package stores

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/zoobzio/astql"
	"github.com/zoobzio/sum"
	"github.com/zoobzio/nestor/models"
)

// Memories provides database access for memory records.
type Memories struct {
	*sum.Database[models.Memory]
}

// NewMemories creates a new memories store.
func NewMemories(db *sqlx.DB, renderer astql.Renderer) (*Memories, error) {
	database, err := sum.NewDatabase[models.Memory](db, "memories", renderer)
	if err != nil {
		return nil, err
	}
	return &Memories{Database: database}, nil
}

// ListByAgentID retrieves all memories belonging to an agent, ordered by creation time.
func (s *Memories) ListByAgentID(ctx context.Context, agentID string) ([]*models.Memory, error) {
	return s.Query().
		Where("agent_id", "=", "agent_id").
		OrderBy("created_at", "ASC").
		Exec(ctx, map[string]any{"agent_id": agentID})
}

// GetByCorrelationID retrieves a memory by its correlation ID.
func (s *Memories) GetByCorrelationID(ctx context.Context, correlationID string) (*models.Memory, error) {
	return s.Select().
		Where("correlation_id", "=", "correlation_id").
		Exec(ctx, map[string]any{"correlation_id": correlationID})
}
