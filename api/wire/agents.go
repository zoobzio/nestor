// Package wire defines request and response types for the public API surface.
package wire

import (
	"time"

	"github.com/zoobzio/check"
)

// AgentResponse is the public API response for agent data.
type AgentResponse struct {
	ID        string    `json:"id" description:"Agent UUID" example:"550e8400-e29b-41d4-a716-446655440001"`
	UserID    string    `json:"user_id" description:"Owning user UUID" example:"550e8400-e29b-41d4-a716-446655440000"`
	Name      string    `json:"name" description:"Agent name" example:"Research Assistant"`
	CreatedAt time.Time `json:"created_at" description:"Agent creation time"`
	UpdatedAt time.Time `json:"updated_at" description:"Last update time"`
}

// Clone returns a shallow copy. All fields are value types.
func (a AgentResponse) Clone() AgentResponse {
	return a
}

// AgentListResponse is the public API response for listing agents.
type AgentListResponse struct {
	Agents []AgentResponse `json:"agents" description:"List of agents"`
}

// Clone returns a deep copy.
func (a AgentListResponse) Clone() AgentListResponse {
	r := a
	if a.Agents != nil {
		r.Agents = make([]AgentResponse, len(a.Agents))
		copy(r.Agents, a.Agents)
	}
	return r
}

// CreateAgentRequest is the request body for creating an agent.
type CreateAgentRequest struct {
	Name string `json:"name" description:"Agent name" example:"Research Assistant"`
}

// Validate validates the request.
func (r *CreateAgentRequest) Validate() error {
	return check.All(
		check.Str(r.Name, "name").Required().MaxLen(255).V(),
	).Err()
}

// Clone returns a shallow copy.
func (r CreateAgentRequest) Clone() CreateAgentRequest {
	return r
}

// UpdateAgentRequest is the request body for updating an agent.
type UpdateAgentRequest struct {
	Name string `json:"name" description:"Agent name" example:"Research Assistant"`
}

// Validate validates the request.
func (r *UpdateAgentRequest) Validate() error {
	return check.All(
		check.Str(r.Name, "name").Required().MaxLen(255).V(),
	).Err()
}

// Clone returns a shallow copy.
func (r UpdateAgentRequest) Clone() UpdateAgentRequest {
	return r
}
