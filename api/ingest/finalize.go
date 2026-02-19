package ingest

import (
	"context"
	"time"

	"github.com/zoobzio/nestor/api/contracts"
	"github.com/zoobzio/nestor/events"
	"github.com/zoobzio/nestor/models"
	"github.com/zoobzio/sum"
)

// finalizeStage marks the memory as ready for search.
//
// Operations:
//  1. Update memory status to "ready"
//  2. Update job status to "completed"
//  3. Set job.CompletedAt
func finalizeStage(ctx context.Context, job *models.IngestJob) (*models.IngestJob, error) {
	events.Ingest.Finalize.Started.Emit(ctx, events.StageEvent{JobID: job.ID})

	job.Stage = models.JobStageFinalize

	if err := finalize(ctx, job); err != nil {
		events.Ingest.Finalize.Failed.Emit(ctx, events.StageFailedEvent{
			JobID: job.ID,
			Error: err.Error(),
		})
		return job, err
	}

	events.Ingest.Finalize.Completed.Emit(ctx, events.StageEvent{JobID: job.ID})
	return job, nil
}

func finalize(ctx context.Context, job *models.IngestJob) error {
	memories := sum.MustUse[contracts.Memories](ctx)

	// Load the memory.
	memory, err := memories.Get(ctx, job.MemoryID)
	if err != nil {
		return err
	}

	// Update status to ready.
	memory.Status = models.MemoryStatusReady
	memory.UpdatedAt = time.Now()

	if err := memories.Set(ctx, memory.ID, memory); err != nil {
		return err
	}

	// Mark job as completed.
	job.Status = models.JobStatusCompleted
	now := time.Now()
	job.CompletedAt = &now

	return nil
}
