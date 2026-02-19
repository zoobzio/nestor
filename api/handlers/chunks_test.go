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

const testChunkID = "550e8400-e29b-41d4-a716-446655440003"

// newTestChunk returns a consistent chunk for handler tests.
func newTestChunk(t *testing.T) *models.Chunk {
	t.Helper()
	now := time.Now()
	return &models.Chunk{
		ID:        testChunkID,
		MemoryID:  testMemoryID,
		Content:   "The user prefers concise responses.",
		Sequence:  0,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// TestListChunks_Success verifies a 200 response with a chunks list.
func TestListChunks_Success(t *testing.T) {
	chunk := newTestChunk(t)
	mock := &nesttest.MockChunks{
		OnListByMemoryID: func(_ context.Context, memoryID string) ([]*models.Chunk, error) {
			return []*models.Chunk{chunk}, nil
		},
	}
	nesttest.SetupRegistry(t, nesttest.WithChunks(mock))

	engine := newTestEngine(t, handlers.ListChunks)
	capture := rtesting.ServeRequest(engine, http.MethodGet, "/agents/"+testAgentID+"/memories/"+testMemoryID+"/chunks", nil)

	rtesting.AssertStatus(t, capture, http.StatusOK)

	var resp struct {
		Chunks []struct {
			ID      string `json:"id"`
			Content string `json:"content"`
		} `json:"chunks"`
	}
	if err := capture.DecodeJSON(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if len(resp.Chunks) != 1 {
		t.Fatalf("expected 1 chunk, got %d", len(resp.Chunks))
	}
	if resp.Chunks[0].ID != testChunkID {
		t.Errorf("chunk ID %q, want %q", resp.Chunks[0].ID, testChunkID)
	}
	if resp.Chunks[0].Content != "The user prefers concise responses." {
		t.Errorf("chunk content %q, want %q", resp.Chunks[0].Content, "The user prefers concise responses.")
	}
}

// TestListChunks_Empty verifies a 200 response with an empty list.
func TestListChunks_Empty(t *testing.T) {
	mock := &nesttest.MockChunks{
		OnListByMemoryID: func(_ context.Context, _ string) ([]*models.Chunk, error) {
			return []*models.Chunk{}, nil
		},
	}
	nesttest.SetupRegistry(t, nesttest.WithChunks(mock))

	engine := newTestEngine(t, handlers.ListChunks)
	capture := rtesting.ServeRequest(engine, http.MethodGet, "/agents/"+testAgentID+"/memories/"+testMemoryID+"/chunks", nil)

	rtesting.AssertStatus(t, capture, http.StatusOK)

	var resp struct {
		Chunks []any `json:"chunks"`
	}
	if err := capture.DecodeJSON(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if len(resp.Chunks) != 0 {
		t.Errorf("expected empty chunks list, got %d", len(resp.Chunks))
	}
}

// TestListChunks_StoreError verifies a 500 response when the store errors.
func TestListChunks_StoreError(t *testing.T) {
	mock := &nesttest.MockChunks{
		OnListByMemoryID: func(_ context.Context, _ string) ([]*models.Chunk, error) {
			return nil, errors.New("database unavailable")
		},
	}
	nesttest.SetupRegistry(t, nesttest.WithChunks(mock))

	engine := newTestEngine(t, handlers.ListChunks)
	capture := rtesting.ServeRequest(engine, http.MethodGet, "/agents/"+testAgentID+"/memories/"+testMemoryID+"/chunks", nil)

	rtesting.AssertStatus(t, capture, http.StatusInternalServerError)
}

// TestCreateChunk_Success verifies a 201 response with the created chunk.
func TestCreateChunk_Success(t *testing.T) {
	var storedChunk *models.Chunk
	mock := &nesttest.MockChunks{
		OnSet: func(_ context.Context, _ string, chunk *models.Chunk) error {
			storedChunk = chunk
			return nil
		},
	}
	nesttest.SetupRegistry(t, nesttest.WithChunks(mock))

	engine := newTestEngine(t, handlers.CreateChunk)

	body := map[string]any{
		"content":   "The user prefers bullet points.",
		"sequence":  1,
		"embedding": []float32{0.1, 0.2, 0.3},
	}
	capture := rtesting.ServeRequest(engine, http.MethodPost, "/agents/"+testAgentID+"/memories/"+testMemoryID+"/chunks", body)

	rtesting.AssertStatus(t, capture, http.StatusCreated)

	if storedChunk == nil {
		t.Fatal("expected chunk to be stored, got nil")
	}
	if storedChunk.Content != "The user prefers bullet points." {
		t.Errorf("stored content %q, want %q", storedChunk.Content, "The user prefers bullet points.")
	}
	if storedChunk.MemoryID != testMemoryID {
		t.Errorf("stored MemoryID %q, want %q", storedChunk.MemoryID, testMemoryID)
	}
	if storedChunk.Sequence != 1 {
		t.Errorf("stored Sequence %d, want 1", storedChunk.Sequence)
	}

	var resp struct {
		Content  string `json:"content"`
		Sequence int    `json:"sequence"`
	}
	if err := capture.DecodeJSON(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if resp.Content != "The user prefers bullet points." {
		t.Errorf("response content %q, want %q", resp.Content, "The user prefers bullet points.")
	}
	if resp.Sequence != 1 {
		t.Errorf("response Sequence %d, want 1", resp.Sequence)
	}
}

// TestCreateChunk_StoreError verifies a 500 response when Set fails.
func TestCreateChunk_StoreError(t *testing.T) {
	mock := &nesttest.MockChunks{
		OnSet: func(_ context.Context, _ string, _ *models.Chunk) error {
			return errors.New("storage failure")
		},
	}
	nesttest.SetupRegistry(t, nesttest.WithChunks(mock))

	engine := newTestEngine(t, handlers.CreateChunk)

	body := map[string]any{"content": "Some chunk.", "sequence": 0}
	capture := rtesting.ServeRequest(engine, http.MethodPost, "/agents/"+testAgentID+"/memories/"+testMemoryID+"/chunks", body)

	rtesting.AssertStatus(t, capture, http.StatusInternalServerError)
}

// TestGetChunk_Success verifies a 200 response with chunk data.
func TestGetChunk_Success(t *testing.T) {
	chunk := newTestChunk(t)
	mock := &nesttest.MockChunks{
		OnGet: func(_ context.Context, key string) (*models.Chunk, error) {
			if key != chunk.ID {
				return nil, errors.New("not found")
			}
			return chunk, nil
		},
	}
	nesttest.SetupRegistry(t, nesttest.WithChunks(mock))

	engine := newTestEngine(t, handlers.GetChunk)
	capture := rtesting.ServeRequest(engine, http.MethodGet, "/agents/"+testAgentID+"/memories/"+testMemoryID+"/chunks/"+testChunkID, nil)

	rtesting.AssertStatus(t, capture, http.StatusOK)

	var resp struct {
		ID      string `json:"id"`
		Content string `json:"content"`
	}
	if err := capture.DecodeJSON(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if resp.ID != testChunkID {
		t.Errorf("response ID %q, want %q", resp.ID, testChunkID)
	}
	if resp.Content != "The user prefers concise responses." {
		t.Errorf("response content %q, want %q", resp.Content, "The user prefers concise responses.")
	}
}

// TestGetChunk_NotFound verifies a 404 response when the chunk does not exist.
func TestGetChunk_NotFound(t *testing.T) {
	mock := &nesttest.MockChunks{
		OnGet: func(_ context.Context, _ string) (*models.Chunk, error) {
			return nil, errors.New("not found")
		},
	}
	nesttest.SetupRegistry(t, nesttest.WithChunks(mock))

	engine := newTestEngine(t, handlers.GetChunk)
	capture := rtesting.ServeRequest(engine, http.MethodGet, "/agents/"+testAgentID+"/memories/"+testMemoryID+"/chunks/does-not-exist", nil)

	rtesting.AssertStatus(t, capture, http.StatusNotFound)
	rtesting.AssertErrorCode(t, capture, "NOT_FOUND")
}

// TestDeleteChunk_Success verifies a 204 response on successful deletion.
func TestDeleteChunk_Success(t *testing.T) {
	var deletedID string
	mock := &nesttest.MockChunks{
		OnDelete: func(_ context.Context, key string) error {
			deletedID = key
			return nil
		},
	}
	nesttest.SetupRegistry(t, nesttest.WithChunks(mock))

	engine := newTestEngine(t, handlers.DeleteChunk)
	capture := rtesting.ServeRequest(engine, http.MethodDelete, "/agents/"+testAgentID+"/memories/"+testMemoryID+"/chunks/"+testChunkID, nil)

	rtesting.AssertStatus(t, capture, http.StatusNoContent)

	if deletedID != testChunkID {
		t.Errorf("deleted ID %q, want %q", deletedID, testChunkID)
	}
}

// TestDeleteChunk_NotFound verifies a 404 when the chunk does not exist.
func TestDeleteChunk_NotFound(t *testing.T) {
	mock := &nesttest.MockChunks{
		OnDelete: func(_ context.Context, _ string) error {
			return errors.New("not found")
		},
	}
	nesttest.SetupRegistry(t, nesttest.WithChunks(mock))

	engine := newTestEngine(t, handlers.DeleteChunk)
	capture := rtesting.ServeRequest(engine, http.MethodDelete, "/agents/"+testAgentID+"/memories/"+testMemoryID+"/chunks/does-not-exist", nil)

	rtesting.AssertStatus(t, capture, http.StatusNotFound)
	rtesting.AssertErrorCode(t, capture, "NOT_FOUND")
}

// TestSearchChunks_Success verifies a 200 response with ordered search results.
func TestSearchChunks_Success(t *testing.T) {
	now := time.Now()
	results := []*models.Chunk{
		{ID: "chunk-1", MemoryID: "memory-1", Content: "Most similar.", Sequence: 0, CreatedAt: now, UpdatedAt: now},
		{ID: "chunk-2", MemoryID: "memory-2", Content: "Less similar.", Sequence: 0, CreatedAt: now, UpdatedAt: now},
	}

	var capturedLimit int
	mock := &nesttest.MockChunks{
		OnSearchByEmbedding: func(_ context.Context, _ []float32, limit int) ([]*models.Chunk, error) {
			capturedLimit = limit
			return results, nil
		},
	}
	nesttest.SetupRegistry(t, nesttest.WithChunks(mock))

	engine := newTestEngine(t, handlers.SearchChunks)

	body := map[string]any{
		"embedding": []float32{0.1, 0.2, 0.3},
		"limit":     5,
	}
	capture := rtesting.ServeRequest(engine, http.MethodPost, "/chunks/search", body)

	rtesting.AssertStatus(t, capture, http.StatusOK)

	if capturedLimit != 5 {
		t.Errorf("captured limit %d, want 5", capturedLimit)
	}

	var resp struct {
		Results []struct {
			ID      string `json:"id"`
			Content string `json:"content"`
		} `json:"results"`
	}
	if err := capture.DecodeJSON(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if len(resp.Results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(resp.Results))
	}
	if resp.Results[0].ID != "chunk-1" {
		t.Errorf("first result ID %q, want %q", resp.Results[0].ID, "chunk-1")
	}
	if resp.Results[0].Content != "Most similar." {
		t.Errorf("first result content %q, want %q", resp.Results[0].Content, "Most similar.")
	}
}

// TestSearchChunks_Empty verifies a 200 response with an empty results list.
func TestSearchChunks_Empty(t *testing.T) {
	mock := &nesttest.MockChunks{
		OnSearchByEmbedding: func(_ context.Context, _ []float32, _ int) ([]*models.Chunk, error) {
			return []*models.Chunk{}, nil
		},
	}
	nesttest.SetupRegistry(t, nesttest.WithChunks(mock))

	engine := newTestEngine(t, handlers.SearchChunks)

	body := map[string]any{
		"embedding": []float32{0.1, 0.2},
		"limit":     10,
	}
	capture := rtesting.ServeRequest(engine, http.MethodPost, "/chunks/search", body)

	rtesting.AssertStatus(t, capture, http.StatusOK)

	var resp struct {
		Results []any `json:"results"`
	}
	if err := capture.DecodeJSON(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if len(resp.Results) != 0 {
		t.Errorf("expected empty results, got %d", len(resp.Results))
	}
}

// TestSearchChunks_StoreError verifies a 500 response when SearchByEmbedding fails.
func TestSearchChunks_StoreError(t *testing.T) {
	mock := &nesttest.MockChunks{
		OnSearchByEmbedding: func(_ context.Context, _ []float32, _ int) ([]*models.Chunk, error) {
			return nil, errors.New("database unavailable")
		},
	}
	nesttest.SetupRegistry(t, nesttest.WithChunks(mock))

	engine := newTestEngine(t, handlers.SearchChunks)

	body := map[string]any{
		"embedding": []float32{0.1, 0.2},
		"limit":     10,
	}
	capture := rtesting.ServeRequest(engine, http.MethodPost, "/chunks/search", body)

	rtesting.AssertStatus(t, capture, http.StatusInternalServerError)
}
