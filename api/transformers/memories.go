// Package transformers provides pure functions for mapping between models and wire types
// on the public API surface.
package transformers

import (
	"encoding/json"

	"github.com/zoobzio/nestor/api/wire"
	"github.com/zoobzio/nestor/models"
)

// MemoryToResponse transforms a Memory model to a public API response.
func MemoryToResponse(m *models.Memory) wire.MemoryResponse {
	var metadata json.RawMessage
	if m.Metadata != nil {
		metadata = json.RawMessage(m.Metadata)
	}
	return wire.MemoryResponse{
		ID:            m.ID,
		AgentID:       m.AgentID,
		Content:       m.Content,
		Status:        string(m.Status),
		CorrelationID: m.CorrelationID,
		Metadata:      metadata,
		CreatedAt:     m.CreatedAt,
		UpdatedAt:     m.UpdatedAt,
	}
}

// MemoriesToList transforms a slice of Memory models to a list response.
func MemoriesToList(memories []*models.Memory) wire.MemoryListResponse {
	resp := wire.MemoryListResponse{
		Memories: make([]wire.MemoryResponse, len(memories)),
	}
	for i, m := range memories {
		resp.Memories[i] = MemoryToResponse(m)
	}
	return resp
}

// ApplyCreateMemory applies a CreateMemoryRequest to a Memory model.
func ApplyCreateMemory(req wire.CreateMemoryRequest, m *models.Memory) {
	m.Content = req.Content
	m.CorrelationID = req.CorrelationID
	if req.Metadata != nil {
		m.Metadata = []byte(req.Metadata)
	}
}

// ApplyUpdateMemory applies an UpdateMemoryRequest to a Memory model.
func ApplyUpdateMemory(req wire.UpdateMemoryRequest, m *models.Memory) {
	m.Content = req.Content
	m.CorrelationID = req.CorrelationID
	if req.Metadata != nil {
		m.Metadata = []byte(req.Metadata)
	}
}
