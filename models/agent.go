// Package models contains domain model types.
package models

import (
	"time"

	"github.com/zoobzio/check"
)

// Agent represents an AI agent belonging to a user.
// Agents are scoped to a single user and own memories.
type Agent struct {
	ID        string    `json:"id" db:"id" constraints:"primarykey" description:"Agent UUID" example:"550e8400-e29b-41d4-a716-446655440001"`
	UserID    string    `json:"user_id" db:"user_id" constraints:"notnull" references:"users(id)" description:"Owning user UUID" example:"550e8400-e29b-41d4-a716-446655440000"`
	Name      string    `json:"name" db:"name" constraints:"notnull" description:"Agent name" example:"Research Assistant"`
	CreatedAt time.Time `json:"created_at" db:"created_at" default:"now()" description:"Agent creation time"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at" default:"now()" description:"Last update time"`
}

// Clone returns a shallow copy of Agent. All fields are value types.
func (a Agent) Clone() Agent {
	return a
}

// Validate validates the Agent model.
func (a Agent) Validate() error {
	return check.All(
		check.Str(a.ID, "id").Required().V(),
		check.Str(a.UserID, "user_id").Required().V(),
		check.Str(a.Name, "name").Required().MaxLen(255).V(),
	).Err()
}
