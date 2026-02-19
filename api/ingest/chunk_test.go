//go:build testing

package ingest

import (
	"context"
	"errors"
	"testing"

	"github.com/zoobzio/nestor/api/contracts"
	"github.com/zoobzio/nestor/models"
	nesttest "github.com/zoobzio/nestor/testing"
)

const testMemoryID = "550e8400-e29b-41d4-a716-446655440002"

// newChunkJob returns an IngestJob ready for the chunk stage (post-store).
func newChunkJob() *models.IngestJob {
	return &models.IngestJob{
		ID:         testJobID,
		CustomerID: testCustomerID,
		AgentID:    testAgentID,
		MemoryID:   testMemoryID,
		Content:    "# Hello\n\nSome markdown content.",
		Stage:      models.JobStageChunk,
		Status:     models.JobStatusRunning,
	}
}

// TestChunkStage_Success_CountsChunks verifies ChunkCount and ChunkIDs are populated.
func TestChunkStage_Success_CountsChunks(t *testing.T) {
	parsed := []contracts.Chunk{
		{Kind: "section", Symbol: "Hello", Content: "Some markdown content.", StartLine: 1, EndLine: 3},
	}
	mockChunker := &nesttest.MockChunker{
		OnChunk: func(_ context.Context, _ []byte) ([]contracts.Chunk, error) {
			return parsed, nil
		},
	}
	mockChunks := &nesttest.MockChunks{}
	nesttest.SetupRegistry(t,
		nesttest.WithChunker(mockChunker),
		nesttest.WithChunks(mockChunks),
	)

	job := newChunkJob()
	result, err := chunkStage(context.Background(), job)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	if result.ChunkCount != 1 {
		t.Errorf("ChunkCount = %d, want 1", result.ChunkCount)
	}
	if len(result.ChunkIDs) != 1 {
		t.Errorf("len(ChunkIDs) = %d, want 1", len(result.ChunkIDs))
	}
	if result.ChunkIDs[0] == "" {
		t.Error("expected non-empty ChunkID, got empty string")
	}
}

// TestChunkStage_Success_ChunkMetadata verifies Kind, Symbol, and Context are mapped correctly.
func TestChunkStage_Success_ChunkMetadata(t *testing.T) {
	symbol := "Introduction"
	ctxHeaders := []string{"Document", "Introduction"}
	parsed := []contracts.Chunk{
		{
			Kind:      "section",
			Symbol:    symbol,
			Context:   ctxHeaders,
			Content:   "Some content here.",
			StartLine: 1,
			EndLine:   5,
		},
	}

	var storedChunk *models.Chunk
	mockChunker := &nesttest.MockChunker{
		OnChunk: func(_ context.Context, _ []byte) ([]contracts.Chunk, error) {
			return parsed, nil
		},
	}
	mockChunks := &nesttest.MockChunks{
		OnSet: func(_ context.Context, _ string, c *models.Chunk) error {
			storedChunk = c
			return nil
		},
	}
	nesttest.SetupRegistry(t,
		nesttest.WithChunker(mockChunker),
		nesttest.WithChunks(mockChunks),
	)

	job := newChunkJob()
	if _, err := chunkStage(context.Background(), job); err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	if storedChunk == nil {
		t.Fatal("expected chunk to be stored, got nil")
	}
	if string(storedChunk.Kind) != "section" {
		t.Errorf("Kind = %q, want %q", storedChunk.Kind, "section")
	}
	if storedChunk.Symbol == nil || *storedChunk.Symbol != symbol {
		t.Errorf("Symbol = %v, want %q", storedChunk.Symbol, symbol)
	}
	if len(storedChunk.Context) != 2 {
		t.Errorf("len(Context) = %d, want 2", len(storedChunk.Context))
	}
	if storedChunk.MemoryID != testMemoryID {
		t.Errorf("MemoryID = %q, want %q", storedChunk.MemoryID, testMemoryID)
	}
}

// TestChunkStage_EmptyContent_ZeroChunks verifies that zero parsed chunks yields ChunkCount of 0.
func TestChunkStage_EmptyContent_ZeroChunks(t *testing.T) {
	mockChunker := &nesttest.MockChunker{
		OnChunk: func(_ context.Context, _ []byte) ([]contracts.Chunk, error) {
			return []contracts.Chunk{}, nil
		},
	}
	mockChunks := &nesttest.MockChunks{}
	nesttest.SetupRegistry(t,
		nesttest.WithChunker(mockChunker),
		nesttest.WithChunks(mockChunks),
	)

	job := newChunkJob()
	job.Content = ""

	result, err := chunkStage(context.Background(), job)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	if result.ChunkCount != 0 {
		t.Errorf("ChunkCount = %d, want 0", result.ChunkCount)
	}
	if len(result.ChunkIDs) != 0 {
		t.Errorf("len(ChunkIDs) = %d, want 0", len(result.ChunkIDs))
	}
}

// TestChunkStage_ChunkerError verifies the chunker error propagates.
func TestChunkStage_ChunkerError(t *testing.T) {
	chunkerErr := errors.New("parse failure")
	mockChunker := &nesttest.MockChunker{
		OnChunk: func(_ context.Context, _ []byte) ([]contracts.Chunk, error) {
			return nil, chunkerErr
		},
	}
	mockChunks := &nesttest.MockChunks{}
	nesttest.SetupRegistry(t,
		nesttest.WithChunker(mockChunker),
		nesttest.WithChunks(mockChunks),
	)

	job := newChunkJob()
	_, err := chunkStage(context.Background(), job)

	if !errors.Is(err, chunkerErr) {
		t.Errorf("error = %v, want %v", err, chunkerErr)
	}
}

// TestChunkStage_StoreError verifies that a chunk store failure propagates.
func TestChunkStage_StoreError(t *testing.T) {
	parsed := []contracts.Chunk{
		{Kind: "section", Content: "Some content.", StartLine: 1, EndLine: 2},
	}
	storeErr := errors.New("storage unavailable")
	mockChunker := &nesttest.MockChunker{
		OnChunk: func(_ context.Context, _ []byte) ([]contracts.Chunk, error) {
			return parsed, nil
		},
	}
	mockChunks := &nesttest.MockChunks{
		OnSet: func(_ context.Context, _ string, _ *models.Chunk) error {
			return storeErr
		},
	}
	nesttest.SetupRegistry(t,
		nesttest.WithChunker(mockChunker),
		nesttest.WithChunks(mockChunks),
	)

	job := newChunkJob()
	_, err := chunkStage(context.Background(), job)

	if !errors.Is(err, storeErr) {
		t.Errorf("error = %v, want %v", err, storeErr)
	}
}
