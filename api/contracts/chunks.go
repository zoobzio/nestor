// Package contracts defines interface boundaries for the public API surface.
package contracts

import (
	"context"

	"github.com/zoobzio/nestor/models"
)

// Chunks defines the contract for chunk operations on the public API surface.
type Chunks interface {
	// Get retrieves a chunk by primary key.
	Get(ctx context.Context, key string) (*models.Chunk, error)
	// Set creates or updates a chunk.
	Set(ctx context.Context, key string, chunk *models.Chunk) error
	// Delete removes a chunk by primary key.
	Delete(ctx context.Context, key string) error
	// ListByMemoryID retrieves all chunks belonging to a memory, ordered by sequence.
	ListByMemoryID(ctx context.Context, memoryID string) ([]*models.Chunk, error)
	// SearchByEmbedding retrieves chunks ordered by cosine distance to the given
	// embedding vector, most similar first.
	SearchByEmbedding(ctx context.Context, embedding []float32, limit int) ([]*models.Chunk, error)
}
