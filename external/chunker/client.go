// Package chunker provides a markdown chunking client backed by chisel.
package chunker

import (
	"context"

	"github.com/zoobzio/chisel/markdown"
	"github.com/zoobzio/nestor/api/contracts"
)

// Client wraps chisel's markdown provider to satisfy contracts.Chunker.
type Client struct {
	provider *markdown.Provider
}

// New creates a new markdown chunking client.
func New() *Client {
	return &Client{
		provider: markdown.New(),
	}
}

// Chunk parses markdown content into semantic chunks.
// Returns an empty slice for empty content.
func (c *Client) Chunk(ctx context.Context, content []byte) ([]contracts.Chunk, error) {
	if len(content) == 0 {
		return []contracts.Chunk{}, nil
	}

	raw, err := c.provider.Chunk(ctx, "", content)
	if err != nil {
		return nil, err
	}

	chunks := make([]contracts.Chunk, len(raw))
	for i, r := range raw {
		chunks[i] = contracts.Chunk{
			Kind:      string(r.Kind),
			Symbol:    r.Symbol,
			Context:   r.Context,
			Content:   r.Content,
			StartLine: r.StartLine,
			EndLine:   r.EndLine,
		}
	}

	return chunks, nil
}
