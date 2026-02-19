// Package models contains domain model types.
package models

import (
	"time"

	"github.com/zoobzio/check"
)

// ChunkKind represents the type of chunk parsed from markdown.
type ChunkKind string

const (
	ChunkKindSection ChunkKind = "section"
)

// Chunk represents a segment of a parent Memory's markdown content.
// Chunks are used for vector search; queries run against chunks but
// results resolve back to the parent Memory.
type Chunk struct {
	ID        string    `json:"id" db:"id" constraints:"primarykey" description:"Chunk UUID" example:"550e8400-e29b-41d4-a716-446655440003"`
	MemoryID  string    `json:"memory_id" db:"memory_id" constraints:"notnull" references:"memories(id)" description:"Parent memory UUID" example:"550e8400-e29b-41d4-a716-446655440002"`
	Content   string    `json:"content" db:"content" constraints:"notnull" type:"text" description:"Segment of parent memory markdown content" example:"The user prefers concise responses."`
	Embedding []float32 `json:"-" db:"embedding" type:"vector(1536)" description:"OpenAI ada-002 embedding vector"`
	Sequence  int       `json:"sequence" db:"sequence" constraints:"notnull" description:"Order of this chunk within the parent memory" example:"0"`
	Kind      ChunkKind `json:"kind" db:"kind" constraints:"notnull" default:"'section'" description:"Chunk type from chisel parser" example:"section"`
	Symbol    *string   `json:"symbol,omitempty" db:"symbol" description:"Header text or identifier for this chunk" example:"Installation"`
	Context   []string  `json:"context,omitempty" db:"context" type:"text[]" description:"Parent header chain for hierarchical context"`
	StartLine int       `json:"start_line" db:"start_line" constraints:"notnull" default:"0" description:"Starting line number in source" example:"1"`
	EndLine   int       `json:"end_line" db:"end_line" constraints:"notnull" default:"0" description:"Ending line number in source" example:"15"`
	CreatedAt time.Time `json:"created_at" db:"created_at" default:"now()" description:"Chunk creation time"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at" default:"now()" description:"Last update time"`
}

// Clone returns a deep copy of Chunk.
func (c Chunk) Clone() Chunk {
	cp := c
	if c.Embedding != nil {
		cp.Embedding = make([]float32, len(c.Embedding))
		copy(cp.Embedding, c.Embedding)
	}
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

// Validate validates the Chunk model.
func (c Chunk) Validate() error {
	return check.All(
		check.Str(c.ID, "id").Required().V(),
		check.Str(c.MemoryID, "memory_id").Required().V(),
		check.Str(c.Content, "content").Required().V(),
		check.Int(c.Sequence, "sequence").Min(0).V(),
	).Err()
}
