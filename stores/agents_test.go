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

// These tests exercise the contracts.Agents interface via MockAgents.
// The concrete Agents store (which embeds sum.Database) requires a live
// database and is covered by integration tests.

func TestAgents_Get_Found(t *testing.T) {
	agent := nesttest.NewAgent(t)

	mock := &nesttest.MockAgents{
		OnGet: func(_ context.Context, key string) (*models.Agent, error) {
			if key != agent.ID {
				return nil, errors.New("not found")
			}
			return agent, nil
		},
	}

	got, err := mock.Get(context.Background(), agent.ID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.ID != agent.ID {
		t.Errorf("got ID %q, want %q", got.ID, agent.ID)
	}
}

func TestAgents_Get_NotFound(t *testing.T) {
	mock := &nesttest.MockAgents{
		OnGet: func(_ context.Context, _ string) (*models.Agent, error) {
			return nil, errors.New("not found")
		},
	}

	_, err := mock.Get(context.Background(), "does-not-exist")
	if err == nil {
		t.Error("expected error for missing agent, got nil")
	}
}

func TestAgents_Set(t *testing.T) {
	var stored *models.Agent

	mock := &nesttest.MockAgents{
		OnSet: func(_ context.Context, _ string, agent *models.Agent) error {
			stored = agent
			return nil
		},
	}

	agent := nesttest.NewAgent(t)
	if err := mock.Set(context.Background(), agent.ID, agent); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if stored == nil {
		t.Fatal("expected agent to be stored, got nil")
	}
	if stored.ID != agent.ID {
		t.Errorf("stored ID %q, want %q", stored.ID, agent.ID)
	}
}

func TestAgents_Set_Error(t *testing.T) {
	mock := &nesttest.MockAgents{
		OnSet: func(_ context.Context, _ string, _ *models.Agent) error {
			return errors.New("storage error")
		},
	}

	agent := nesttest.NewAgent(t)
	if err := mock.Set(context.Background(), agent.ID, agent); err == nil {
		t.Error("expected error from Set, got nil")
	}
}

func TestAgents_Delete(t *testing.T) {
	var deletedKey string

	mock := &nesttest.MockAgents{
		OnDelete: func(_ context.Context, key string) error {
			deletedKey = key
			return nil
		},
	}

	if err := mock.Delete(context.Background(), "agent-id"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if deletedKey != "agent-id" {
		t.Errorf("deleted key %q, want %q", deletedKey, "agent-id")
	}
}

func TestAgents_Delete_Error(t *testing.T) {
	mock := &nesttest.MockAgents{
		OnDelete: func(_ context.Context, _ string) error {
			return errors.New("not found")
		},
	}

	if err := mock.Delete(context.Background(), "missing-id"); err == nil {
		t.Error("expected error from Delete, got nil")
	}
}

func TestAgents_ListByUserID(t *testing.T) {
	now := time.Now()
	userID := "550e8400-e29b-41d4-a716-446655440000"
	expected := []*models.Agent{
		{ID: "agent-1", UserID: userID, Name: "First", CreatedAt: now, UpdatedAt: now},
		{ID: "agent-2", UserID: userID, Name: "Second", CreatedAt: now, UpdatedAt: now},
	}

	mock := &nesttest.MockAgents{
		OnListByUserID: func(_ context.Context, id string) ([]*models.Agent, error) {
			if id != userID {
				return []*models.Agent{}, nil
			}
			return expected, nil
		},
	}

	list, err := mock.ListByUserID(context.Background(), userID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(list) != 2 {
		t.Fatalf("expected 2 agents, got %d", len(list))
	}
	if list[0].ID != "agent-1" {
		t.Errorf("first agent ID %q, want %q", list[0].ID, "agent-1")
	}
	if list[1].ID != "agent-2" {
		t.Errorf("second agent ID %q, want %q", list[1].ID, "agent-2")
	}
}

func TestAgents_ListByUserID_Empty(t *testing.T) {
	mock := &nesttest.MockAgents{
		OnListByUserID: func(_ context.Context, _ string) ([]*models.Agent, error) {
			return []*models.Agent{}, nil
		},
	}

	list, err := mock.ListByUserID(context.Background(), "user-with-no-agents")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(list) != 0 {
		t.Errorf("expected empty list, got %d agents", len(list))
	}
}

func TestAgents_ListByUserID_Error(t *testing.T) {
	mock := &nesttest.MockAgents{
		OnListByUserID: func(_ context.Context, _ string) ([]*models.Agent, error) {
			return nil, errors.New("database error")
		},
	}

	_, err := mock.ListByUserID(context.Background(), "any-user")
	if err == nil {
		t.Error("expected error from ListByUserID, got nil")
	}
}
