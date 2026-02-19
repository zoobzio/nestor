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

// These tests exercise the contracts.IngestJobs interface via MockIngestJobs.
// The concrete IngestJobs store (which embeds sum.Database) requires a live
// database and is covered by integration tests.

func TestIngestJobs_Get_Found(t *testing.T) {
	job := nesttest.NewIngestJob(t)

	mock := &nesttest.MockIngestJobs{
		OnGet: func(_ context.Context, id string) (*models.IngestJob, error) {
			if id != job.ID {
				return nil, errors.New("not found")
			}
			return job, nil
		},
	}

	got, err := mock.Get(context.Background(), job.ID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.ID != job.ID {
		t.Errorf("got ID %q, want %q", got.ID, job.ID)
	}
}

func TestIngestJobs_Get_NotFound(t *testing.T) {
	mock := &nesttest.MockIngestJobs{
		OnGet: func(_ context.Context, _ string) (*models.IngestJob, error) {
			return nil, errors.New("not found")
		},
	}

	_, err := mock.Get(context.Background(), "does-not-exist")
	if err == nil {
		t.Error("expected error for missing job, got nil")
	}
}

func TestIngestJobs_Save_Insert(t *testing.T) {
	var stored *models.IngestJob

	mock := &nesttest.MockIngestJobs{
		OnSave: func(_ context.Context, job *models.IngestJob) error {
			stored = job
			return nil
		},
	}

	job := nesttest.NewIngestJob(t)
	if err := mock.Save(context.Background(), job); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if stored == nil {
		t.Fatal("expected job to be stored, got nil")
	}
	if stored.ID != job.ID {
		t.Errorf("stored ID %q, want %q", stored.ID, job.ID)
	}
	if stored.Status != models.JobStatusPending {
		t.Errorf("stored Status %q, want %q", stored.Status, models.JobStatusPending)
	}
}

func TestIngestJobs_Save_Update(t *testing.T) {
	now := time.Now()
	job := nesttest.NewIngestJob(t)
	job.Status = models.JobStatusRunning
	job.Stage = models.JobStageChunk
	job.StartedAt = &now

	var stored *models.IngestJob

	mock := &nesttest.MockIngestJobs{
		OnSave: func(_ context.Context, j *models.IngestJob) error {
			stored = j
			return nil
		},
	}

	if err := mock.Save(context.Background(), job); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if stored == nil {
		t.Fatal("expected job to be stored, got nil")
	}
	if stored.Status != models.JobStatusRunning {
		t.Errorf("stored Status %q, want %q", stored.Status, models.JobStatusRunning)
	}
	if stored.Stage != models.JobStageChunk {
		t.Errorf("stored Stage %q, want %q", stored.Stage, models.JobStageChunk)
	}
	if stored.StartedAt == nil {
		t.Error("expected StartedAt to be set, got nil")
	}
}

func TestIngestJobs_Save_Error(t *testing.T) {
	mock := &nesttest.MockIngestJobs{
		OnSave: func(_ context.Context, _ *models.IngestJob) error {
			return errors.New("storage error")
		},
	}

	job := nesttest.NewIngestJob(t)
	if err := mock.Save(context.Background(), job); err == nil {
		t.Error("expected error from Save, got nil")
	}
}

func TestIngestJobs_ListByStatus_Pending(t *testing.T) {
	now := time.Now()
	expected := []*models.IngestJob{
		{
			ID: "job-1", CustomerID: "cust-1", AgentID: "agent-1", MemoryID: "mem-1",
			Content: "First job.", Stage: models.JobStageValidate, Status: models.JobStatusPending, CreatedAt: now,
		},
		{
			ID: "job-2", CustomerID: "cust-1", AgentID: "agent-1", MemoryID: "mem-2",
			Content: "Second job.", Stage: models.JobStageValidate, Status: models.JobStatusPending, CreatedAt: now,
		},
	}

	mock := &nesttest.MockIngestJobs{
		OnListByStatus: func(_ context.Context, status models.JobStatus) ([]*models.IngestJob, error) {
			if status != models.JobStatusPending {
				return []*models.IngestJob{}, nil
			}
			return expected, nil
		},
	}

	list, err := mock.ListByStatus(context.Background(), models.JobStatusPending)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(list) != 2 {
		t.Fatalf("expected 2 jobs, got %d", len(list))
	}
	if list[0].ID != "job-1" {
		t.Errorf("first job ID %q, want %q", list[0].ID, "job-1")
	}
	if list[1].ID != "job-2" {
		t.Errorf("second job ID %q, want %q", list[1].ID, "job-2")
	}
}

func TestIngestJobs_ListByStatus_Failed(t *testing.T) {
	now := time.Now()
	errMsg := "pipeline error"
	expected := []*models.IngestJob{
		{
			ID: "job-3", CustomerID: "cust-2", AgentID: "agent-2", MemoryID: "mem-3",
			Content: "Failed job.", Stage: models.JobStageEmbed, Status: models.JobStatusFailed,
			Error: &errMsg, CreatedAt: now,
		},
	}

	mock := &nesttest.MockIngestJobs{
		OnListByStatus: func(_ context.Context, status models.JobStatus) ([]*models.IngestJob, error) {
			if status != models.JobStatusFailed {
				return []*models.IngestJob{}, nil
			}
			return expected, nil
		},
	}

	list, err := mock.ListByStatus(context.Background(), models.JobStatusFailed)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(list) != 1 {
		t.Fatalf("expected 1 job, got %d", len(list))
	}
	if list[0].Status != models.JobStatusFailed {
		t.Errorf("job Status %q, want %q", list[0].Status, models.JobStatusFailed)
	}
	if list[0].Error == nil || *list[0].Error != errMsg {
		t.Errorf("job Error %v, want %q", list[0].Error, errMsg)
	}
}

func TestIngestJobs_ListByStatus_Empty(t *testing.T) {
	mock := &nesttest.MockIngestJobs{
		OnListByStatus: func(_ context.Context, _ models.JobStatus) ([]*models.IngestJob, error) {
			return []*models.IngestJob{}, nil
		},
	}

	list, err := mock.ListByStatus(context.Background(), models.JobStatusCompleted)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(list) != 0 {
		t.Errorf("expected empty list, got %d jobs", len(list))
	}
}

func TestIngestJobs_ListByStatus_Error(t *testing.T) {
	mock := &nesttest.MockIngestJobs{
		OnListByStatus: func(_ context.Context, _ models.JobStatus) ([]*models.IngestJob, error) {
			return nil, errors.New("database error")
		},
	}

	_, err := mock.ListByStatus(context.Background(), models.JobStatusRunning)
	if err == nil {
		t.Error("expected error from ListByStatus, got nil")
	}
}
