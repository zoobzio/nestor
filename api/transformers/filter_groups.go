// Package transformers provides pure functions for mapping between models and wire types
// on the public API surface.
package transformers

import (
	"encoding/json"

	"github.com/zoobzio/nestor/api/wire"
	"github.com/zoobzio/nestor/models"
)

// FilterGroupToResponse transforms a FilterGroup model to a public API response.
func FilterGroupToResponse(f *models.FilterGroup) wire.FilterGroupResponse {
	var criteria json.RawMessage
	if f.Criteria != nil {
		criteria = json.RawMessage(f.Criteria)
	}
	return wire.FilterGroupResponse{
		ID:             f.ID,
		AgentID:        f.AgentID,
		Name:           f.Name,
		Criteria:       criteria,
		PromptTemplate: f.PromptTemplate,
		CreatedAt:      f.CreatedAt,
		UpdatedAt:      f.UpdatedAt,
	}
}

// FilterGroupsToList transforms a slice of FilterGroup models to a list response.
func FilterGroupsToList(filterGroups []*models.FilterGroup) wire.FilterGroupListResponse {
	resp := wire.FilterGroupListResponse{
		FilterGroups: make([]wire.FilterGroupResponse, len(filterGroups)),
	}
	for i, f := range filterGroups {
		resp.FilterGroups[i] = FilterGroupToResponse(f)
	}
	return resp
}

// ApplyCreateFilterGroup applies a CreateFilterGroupRequest to a FilterGroup model.
func ApplyCreateFilterGroup(req wire.CreateFilterGroupRequest, f *models.FilterGroup) {
	f.Name = req.Name
	f.PromptTemplate = req.PromptTemplate
	if req.Criteria != nil {
		f.Criteria = []byte(req.Criteria)
	}
}

// ApplyUpdateFilterGroup applies an UpdateFilterGroupRequest to a FilterGroup model.
func ApplyUpdateFilterGroup(req wire.UpdateFilterGroupRequest, f *models.FilterGroup) {
	f.Name = req.Name
	f.PromptTemplate = req.PromptTemplate
	if req.Criteria != nil {
		f.Criteria = []byte(req.Criteria)
	}
}
