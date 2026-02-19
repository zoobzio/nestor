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

// setupFullPipeline registers all dependencies required for a successful end-to-end run.
// Returns the job submitted. The in-memory stores allow realistic state threading
// between stages without a real database.
func setupFullPipeline(t *testing.T) {
	t.Helper()

	agent := &models.Agent{ID: testAgentID}
	memory := &models.Memory{
		ID:      testMemoryID,
		AgentID: testAgentID,
		Status:  models.MemoryStatusProcessing,
	}

	// In-memory chunk store keyed by ID.
	chunkStore := map[string]*models.Chunk{}

	mockAgents := &nesttest.MockAgents{
		OnGet: func(_ context.Context, _ string) (*models.Agent, error) {
			return agent, nil
		},
	}

	mockMemories := &nesttest.MockMemories{
		OnSet: func(_ context.Context, key string, m *models.Memory) error {
			memory = m
			memory.ID = key
			return nil
		},
		OnGet: func(_ context.Context, _ string) (*models.Memory, error) {
			return memory, nil
		},
	}

	mockChunks := &nesttest.MockChunks{
		OnSet: func(_ context.Context, key string, c *models.Chunk) error {
			chunkStore[key] = c
			return nil
		},
		OnGet: func(_ context.Context, key string) (*models.Chunk, error) {
			return chunkStore[key], nil
		},
	}

	mockChunker := &nesttest.MockChunker{
		OnChunk: func(_ context.Context, _ []byte) ([]contracts.Chunk, error) {
			return []contracts.Chunk{
				{Kind: "section", Symbol: "Hello", Content: "Some content.", StartLine: 1, EndLine: 3},
			}, nil
		},
	}

	mockEmbedder := &nesttest.MockEmbedder{}

	nesttest.SetupRegistry(t,
		nesttest.WithAgents(mockAgents),
		nesttest.WithMemories(mockMemories),
		nesttest.WithChunks(mockChunks),
		nesttest.WithChunker(mockChunker),
		nesttest.WithEmbedder(mockEmbedder),
	)
}

// newPipelineJob returns a job that should traverse the full pipeline successfully.
func newPipelineJob() *models.IngestJob {
	return &models.IngestJob{
		ID:         testJobID,
		CustomerID: testCustomerID,
		AgentID:    testAgentID,
		Content:    "# Hello\n\nSome markdown content.",
		Stage:      models.JobStageValidate,
		Status:     models.JobStatusPending,
	}
}

// TestPipeline_FullSuccess verifies all five stages run and the job reaches "completed" status.
func TestPipeline_FullSuccess(t *testing.T) {
	setupFullPipeline(t)

	pipeline := NewPipeline()
	defer pipeline.Close()

	job := newPipelineJob()
	result, err := pipeline.Process(context.Background(), job)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	if result.Status != models.JobStatusCompleted {
		t.Errorf("status = %q, want %q", result.Status, models.JobStatusCompleted)
	}
	if result.MemoryID == "" {
		t.Error("expected MemoryID to be set, got empty string")
	}
	if result.ChunkCount == 0 {
		t.Error("expected ChunkCount > 0, got 0")
	}
	if result.CompletedAt == nil {
		t.Error("expected CompletedAt to be set, got nil")
	}
}

// TestPipeline_StopsOnValidateFailure verifies the pipeline halts when validate fails.
func TestPipeline_StopsOnValidateFailure(t *testing.T) {
	// Agent not found — validate will return ErrAgentNotFound.
	mockAgents := &nesttest.MockAgents{
		OnGet: func(_ context.Context, _ string) (*models.Agent, error) {
			return nil, nil
		},
	}
	nesttest.SetupRegistry(t, nesttest.WithAgents(mockAgents))

	pipeline := NewPipeline()
	defer pipeline.Close()

	job := newPipelineJob()
	_, err := pipeline.Process(context.Background(), job)

	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, ErrAgentNotFound) {
		t.Errorf("error = %v, want ErrAgentNotFound", err)
	}
}

// TestPipeline_StopsOnStoreFailure verifies the pipeline halts when the store stage fails.
func TestPipeline_StopsOnStoreFailure(t *testing.T) {
	agent := &models.Agent{ID: testAgentID}
	storeErr := errors.New("disk full")

	mockAgents := &nesttest.MockAgents{
		OnGet: func(_ context.Context, _ string) (*models.Agent, error) {
			return agent, nil
		},
	}
	mockMemories := &nesttest.MockMemories{
		OnSet: func(_ context.Context, _ string, _ *models.Memory) error {
			return storeErr
		},
	}

	nesttest.SetupRegistry(t,
		nesttest.WithAgents(mockAgents),
		nesttest.WithMemories(mockMemories),
	)

	pipeline := NewPipeline()
	defer pipeline.Close()

	job := newPipelineJob()
	_, err := pipeline.Process(context.Background(), job)

	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, storeErr) {
		t.Errorf("error = %v, want %v", err, storeErr)
	}
}
