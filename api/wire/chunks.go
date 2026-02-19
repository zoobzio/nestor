// Package wire defines request and response types for the public API surface.
package wire

import (
	"time"

	"github.com/zoobzio/check"
)

// ChunkResponse is the public API response for chunk data.
// Embeddings are intentionally omitted from the response.
type ChunkResponse struct {
	ID        string    `json:"id" description:"Chunk UUID" example:"550e8400-e29b-41d4-a716-446655440003"`
	MemoryID  string    `json:"memory_id" description:"Parent memory UUID" example:"550e8400-e29b-41d4-a716-446655440002"`
	Content   string    `json:"content" description:"Segment of parent memory markdown content" example:"The user prefers concise responses."`
	Sequence  int       `json:"sequence" description:"Order of this chunk within the parent memory" example:"0"`
	Kind      string    `json:"kind" description:"Chunk type from chisel parser" example:"section"`
	Symbol    *string   `json:"symbol,omitempty" description:"Header text or identifier for this chunk" example:"Installation"`
	Context   []string  `json:"context,omitempty" description:"Parent header chain for hierarchical context"`
	StartLine int       `json:"start_line" description:"Starting line number in source" example:"1"`
	EndLine   int       `json:"end_line" description:"Ending line number in source" example:"15"`
	CreatedAt time.Time `json:"created_at" description:"Chunk creation time"`
	UpdatedAt time.Time `json:"updated_at" description:"Last update time"`
}

// Clone returns a deep copy of ChunkResponse.
func (c ChunkResponse) Clone() ChunkResponse {
	cp := c
	if c.Symbol != nil {
		s := *c.Symbol
		cp.Symbol = &s
	}
	if c.Context != nil {
		cp.Context = make([]string, len(c.Context))
		copy(cp.Context, c.Context)
	}
	return cp
}

// ChunkListResponse is the public API response for listing chunks.
type ChunkListResponse struct {
	Chunks []ChunkResponse `json:"chunks" description:"List of chunks ordered by sequence"`
}

// Clone returns a deep copy of ChunkListResponse.
func (c ChunkListResponse) Clone() ChunkListResponse {
	cp := c
	if c.Chunks != nil {
		cp.Chunks = make([]ChunkResponse, len(c.Chunks))
		copy(cp.Chunks, c.Chunks)
	}
	return cp
}

// ChunkSearchRequest is the request body for searching chunks by embedding.
type ChunkSearchRequest struct {
	Embedding []float32 `json:"embedding" description:"Query embedding vector (1536 dimensions for OpenAI ada-002)"`
	Limit     int       `json:"limit" description:"Maximum number of results to return" example:"10"`
}

// Validate validates the request.
func (r *ChunkSearchRequest) Validate() error {
	return check.All(
		check.Int(r.Limit, "limit").Positive().V(),
	).Err()
}

// Clone returns a deep copy of ChunkSearchRequest.
func (r ChunkSearchRequest) Clone() ChunkSearchRequest {
	c := r
	if r.Embedding != nil {
		c.Embedding = make([]float32, len(r.Embedding))
		copy(c.Embedding, r.Embedding)
	}
	return c
}

// ChunkSearchResponse is the public API response for chunk similarity search.
// Results are ordered by cosine similarity, most similar first.
type ChunkSearchResponse struct {
	Results []ChunkResponse `json:"results" description:"Chunks ordered by similarity to the query embedding"`
}

// Clone returns a deep copy of ChunkSearchResponse.
func (c ChunkSearchResponse) Clone() ChunkSearchResponse {
	cp := c
	if c.Results != nil {
		cp.Results = make([]ChunkResponse, len(c.Results))
		copy(cp.Results, c.Results)
	}
	return cp
}

// CreateChunkRequest is the request body for creating a chunk.
type CreateChunkRequest struct {
	Content   string    `json:"content" description:"Segment of parent memory markdown content" example:"The user prefers concise responses."`
	Embedding []float32 `json:"embedding" description:"Embedding vector for this chunk (1536 dimensions for OpenAI ada-002)"`
	Sequence  int       `json:"sequence" description:"Order of this chunk within the parent memory" example:"0"`
}

// Validate validates the request.
func (r *CreateChunkRequest) Validate() error {
	return check.All(
		check.Str(r.Content, "content").Required().V(),
	).Err()
}

// Clone returns a deep copy of CreateChunkRequest.
func (r CreateChunkRequest) Clone() CreateChunkRequest {
	c := r
	if r.Embedding != nil {
		c.Embedding = make([]float32, len(r.Embedding))
		copy(c.Embedding, r.Embedding)
	}
	return c
}
