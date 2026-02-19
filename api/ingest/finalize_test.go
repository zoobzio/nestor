//go:build testing

package ingest

import (
	"context"
	"errors"
	"testing"

	"github.com/zoobzio/nestor/models"
	nesttest "github.com/zoobzio/nestor/testing"
)

// newFinalizeJob returns an IngestJob ready for the finalize stage (post-embed).
func newFinalizeJob() *models.IngestJob {
	return &models.IngestJob{
		ID:            testJobID,
		CustomerID:    testCustomerID,
		AgentID:       testAgentID,
		MemoryID:      testMemoryID,
		ChunkIDs:      []string{testChunkID1},
		ChunkCount:    1,
		EmbeddedCount: 1,
		Content:       "# Hello\n\nSome markdown content.",
		Stage:         models.JobStageFinalize,
		Status:        models.JobStatusRunning,
	}
}

// TestFinalizeStage_Success_StatusReady verifies the memory status is set to "ready".
func TestFinalizeStage_Success_StatusReady(t *testing.T) {
	memory := &models.Memory{
		ID:      testMemoryID,
		AgentID: testAgentID,
		Content: "Some content.",
		Status:  models.MemoryStatusProcessing,
	}

	var updated *models.Memory
	mock := &nesttest.MockMemories{
		OnGet: func(_ context.Context, key string) (*models.Memory, error) {
			if key == testMemoryID {
				return memory, nil
			}
			return nil, nil
		},
		OnSet: func(_ context.Context, _ string, m *models.Memory) error {
			updated = m
			return nil
		},
	}
	nesttest.SetupRegistry(t, nesttest.WithMemories(mock))

	job := newFinalizeJob()
	if _, err := finalizeStage(context.Background(), job); err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	if updated == nil {
		t.Fatal("expected memory to be updated, got nil")
	}
	if updated.Status != models.MemoryStatusReady {
		t.Errorf("status = %q, want %q", updated.Status, models.MemoryStatusReady)
	}
}

// TestFinalizeStage_Success_SetsCompletedAt verifies CompletedAt is populated and Status is completed.
func TestFinalizeStage_Success_SetsCompletedAt(t *testing.T) {
	memory := &models.Memory{
		ID:      testMemoryID,
		AgentID: testAgentID,
		Content: "Some content.",
		Status:  models.MemoryStatusProcessing,
	}
	mock := &nesttest.MockMemories{
		OnGet: func(_ context.Context, _ string) (*models.Memory, error) {
			return memory, nil
		},
	}
	nesttest.SetupRegistry(t, nesttest.WithMemories(mock))

	job := newFinalizeJob()
	result, err := finalizeStage(context.Background(), job)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	if result.CompletedAt == nil {
		t.Error("expected CompletedAt to be set, got nil")
	}
	if result.Status != models.JobStatusCompleted {
		t.Errorf("status = %q, want %q", result.Status, models.JobStatusCompleted)
	}
}

// TestFinalizeStage_GetMemoryError verifies a Get failure propagates.
func TestFinalizeStage_GetMemoryError(t *testing.T) {
	getErr := errors.New("database unavailable")
	mock := &nesttest.MockMemories{
		OnGet: func(_ context.Context, _ string) (*models.Memory, error) {
			return nil, getErr
		},
	}
	nesttest.SetupRegistry(t, nesttest.WithMemories(mock))

	job := newFinalizeJob()
	_, err := finalizeStage(context.Background(), job)

	if !errors.Is(err, getErr) {
		t.Errorf("error = %v, want %v", err, getErr)
	}
}

// TestFinalizeStage_SetMemoryError verifies a Set failure propagates.
func TestFinalizeStage_SetMemoryError(t *testing.T) {
	memory := &models.Memory{
		ID:      testMemoryID,
		AgentID: testAgentID,
		Content: "Some content.",
		Status:  models.MemoryStatusProcessing,
	}
	setErr := errors.New("disk full")
	mock := &nesttest.MockMemories{
		OnGet: func(_ context.Context, _ string) (*models.Memory, error) {
			return memory, nil
		},
		OnSet: func(_ context.Context, _ string, _ *models.Memory) error {
			return setErr
		},
	}
	nesttest.SetupRegistry(t, nesttest.WithMemories(mock))

	job := newFinalizeJob()
	_, err := finalizeStage(context.Background(), job)

	if !errors.Is(err, setErr) {
		t.Errorf("error = %v, want %v", err, setErr)
	}
}
