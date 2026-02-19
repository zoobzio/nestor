// Package capacitors provides hot-reload runtime configuration via flux capacitors.
package capacitors

import (
	"context"
	"log"

	"github.com/zoobzio/check"
	"github.com/zoobzio/flux"
)

// IngestConfig holds operational settings for the ingest pipeline.
// These values can be changed at runtime without restarting the application.
type IngestConfig struct {
	MaxContentSize int64 `json:"max_content_size"` // Max memory content in bytes
	WorkerCount    int   `json:"worker_count"`     // Pipeline workers
}

// Validate checks that IngestConfig values are within acceptable bounds.
func (c IngestConfig) Validate() error {
	return check.All(
		check.Int(int(c.MaxContentSize), "max_content_size").Positive().V(),
		check.Int(c.WorkerCount, "worker_count").Positive().Max(256).V(),
	).Err()
}

// DefaultIngestConfig returns sensible defaults for the ingest pipeline.
func DefaultIngestConfig() IngestConfig {
	return IngestConfig{
		MaxContentSize: 10 * 1024 * 1024, // 10MB
		WorkerCount:    4,
	}
}

// applyIngestConfig applies updated ingest configuration to the running system.
// Wire this to the ingest pipeline once it is implemented.
func applyIngestConfig(cfg IngestConfig) {
	// TODO: wire to ingest pipeline
	// e.g., ingest.SetWorkerCount(cfg.WorkerCount)
	//       ingest.SetMaxContentSize(cfg.MaxContentSize)
	log.Printf("ingest config updated: max_content_size=%d worker_count=%d", cfg.MaxContentSize, cfg.WorkerCount)
}

// InitIngest initializes the IngestConfig capacitor.
// It applies defaults immediately, then watches the provided watcher for live changes.
func InitIngest(ctx context.Context, watcher flux.Watcher) error {
	applyIngestConfig(DefaultIngestConfig())

	c := flux.New[IngestConfig](
		watcher,
		func(_ context.Context, _, curr IngestConfig) error {
			applyIngestConfig(curr)
			return nil
		},
	)

	go func() {
		if err := c.Start(ctx); err != nil {
			log.Printf("ingest capacitor error: %v", err)
		}
	}()

	return nil
}
