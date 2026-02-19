package stores

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/zoobzio/astql"
	"github.com/zoobzio/sum"
	"github.com/zoobzio/nestor/models"
)

// FilterGroups provides database access for filter group records.
type FilterGroups struct {
	*sum.Database[models.FilterGroup]
}

// NewFilterGroups creates a new filter groups store.
func NewFilterGroups(db *sqlx.DB, renderer astql.Renderer) (*FilterGroups, error) {
	database, err := sum.NewDatabase[models.FilterGroup](db, "filter_groups", renderer)
	if err != nil {
		return nil, err
	}
	return &FilterGroups{Database: database}, nil
}

// ListByAgentID retrieves all filter groups belonging to an agent, ordered by creation time.
func (s *FilterGroups) ListByAgentID(ctx context.Context, agentID string) ([]*models.FilterGroup, error) {
	return s.Query().
		Where("agent_id", "=", "agent_id").
		OrderBy("created_at", "ASC").
		Exec(ctx, map[string]any{"agent_id": agentID})
}
