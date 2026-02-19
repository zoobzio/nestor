package stores

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/zoobzio/astql"
	"github.com/zoobzio/sum"
	"github.com/zoobzio/nestor/models"
)

// Agents provides database access for agent records.
type Agents struct {
	*sum.Database[models.Agent]
}

// NewAgents creates a new agents store.
func NewAgents(db *sqlx.DB, renderer astql.Renderer) (*Agents, error) {
	database, err := sum.NewDatabase[models.Agent](db, "agents", renderer)
	if err != nil {
		return nil, err
	}
	return &Agents{Database: database}, nil
}

// ListByUserID retrieves all agents belonging to a user.
func (s *Agents) ListByUserID(ctx context.Context, userID string) ([]*models.Agent, error) {
	return s.Query().
		Where("user_id", "=", "user_id").
		OrderBy("created_at", "ASC").
		Exec(ctx, map[string]any{"user_id": userID})
}
