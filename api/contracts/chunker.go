// Package contracts defines interface boundaries for the public API surface.
package contracts

import "context"

// Chunk represents a parsed segment from chisel.
type Chunk struct {
	Kind      string   // e.g., "section"
	Symbol    string   // Header text
	Context   []string // Parent header chain
	Content   string   // Raw text of the section
	StartLine int
	EndLine   int
}

// Chunker parses content into semantic chunks.
type Chunker interface {
	// Chunk parses markdown content into semantic chunks using chisel.
	// Returns chunks with Kind, Symbol, Context, StartLine, EndLine, Content.
	Chunk(ctx context.Context, content []byte) ([]Chunk, error)
}
