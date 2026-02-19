package stores

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/zoobzio/astql"
	"github.com/zoobzio/sum"
	"github.com/zoobzio/nestor/models"
)

// IngestJobs provides database access for ingest job records.
type IngestJobs struct {
	*sum.Database[models.IngestJob]
}

// NewIngestJobs creates a new ingest jobs store.
func NewIngestJobs(db *sqlx.DB, renderer astql.Renderer) (*IngestJobs, error) {
	database, err := sum.NewDatabase[models.IngestJob](db, "ingest_jobs", renderer)
	if err != nil {
		return nil, err
	}
	return &IngestJobs{Database: database}, nil
}

// Save persists a job by inserting or updating it by primary key.
func (s *IngestJobs) Save(ctx context.Context, job *models.IngestJob) error {
	return s.Set(ctx, job.ID, job)
}

// ListByStatus returns all jobs with the given status, ordered by creation time ascending.
func (s *IngestJobs) ListByStatus(ctx context.Context, status models.JobStatus) ([]*models.IngestJob, error) {
	return s.Query().
		Where("status", "=", "status").
		OrderBy("created_at", "ASC").
		Exec(ctx, map[string]any{"status": status})
}
