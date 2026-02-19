// Package main is the entry point for the application.
package main

import (
	"context"
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/zoobzio/aperture"
	"github.com/zoobzio/astql/postgres"
	"github.com/zoobzio/capitan"
	"github.com/zoobzio/sum"
	"github.com/zoobzio/nestor/api/contracts"
	"github.com/zoobzio/nestor/api/handlers"
	"github.com/zoobzio/nestor/api/ingest"
	"github.com/zoobzio/nestor/config"
	"github.com/zoobzio/nestor/events"
	"github.com/zoobzio/nestor/external/chunker"
	"github.com/zoobzio/nestor/external/embedding"
	intotel "github.com/zoobzio/nestor/internal/otel"
	"github.com/zoobzio/nestor/stores"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	log.Println("starting...")
	ctx := context.Background()

	// Initialize sum service and registry.
	svc := sum.New()
	k := sum.Start()

	// =========================================================================
	// 1. Load Configuration
	// =========================================================================

	if err := sum.Config[config.App](ctx, k, nil); err != nil {
		return fmt.Errorf("failed to load app config: %w", err)
	}
	if err := sum.Config[config.Database](ctx, k, nil); err != nil {
		return fmt.Errorf("failed to load database config: %w", err)
	}
	if err := sum.Config[config.Embedding](ctx, k, nil); err != nil {
		return fmt.Errorf("failed to load embedding config: %w", err)
	}

	// =========================================================================
	// 2. Connect to Infrastructure
	// =========================================================================

	dbCfg := sum.MustUse[config.Database](ctx)
	db, err := sqlx.Connect("postgres", dbCfg.DSN)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}
	defer func() { _ = db.Close() }()
	log.Println("database connected")
	capitan.Emit(ctx, events.StartupDatabaseConnected)

	renderer := postgres.New()

	// =========================================================================
	// 3. Create and Register Stores
	// =========================================================================

	allStores, err := stores.New(db, renderer)
	if err != nil {
		return fmt.Errorf("failed to create stores: %w", err)
	}

	sum.Register[contracts.Agents](k, allStores.Agents)
	sum.Register[contracts.Memories](k, allStores.Memories)
	sum.Register[contracts.Chunks](k, allStores.Chunks)
	sum.Register[contracts.IngestJobs](k, allStores.IngestJobs)
	sum.Register[contracts.FilterGroups](k, allStores.FilterGroups)

	// =========================================================================
	// 4. Register External Clients
	// =========================================================================

	chunkerClient := chunker.New()
	sum.Register[contracts.Chunker](k, chunkerClient)

	embedCfg := sum.MustUse[config.Embedding](ctx)
	embedClient, err := embedding.NewClient(embedCfg.Provider, embedCfg.Model, embedCfg.APIKey, embedCfg.Dimensions)
	if err != nil {
		return fmt.Errorf("failed to create embedding client: %w", err)
	}
	sum.Register[contracts.Embedder](k, embedClient)

	// =========================================================================
	// 5. Freeze Registry
	// =========================================================================

	sum.Freeze(k)
	capitan.Emit(ctx, events.StartupServicesReady)

	// =========================================================================
	// 6. Initialize Observability (OTEL + Aperture)
	// =========================================================================

	otelEndpoint := "localhost:4318"
	serviceName := "nestor"

	otelProviders, err := intotel.New(ctx, intotel.Config{
		Endpoint:    otelEndpoint,
		ServiceName: serviceName,
	})
	if err != nil {
		return fmt.Errorf("failed to create otel providers: %w", err)
	}
	defer func() { _ = otelProviders.Shutdown(ctx) }()
	log.Println("observability initialized")
	capitan.Emit(ctx, events.StartupOTELReady)

	ap, err := aperture.New(
		capitan.Default(),
		otelProviders.Log,
		otelProviders.Metric,
		otelProviders.Trace,
	)
	if err != nil {
		return fmt.Errorf("failed to create aperture: %w", err)
	}
	defer ap.Close()
	capitan.Emit(ctx, events.StartupApertureReady)

	// =========================================================================
	// 7. Start Pipeline Worker
	// =========================================================================

	worker := ingest.NewWorker()
	worker.Start(ctx)
	defer func() { _ = worker.Stop() }()

	// =========================================================================
	// 8. Register Handlers and Run
	// =========================================================================

	svc.Handle(handlers.All()...)

	appCfg := sum.MustUse[config.App](ctx)
	capitan.Emit(ctx, events.StartupServerListening, events.StartupPortKey.Field(appCfg.Port))
	log.Printf("starting server on port %d...", appCfg.Port)
	return svc.Run("", appCfg.Port)
}
