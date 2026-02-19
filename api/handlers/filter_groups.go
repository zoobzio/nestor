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

// ListFilterGroups returns all filter groups for the given agent.
var ListFilterGroups = rocco.GET("/agents/{agent_id}/filter-groups", func(req *rocco.Request[rocco.NoBody]) (wire.FilterGroupListResponse, error) {
	filterGroups := sum.MustUse[contracts.FilterGroups](req.Context)

	agentID := req.Params.Path["agent_id"]

	list, err := filterGroups.ListByAgentID(req.Context, agentID)
	if err != nil {
		return wire.FilterGroupListResponse{}, err
	}

	return transformers.FilterGroupsToList(list), nil
}).WithPathParams("agent_id").
	WithSummary("List filter groups").
	WithDescription("Returns all filter groups belonging to the given agent.").
	WithTags("FilterGroups").
	WithAuthentication()

// CreateFilterGroup creates a new filter group for the given agent.
var CreateFilterGroup = rocco.POST("/agents/{agent_id}/filter-groups", func(req *rocco.Request[wire.CreateFilterGroupRequest]) (wire.FilterGroupResponse, error) {
	filterGroups := sum.MustUse[contracts.FilterGroups](req.Context)

	agentID := req.Params.Path["agent_id"]

	fg := &models.FilterGroup{
		ID:        newUUID(),
		AgentID:   agentID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	transformers.ApplyCreateFilterGroup(req.Body, fg)

	if err := filterGroups.Set(req.Context, fg.ID, fg); err != nil {
		return wire.FilterGroupResponse{}, err
	}

	return transformers.FilterGroupToResponse(fg), nil
}).WithPathParams("agent_id").
	WithSummary("Create filter group").
	WithDescription("Creates a new filter group for the given agent.").
	WithTags("FilterGroups").
	WithAuthentication().
	WithSuccessStatus(201)

// GetFilterGroup returns a single filter group by ID.
var GetFilterGroup = rocco.GET("/agents/{agent_id}/filter-groups/{id}", func(req *rocco.Request[rocco.NoBody]) (wire.FilterGroupResponse, error) {
	filterGroups := sum.MustUse[contracts.FilterGroups](req.Context)

	id := req.Params.Path["id"]

	fg, err := filterGroups.Get(req.Context, id)
	if err != nil {
		return wire.FilterGroupResponse{}, ErrFilterGroupNotFound
	}

	return transformers.FilterGroupToResponse(fg), nil
}).WithPathParams("agent_id", "id").
	WithSummary("Get filter group").
	WithDescription("Returns a filter group by ID.").
	WithTags("FilterGroups").
	WithErrors(ErrFilterGroupNotFound).
	WithAuthentication()

// UpdateFilterGroup updates an existing filter group.
var UpdateFilterGroup = rocco.PATCH("/agents/{agent_id}/filter-groups/{id}", func(req *rocco.Request[wire.UpdateFilterGroupRequest]) (wire.FilterGroupResponse, error) {
	filterGroups := sum.MustUse[contracts.FilterGroups](req.Context)

	id := req.Params.Path["id"]

	fg, err := filterGroups.Get(req.Context, id)
	if err != nil {
		return wire.FilterGroupResponse{}, ErrFilterGroupNotFound
	}

	transformers.ApplyUpdateFilterGroup(req.Body, fg)
	fg.UpdatedAt = time.Now()

	if err := filterGroups.Set(req.Context, id, fg); err != nil {
		return wire.FilterGroupResponse{}, err
	}

	return transformers.FilterGroupToResponse(fg), nil
}).WithPathParams("agent_id", "id").
	WithSummary("Update filter group").
	WithDescription("Updates a filter group's details.").
	WithTags("FilterGroups").
	WithErrors(ErrFilterGroupNotFound).
	WithAuthentication()

// DeleteFilterGroup removes a filter group by ID.
var DeleteFilterGroup = rocco.DELETE("/agents/{agent_id}/filter-groups/{id}", func(req *rocco.Request[rocco.NoBody]) (rocco.NoBody, error) {
	filterGroups := sum.MustUse[contracts.FilterGroups](req.Context)

	id := req.Params.Path["id"]

	if err := filterGroups.Delete(req.Context, id); err != nil {
		return rocco.NoBody{}, ErrFilterGroupNotFound
	}

	return rocco.NoBody{}, nil
}).WithPathParams("agent_id", "id").
	WithSummary("Delete filter group").
	WithDescription("Deletes a filter group by ID.").
	WithTags("FilterGroups").
	WithErrors(ErrFilterGroupNotFound).
	WithAuthentication().
	WithSuccessStatus(204)
