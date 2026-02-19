// Package models contains domain model types.
package models

import (
	"time"

	"github.com/zoobzio/check"
)

// MemoryStatus represents the processing state of a memory.
type MemoryStatus string

const (
	MemoryStatusProcessing MemoryStatus = "processing"
	MemoryStatusReady      MemoryStatus = "ready"
	MemoryStatusFailed     MemoryStatus = "failed"
)

// Memory represents a semantic memory belonging to an agent.
// Memories store markdown content with optional correlation to external events
// and arbitrary agent-defined metadata.
type Memory struct {
	ID            string       `json:"id" db:"id" constraints:"primarykey" description:"Memory UUID" example:"550e8400-e29b-41d4-a716-446655440002"`
	AgentID       string       `json:"agent_id" db:"agent_id" constraints:"notnull" references:"agents(id)" description:"Owning agent UUID" example:"550e8400-e29b-41d4-a716-446655440001"`
	Content       string       `json:"content" db:"content" constraints:"notnull" type:"text" description:"Memory content in markdown" example:"The user prefers concise responses."`
	Status        MemoryStatus `json:"status" db:"status" constraints:"notnull" default:"'processing'" description:"Processing status: processing, ready, failed" example:"ready"`
	CorrelationID *string      `json:"correlation_id,omitempty" db:"correlation_id" description:"Optional correlation ID for linking to external events" example:"evt_abc123"`
	Metadata      []byte       `json:"metadata,omitempty" db:"metadata" type:"jsonb" description:"Arbitrary agent-defined metadata"`
	CreatedAt     time.Time    `json:"created_at" db:"created_at" default:"now()" description:"Memory creation time"`
	UpdatedAt     time.Time    `json:"updated_at" db:"updated_at" default:"now()" description:"Last update time"`
}

// Clone returns a deep copy of Memory.
func (m Memory) Clone() Memory {
	c := m
	if m.CorrelationID != nil {
		s := *m.CorrelationID
		c.CorrelationID = &s
	}
	if m.Metadata != nil {
		c.Metadata = make([]byte, len(m.Metadata))
		copy(c.Metadata, m.Metadata)
	}
	return c
}

// Validate validates the Memory model.
func (m Memory) Validate() error {
	return check.All(
		check.Str(m.ID, "id").Required().V(),
		check.Str(m.AgentID, "agent_id").Required().V(),
		check.Str(m.Content, "content").Required().V(),
	).Err()
}
