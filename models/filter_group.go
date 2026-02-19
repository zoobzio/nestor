// Package models contains domain model types.
package models

import (
	"time"

	"github.com/zoobzio/check"
)

// FilterGroup represents a named set of filter criteria belonging to an agent.
// FilterGroups define conditions and prompt templates used during event processing
// to generate targeted memories.
type FilterGroup struct {
	ID             string    `json:"id" db:"id" constraints:"primarykey" description:"FilterGroup UUID" example:"550e8400-e29b-41d4-a716-446655440010"`
	AgentID        string    `json:"agent_id" db:"agent_id" constraints:"notnull" references:"agents(id)" description:"Owning agent UUID" example:"550e8400-e29b-41d4-a716-446655440001"`
	Name           string    `json:"name" db:"name" constraints:"notnull" description:"FilterGroup name" example:"Support Escalations"`
	Criteria       []byte    `json:"criteria,omitempty" db:"criteria" type:"jsonb" description:"Flexible filter conditions as JSON"`
	PromptTemplate string    `json:"prompt_template" db:"prompt_template" constraints:"notnull" type:"text" description:"Instructions for memory generation" example:"Summarize this event focusing on resolution steps."`
	CreatedAt      time.Time `json:"created_at" db:"created_at" default:"now()" description:"FilterGroup creation time"`
	UpdatedAt      time.Time `json:"updated_at" db:"updated_at" default:"now()" description:"Last update time"`
}

// Clone returns a deep copy of FilterGroup.
func (f FilterGroup) Clone() FilterGroup {
	c := f
	if f.Criteria != nil {
		c.Criteria = make([]byte, len(f.Criteria))
		copy(c.Criteria, f.Criteria)
	}
	return c
}

// Validate validates the FilterGroup model.
func (f FilterGroup) Validate() error {
	return check.All(
		check.Str(f.ID, "id").Required().V(),
		check.Str(f.AgentID, "agent_id").Required().V(),
		check.Str(f.Name, "name").Required().MaxLen(255).V(),
		check.Str(f.PromptTemplate, "prompt_template").Required().V(),
	).Err()
}
