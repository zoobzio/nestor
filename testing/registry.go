//go:build testing

// Package testing provides test infrastructure: fixtures, mocks, and helpers.
package testing

import (
	"testing"

	"github.com/zoobzio/sum"
	"github.com/zoobzio/nestor/api/contracts"
)

// RegistryOption configures the sum registry for a test.
type RegistryOption func(k sum.Key)

// WithAgents registers a mock Agents implementation.
func WithAgents(a contracts.Agents) RegistryOption {
	return func(k sum.Key) {
		sum.Register[contracts.Agents](k, a)
	}
}

// WithMemories registers a mock Memories implementation.
func WithMemories(m contracts.Memories) RegistryOption {
	return func(k sum.Key) {
		sum.Register[contracts.Memories](k, m)
	}
}

// WithChunks registers a mock Chunks implementation.
func WithChunks(c contracts.Chunks) RegistryOption {
	return func(k sum.Key) {
		sum.Register[contracts.Chunks](k, c)
	}
}

// WithIngestJobs registers a mock IngestJobs implementation.
func WithIngestJobs(j contracts.IngestJobs) RegistryOption {
	return func(k sum.Key) {
		sum.Register[contracts.IngestJobs](k, j)
	}
}

// WithChunker registers a mock Chunker implementation.
func WithChunker(c contracts.Chunker) RegistryOption {
	return func(k sum.Key) {
		sum.Register[contracts.Chunker](k, c)
	}
}

// WithEmbedder registers a mock Embedder implementation.
func WithEmbedder(e contracts.Embedder) RegistryOption {
	return func(k sum.Key) {
		sum.Register[contracts.Embedder](k, e)
	}
}

// SetupRegistry initialises a fresh sum registry, applies all options, freezes
// it, and registers cleanup. It returns a test context suitable for use
// in handler tests.
func SetupRegistry(t *testing.T, opts ...RegistryOption) {
	t.Helper()
	sum.Reset()
	k := sum.Start()
	for _, opt := range opts {
		opt(k)
	}
	sum.Freeze(k)
	t.Cleanup(sum.Reset)
}
