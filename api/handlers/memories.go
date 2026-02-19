package handlers

import (
	"time"

	"github.com/zoobzio/rocco"
	"github.com/zoobzio/sum"
	"github.com/zoobzio/nestor/api/contracts"
	"github.com/zoobzio/nestor/api/transformers"
	"github.com/zoobzio/nestor/api/wire"
	"github.com/zoobzio/nestor/models"
)

// ListMemories returns all memories for a given agent.
var ListMemories = rocco.GET("/agents/{agent_id}/memories", func(req *rocco.Request[rocco.NoBody]) (wire.MemoryListResponse, error) {
	memories := sum.MustUse[contracts.Memories](req.Context)

	agentID := req.Params.Path["agent_id"]

	list, err := memories.ListByAgentID(req.Context, agentID)
	if err != nil {
		return wire.MemoryListResponse{}, err
	}

	return transformers.MemoriesToList(list), nil
}).WithPathParams("agent_id").
	WithSummary("List memories").
	WithDescription("Returns all memories belonging to the specified agent.").
	WithTags("Memories").
	WithAuthentication()

// CreateMemory creates a new memory for the specified agent.
var CreateMemory = rocco.POST("/agents/{agent_id}/memories", func(req *rocco.Request[wire.CreateMemoryRequest]) (wire.MemoryResponse, error) {
	memories := sum.MustUse[contracts.Memories](req.Context)

	memory := &models.Memory{
		ID:        newUUID(),
		AgentID:   req.Params.Path["agent_id"],
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	transformers.ApplyCreateMemory(req.Body, memory)

	if err := memories.Set(req.Context, memory.ID, memory); err != nil {
		return wire.MemoryResponse{}, err
	}

	return transformers.MemoryToResponse(memory), nil
}).WithPathParams("agent_id").
	WithSummary("Create memory").
	WithDescription("Creates a new memory for the specified agent.").
	WithTags("Memories").
	WithAuthentication().
	WithSuccessStatus(201)

// GetMemory returns a single memory by ID.
var GetMemory = rocco.GET("/agents/{agent_id}/memories/{id}", func(req *rocco.Request[rocco.NoBody]) (wire.MemoryResponse, error) {
	memories := sum.MustUse[contracts.Memories](req.Context)

	id := req.Params.Path["id"]

	memory, err := memories.Get(req.Context, id)
	if err != nil {
		return wire.MemoryResponse{}, ErrMemoryNotFound
	}

	return transformers.MemoryToResponse(memory), nil
}).WithPathParams("agent_id", "id").
	WithSummary("Get memory").
	WithDescription("Returns a memory by ID.").
	WithTags("Memories").
	WithErrors(ErrMemoryNotFound).
	WithAuthentication()

// UpdateMemory updates an existing memory.
var UpdateMemory = rocco.PATCH("/agents/{agent_id}/memories/{id}", func(req *rocco.Request[wire.UpdateMemoryRequest]) (wire.MemoryResponse, error) {
	memories := sum.MustUse[contracts.Memories](req.Context)

	id := req.Params.Path["id"]

	memory, err := memories.Get(req.Context, id)
	if err != nil {
		return wire.MemoryResponse{}, ErrMemoryNotFound
	}

	transformers.ApplyUpdateMemory(req.Body, memory)
	memory.UpdatedAt = time.Now()

	if err := memories.Set(req.Context, id, memory); err != nil {
		return wire.MemoryResponse{}, err
	}

	return transformers.MemoryToResponse(memory), nil
}).WithPathParams("agent_id", "id").
	WithSummary("Update memory").
	WithDescription("Updates a memory's content, correlation ID, or metadata.").
	WithTags("Memories").
	WithErrors(ErrMemoryNotFound).
	WithAuthentication()

// DeleteMemory removes a memory by ID.
var DeleteMemory = rocco.DELETE("/agents/{agent_id}/memories/{id}", func(req *rocco.Request[rocco.NoBody]) (rocco.NoBody, error) {
	memories := sum.MustUse[contracts.Memories](req.Context)

	id := req.Params.Path["id"]

	if err := memories.Delete(req.Context, id); err != nil {
		return rocco.NoBody{}, ErrMemoryNotFound
	}

	return rocco.NoBody{}, nil
}).WithPathParams("agent_id", "id").
	WithSummary("Delete memory").
	WithDescription("Deletes a memory by ID.").
	WithTags("Memories").
	WithErrors(ErrMemoryNotFound).
	WithAuthentication().
	WithSuccessStatus(204)

// GetMemoryByCorrelationID retrieves a memory by its correlation ID.
var GetMemoryByCorrelationID = rocco.GET("/agents/{agent_id}/memories/by-correlation/{correlation_id}", func(req *rocco.Request[rocco.NoBody]) (wire.MemoryResponse, error) {
	memories := sum.MustUse[contracts.Memories](req.Context)

	correlationID := req.Params.Path["correlation_id"]

	memory, err := memories.GetByCorrelationID(req.Context, correlationID)
	if err != nil {
		return wire.MemoryResponse{}, ErrMemoryNotFound
	}

	return transformers.MemoryToResponse(memory), nil
}).WithPathParams("agent_id", "correlation_id").
	WithSummary("Get memory by correlation ID").
	WithDescription("Returns a memory by its correlation ID.").
	WithTags("Memories").
	WithErrors(ErrMemoryNotFound).
	WithAuthentication()
