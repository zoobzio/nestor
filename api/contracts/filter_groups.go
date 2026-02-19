// Package contracts defines interface boundaries for the public API surface.
package contracts

import (
	"context"

	"github.com/zoobzio/nestor/models"
)

// FilterGroups defines the contract for filter group operations on the public API surface.
type FilterGroups interface {
	// Get retrieves a filter group by primary key.
	Get(ctx context.Context, key string) (*models.FilterGroup, error)
	// Set creates or updates a filter group.
	Set(ctx context.Context, key string, filterGroup *models.FilterGroup) error
	// Delete removes a filter group by primary key.
	Delete(ctx context.Context, key string) error
	// ListByAgentID retrieves all filter groups belonging to an agent.
	ListByAgentID(ctx context.Context, agentID string) ([]*models.FilterGroup, error)
}
