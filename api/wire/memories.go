// Package wire defines request and response types for the public API surface.
package wire

import (
	"encoding/json"
	"time"

	"github.com/zoobzio/check"
)

// MemoryResponse is the public API response for memory data.
type MemoryResponse struct {
	ID            string          `json:"id" description:"Memory UUID" example:"550e8400-e29b-41d4-a716-446655440002"`
	AgentID       string          `json:"agent_id" description:"Owning agent UUID" example:"550e8400-e29b-41d4-a716-446655440001"`
	Content       string          `json:"content" description:"Memory content in markdown" example:"The user prefers concise responses."`
	Status        string          `json:"status" description:"Processing status: processing, ready, failed" example:"ready"`
	CorrelationID *string         `json:"correlation_id,omitempty" description:"Optional correlation ID for linking to external events" example:"evt_abc123"`
	Metadata      json.RawMessage `json:"metadata,omitempty" description:"Arbitrary agent-defined metadata"`
	CreatedAt     time.Time       `json:"created_at" description:"Memory creation time"`
	UpdatedAt     time.Time       `json:"updated_at" description:"Last update time"`
}

// Clone returns a deep copy of MemoryResponse.
func (m MemoryResponse) Clone() MemoryResponse {
	c := m
	if m.CorrelationID != nil {
		s := *m.CorrelationID
		c.CorrelationID = &s
	}
	if m.Metadata != nil {
		c.Metadata = make(json.RawMessage, len(m.Metadata))
		copy(c.Metadata, m.Metadata)
	}
	return c
}

// MemoryListResponse is the public API response for listing memories.
type MemoryListResponse struct {
	Memories []MemoryResponse `json:"memories" description:"List of memories"`
}

// Clone returns a deep copy of MemoryListResponse.
func (m MemoryListResponse) Clone() MemoryListResponse {
	c := m
	if m.Memories != nil {
		c.Memories = make([]MemoryResponse, len(m.Memories))
		for i, mem := range m.Memories {
			c.Memories[i] = mem.Clone()
		}
	}
	return c
}

// CreateMemoryRequest is the request body for creating a memory.
type CreateMemoryRequest struct {
	Content       string          `json:"content" description:"Memory content in markdown" example:"The user prefers concise responses."`
	CorrelationID *string         `json:"correlation_id,omitempty" description:"Optional correlation ID for linking to external events" example:"evt_abc123"`
	Metadata      json.RawMessage `json:"metadata,omitempty" description:"Arbitrary agent-defined metadata"`
}

// Validate validates the request.
func (r *CreateMemoryRequest) Validate() error {
	return check.All(
		check.Str(r.Content, "content").Required().V(),
	).Err()
}

// Clone returns a deep copy of CreateMemoryRequest.
func (r CreateMemoryRequest) Clone() CreateMemoryRequest {
	c := r
	if r.CorrelationID != nil {
		s := *r.CorrelationID
		c.CorrelationID = &s
	}
	if r.Metadata != nil {
		c.Metadata = make(json.RawMessage, len(r.Metadata))
		copy(c.Metadata, r.Metadata)
	}
	return c
}

// UpdateMemoryRequest is the request body for updating a memory.
type UpdateMemoryRequest struct {
	Content       string          `json:"content" description:"Memory content in markdown" example:"The user prefers concise responses."`
	CorrelationID *string         `json:"correlation_id,omitempty" description:"Optional correlation ID for linking to external events" example:"evt_abc123"`
	Metadata      json.RawMessage `json:"metadata,omitempty" description:"Arbitrary agent-defined metadata"`
}

// Validate validates the request.
func (r *UpdateMemoryRequest) Validate() error {
	return check.All(
		check.Str(r.Content, "content").Required().V(),
	).Err()
}

// Clone returns a deep copy of UpdateMemoryRequest.
func (r UpdateMemoryRequest) Clone() UpdateMemoryRequest {
	c := r
	if r.CorrelationID != nil {
		s := *r.CorrelationID
		c.CorrelationID = &s
	}
	if r.Metadata != nil {
		c.Metadata = make(json.RawMessage, len(r.Metadata))
		copy(c.Metadata, r.Metadata)
	}
	return c
}
