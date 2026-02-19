// Package main is the entry point for the admin API.
package main

import (
	"context"
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/zoobzio/astql/postgres"
	"github.com/zoobzio/capitan"
	"github.com/zoobzio/sum"
	admincontracts "github.com/zoobzio/nestor/admin/contracts"
	adminhandlers "github.com/zoobzio/nestor/admin/handlers"
	"github.com/zoobzio/nestor/config"
	"github.com/zoobzio/nestor/events"
	"github.com/zoobzio/nestor/stores"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	log.Println("starting admin...")
	ctx := context.Background()

	// Initialize sum service and registry.
	svc := sum.New()
	k := sum.Start()

	// =========================================================================
	// 1. Load Configuration
	// =========================================================================

	if err := sum.Config[config.Admin](ctx, k, nil); err != nil {
		return fmt.Errorf("failed to load admin config: %w", err)
	}
	if err := sum.Config[config.Database](ctx, k, nil); err != nil {
		return fmt.Errorf("failed to load database config: %w", err)
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

	sum.Register[admincontracts.Users](k, allStores.Users)

	// =========================================================================
	// 4. Freeze Registry
	// =========================================================================

	sum.Freeze(k)
	capitan.Emit(ctx, events.StartupServicesReady)

	// =========================================================================
	// 5. Register Handlers and Run
	// =========================================================================

	svc.Handle(adminhandlers.All()...)

	adminCfg := sum.MustUse[config.Admin](ctx)
	capitan.Emit(ctx, events.StartupServerListening, events.StartupPortKey.Field(adminCfg.Port))
	log.Printf("starting admin server on port %d...", adminCfg.Port)

	return svc.Run("", adminCfg.Port)
}
