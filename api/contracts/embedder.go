// Package contracts defines interface boundaries for the public API surface.
package contracts

import "context"

// Embedder generates vector embeddings for text.
type Embedder interface {
	// Embed generates embeddings for the given texts.
	// Returns one embedding vector per input text.
	// Embedding dimension depends on the underlying model (e.g., 1536 for OpenAI ada-002).
	Embed(ctx context.Context, texts []string) ([][]float32, error)
}
