//go:build testing

package models

import (
	"testing"
	"time"
)

func TestMemoryValidate_Valid(t *testing.T) {
	m := Memory{
		ID:        "550e8400-e29b-41d4-a716-446655440002",
		AgentID:   "550e8400-e29b-41d4-a716-446655440001",
		Content:   "The user prefers concise responses.",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	if err := m.Validate(); err != nil {
		t.Errorf("expected no error, got: %v", err)
	}
}

func TestMemoryValidate_MissingID(t *testing.T) {
	m := Memory{
		AgentID: "550e8400-e29b-41d4-a716-446655440001",
		Content: "The user prefers concise responses.",
	}
	if err := m.Validate(); err == nil {
		t.Error("expected error for missing id, got nil")
	}
}

func TestMemoryValidate_MissingAgentID(t *testing.T) {
	m := Memory{
		ID:      "550e8400-e29b-41d4-a716-446655440002",
		Content: "The user prefers concise responses.",
	}
	if err := m.Validate(); err == nil {
		t.Error("expected error for missing agent_id, got nil")
	}
}

func TestMemoryValidate_MissingContent(t *testing.T) {
	m := Memory{
		ID:      "550e8400-e29b-41d4-a716-446655440002",
		AgentID: "550e8400-e29b-41d4-a716-446655440001",
	}
	if err := m.Validate(); err == nil {
		t.Error("expected error for missing content, got nil")
	}
}

func TestMemoryClone_Basic(t *testing.T) {
	now := time.Now()
	original := Memory{
		ID:        "550e8400-e29b-41d4-a716-446655440002",
		AgentID:   "550e8400-e29b-41d4-a716-446655440001",
		Content:   "The user prefers concise responses.",
		CreatedAt: now,
		UpdatedAt: now,
	}

	clone := original.Clone()

	if clone.ID != original.ID {
		t.Errorf("Clone ID mismatch: got %q, want %q", clone.ID, original.ID)
	}
	if clone.AgentID != original.AgentID {
		t.Errorf("Clone AgentID mismatch: got %q, want %q", clone.AgentID, original.AgentID)
	}
	if clone.Content != original.Content {
		t.Errorf("Clone Content mismatch: got %q, want %q", clone.Content, original.Content)
	}
	if !clone.CreatedAt.Equal(original.CreatedAt) {
		t.Errorf("Clone CreatedAt mismatch: got %v, want %v", clone.CreatedAt, original.CreatedAt)
	}
	if !clone.UpdatedAt.Equal(original.UpdatedAt) {
		t.Errorf("Clone UpdatedAt mismatch: got %v, want %v", clone.UpdatedAt, original.UpdatedAt)
	}
	if clone.CorrelationID != nil {
		t.Error("expected nil CorrelationID in clone, got non-nil")
	}
	if clone.Metadata != nil {
		t.Error("expected nil Metadata in clone, got non-nil")
	}
}

func TestMemoryClone_NilCorrelationID(t *testing.T) {
	original := Memory{
		ID:      "550e8400-e29b-41d4-a716-446655440002",
		AgentID: "550e8400-e29b-41d4-a716-446655440001",
		Content: "Some memory.",
	}

	clone := original.Clone()

	if clone.CorrelationID != nil {
		t.Error("expected nil CorrelationID in clone, got non-nil")
	}
}

func TestMemoryClone_WithCorrelationID_IsDeepCopy(t *testing.T) {
	corrID := "evt_abc123"
	original := Memory{
		ID:            "550e8400-e29b-41d4-a716-446655440002",
		AgentID:       "550e8400-e29b-41d4-a716-446655440001",
		Content:       "Some memory.",
		CorrelationID: &corrID,
	}

	clone := original.Clone()

	if clone.CorrelationID == nil {
		t.Fatal("expected non-nil CorrelationID in clone, got nil")
	}
	if *clone.CorrelationID != corrID {
		t.Errorf("Clone CorrelationID %q, want %q", *clone.CorrelationID, corrID)
	}
	// Mutating the clone's correlation ID must not affect the original.
	*clone.CorrelationID = "mutated"
	if *original.CorrelationID == "mutated" {
		t.Error("mutating clone CorrelationID affected original")
	}
}

func TestMemoryClone_WithMetadata_IsDeepCopy(t *testing.T) {
	original := Memory{
		ID:       "550e8400-e29b-41d4-a716-446655440002",
		AgentID:  "550e8400-e29b-41d4-a716-446655440001",
		Content:  "Some memory.",
		Metadata: []byte(`{"key":"value"}`),
	}

	clone := original.Clone()

	if string(clone.Metadata) != string(original.Metadata) {
		t.Errorf("Clone Metadata %q, want %q", clone.Metadata, original.Metadata)
	}
	// Mutating the clone's metadata must not affect the original.
	clone.Metadata[0] = 'X'
	if original.Metadata[0] == 'X' {
		t.Error("mutating clone Metadata affected original")
	}
}
