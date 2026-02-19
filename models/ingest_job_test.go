//go:build testing

package models

import (
	"testing"
	"time"
)

func TestIngestJobValidate_Valid(t *testing.T) {
	j := IngestJob{
		ID:         "550e8400-e29b-41d4-a716-446655440000",
		CustomerID: "550e8400-e29b-41d4-a716-446655440001",
		AgentID:    "550e8400-e29b-41d4-a716-446655440002",
		Content:    "Raw content for ingestion.",
	}
	if err := j.Validate(); err != nil {
		t.Errorf("expected no error, got: %v", err)
	}
}

func TestIngestJobValidate_MissingID(t *testing.T) {
	j := IngestJob{
		CustomerID: "550e8400-e29b-41d4-a716-446655440001",
		AgentID:    "550e8400-e29b-41d4-a716-446655440002",
		Content:    "Raw content for ingestion.",
	}
	if err := j.Validate(); err == nil {
		t.Error("expected error for missing id, got nil")
	}
}

func TestIngestJobValidate_MissingCustomerID(t *testing.T) {
	j := IngestJob{
		ID:      "550e8400-e29b-41d4-a716-446655440000",
		AgentID: "550e8400-e29b-41d4-a716-446655440002",
		Content: "Raw content for ingestion.",
	}
	if err := j.Validate(); err == nil {
		t.Error("expected error for missing customer_id, got nil")
	}
}

func TestIngestJobValidate_MissingAgentID(t *testing.T) {
	j := IngestJob{
		ID:         "550e8400-e29b-41d4-a716-446655440000",
		CustomerID: "550e8400-e29b-41d4-a716-446655440001",
		Content:    "Raw content for ingestion.",
	}
	if err := j.Validate(); err == nil {
		t.Error("expected error for missing agent_id, got nil")
	}
}

func TestIngestJobValidate_MissingContent(t *testing.T) {
	j := IngestJob{
		ID:         "550e8400-e29b-41d4-a716-446655440000",
		CustomerID: "550e8400-e29b-41d4-a716-446655440001",
		AgentID:    "550e8400-e29b-41d4-a716-446655440002",
	}
	if err := j.Validate(); err == nil {
		t.Error("expected error for missing content, got nil")
	}
}

func TestIngestJobClone_Basic(t *testing.T) {
	now := time.Now()
	original := IngestJob{
		ID:         "550e8400-e29b-41d4-a716-446655440000",
		CustomerID: "550e8400-e29b-41d4-a716-446655440001",
		AgentID:    "550e8400-e29b-41d4-a716-446655440002",
		MemoryID:   "550e8400-e29b-41d4-a716-446655440003",
		Content:    "Raw content for ingestion.",
		Stage:      JobStageValidate,
		Status:     JobStatusPending,
		CreatedAt:  now,
	}

	clone := original.Clone()

	if clone.ID != original.ID {
		t.Errorf("Clone ID mismatch: got %q, want %q", clone.ID, original.ID)
	}
	if clone.CustomerID != original.CustomerID {
		t.Errorf("Clone CustomerID mismatch: got %q, want %q", clone.CustomerID, original.CustomerID)
	}
	if clone.AgentID != original.AgentID {
		t.Errorf("Clone AgentID mismatch: got %q, want %q", clone.AgentID, original.AgentID)
	}
	if clone.MemoryID != original.MemoryID {
		t.Errorf("Clone MemoryID mismatch: got %q, want %q", clone.MemoryID, original.MemoryID)
	}
	if clone.Content != original.Content {
		t.Errorf("Clone Content mismatch: got %q, want %q", clone.Content, original.Content)
	}
	if clone.Stage != original.Stage {
		t.Errorf("Clone Stage mismatch: got %q, want %q", clone.Stage, original.Stage)
	}
	if clone.Status != original.Status {
		t.Errorf("Clone Status mismatch: got %q, want %q", clone.Status, original.Status)
	}
	if !clone.CreatedAt.Equal(original.CreatedAt) {
		t.Errorf("Clone CreatedAt mismatch: got %v, want %v", clone.CreatedAt, original.CreatedAt)
	}
}

func TestIngestJobClone_ChunkIDs_IsDeepCopy(t *testing.T) {
	original := IngestJob{
		ID:         "550e8400-e29b-41d4-a716-446655440000",
		CustomerID: "550e8400-e29b-41d4-a716-446655440001",
		AgentID:    "550e8400-e29b-41d4-a716-446655440002",
		Content:    "Raw content for ingestion.",
		ChunkIDs:   []string{"chunk-1", "chunk-2", "chunk-3"},
	}

	clone := original.Clone()

	if len(clone.ChunkIDs) != len(original.ChunkIDs) {
		t.Fatalf("Clone ChunkIDs length mismatch: got %d, want %d", len(clone.ChunkIDs), len(original.ChunkIDs))
	}
	for i, id := range original.ChunkIDs {
		if clone.ChunkIDs[i] != id {
			t.Errorf("Clone ChunkIDs[%d] mismatch: got %q, want %q", i, clone.ChunkIDs[i], id)
		}
	}
	// Mutating the clone's ChunkIDs must not affect the original.
	clone.ChunkIDs[0] = "mutated"
	if original.ChunkIDs[0] == "mutated" {
		t.Error("mutating clone ChunkIDs affected original")
	}
}

func TestIngestJobClone_NilChunkIDs(t *testing.T) {
	original := IngestJob{
		ID:         "550e8400-e29b-41d4-a716-446655440000",
		CustomerID: "550e8400-e29b-41d4-a716-446655440001",
		AgentID:    "550e8400-e29b-41d4-a716-446655440002",
		Content:    "Raw content for ingestion.",
	}

	clone := original.Clone()

	if clone.ChunkIDs != nil {
		t.Error("expected nil ChunkIDs in clone, got non-nil")
	}
}

func TestIngestJobClone_Error_IsDeepCopy(t *testing.T) {
	errMsg := "something went wrong"
	original := IngestJob{
		ID:         "550e8400-e29b-41d4-a716-446655440000",
		CustomerID: "550e8400-e29b-41d4-a716-446655440001",
		AgentID:    "550e8400-e29b-41d4-a716-446655440002",
		Content:    "Raw content for ingestion.",
		Error:      &errMsg,
	}

	clone := original.Clone()

	if clone.Error == nil {
		t.Fatal("expected non-nil Error in clone, got nil")
	}
	if *clone.Error != errMsg {
		t.Errorf("Clone Error mismatch: got %q, want %q", *clone.Error, errMsg)
	}
	// Mutating the clone's Error must not affect the original.
	*clone.Error = "mutated"
	if *original.Error == "mutated" {
		t.Error("mutating clone Error affected original")
	}
}

func TestIngestJobClone_NilError(t *testing.T) {
	original := IngestJob{
		ID:         "550e8400-e29b-41d4-a716-446655440000",
		CustomerID: "550e8400-e29b-41d4-a716-446655440001",
		AgentID:    "550e8400-e29b-41d4-a716-446655440002",
		Content:    "Raw content for ingestion.",
	}

	clone := original.Clone()

	if clone.Error != nil {
		t.Error("expected nil Error in clone, got non-nil")
	}
}

func TestIngestJobClone_StartedAt_IsDeepCopy(t *testing.T) {
	started := time.Now()
	original := IngestJob{
		ID:         "550e8400-e29b-41d4-a716-446655440000",
		CustomerID: "550e8400-e29b-41d4-a716-446655440001",
		AgentID:    "550e8400-e29b-41d4-a716-446655440002",
		Content:    "Raw content for ingestion.",
		StartedAt:  &started,
	}

	clone := original.Clone()

	if clone.StartedAt == nil {
		t.Fatal("expected non-nil StartedAt in clone, got nil")
	}
	if !clone.StartedAt.Equal(*original.StartedAt) {
		t.Errorf("Clone StartedAt mismatch: got %v, want %v", clone.StartedAt, original.StartedAt)
	}
	// Mutating the clone's StartedAt must not affect the original.
	*clone.StartedAt = time.Time{}
	if original.StartedAt.IsZero() {
		t.Error("mutating clone StartedAt affected original")
	}
}

func TestIngestJobClone_NilStartedAt(t *testing.T) {
	original := IngestJob{
		ID:         "550e8400-e29b-41d4-a716-446655440000",
		CustomerID: "550e8400-e29b-41d4-a716-446655440001",
		AgentID:    "550e8400-e29b-41d4-a716-446655440002",
		Content:    "Raw content for ingestion.",
	}

	clone := original.Clone()

	if clone.StartedAt != nil {
		t.Error("expected nil StartedAt in clone, got non-nil")
	}
}

func TestIngestJobClone_CompletedAt_IsDeepCopy(t *testing.T) {
	completed := time.Now()
	original := IngestJob{
		ID:          "550e8400-e29b-41d4-a716-446655440000",
		CustomerID:  "550e8400-e29b-41d4-a716-446655440001",
		AgentID:     "550e8400-e29b-41d4-a716-446655440002",
		Content:     "Raw content for ingestion.",
		CompletedAt: &completed,
	}

	clone := original.Clone()

	if clone.CompletedAt == nil {
		t.Fatal("expected non-nil CompletedAt in clone, got nil")
	}
	if !clone.CompletedAt.Equal(*original.CompletedAt) {
		t.Errorf("Clone CompletedAt mismatch: got %v, want %v", clone.CompletedAt, original.CompletedAt)
	}
	// Mutating the clone's CompletedAt must not affect the original.
	*clone.CompletedAt = time.Time{}
	if original.CompletedAt.IsZero() {
		t.Error("mutating clone CompletedAt affected original")
	}
}

func TestIngestJobClone_NilCompletedAt(t *testing.T) {
	original := IngestJob{
		ID:         "550e8400-e29b-41d4-a716-446655440000",
		CustomerID: "550e8400-e29b-41d4-a716-446655440001",
		AgentID:    "550e8400-e29b-41d4-a716-446655440002",
		Content:    "Raw content for ingestion.",
	}

	clone := original.Clone()

	if clone.CompletedAt != nil {
		t.Error("expected nil CompletedAt in clone, got non-nil")
	}
}
