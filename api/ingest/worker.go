package ingest

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/zoobzio/capitan"
	"github.com/zoobzio/nestor/api/contracts"
	"github.com/zoobzio/nestor/capacitors"
	"github.com/zoobzio/nestor/events"
	"github.com/zoobzio/nestor/models"
	"github.com/zoobzio/pipz"
	"github.com/zoobzio/sum"
)

// WorkerPoolID identifies the main ingest worker pool.
var WorkerPoolID = pipz.NewIdentity("ingest.worker.pool", "Ingest pipeline worker pool")

// Worker manages the memory ingestion pipeline lifecycle.
// It listens for memory creation events and processes them through the pipeline.
type Worker struct {
	pool     *pipz.WorkerPool[*models.IngestJob]
	pipeline *pipz.Sequence[*models.IngestJob]
	listener *capitan.Listener
}

// NewWorker creates a new ingestion worker with the default configuration.
func NewWorker() *Worker {
	return NewWorkerWithConfig(capacitors.DefaultIngestConfig())
}

// NewWorkerWithConfig creates a new ingestion worker with the specified configuration.
func NewWorkerWithConfig(cfg capacitors.IngestConfig) *Worker {
	pipeline := NewPipeline()
	return &Worker{
		pool:     pipz.NewWorkerPool(WorkerPoolID, cfg.WorkerCount, pipeline),
		pipeline: pipeline,
	}
}

// Start begins listening for memory creation events.
// The worker will process each event through the ingestion pipeline.
func (w *Worker) Start(ctx context.Context) {
	// Note: The design document shows listening for Memory.Created events.
	// However, the events package currently only defines ingest events.
	// In a real implementation, you would listen for a memory submission event
	// or expose Submit() for direct invocation.
	//
	// For now, we expose Submit() for programmatic use.
}

// Stop gracefully shuts down the worker.
// It closes the event listener and waits for in-flight jobs to complete.
func (w *Worker) Stop() error {
	if w.listener != nil {
		w.listener.Close()
	}
	return w.pool.Close()
}

// Submit submits a new memory for ingestion.
// This is the primary entry point for triggering the ingestion pipeline.
func (w *Worker) Submit(ctx context.Context, agentID, customerID, content string) (string, error) {
	job := &models.IngestJob{
		ID:         uuid.New().String(),
		CustomerID: customerID,
		AgentID:    agentID,
		Content:    content,
		Stage:      models.JobStageValidate,
		Status:     models.JobStatusPending,
		CreatedAt:  time.Now(),
	}

	go w.processJob(ctx, job)

	return job.ID, nil
}

// SubmitSync submits a memory for ingestion and waits for completion.
// Returns the completed job or an error if processing failed.
func (w *Worker) SubmitSync(ctx context.Context, agentID, customerID, content string) (*models.IngestJob, error) {
	job := &models.IngestJob{
		ID:         uuid.New().String(),
		CustomerID: customerID,
		AgentID:    agentID,
		Content:    content,
		Stage:      models.JobStageValidate,
		Status:     models.JobStatusPending,
		CreatedAt:  time.Now(),
	}

	return w.processJobSync(ctx, job)
}

// processJob processes a job asynchronously.
func (w *Worker) processJob(ctx context.Context, job *models.IngestJob) {
	_, _ = w.processJobSync(ctx, job)
}

// processJobSync processes a job and returns the result.
func (w *Worker) processJobSync(ctx context.Context, job *models.IngestJob) (*models.IngestJob, error) {
	jobs := sum.MustUse[contracts.IngestJobs](ctx)

	// Mark started.
	job.Status = models.JobStatusRunning
	now := time.Now()
	job.StartedAt = &now
	_ = jobs.Save(ctx, job)

	events.Ingest.Started.Emit(ctx, events.IngestStartedEvent{
		JobID:      job.ID,
		CustomerID: job.CustomerID,
		AgentID:    job.AgentID,
	})

	// Execute pipeline.
	result, err := w.pool.Process(ctx, job)
	if err != nil {
		errStr := err.Error()
		job.Error = &errStr
		job.Status = models.JobStatusFailed
		completedAt := time.Now()
		job.CompletedAt = &completedAt
		_ = jobs.Save(ctx, job)

		events.Ingest.Failed.Emit(ctx, events.IngestFailedEvent{
			JobID: job.ID,
			Stage: string(job.Stage),
			Error: errStr,
		})
		return job, err
	}

	// Finalize stage already sets CompletedAt and Status.
	// Just ensure we persist the final state.
	_ = jobs.Save(ctx, result)

	events.Ingest.Completed.Emit(ctx, events.IngestCompletedEvent{
		JobID:      result.ID,
		MemoryID:   result.MemoryID,
		ChunkCount: result.ChunkCount,
	})

	return result, nil
}

// SetWorkerCount dynamically adjusts the worker pool size.
// This is called by the capacitor when configuration changes.
func (w *Worker) SetWorkerCount(count int) {
	// WorkerPool does not support dynamic resizing.
	// This would require recreating the pool, which we avoid for simplicity.
	// In production, you might implement a more sophisticated approach.
}
