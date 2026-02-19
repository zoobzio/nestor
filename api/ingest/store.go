package ingest

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/zoobzio/nestor/api/contracts"
	"github.com/zoobzio/nestor/events"
	"github.com/zoobzio/nestor/models"
	"github.com/zoobzio/sum"
)

// storeStage persists the raw memory before further processing.
//
// Operations:
//  1. Create Memory model with status "processing"
//  2. Store to database
//  3. Update job.MemoryID with created ID
func storeStage(ctx context.Context, job *models.IngestJob) (*models.IngestJob, error) {
	events.Ingest.Store.Started.Emit(ctx, events.StageEvent{JobID: job.ID})

	job.Stage = models.JobStageStore

	if err := store(ctx, job); err != nil {
		events.Ingest.Store.Failed.Emit(ctx, events.StageFailedEvent{
			JobID: job.ID,
			Error: err.Error(),
		})
		return job, err
	}

	events.Ingest.Store.Completed.Emit(ctx, events.StageEvent{JobID: job.ID})
	return job, nil
}

func store(ctx context.Context, job *models.IngestJob) error {
	memories := sum.MustUse[contracts.Memories](ctx)

	// Generate memory ID if not already set.
	memoryID := job.MemoryID
	if memoryID == "" {
		memoryID = uuid.New().String()
	}

	now := time.Now()
	memory := &models.Memory{
		ID:        memoryID,
		AgentID:   job.AgentID,
		Content:   job.Content,
		Status:    models.MemoryStatusProcessing,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := memories.Set(ctx, memoryID, memory); err != nil {
		return err
	}

	job.MemoryID = memoryID
	return nil
}
