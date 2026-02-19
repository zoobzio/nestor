// Package models contains domain model types.
package models

import (
	"time"

	"github.com/zoobzio/check"
)

// JobStage represents the current processing stage of an ingest job.
type JobStage string

const (
	JobStageValidate JobStage = "validate"
	JobStageStore    JobStage = "store"
	JobStageChunk    JobStage = "chunk"
	JobStageEmbed    JobStage = "embed"
	JobStageFinalize JobStage = "finalize"
)

// JobStatus represents the overall status of an ingest job.
type JobStatus string

const (
	JobStatusPending   JobStatus = "pending"
	JobStatusRunning   JobStatus = "running"
	JobStatusCompleted JobStatus = "completed"
	JobStatusFailed    JobStatus = "failed"
)

// IngestJob is the carrier struct that flows through the memory ingestion pipeline.
// It holds all state required by each pipeline stage as the job progresses from
// validation through storage, chunking, embedding, and finalization.
type IngestJob struct {
	ID            string     `json:"id" db:"id" constraints:"primarykey" description:"Job UUID" example:"550e8400-e29b-41d4-a716-446655440000"`
	CustomerID    string     `json:"customer_id" db:"customer_id" constraints:"notnull" description:"Customer UUID" example:"550e8400-e29b-41d4-a716-446655440001"`
	AgentID       string     `json:"agent_id" db:"agent_id" constraints:"notnull" description:"Agent UUID" example:"550e8400-e29b-41d4-a716-446655440002"`
	MemoryID      string     `json:"memory_id" db:"memory_id" constraints:"notnull" description:"Memory UUID" example:"550e8400-e29b-41d4-a716-446655440003"`
	Content       string     `json:"content" db:"content" constraints:"notnull" type:"text" description:"Raw content submitted for ingestion"`
	Stage         JobStage   `json:"stage" db:"stage" constraints:"notnull" default:"'validate'" description:"Current pipeline stage: validate, store, chunk, embed, finalize" example:"chunk"`
	Status        JobStatus  `json:"status" db:"status" constraints:"notnull" default:"'pending'" description:"Overall job status: pending, running, completed, failed" example:"running"`
	ChunkIDs      []string   `json:"chunk_ids" db:"chunk_ids" type:"text[]" description:"IDs of chunks produced by the chunking stage"`
	ChunkCount    int        `json:"chunk_count" db:"chunk_count" constraints:"notnull" default:"0" description:"Total number of chunks produced" example:"12"`
	EmbeddedCount int        `json:"embedded_count" db:"embedded_count" constraints:"notnull" default:"0" description:"Number of chunks successfully embedded" example:"10"`
	Error         *string    `json:"error,omitempty" db:"error" description:"Error message if the job failed"`
	CreatedAt     time.Time  `json:"created_at" db:"created_at" constraints:"notnull" default:"now()" description:"Job creation time"`
	StartedAt     *time.Time `json:"started_at,omitempty" db:"started_at" description:"Time the job began processing"`
	CompletedAt   *time.Time `json:"completed_at,omitempty" db:"completed_at" description:"Time the job finished (success or failure)"`
}

// Clone returns a deep copy of IngestJob, suitable for worker pool use.
// Pointer fields and slices are copied independently.
func (j *IngestJob) Clone() *IngestJob {
	c := *j

	if j.ChunkIDs != nil {
		c.ChunkIDs = make([]string, len(j.ChunkIDs))
		copy(c.ChunkIDs, j.ChunkIDs)
	}

	if j.Error != nil {
		s := *j.Error
		c.Error = &s
	}

	if j.StartedAt != nil {
		t := *j.StartedAt
		c.StartedAt = &t
	}

	if j.CompletedAt != nil {
		t := *j.CompletedAt
		c.CompletedAt = &t
	}

	return &c
}

// Validate validates the IngestJob model.
func (j IngestJob) Validate() error {
	return check.All(
		check.Str(j.ID, "id").Required().V(),
		check.Str(j.CustomerID, "customer_id").Required().V(),
		check.Str(j.AgentID, "agent_id").Required().V(),
		check.Str(j.Content, "content").Required().V(),
	).Err()
}
