// Package models contains domain model types.
package models

import (
	"time"

	"github.com/zoobzio/check"
)

// User represents a human signing up for the service.
// Users are managed by admin staff. All agents are scoped to a user.
type User struct {
	ID        string    `json:"id" db:"id" constraints:"primarykey" description:"User UUID" example:"550e8400-e29b-41d4-a716-446655440000"`
	Name      string    `json:"name" db:"name" constraints:"notnull" description:"User name" example:"Jane Doe"`
	CreatedAt time.Time `json:"created_at" db:"created_at" default:"now()" description:"Account creation time"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at" default:"now()" description:"Last update time"`
}

// Clone returns a shallow copy of User. All fields are value types.
func (u User) Clone() User {
	return u
}

// Validate validates the User model.
func (u User) Validate() error {
	return check.All(
		check.Str(u.ID, "id").Required().V(),
		check.Str(u.Name, "name").Required().MaxLen(255).V(),
	).Err()
}
