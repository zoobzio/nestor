package stores

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/zoobzio/astql"
	"github.com/zoobzio/sum"
	"github.com/zoobzio/nestor/models"
)

// Users provides database access for user records.
type Users struct {
	*sum.Database[models.User]
}

// NewUsers creates a new users store.
func NewUsers(db *sqlx.DB, renderer astql.Renderer) (*Users, error) {
	database, err := sum.NewDatabase[models.User](db, "users", renderer)
	if err != nil {
		return nil, err
	}
	return &Users{Database: database}, nil
}

// List retrieves all users.
func (s *Users) List(ctx context.Context) ([]*models.User, error) {
	return s.Query().
		OrderBy("created_at", "ASC").
		Exec(ctx, nil)
}
