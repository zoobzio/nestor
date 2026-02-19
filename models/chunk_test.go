//go:build testing

package models

import (
	"testing"
	"time"
)

func TestChunkValidate_Valid(t *testing.T) {
	c := Chunk{
		ID:        "550e8400-e29b-41d4-a716-446655440003",
		MemoryID:  "550e8400-e29b-41d4-a716-446655440002",
		Content:   "The user prefers concise responses.",
		Sequence:  0,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	if err := c.Validate(); err != nil {
		t.Errorf("expected no error, got: %v", err)
	}
}

func TestChunkValidate_MissingID(t *testing.T) {
	c := Chunk{
		MemoryID: "550e8400-e29b-41d4-a716-446655440002",
		Content:  "The user prefers concise responses.",
		Sequence: 0,
	}
	if err := c.Validate(); err == nil {
		t.Error("expected error for missing id, got nil")
	}
}

func TestChunkValidate_MissingMemoryID(t *testing.T) {
	c := Chunk{
		ID:       "550e8400-e29b-41d4-a716-446655440003",
		Content:  "The user prefers concise responses.",
		Sequence: 0,
	}
	if err := c.Validate(); err == nil {
		t.Error("expected error for missing memory_id, got nil")
	}
}

func TestChunkValidate_MissingContent(t *testing.T) {
	c := Chunk{
		ID:       "550e8400-e29b-41d4-a716-446655440003",
		MemoryID: "550e8400-e29b-41d4-a716-446655440002",
		Sequence: 0,
	}
	if err := c.Validate(); err == nil {
		t.Error("expected error for missing content, got nil")
	}
}

func TestChunkValidate_NegativeSequence(t *testing.T) {
	c := Chunk{
		ID:       "550e8400-e29b-41d4-a716-446655440003",
		MemoryID: "550e8400-e29b-41d4-a716-446655440002",
		Content:  "The user prefers concise responses.",
		Sequence: -1,
	}
	if err := c.Validate(); err == nil {
		t.Error("expected error for negative sequence, got nil")
	}
}

func TestChunkClone_Basic(t *testing.T) {
	now := time.Now()
	original := Chunk{
		ID:        "550e8400-e29b-41d4-a716-446655440003",
		MemoryID:  "550e8400-e29b-41d4-a716-446655440002",
		Content:   "The user prefers concise responses.",
		Sequence:  2,
		CreatedAt: now,
		UpdatedAt: now,
	}

	clone := original.Clone()

	if clone.ID != original.ID {
		t.Errorf("Clone ID %q, want %q", clone.ID, original.ID)
	}
	if clone.MemoryID != original.MemoryID {
		t.Errorf("Clone MemoryID %q, want %q", clone.MemoryID, original.MemoryID)
	}
	if clone.Content != original.Content {
		t.Errorf("Clone Content %q, want %q", clone.Content, original.Content)
	}
	if clone.Sequence != original.Sequence {
		t.Errorf("Clone Sequence %d, want %d", clone.Sequence, original.Sequence)
	}
	if !clone.CreatedAt.Equal(original.CreatedAt) {
		t.Errorf("Clone CreatedAt %v, want %v", clone.CreatedAt, original.CreatedAt)
	}
	if !clone.UpdatedAt.Equal(original.UpdatedAt) {
		t.Errorf("Clone UpdatedAt %v, want %v", clone.UpdatedAt, original.UpdatedAt)
	}
	if clone.Embedding != nil {
		t.Error("expected nil Embedding in clone, got non-nil")
	}
}

func TestChunkClone_NilEmbedding(t *testing.T) {
	original := Chunk{
		ID:       "550e8400-e29b-41d4-a716-446655440003",
		MemoryID: "550e8400-e29b-41d4-a716-446655440002",
		Content:  "Some chunk.",
	}

	clone := original.Clone()

	if clone.Embedding != nil {
		t.Error("expected nil Embedding in clone, got non-nil")
	}
}

func TestChunkClone_WithEmbedding_IsDeepCopy(t *testing.T) {
	original := Chunk{
		ID:        "550e8400-e29b-41d4-a716-446655440003",
		MemoryID:  "550e8400-e29b-41d4-a716-446655440002",
		Content:   "Some chunk.",
		Embedding: []float32{0.1, 0.2, 0.3},
	}

	clone := original.Clone()

	if len(clone.Embedding) != len(original.Embedding) {
		t.Fatalf("Clone Embedding length %d, want %d", len(clone.Embedding), len(original.Embedding))
	}
	for i, v := range original.Embedding {
		if clone.Embedding[i] != v {
			t.Errorf("Clone Embedding[%d] = %f, want %f", i, clone.Embedding[i], v)
		}
	}
	// Mutating the clone's embedding must not affect the original.
	clone.Embedding[0] = 9.9
	if original.Embedding[0] == 9.9 {
		t.Error("mutating clone Embedding affected original")
	}
}
