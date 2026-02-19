package handlers

import (
	"crypto/rand"
	"fmt"
	"time"

	"github.com/zoobzio/rocco"
	"github.com/zoobzio/sum"
	"github.com/zoobzio/nestor/api/contracts"
	"github.com/zoobzio/nestor/api/transformers"
	"github.com/zoobzio/nestor/api/wire"
	"github.com/zoobzio/nestor/models"
)

// newUUID generates a random UUID v4.
func newUUID() string {
	var b [16]byte
	_, _ = rand.Read(b[:])
	b[6] = (b[6] & 0x0f) | 0x40
	b[8] = (b[8] & 0x3f) | 0x80
	return fmt.Sprintf("%08x-%04x-%04x-%04x-%012x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
}

// ListAgents returns all agents for the authenticated user.
var ListAgents = rocco.GET("/agents", func(req *rocco.Request[rocco.NoBody]) (wire.AgentListResponse, error) {
	agents := sum.MustUse[contracts.Agents](req.Context)

	userID := req.Identity.ID()

	list, err := agents.ListByUserID(req.Context, userID)
	if err != nil {
		return wire.AgentListResponse{}, err
	}

	return transformers.AgentsToList(list), nil
}).WithSummary("List agents").
	WithDescription("Returns all agents belonging to the authenticated user.").
	WithTags("Agents").
	WithAuthentication()

// CreateAgent creates a new agent for the authenticated user.
var CreateAgent = rocco.POST("/agents", func(req *rocco.Request[wire.CreateAgentRequest]) (wire.AgentResponse, error) {
	agents := sum.MustUse[contracts.Agents](req.Context)

	agent := &models.Agent{
		ID:        newUUID(),
		UserID:    req.Identity.ID(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	transformers.ApplyCreateAgent(req.Body, agent)

	if err := agents.Set(req.Context, agent.ID, agent); err != nil {
		return wire.AgentResponse{}, err
	}

	return transformers.AgentToResponse(agent), nil
}).WithSummary("Create agent").
	WithDescription("Creates a new agent for the authenticated user.").
	WithTags("Agents").
	WithAuthentication().
	WithSuccessStatus(201)

// GetAgent returns a single agent by ID.
var GetAgent = rocco.GET("/agents/{id}", func(req *rocco.Request[rocco.NoBody]) (wire.AgentResponse, error) {
	agents := sum.MustUse[contracts.Agents](req.Context)

	id := req.Params.Path["id"]

	agent, err := agents.Get(req.Context, id)
	if err != nil {
		return wire.AgentResponse{}, ErrAgentNotFound
	}

	return transformers.AgentToResponse(agent), nil
}).WithPathParams("id").
	WithSummary("Get agent").
	WithDescription("Returns an agent by ID.").
	WithTags("Agents").
	WithErrors(ErrAgentNotFound).
	WithAuthentication()

// UpdateAgent updates an existing agent.
var UpdateAgent = rocco.PATCH("/agents/{id}", func(req *rocco.Request[wire.UpdateAgentRequest]) (wire.AgentResponse, error) {
	agents := sum.MustUse[contracts.Agents](req.Context)

	id := req.Params.Path["id"]

	agent, err := agents.Get(req.Context, id)
	if err != nil {
		return wire.AgentResponse{}, ErrAgentNotFound
	}

	transformers.ApplyUpdateAgent(req.Body, agent)
	agent.UpdatedAt = time.Now()

	if err := agents.Set(req.Context, id, agent); err != nil {
		return wire.AgentResponse{}, err
	}

	return transformers.AgentToResponse(agent), nil
}).WithPathParams("id").
	WithSummary("Update agent").
	WithDescription("Updates an agent's details.").
	WithTags("Agents").
	WithErrors(ErrAgentNotFound).
	WithAuthentication()

// DeleteAgent removes an agent by ID.
var DeleteAgent = rocco.DELETE("/agents/{id}", func(req *rocco.Request[rocco.NoBody]) (rocco.NoBody, error) {
	agents := sum.MustUse[contracts.Agents](req.Context)

	id := req.Params.Path["id"]

	if err := agents.Delete(req.Context, id); err != nil {
		return rocco.NoBody{}, ErrAgentNotFound
	}

	return rocco.NoBody{}, nil
}).WithPathParams("id").
	WithSummary("Delete agent").
	WithDescription("Deletes an agent by ID.").
	WithTags("Agents").
	WithErrors(ErrAgentNotFound).
	WithAuthentication().
	WithSuccessStatus(204)
