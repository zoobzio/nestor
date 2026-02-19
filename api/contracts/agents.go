// Package contracts defines interface boundaries for the public API surface.
package contracts

import (
	"context"

	"github.com/zoobzio/nestor/models"
)

// Agents defines the contract for agent operations on the public API surface.
type Agents interface {
	// Get retrieves an agent by primary key.
	Get(ctx context.Context, key string) (*models.Agent, error)
	// Set creates or updates an agent.
	Set(ctx context.Context, key string, agent *models.Agent) error
	// Delete removes an agent by primary key.
	Delete(ctx context.Context, key string) error
	// ListByUserID retrieves all agents belonging to a user.
	ListByUserID(ctx context.Context, userID string) ([]*models.Agent, error)
}
