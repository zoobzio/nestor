package capacitors

import (
	"context"
	"fmt"

	"github.com/zoobzio/flux"
	fluxfile "github.com/zoobzio/flux/file"
)

// Config holds the watcher sources for all capacitors.
type Config struct {
	// IngestConfigPath is the filesystem path to the IngestConfig JSON file.
	// Example: "/etc/nestor/ingest.json"
	IngestConfigPath string

	// EmbedConfigPath is the filesystem path to the EmbedConfig JSON file.
	// Example: "/etc/nestor/embed.json"
	EmbedConfigPath string
}

// Init initializes all capacitors with their respective watchers.
// Call this after infrastructure is connected but before sum.Freeze().
func Init(ctx context.Context, cfg Config) error {
	var ingestWatcher flux.Watcher
	if cfg.IngestConfigPath != "" {
		ingestWatcher = fluxfile.New(cfg.IngestConfigPath)
	}

	if err := InitIngest(ctx, ingestWatcher); err != nil {
		return fmt.Errorf("failed to initialize ingest capacitor: %w", err)
	}

	var embedWatcher flux.Watcher
	if cfg.EmbedConfigPath != "" {
		embedWatcher = fluxfile.New(cfg.EmbedConfigPath)
	}

	if err := InitEmbed(ctx, embedWatcher); err != nil {
		return fmt.Errorf("failed to initialize embed capacitor: %w", err)
	}

	return nil
}
