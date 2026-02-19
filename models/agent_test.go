//go:build testing

package models

import (
	"strings"
	"testing"
	"time"
)

func TestValidate_Valid(t *testing.T) {
	a := Agent{
		ID:        "550e8400-e29b-41d4-a716-446655440001",
		UserID:    "550e8400-e29b-41d4-a716-446655440000",
		Name:      "Research Assistant",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	if err := a.Validate(); err != nil {
		t.Errorf("expected no error, got: %v", err)
	}
}

func TestValidate_MissingID(t *testing.T) {
	a := Agent{
		UserID: "550e8400-e29b-41d4-a716-446655440000",
		Name:   "Research Assistant",
	}
	if err := a.Validate(); err == nil {
		t.Error("expected error for missing id, got nil")
	}
}

func TestValidate_MissingUserID(t *testing.T) {
	a := Agent{
		ID:   "550e8400-e29b-41d4-a716-446655440001",
		Name: "Research Assistant",
	}
	if err := a.Validate(); err == nil {
		t.Error("expected error for missing user_id, got nil")
	}
}

func TestValidate_MissingName(t *testing.T) {
	a := Agent{
		ID:     "550e8400-e29b-41d4-a716-446655440001",
		UserID: "550e8400-e29b-41d4-a716-446655440000",
	}
	if err := a.Validate(); err == nil {
		t.Error("expected error for missing name, got nil")
	}
}

func TestValidate_NameTooLong(t *testing.T) {
	a := Agent{
		ID:     "550e8400-e29b-41d4-a716-446655440001",
		UserID: "550e8400-e29b-41d4-a716-446655440000",
		Name:   strings.Repeat("a", 256),
	}
	if err := a.Validate(); err == nil {
		t.Error("expected error for name exceeding max length, got nil")
	}
}

func TestClone(t *testing.T) {
	now := time.Now()
	original := Agent{
		ID:        "550e8400-e29b-41d4-a716-446655440001",
		UserID:    "550e8400-e29b-41d4-a716-446655440000",
		Name:      "Research Assistant",
		CreatedAt: now,
		UpdatedAt: now,
	}

	clone := original.Clone()

	if clone.ID != original.ID {
		t.Errorf("Clone ID mismatch: got %q, want %q", clone.ID, original.ID)
	}
	if clone.UserID != original.UserID {
		t.Errorf("Clone UserID mismatch: got %q, want %q", clone.UserID, original.UserID)
	}
	if clone.Name != original.Name {
		t.Errorf("Clone Name mismatch: got %q, want %q", clone.Name, original.Name)
	}
	if !clone.CreatedAt.Equal(original.CreatedAt) {
		t.Errorf("Clone CreatedAt mismatch: got %v, want %v", clone.CreatedAt, original.CreatedAt)
	}
	if !clone.UpdatedAt.Equal(original.UpdatedAt) {
		t.Errorf("Clone UpdatedAt mismatch: got %v, want %v", clone.UpdatedAt, original.UpdatedAt)
	}

	// Mutating the clone must not affect the original.
	clone.Name = "Mutated"
	if original.Name == "Mutated" {
		t.Error("mutating clone affected original")
	}
}
