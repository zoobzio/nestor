// Package contracts defines interface boundaries for the public API surface.
package contracts

import (
	"context"

	"github.com/zoobzio/nestor/models"
)

// Memories defines the contract for memory operations on the public API surface.
type Memories interface {
	// Get retrieves a memory by primary key.
	Get(ctx context.Context, key string) (*models.Memory, error)
	// Set creates or updates a memory.
	Set(ctx context.Context, key string, memory *models.Memory) error
	// Delete removes a memory by primary key.
	Delete(ctx context.Context, key string) error
	// ListByAgentID retrieves all memories belonging to an agent.
	ListByAgentID(ctx context.Context, agentID string) ([]*models.Memory, error)
	// GetByCorrelationID retrieves a memory by its correlation ID.
	GetByCorrelationID(ctx context.Context, correlationID string) (*models.Memory, error)
}
