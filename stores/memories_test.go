//go:build testing

package stores_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/zoobzio/nestor/models"
	nesttest "github.com/zoobzio/nestor/testing"
)

// These tests exercise the contracts.Memories interface via MockMemories.
// The concrete Memories store (which embeds sum.Database) requires a live
// database and is covered by integration tests.

func TestMemories_Get_Found(t *testing.T) {
	memory := nesttest.NewMemory(t)

	mock := &nesttest.MockMemories{
		OnGet: func(_ context.Context, key string) (*models.Memory, error) {
			if key != memory.ID {
				return nil, errors.New("not found")
			}
			return memory, nil
		},
	}

	got, err := mock.Get(context.Background(), memory.ID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.ID != memory.ID {
		t.Errorf("got ID %q, want %q", got.ID, memory.ID)
	}
}

func TestMemories_Get_NotFound(t *testing.T) {
	mock := &nesttest.MockMemories{
		OnGet: func(_ context.Context, _ string) (*models.Memory, error) {
			return nil, errors.New("not found")
		},
	}

	_, err := mock.Get(context.Background(), "does-not-exist")
	if err == nil {
		t.Error("expected error for missing memory, got nil")
	}
}

func TestMemories_Set(t *testing.T) {
	var stored *models.Memory

	mock := &nesttest.MockMemories{
		OnSet: func(_ context.Context, _ string, memory *models.Memory) error {
			stored = memory
			return nil
		},
	}

	memory := nesttest.NewMemory(t)
	if err := mock.Set(context.Background(), memory.ID, memory); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if stored == nil {
		t.Fatal("expected memory to be stored, got nil")
	}
	if stored.ID != memory.ID {
		t.Errorf("stored ID %q, want %q", stored.ID, memory.ID)
	}
}

func TestMemories_Set_Error(t *testing.T) {
	mock := &nesttest.MockMemories{
		OnSet: func(_ context.Context, _ string, _ *models.Memory) error {
			return errors.New("storage error")
		},
	}

	memory := nesttest.NewMemory(t)
	if err := mock.Set(context.Background(), memory.ID, memory); err == nil {
		t.Error("expected error from Set, got nil")
	}
}

func TestMemories_Delete(t *testing.T) {
	var deletedKey string

	mock := &nesttest.MockMemories{
		OnDelete: func(_ context.Context, key string) error {
			deletedKey = key
			return nil
		},
	}

	if err := mock.Delete(context.Background(), "memory-id"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if deletedKey != "memory-id" {
		t.Errorf("deleted key %q, want %q", deletedKey, "memory-id")
	}
}

func TestMemories_Delete_Error(t *testing.T) {
	mock := &nesttest.MockMemories{
		OnDelete: func(_ context.Context, _ string) error {
			return errors.New("not found")
		},
	}

	if err := mock.Delete(context.Background(), "missing-id"); err == nil {
		t.Error("expected error from Delete, got nil")
	}
}

func TestMemories_ListByAgentID(t *testing.T) {
	now := time.Now()
	agentID := "550e8400-e29b-41d4-a716-446655440001"
	expected := []*models.Memory{
		{ID: "memory-1", AgentID: agentID, Content: "First memory.", CreatedAt: now, UpdatedAt: now},
		{ID: "memory-2", AgentID: agentID, Content: "Second memory.", CreatedAt: now, UpdatedAt: now},
	}

	mock := &nesttest.MockMemories{
		OnListByAgentID: func(_ context.Context, id string) ([]*models.Memory, error) {
			if id != agentID {
				return []*models.Memory{}, nil
			}
			return expected, nil
		},
	}

	list, err := mock.ListByAgentID(context.Background(), agentID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(list) != 2 {
		t.Fatalf("expected 2 memories, got %d", len(list))
	}
	if list[0].ID != "memory-1" {
		t.Errorf("first memory ID %q, want %q", list[0].ID, "memory-1")
	}
	if list[1].ID != "memory-2" {
		t.Errorf("second memory ID %q, want %q", list[1].ID, "memory-2")
	}
}

func TestMemories_ListByAgentID_Empty(t *testing.T) {
	mock := &nesttest.MockMemories{
		OnListByAgentID: func(_ context.Context, _ string) ([]*models.Memory, error) {
			return []*models.Memory{}, nil
		},
	}

	list, err := mock.ListByAgentID(context.Background(), "agent-with-no-memories")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(list) != 0 {
		t.Errorf("expected empty list, got %d memories", len(list))
	}
}

func TestMemories_ListByAgentID_Error(t *testing.T) {
	mock := &nesttest.MockMemories{
		OnListByAgentID: func(_ context.Context, _ string) ([]*models.Memory, error) {
			return nil, errors.New("database error")
		},
	}

	_, err := mock.ListByAgentID(context.Background(), "any-agent")
	if err == nil {
		t.Error("expected error from ListByAgentID, got nil")
	}
}

func TestMemories_GetByCorrelationID_Found(t *testing.T) {
	corrID := "evt_abc123"
	memory := nesttest.NewMemory(t)
	memory.CorrelationID = &corrID

	mock := &nesttest.MockMemories{
		OnGetByCorrelationID: func(_ context.Context, id string) (*models.Memory, error) {
			if id != corrID {
				return nil, errors.New("not found")
			}
			return memory, nil
		},
	}

	got, err := mock.GetByCorrelationID(context.Background(), corrID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.ID != memory.ID {
		t.Errorf("got ID %q, want %q", got.ID, memory.ID)
	}
	if got.CorrelationID == nil || *got.CorrelationID != corrID {
		t.Errorf("got CorrelationID %v, want %q", got.CorrelationID, corrID)
	}
}

func TestMemories_GetByCorrelationID_NotFound(t *testing.T) {
	mock := &nesttest.MockMemories{
		OnGetByCorrelationID: func(_ context.Context, _ string) (*models.Memory, error) {
			return nil, errors.New("not found")
		},
	}

	_, err := mock.GetByCorrelationID(context.Background(), "does-not-exist")
	if err == nil {
		t.Error("expected error for missing correlation ID, got nil")
	}
}
