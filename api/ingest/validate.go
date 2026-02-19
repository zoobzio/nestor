package ingest

import (
	"context"
	"errors"
	"unicode/utf8"

	"github.com/zoobzio/nestor/api/contracts"
	"github.com/zoobzio/nestor/capacitors"
	"github.com/zoobzio/nestor/events"
	"github.com/zoobzio/nestor/models"
	"github.com/zoobzio/sum"
)

// Validation errors.
var (
	ErrAgentNotFound    = errors.New("agent not found")
	ErrContentEmpty     = errors.New("content is empty")
	ErrContentTooLarge  = errors.New("content exceeds maximum size")
	ErrContentInvalidUTF8 = errors.New("content is not valid UTF-8")
)

// validateStage ensures the job has valid inputs before expensive operations.
//
// Checks performed:
//   - Agent exists
//   - Content is non-empty
//   - Content is valid UTF-8
//   - Content size is within limits (configurable via capacitor)
func validateStage(ctx context.Context, job *models.IngestJob) (*models.IngestJob, error) {
	events.Ingest.Validate.Started.Emit(ctx, events.StageEvent{JobID: job.ID})

	job.Stage = models.JobStageValidate

	if err := validate(ctx, job); err != nil {
		events.Ingest.Validate.Failed.Emit(ctx, events.StageFailedEvent{
			JobID: job.ID,
			Error: err.Error(),
		})
		return job, err
	}

	events.Ingest.Validate.Completed.Emit(ctx, events.StageEvent{JobID: job.ID})
	return job, nil
}

func validate(ctx context.Context, job *models.IngestJob) error {
	agents := sum.MustUse[contracts.Agents](ctx)

	// Verify agent exists.
	agent, err := agents.Get(ctx, job.AgentID)
	if err != nil {
		return err
	}
	if agent == nil {
		return ErrAgentNotFound
	}

	// Validate content is non-empty.
	if len(job.Content) == 0 {
		return ErrContentEmpty
	}

	// Validate content is valid UTF-8.
	if !utf8.ValidString(job.Content) {
		return ErrContentInvalidUTF8
	}

	// Check content size limits from capacitor.
	cfg := capacitors.DefaultIngestConfig()
	if int64(len(job.Content)) > cfg.MaxContentSize {
		return ErrContentTooLarge
	}

	return nil
}
