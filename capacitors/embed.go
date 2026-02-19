package capacitors

import (
	"context"
	"log"
	"time"

	"github.com/zoobzio/check"
	"github.com/zoobzio/flux"
)

// EmbedConfig holds operational settings for the embedding pipeline.
// These values can be changed at runtime without restarting the application.
type EmbedConfig struct {
	Workers   int           `json:"workers"`    // Concurrent embedding batches
	BatchSize int           `json:"batch_size"` // Chunks per batch
	Timeout   time.Duration `json:"timeout"`    // Per-batch timeout
}

// Validate checks that EmbedConfig values are within acceptable bounds.
func (c EmbedConfig) Validate() error {
	return check.All(
		check.Int(c.Workers, "workers").Positive().Max(64).V(),
		check.Int(c.BatchSize, "batch_size").Positive().Max(1024).V(),
		check.DurationPositive(c.Timeout, "timeout"),
		check.DurationMax(c.Timeout, 5*time.Minute, "timeout"),
	).Err()
}

// DefaultEmbedConfig returns sensible defaults for the embedding pipeline.
func DefaultEmbedConfig() EmbedConfig {
	return EmbedConfig{
		Workers:   4,
		BatchSize: 128,
		Timeout:   30 * time.Second,
	}
}

// applyEmbedConfig applies updated embedding configuration to the running system.
// Wire this to the embedding pipeline once it is implemented.
func applyEmbedConfig(cfg EmbedConfig) {
	// TODO: wire to embedding pipeline
	// e.g., embed.SetWorkers(cfg.Workers)
	//       embed.SetBatchSize(cfg.BatchSize)
	//       embed.SetTimeout(cfg.Timeout)
	log.Printf("embed config updated: workers=%d batch_size=%d timeout=%s", cfg.Workers, cfg.BatchSize, cfg.Timeout)
}

// InitEmbed initializes the EmbedConfig capacitor.
// It applies defaults immediately, then watches the provided watcher for live changes.
func InitEmbed(ctx context.Context, watcher flux.Watcher) error {
	applyEmbedConfig(DefaultEmbedConfig())

	c := flux.New[EmbedConfig](
		watcher,
		func(_ context.Context, _, curr EmbedConfig) error {
			applyEmbedConfig(curr)
			return nil
		},
	)

	go func() {
		if err := c.Start(ctx); err != nil {
			log.Printf("embed capacitor error: %v", err)
		}
	}()

	return nil
}
