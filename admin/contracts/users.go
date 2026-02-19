// Package contracts defines interface boundaries for the admin API surface.
package contracts

import (
	"context"

	"github.com/zoobzio/nestor/models"
)

// Users defines the contract for user operations on the admin API surface.
type Users interface {
	// Get retrieves a user by primary key.
	Get(ctx context.Context, key string) (*models.User, error)
	// Set creates or updates a user.
	Set(ctx context.Context, key string, user *models.User) error
	// Delete removes a user by primary key.
	Delete(ctx context.Context, key string) error
	// List retrieves all users.
	List(ctx context.Context) ([]*models.User, error)
}
