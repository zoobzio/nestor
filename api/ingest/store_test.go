//go:build testing

package ingest

import (
	"context"
	"errors"
	"testing"

	"github.com/zoobzio/nestor/models"
	nesttest "github.com/zoobzio/nestor/testing"
)

// newStoreJob returns an IngestJob ready for the store stage (post-validate).
func newStoreJob() *models.IngestJob {
	return &models.IngestJob{
		ID:         testJobID,
		CustomerID: testCustomerID,
		AgentID:    testAgentID,
		Content:    "# Hello\n\nSome markdown content.",
		Stage:      models.JobStageStore,
		Status:     models.JobStatusRunning,
	}
}

// TestStoreStage_Success_SetsMemoryID verifies that MemoryID is populated after the stage runs.
func TestStoreStage_Success_SetsMemoryID(t *testing.T) {
	mock := &nesttest.MockMemories{}
	nesttest.SetupRegistry(t, nesttest.WithMemories(mock))

	job := newStoreJob()
	job.MemoryID = "" // ensure no pre-existing ID

	result, err := storeStage(context.Background(), job)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if result.MemoryID == "" {
		t.Error("expected MemoryID to be set, got empty string")
	}
}

// TestStoreStage_Success_StatusProcessing verifies the stored memory has status "processing".
func TestStoreStage_Success_StatusProcessing(t *testing.T) {
	var stored *models.Memory
	mock := &nesttest.MockMemories{
		OnSet: func(_ context.Context, _ string, m *models.Memory) error {
			stored = m
			return nil
		},
	}
	nesttest.SetupRegistry(t, nesttest.WithMemories(mock))

	job := newStoreJob()
	if _, err := storeStage(context.Background(), job); err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	if stored == nil {
		t.Fatal("expected memory to be stored, got nil")
	}
	if stored.Status != models.MemoryStatusProcessing {
		t.Errorf("status = %q, want %q", stored.Status, models.MemoryStatusProcessing)
	}
}

// TestStoreStage_PreservesExistingMemoryID verifies that a pre-set MemoryID is preserved.
func TestStoreStage_PreservesExistingMemoryID(t *testing.T) {
	const existingID = "550e8400-e29b-41d4-a716-446655440099"
	var storedKey string
	mock := &nesttest.MockMemories{
		OnSet: func(_ context.Context, key string, _ *models.Memory) error {
			storedKey = key
			return nil
		},
	}
	nesttest.SetupRegistry(t, nesttest.WithMemories(mock))

	job := newStoreJob()
	job.MemoryID = existingID

	result, err := storeStage(context.Background(), job)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if result.MemoryID != existingID {
		t.Errorf("MemoryID = %q, want %q", result.MemoryID, existingID)
	}
	if storedKey != existingID {
		t.Errorf("stored key = %q, want %q", storedKey, existingID)
	}
}

// TestStoreStage_StoreError verifies the store error propagates.
func TestStoreStage_StoreError(t *testing.T) {
	storeErr := errors.New("disk full")
	mock := &nesttest.MockMemories{
		OnSet: func(_ context.Context, _ string, _ *models.Memory) error {
			return storeErr
		},
	}
	nesttest.SetupRegistry(t, nesttest.WithMemories(mock))

	job := newStoreJob()
	_, err := storeStage(context.Background(), job)

	if !errors.Is(err, storeErr) {
		t.Errorf("error = %v, want %v", err, storeErr)
	}
}
