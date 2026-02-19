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

// These tests exercise the contracts.Chunks interface via MockChunks.
// The concrete Chunks store (which embeds sum.Database) requires a live
// database and is covered by integration tests.

func TestChunks_Get_Found(t *testing.T) {
	chunk := nesttest.NewChunk(t)

	mock := &nesttest.MockChunks{
		OnGet: func(_ context.Context, key string) (*models.Chunk, error) {
			if key != chunk.ID {
				return nil, errors.New("not found")
			}
			return chunk, nil
		},
	}

	got, err := mock.Get(context.Background(), chunk.ID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.ID != chunk.ID {
		t.Errorf("got ID %q, want %q", got.ID, chunk.ID)
	}
}

func TestChunks_Get_NotFound(t *testing.T) {
	mock := &nesttest.MockChunks{
		OnGet: func(_ context.Context, _ string) (*models.Chunk, error) {
			return nil, errors.New("not found")
		},
	}

	_, err := mock.Get(context.Background(), "does-not-exist")
	if err == nil {
		t.Error("expected error for missing chunk, got nil")
	}
}

func TestChunks_Set(t *testing.T) {
	var stored *models.Chunk

	mock := &nesttest.MockChunks{
		OnSet: func(_ context.Context, _ string, chunk *models.Chunk) error {
			stored = chunk
			return nil
		},
	}

	chunk := nesttest.NewChunk(t)
	if err := mock.Set(context.Background(), chunk.ID, chunk); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if stored == nil {
		t.Fatal("expected chunk to be stored, got nil")
	}
	if stored.ID != chunk.ID {
		t.Errorf("stored ID %q, want %q", stored.ID, chunk.ID)
	}
}

func TestChunks_Set_Error(t *testing.T) {
	mock := &nesttest.MockChunks{
		OnSet: func(_ context.Context, _ string, _ *models.Chunk) error {
			return errors.New("storage error")
		},
	}

	chunk := nesttest.NewChunk(t)
	if err := mock.Set(context.Background(), chunk.ID, chunk); err == nil {
		t.Error("expected error from Set, got nil")
	}
}

func TestChunks_Delete(t *testing.T) {
	var deletedKey string

	mock := &nesttest.MockChunks{
		OnDelete: func(_ context.Context, key string) error {
			deletedKey = key
			return nil
		},
	}

	if err := mock.Delete(context.Background(), "chunk-id"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if deletedKey != "chunk-id" {
		t.Errorf("deleted key %q, want %q", deletedKey, "chunk-id")
	}
}

func TestChunks_Delete_Error(t *testing.T) {
	mock := &nesttest.MockChunks{
		OnDelete: func(_ context.Context, _ string) error {
			return errors.New("not found")
		},
	}

	if err := mock.Delete(context.Background(), "missing-id"); err == nil {
		t.Error("expected error from Delete, got nil")
	}
}

func TestChunks_ListByMemoryID(t *testing.T) {
	now := time.Now()
	memoryID := "550e8400-e29b-41d4-a716-446655440002"
	expected := []*models.Chunk{
		{ID: "chunk-1", MemoryID: memoryID, Content: "First chunk.", Sequence: 0, CreatedAt: now, UpdatedAt: now},
		{ID: "chunk-2", MemoryID: memoryID, Content: "Second chunk.", Sequence: 1, CreatedAt: now, UpdatedAt: now},
	}

	mock := &nesttest.MockChunks{
		OnListByMemoryID: func(_ context.Context, id string) ([]*models.Chunk, error) {
			if id != memoryID {
				return []*models.Chunk{}, nil
			}
			return expected, nil
		},
	}

	list, err := mock.ListByMemoryID(context.Background(), memoryID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(list) != 2 {
		t.Fatalf("expected 2 chunks, got %d", len(list))
	}
	if list[0].ID != "chunk-1" {
		t.Errorf("first chunk ID %q, want %q", list[0].ID, "chunk-1")
	}
	if list[1].ID != "chunk-2" {
		t.Errorf("second chunk ID %q, want %q", list[1].ID, "chunk-2")
	}
	if list[0].Sequence != 0 {
		t.Errorf("first chunk Sequence %d, want 0", list[0].Sequence)
	}
	if list[1].Sequence != 1 {
		t.Errorf("second chunk Sequence %d, want 1", list[1].Sequence)
	}
}

func TestChunks_ListByMemoryID_Empty(t *testing.T) {
	mock := &nesttest.MockChunks{
		OnListByMemoryID: func(_ context.Context, _ string) ([]*models.Chunk, error) {
			return []*models.Chunk{}, nil
		},
	}

	list, err := mock.ListByMemoryID(context.Background(), "memory-with-no-chunks")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(list) != 0 {
		t.Errorf("expected empty list, got %d chunks", len(list))
	}
}

func TestChunks_ListByMemoryID_Error(t *testing.T) {
	mock := &nesttest.MockChunks{
		OnListByMemoryID: func(_ context.Context, _ string) ([]*models.Chunk, error) {
			return nil, errors.New("database error")
		},
	}

	_, err := mock.ListByMemoryID(context.Background(), "any-memory")
	if err == nil {
		t.Error("expected error from ListByMemoryID, got nil")
	}
}

func TestChunks_SearchByEmbedding(t *testing.T) {
	now := time.Now()
	query := []float32{0.1, 0.2, 0.3}
	expected := []*models.Chunk{
		{ID: "chunk-1", MemoryID: "memory-1", Content: "Most similar.", Sequence: 0, CreatedAt: now, UpdatedAt: now},
		{ID: "chunk-2", MemoryID: "memory-2", Content: "Less similar.", Sequence: 0, CreatedAt: now, UpdatedAt: now},
	}

	var capturedEmbedding []float32
	var capturedLimit int

	mock := &nesttest.MockChunks{
		OnSearchByEmbedding: func(_ context.Context, embedding []float32, limit int) ([]*models.Chunk, error) {
			capturedEmbedding = embedding
			capturedLimit = limit
			return expected, nil
		},
	}

	results, err := mock.SearchByEmbedding(context.Background(), query, 5)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
	if results[0].ID != "chunk-1" {
		t.Errorf("first result ID %q, want %q", results[0].ID, "chunk-1")
	}
	if len(capturedEmbedding) != len(query) {
		t.Errorf("captured embedding length %d, want %d", len(capturedEmbedding), len(query))
	}
	if capturedLimit != 5 {
		t.Errorf("captured limit %d, want 5", capturedLimit)
	}
}

func TestChunks_SearchByEmbedding_Empty(t *testing.T) {
	mock := &nesttest.MockChunks{
		OnSearchByEmbedding: func(_ context.Context, _ []float32, _ int) ([]*models.Chunk, error) {
			return []*models.Chunk{}, nil
		},
	}

	results, err := mock.SearchByEmbedding(context.Background(), []float32{0.1, 0.2}, 10)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 0 {
		t.Errorf("expected empty results, got %d chunks", len(results))
	}
}

func TestChunks_SearchByEmbedding_Error(t *testing.T) {
	mock := &nesttest.MockChunks{
		OnSearchByEmbedding: func(_ context.Context, _ []float32, _ int) ([]*models.Chunk, error) {
			return nil, errors.New("database error")
		},
	}

	_, err := mock.SearchByEmbedding(context.Background(), []float32{0.1}, 10)
	if err == nil {
		t.Error("expected error from SearchByEmbedding, got nil")
	}
}
