//go:build testing

package handlers_test

import (
	"context"
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/zoobzio/nestor/api/handlers"
	"github.com/zoobzio/nestor/models"
	nesttest "github.com/zoobzio/nestor/testing"
	rtesting "github.com/zoobzio/rocco/testing"
)

const testMemoryID = "550e8400-e29b-41d4-a716-446655440002"

// newTestMemory returns a consistent memory for handler tests.
func newTestMemory(t *testing.T) *models.Memory {
	t.Helper()
	now := time.Now()
	return &models.Memory{
		ID:        testMemoryID,
		AgentID:   testAgentID,
		Content:   "The user prefers concise responses.",
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// TestListMemories_Success verifies a 200 response with a memories list.
func TestListMemories_Success(t *testing.T) {
	memory := newTestMemory(t)
	mock := &nesttest.MockMemories{
		OnListByAgentID: func(_ context.Context, agentID string) ([]*models.Memory, error) {
			return []*models.Memory{memory}, nil
		},
	}
	nesttest.SetupRegistry(t, nesttest.WithMemories(mock))

	engine := newTestEngine(t, handlers.ListMemories)
	capture := rtesting.ServeRequest(engine, http.MethodGet, "/agents/"+testAgentID+"/memories", nil)

	rtesting.AssertStatus(t, capture, http.StatusOK)

	var resp struct {
		Memories []struct {
			ID      string `json:"id"`
			Content string `json:"content"`
		} `json:"memories"`
	}
	if err := capture.DecodeJSON(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if len(resp.Memories) != 1 {
		t.Fatalf("expected 1 memory, got %d", len(resp.Memories))
	}
	if resp.Memories[0].ID != testMemoryID {
		t.Errorf("memory ID %q, want %q", resp.Memories[0].ID, testMemoryID)
	}
	if resp.Memories[0].Content != "The user prefers concise responses." {
		t.Errorf("memory content %q, want %q", resp.Memories[0].Content, "The user prefers concise responses.")
	}
}

// TestListMemories_Empty verifies a 200 response with an empty list.
func TestListMemories_Empty(t *testing.T) {
	mock := &nesttest.MockMemories{
		OnListByAgentID: func(_ context.Context, _ string) ([]*models.Memory, error) {
			return []*models.Memory{}, nil
		},
	}
	nesttest.SetupRegistry(t, nesttest.WithMemories(mock))

	engine := newTestEngine(t, handlers.ListMemories)
	capture := rtesting.ServeRequest(engine, http.MethodGet, "/agents/"+testAgentID+"/memories", nil)

	rtesting.AssertStatus(t, capture, http.StatusOK)

	var resp struct {
		Memories []any `json:"memories"`
	}
	if err := capture.DecodeJSON(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if len(resp.Memories) != 0 {
		t.Errorf("expected empty memories list, got %d", len(resp.Memories))
	}
}

// TestListMemories_StoreError verifies a 500 response when the store errors.
func TestListMemories_StoreError(t *testing.T) {
	mock := &nesttest.MockMemories{
		OnListByAgentID: func(_ context.Context, _ string) ([]*models.Memory, error) {
			return nil, errors.New("database unavailable")
		},
	}
	nesttest.SetupRegistry(t, nesttest.WithMemories(mock))

	engine := newTestEngine(t, handlers.ListMemories)
	capture := rtesting.ServeRequest(engine, http.MethodGet, "/agents/"+testAgentID+"/memories", nil)

	rtesting.AssertStatus(t, capture, http.StatusInternalServerError)
}

// TestCreateMemory_Success verifies a 201 response with the created memory.
func TestCreateMemory_Success(t *testing.T) {
	var storedMemory *models.Memory
	mock := &nesttest.MockMemories{
		OnSet: func(_ context.Context, _ string, memory *models.Memory) error {
			storedMemory = memory
			return nil
		},
	}
	nesttest.SetupRegistry(t, nesttest.WithMemories(mock))

	engine := newTestEngine(t, handlers.CreateMemory)

	body := map[string]string{"content": "The user prefers bullet points."}
	capture := rtesting.ServeRequest(engine, http.MethodPost, "/agents/"+testAgentID+"/memories", body)

	rtesting.AssertStatus(t, capture, http.StatusCreated)

	if storedMemory == nil {
		t.Fatal("expected memory to be stored, got nil")
	}
	if storedMemory.Content != "The user prefers bullet points." {
		t.Errorf("stored content %q, want %q", storedMemory.Content, "The user prefers bullet points.")
	}
	if storedMemory.AgentID != testAgentID {
		t.Errorf("stored AgentID %q, want %q", storedMemory.AgentID, testAgentID)
	}

	var resp struct {
		Content string `json:"content"`
	}
	if err := capture.DecodeJSON(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if resp.Content != "The user prefers bullet points." {
		t.Errorf("response content %q, want %q", resp.Content, "The user prefers bullet points.")
	}
}

// TestCreateMemory_StoreError verifies a 500 response when Set fails.
func TestCreateMemory_StoreError(t *testing.T) {
	mock := &nesttest.MockMemories{
		OnSet: func(_ context.Context, _ string, _ *models.Memory) error {
			return errors.New("storage failure")
		},
	}
	nesttest.SetupRegistry(t, nesttest.WithMemories(mock))

	engine := newTestEngine(t, handlers.CreateMemory)

	body := map[string]string{"content": "Some memory."}
	capture := rtesting.ServeRequest(engine, http.MethodPost, "/agents/"+testAgentID+"/memories", body)

	rtesting.AssertStatus(t, capture, http.StatusInternalServerError)
}

// TestGetMemory_Success verifies a 200 response with memory data.
func TestGetMemory_Success(t *testing.T) {
	memory := newTestMemory(t)
	mock := &nesttest.MockMemories{
		OnGet: func(_ context.Context, key string) (*models.Memory, error) {
			if key != memory.ID {
				return nil, errors.New("not found")
			}
			return memory, nil
		},
	}
	nesttest.SetupRegistry(t, nesttest.WithMemories(mock))

	engine := newTestEngine(t, handlers.GetMemory)
	capture := rtesting.ServeRequest(engine, http.MethodGet, "/agents/"+testAgentID+"/memories/"+testMemoryID, nil)

	rtesting.AssertStatus(t, capture, http.StatusOK)

	var resp struct {
		ID      string `json:"id"`
		Content string `json:"content"`
	}
	if err := capture.DecodeJSON(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if resp.ID != testMemoryID {
		t.Errorf("response ID %q, want %q", resp.ID, testMemoryID)
	}
	if resp.Content != "The user prefers concise responses." {
		t.Errorf("response content %q, want %q", resp.Content, "The user prefers concise responses.")
	}
}

// TestGetMemory_NotFound verifies a 404 response when the memory does not exist.
func TestGetMemory_NotFound(t *testing.T) {
	mock := &nesttest.MockMemories{
		OnGet: func(_ context.Context, _ string) (*models.Memory, error) {
			return nil, errors.New("not found")
		},
	}
	nesttest.SetupRegistry(t, nesttest.WithMemories(mock))

	engine := newTestEngine(t, handlers.GetMemory)
	capture := rtesting.ServeRequest(engine, http.MethodGet, "/agents/"+testAgentID+"/memories/does-not-exist", nil)

	rtesting.AssertStatus(t, capture, http.StatusNotFound)
	rtesting.AssertErrorCode(t, capture, "NOT_FOUND")
}

// TestUpdateMemory_Success verifies a 200 response with updated memory data.
func TestUpdateMemory_Success(t *testing.T) {
	memory := newTestMemory(t)
	mock := &nesttest.MockMemories{
		OnGet: func(_ context.Context, key string) (*models.Memory, error) {
			if key != memory.ID {
				return nil, errors.New("not found")
			}
			return memory, nil
		},
		OnSet: func(_ context.Context, _ string, updated *models.Memory) error {
			memory = updated
			return nil
		},
	}
	nesttest.SetupRegistry(t, nesttest.WithMemories(mock))

	engine := newTestEngine(t, handlers.UpdateMemory)

	body := map[string]string{"content": "Updated memory content."}
	capture := rtesting.ServeRequest(engine, http.MethodPatch, "/agents/"+testAgentID+"/memories/"+testMemoryID, body)

	rtesting.AssertStatus(t, capture, http.StatusOK)

	var resp struct {
		ID      string `json:"id"`
		Content string `json:"content"`
	}
	if err := capture.DecodeJSON(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if resp.Content != "Updated memory content." {
		t.Errorf("response content %q, want %q", resp.Content, "Updated memory content.")
	}
}

// TestUpdateMemory_NotFound verifies a 404 when the memory does not exist.
func TestUpdateMemory_NotFound(t *testing.T) {
	mock := &nesttest.MockMemories{
		OnGet: func(_ context.Context, _ string) (*models.Memory, error) {
			return nil, errors.New("not found")
		},
	}
	nesttest.SetupRegistry(t, nesttest.WithMemories(mock))

	engine := newTestEngine(t, handlers.UpdateMemory)

	body := map[string]string{"content": "Updated content."}
	capture := rtesting.ServeRequest(engine, http.MethodPatch, "/agents/"+testAgentID+"/memories/does-not-exist", body)

	rtesting.AssertStatus(t, capture, http.StatusNotFound)
	rtesting.AssertErrorCode(t, capture, "NOT_FOUND")
}

// TestUpdateMemory_SetError verifies a 500 when Set fails after a successful Get.
func TestUpdateMemory_SetError(t *testing.T) {
	memory := newTestMemory(t)
	mock := &nesttest.MockMemories{
		OnGet: func(_ context.Context, key string) (*models.Memory, error) {
			if key != memory.ID {
				return nil, errors.New("not found")
			}
			return memory, nil
		},
		OnSet: func(_ context.Context, _ string, _ *models.Memory) error {
			return errors.New("storage failure")
		},
	}
	nesttest.SetupRegistry(t, nesttest.WithMemories(mock))

	engine := newTestEngine(t, handlers.UpdateMemory)

	body := map[string]string{"content": "Updated content."}
	capture := rtesting.ServeRequest(engine, http.MethodPatch, "/agents/"+testAgentID+"/memories/"+testMemoryID, body)

	rtesting.AssertStatus(t, capture, http.StatusInternalServerError)
}

// TestDeleteMemory_Success verifies a 204 response on successful deletion.
func TestDeleteMemory_Success(t *testing.T) {
	var deletedID string
	mock := &nesttest.MockMemories{
		OnDelete: func(_ context.Context, key string) error {
			deletedID = key
			return nil
		},
	}
	nesttest.SetupRegistry(t, nesttest.WithMemories(mock))

	engine := newTestEngine(t, handlers.DeleteMemory)
	capture := rtesting.ServeRequest(engine, http.MethodDelete, "/agents/"+testAgentID+"/memories/"+testMemoryID, nil)

	rtesting.AssertStatus(t, capture, http.StatusNoContent)

	if deletedID != testMemoryID {
		t.Errorf("deleted ID %q, want %q", deletedID, testMemoryID)
	}
}

// TestDeleteMemory_NotFound verifies a 404 when the memory does not exist.
func TestDeleteMemory_NotFound(t *testing.T) {
	mock := &nesttest.MockMemories{
		OnDelete: func(_ context.Context, _ string) error {
			return errors.New("not found")
		},
	}
	nesttest.SetupRegistry(t, nesttest.WithMemories(mock))

	engine := newTestEngine(t, handlers.DeleteMemory)
	capture := rtesting.ServeRequest(engine, http.MethodDelete, "/agents/"+testAgentID+"/memories/does-not-exist", nil)

	rtesting.AssertStatus(t, capture, http.StatusNotFound)
	rtesting.AssertErrorCode(t, capture, "NOT_FOUND")
}

// TestGetMemoryByCorrelationID_Success verifies a 200 response with memory data.
func TestGetMemoryByCorrelationID_Success(t *testing.T) {
	corrID := "evt_abc123"
	memory := newTestMemory(t)
	memory.CorrelationID = &corrID

	mock := &nesttest.MockMemories{
		OnGetByCorrelationID: func(_ context.Context, id string) (*models.Memory, error) {
			if id != corrID {
				return nil, errors.New("not found")
			}
			return memory, nil
		},
	}
	nesttest.SetupRegistry(t, nesttest.WithMemories(mock))

	engine := newTestEngine(t, handlers.GetMemoryByCorrelationID)
	capture := rtesting.ServeRequest(engine, http.MethodGet, "/agents/"+testAgentID+"/memories/by-correlation/"+corrID, nil)

	rtesting.AssertStatus(t, capture, http.StatusOK)

	var resp struct {
		ID            string  `json:"id"`
		CorrelationID *string `json:"correlation_id"`
	}
	if err := capture.DecodeJSON(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if resp.ID != testMemoryID {
		t.Errorf("response ID %q, want %q", resp.ID, testMemoryID)
	}
	if resp.CorrelationID == nil || *resp.CorrelationID != corrID {
		t.Errorf("response CorrelationID %v, want %q", resp.CorrelationID, corrID)
	}
}

// TestGetMemoryByCorrelationID_NotFound verifies a 404 when the correlation ID does not exist.
func TestGetMemoryByCorrelationID_NotFound(t *testing.T) {
	mock := &nesttest.MockMemories{
		OnGetByCorrelationID: func(_ context.Context, _ string) (*models.Memory, error) {
			return nil, errors.New("not found")
		},
	}
	nesttest.SetupRegistry(t, nesttest.WithMemories(mock))

	engine := newTestEngine(t, handlers.GetMemoryByCorrelationID)
	capture := rtesting.ServeRequest(engine, http.MethodGet, "/agents/"+testAgentID+"/memories/by-correlation/does-not-exist", nil)

	rtesting.AssertStatus(t, capture, http.StatusNotFound)
	rtesting.AssertErrorCode(t, capture, "NOT_FOUND")
}
