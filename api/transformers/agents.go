// Package transformers provides pure functions for mapping between models and wire types
// on the public API surface.
package transformers

import (
	"github.com/zoobzio/nestor/api/wire"
	"github.com/zoobzio/nestor/models"
)

// AgentToResponse transforms an Agent model to a public API response.
func AgentToResponse(a *models.Agent) wire.AgentResponse {
	return wire.AgentResponse{
		ID:        a.ID,
		UserID:    a.UserID,
		Name:      a.Name,
		CreatedAt: a.CreatedAt,
		UpdatedAt: a.UpdatedAt,
	}
}

// AgentsToList transforms a slice of Agent models to a list response.
func AgentsToList(agents []*models.Agent) wire.AgentListResponse {
	resp := wire.AgentListResponse{
		Agents: make([]wire.AgentResponse, len(agents)),
	}
	for i, a := range agents {
		resp.Agents[i] = AgentToResponse(a)
	}
	return resp
}

// ApplyCreateAgent applies a CreateAgentRequest to an Agent model.
func ApplyCreateAgent(req wire.CreateAgentRequest, a *models.Agent) {
	a.Name = req.Name
}

// ApplyUpdateAgent applies an UpdateAgentRequest to an Agent model.
func ApplyUpdateAgent(req wire.UpdateAgentRequest, a *models.Agent) {
	a.Name = req.Name
}
