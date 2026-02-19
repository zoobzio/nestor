// Package transformers provides pure functions for mapping between models and wire types
// on the public API surface.
package transformers

import (
	"github.com/zoobzio/nestor/api/wire"
	"github.com/zoobzio/nestor/models"
)

// ChunkToResponse transforms a Chunk model to a public API response.
// Embeddings are excluded from the response.
func ChunkToResponse(c *models.Chunk) wire.ChunkResponse {
	return wire.ChunkResponse{
		ID:        c.ID,
		MemoryID:  c.MemoryID,
		Content:   c.Content,
		Sequence:  c.Sequence,
		Kind:      string(c.Kind),
		Symbol:    c.Symbol,
		Context:   c.Context,
		StartLine: c.StartLine,
		EndLine:   c.EndLine,
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt,
	}
}

// ChunksToList transforms a slice of Chunk models to a list response.
func ChunksToList(chunks []*models.Chunk) wire.ChunkListResponse {
	resp := wire.ChunkListResponse{
		Chunks: make([]wire.ChunkResponse, len(chunks)),
	}
	for i, c := range chunks {
		resp.Chunks[i] = ChunkToResponse(c)
	}
	return resp
}

// ChunksToSearchResponse transforms a slice of Chunk models to a search response.
func ChunksToSearchResponse(chunks []*models.Chunk) wire.ChunkSearchResponse {
	resp := wire.ChunkSearchResponse{
		Results: make([]wire.ChunkResponse, len(chunks)),
	}
	for i, c := range chunks {
		resp.Results[i] = ChunkToResponse(c)
	}
	return resp
}

// ApplyCreateChunk applies a CreateChunkRequest to a Chunk model.
func ApplyCreateChunk(req wire.CreateChunkRequest, c *models.Chunk) {
	c.Content = req.Content
	c.Sequence = req.Sequence
	if req.Embedding != nil {
		c.Embedding = make([]float32, len(req.Embedding))
		copy(c.Embedding, req.Embedding)
	}
}
