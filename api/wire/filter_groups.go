// Package wire defines request and response types for the public API surface.
package wire

import (
	"encoding/json"
	"time"

	"github.com/zoobzio/check"
)

// FilterGroupResponse is the public API response for filter group data.
type FilterGroupResponse struct {
	ID             string          `json:"id" description:"FilterGroup UUID" example:"550e8400-e29b-41d4-a716-446655440010"`
	AgentID        string          `json:"agent_id" description:"Owning agent UUID" example:"550e8400-e29b-41d4-a716-446655440001"`
	Name           string          `json:"name" description:"FilterGroup name" example:"Support Escalations"`
	Criteria       json.RawMessage `json:"criteria,omitempty" description:"Flexible filter conditions as JSON"`
	PromptTemplate string          `json:"prompt_template" description:"Instructions for memory generation" example:"Summarize this event focusing on resolution steps."`
	CreatedAt      time.Time       `json:"created_at" description:"FilterGroup creation time"`
	UpdatedAt      time.Time       `json:"updated_at" description:"Last update time"`
}

// Clone returns a deep copy of FilterGroupResponse.
func (f FilterGroupResponse) Clone() FilterGroupResponse {
	c := f
	if f.Criteria != nil {
		c.Criteria = make(json.RawMessage, len(f.Criteria))
		copy(c.Criteria, f.Criteria)
	}
	return c
}

// FilterGroupListResponse is the public API response for listing filter groups.
type FilterGroupListResponse struct {
	FilterGroups []FilterGroupResponse `json:"filter_groups" description:"List of filter groups"`
}

// Clone returns a deep copy of FilterGroupListResponse.
func (f FilterGroupListResponse) Clone() FilterGroupListResponse {
	c := f
	if f.FilterGroups != nil {
		c.FilterGroups = make([]FilterGroupResponse, len(f.FilterGroups))
		for i, fg := range f.FilterGroups {
			c.FilterGroups[i] = fg.Clone()
		}
	}
	return c
}

// CreateFilterGroupRequest is the request body for creating a filter group.
type CreateFilterGroupRequest struct {
	Name           string          `json:"name" description:"FilterGroup name" example:"Support Escalations"`
	Criteria       json.RawMessage `json:"criteria,omitempty" description:"Flexible filter conditions as JSON"`
	PromptTemplate string          `json:"prompt_template" description:"Instructions for memory generation" example:"Summarize this event focusing on resolution steps."`
}

// Validate validates the request.
func (r *CreateFilterGroupRequest) Validate() error {
	return check.All(
		check.Str(r.Name, "name").Required().MaxLen(255).V(),
		check.Str(r.PromptTemplate, "prompt_template").Required().V(),
	).Err()
}

// Clone returns a deep copy of CreateFilterGroupRequest.
func (r CreateFilterGroupRequest) Clone() CreateFilterGroupRequest {
	c := r
	if r.Criteria != nil {
		c.Criteria = make(json.RawMessage, len(r.Criteria))
		copy(c.Criteria, r.Criteria)
	}
	return c
}

// UpdateFilterGroupRequest is the request body for updating a filter group.
type UpdateFilterGroupRequest struct {
	Name           string          `json:"name" description:"FilterGroup name" example:"Support Escalations"`
	Criteria       json.RawMessage `json:"criteria,omitempty" description:"Flexible filter conditions as JSON"`
	PromptTemplate string          `json:"prompt_template" description:"Instructions for memory generation" example:"Summarize this event focusing on resolution steps."`
}

// Validate validates the request.
func (r *UpdateFilterGroupRequest) Validate() error {
	return check.All(
		check.Str(r.Name, "name").Required().MaxLen(255).V(),
		check.Str(r.PromptTemplate, "prompt_template").Required().V(),
	).Err()
}

// Clone returns a deep copy of UpdateFilterGroupRequest.
func (r UpdateFilterGroupRequest) Clone() UpdateFilterGroupRequest {
	c := r
	if r.Criteria != nil {
		c.Criteria = make(json.RawMessage, len(r.Criteria))
		copy(c.Criteria, r.Criteria)
	}
	return c
}
