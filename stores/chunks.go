package stores

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/zoobzio/astql"
	"github.com/zoobzio/sum"
	"github.com/zoobzio/nestor/models"
)

// Chunks provides database access for chunk records.
type Chunks struct {
	*sum.Database[models.Chunk]
}

// NewChunks creates a new chunks store.
func NewChunks(db *sqlx.DB, renderer astql.Renderer) (*Chunks, error) {
	database, err := sum.NewDatabase[models.Chunk](db, "chunks", renderer)
	if err != nil {
		return nil, err
	}
	return &Chunks{Database: database}, nil
}

// ListByMemoryID retrieves all chunks belonging to a memory, ordered by sequence.
func (s *Chunks) ListByMemoryID(ctx context.Context, memoryID string) ([]*models.Chunk, error) {
	return s.Query().
		Where("memory_id", "=", "memory_id").
		OrderBy("sequence", "ASC").
		Exec(ctx, map[string]any{"memory_id": memoryID})
}

// SearchByEmbedding retrieves chunks ordered by cosine distance to the given
// embedding vector, returning the most similar results first. Limit controls
// the maximum number of results returned.
func (s *Chunks) SearchByEmbedding(ctx context.Context, embedding []float32, limit int) ([]*models.Chunk, error) {
	return s.Query().
		OrderByExpr("embedding", "<=>", "query_vec", "ASC").
		Limit(limit).
		Exec(ctx, map[string]any{"query_vec": embedding})
}

// SearchByEmbeddingForMemory retrieves chunks for a specific memory ordered by
// cosine distance, useful for re-ranking within a single memory's chunks.
func (s *Chunks) SearchByEmbeddingForMemory(ctx context.Context, memoryID string, embedding []float32, limit int) ([]*models.Chunk, error) {
	return s.Query().
		Where("memory_id", "=", "memory_id").
		OrderByExpr("embedding", "<=>", "query_vec", "ASC").
		Limit(limit).
		Exec(ctx, map[string]any{
			"memory_id": memoryID,
			"query_vec": embedding,
		})
}
