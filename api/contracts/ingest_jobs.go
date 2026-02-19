package contracts

import (
	"context"

	"github.com/zoobzio/nestor/models"
)

// IngestJobs tracks ingestion job state for observability and recovery.
type IngestJobs interface {
	// Get retrieves a job by ID.
	Get(ctx context.Context, id string) (*models.IngestJob, error)

	// Save persists a job (insert or update).
	Save(ctx context.Context, job *models.IngestJob) error

	// ListByStatus returns jobs with the given status.
	ListByStatus(ctx context.Context, status models.JobStatus) ([]*models.IngestJob, error)
}
