//go:build testing

package ingest

import (
	"context"
	"errors"
	"testing"

	"github.com/zoobzio/nestor/models"
	nesttest "github.com/zoobzio/nestor/testing"
)

const (
	testChunkID1 = "550e8400-e29b-41d4-a716-446655440003"
	testChunkID2 = "550e8400-e29b-41d4-a716-446655440004"
)

// newEmbedJob returns an IngestJob ready for the embed stage (post-chunk).
func newEmbedJob(chunkIDs []string) *models.IngestJob {
	return &models.IngestJob{
		ID:         testJobID,
		CustomerID: testCustomerID,
		AgentID:    testAgentID,
		MemoryID:   testMemoryID,
		ChunkIDs:   chunkIDs,
		ChunkCount: len(chunkIDs),
		Content:    "# Hello\n\nSome markdown content.",
		Stage:      models.JobStageEmbed,
		Status:     models.JobStatusRunning,
	}
}

// TestEmbedStage_NoChunks_Skips verifies that an empty ChunkIDs list skips embedding entirely.
func TestEmbedStage_NoChunks_Skips(t *testing.T) {
	mockEmbedder := &nesttest.MockEmbedder{}
	mockChunks := &nesttest.MockChunks{}
	nesttest.SetupRegistry(t,
		nesttest.WithEmbedder(mockEmbedder),
		nesttest.WithChunks(mockChunks),
	)

	job := newEmbedJob([]string{})
	result, err := embedStage(context.Background(), job)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if result.EmbeddedCount != 0 {
		t.Errorf("EmbeddedCount = %d, want 0", result.EmbeddedCount)
	}
}

// TestEmbedStage_Success_UpdatesEmbeddedCount verifies EmbeddedCount equals the number of chunks embedded.
func TestEmbedStage_Success_UpdatesEmbeddedCount(t *testing.T) {
	chunk := &models.Chunk{
		ID:       testChunkID1,
		MemoryID: testMemoryID,
		Content:  "Some content.",
		Sequence: 0,
	}

	mockChunks := &nesttest.MockChunks{
		OnGet: func(_ context.Context, key string) (*models.Chunk, error) {
			if key == testChunkID1 {
				return chunk, nil
			}
			return nil, nil
		},
	}
	mockEmbedder := &nesttest.MockEmbedder{}
	nesttest.SetupRegistry(t,
		nesttest.WithEmbedder(mockEmbedder),
		nesttest.WithChunks(mockChunks),
	)

	job := newEmbedJob([]string{testChunkID1})
	result, err := embedStage(context.Background(), job)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if result.EmbeddedCount != 1 {
		t.Errorf("EmbeddedCount = %d, want 1", result.EmbeddedCount)
	}
}

// TestEmbedStage_BatchProcessing verifies multiple chunks are all embedded.
func TestEmbedStage_BatchProcessing(t *testing.T) {
	chunk1 := &models.Chunk{ID: testChunkID1, MemoryID: testMemoryID, Content: "Chunk one.", Sequence: 0}
	chunk2 := &models.Chunk{ID: testChunkID2, MemoryID: testMemoryID, Content: "Chunk two.", Sequence: 1}

	chunkMap := map[string]*models.Chunk{
		testChunkID1: chunk1,
		testChunkID2: chunk2,
	}

	mockChunks := &nesttest.MockChunks{
		OnGet: func(_ context.Context, key string) (*models.Chunk, error) {
			return chunkMap[key], nil
		},
	}
	mockEmbedder := &nesttest.MockEmbedder{}
	nesttest.SetupRegistry(t,
		nesttest.WithEmbedder(mockEmbedder),
		nesttest.WithChunks(mockChunks),
	)

	job := newEmbedJob([]string{testChunkID1, testChunkID2})
	result, err := embedStage(context.Background(), job)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if result.EmbeddedCount != 2 {
		t.Errorf("EmbeddedCount = %d, want 2", result.EmbeddedCount)
	}
}

// TestEmbedStage_EmbedderError verifies the embedder error propagates.
func TestEmbedStage_EmbedderError(t *testing.T) {
	chunk := &models.Chunk{
		ID:       testChunkID1,
		MemoryID: testMemoryID,
		Content:  "Some content.",
		Sequence: 0,
	}
	embedErr := errors.New("embedding service unavailable")

	mockChunks := &nesttest.MockChunks{
		OnGet: func(_ context.Context, key string) (*models.Chunk, error) {
			if key == testChunkID1 {
				return chunk, nil
			}
			return nil, nil
		},
	}
	mockEmbedder := &nesttest.MockEmbedder{
		OnEmbed: func(_ context.Context, _ []string) ([][]float32, error) {
			return nil, embedErr
		},
	}
	nesttest.SetupRegistry(t,
		nesttest.WithEmbedder(mockEmbedder),
		nesttest.WithChunks(mockChunks),
	)

	job := newEmbedJob([]string{testChunkID1})
	_, err := embedStage(context.Background(), job)

	if err == nil {
		t.Fatal("expected error, got nil")
	}
}
