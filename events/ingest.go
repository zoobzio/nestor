// Package events provides event definitions for the application.
package events

import (
	"github.com/zoobzio/capitan"
	"github.com/zoobzio/sum"
)

// IngestStartedEvent carries data for the ingest pipeline start.
type IngestStartedEvent struct {
	JobID      string `json:"job_id"`
	CustomerID string `json:"customer_id"`
	AgentID    string `json:"agent_id"`
}

// IngestCompletedEvent carries data for a successful ingest completion.
type IngestCompletedEvent struct {
	JobID      string `json:"job_id"`
	MemoryID   string `json:"memory_id"`
	ChunkCount int    `json:"chunk_count"`
}

// IngestFailedEvent carries data for a failed ingest pipeline.
type IngestFailedEvent struct {
	JobID string `json:"job_id"`
	Stage string `json:"stage"`
	Error string `json:"error"`
}

// StageEvent carries data for a pipeline stage lifecycle event.
type StageEvent struct {
	JobID string `json:"job_id"`
}

// StageFailedEvent carries data for a failed pipeline stage.
type StageFailedEvent struct {
	JobID string `json:"job_id"`
	Error string `json:"error"`
}

// ChunkCompletedEvent carries data for a completed chunking stage.
type ChunkCompletedEvent struct {
	JobID      string `json:"job_id"`
	ChunkCount int    `json:"chunk_count"`
}

// EmbedCompletedEvent carries data for a completed embedding stage.
type EmbedCompletedEvent struct {
	JobID         string `json:"job_id"`
	EmbeddedCount int    `json:"embedded_count"`
}

// EmbedBatchEvent carries data for a completed embedding batch.
type EmbedBatchEvent struct {
	JobID    string `json:"job_id"`
	BatchIdx int    `json:"batch_idx"`
	Count    int    `json:"count"`
}

// Ingest lifecycle signals.
var (
	ingestStartedSignal   = capitan.NewSignal("nestor.ingest.started", "Memory ingestion pipeline started")
	ingestCompletedSignal = capitan.NewSignal("nestor.ingest.completed", "Memory ingestion pipeline completed")
	ingestFailedSignal    = capitan.NewSignal("nestor.ingest.failed", "Memory ingestion pipeline failed")
)

// Ingest validate stage signals.
var (
	ingestValidateStartedSignal   = capitan.NewSignal("nestor.ingest.validate.started", "Ingest validate stage started")
	ingestValidateCompletedSignal = capitan.NewSignal("nestor.ingest.validate.completed", "Ingest validate stage completed")
	ingestValidateFailedSignal    = capitan.NewSignal("nestor.ingest.validate.failed", "Ingest validate stage failed")
)

// Ingest store stage signals.
var (
	ingestStoreStartedSignal   = capitan.NewSignal("nestor.ingest.store.started", "Ingest store stage started")
	ingestStoreCompletedSignal = capitan.NewSignal("nestor.ingest.store.completed", "Ingest store stage completed")
	ingestStoreFailedSignal    = capitan.NewSignal("nestor.ingest.store.failed", "Ingest store stage failed")
)

// Ingest chunk stage signals.
var (
	ingestChunkStartedSignal   = capitan.NewSignal("nestor.ingest.chunk.started", "Ingest chunk stage started")
	ingestChunkCompletedSignal = capitan.NewSignal("nestor.ingest.chunk.completed", "Ingest chunk stage completed")
	ingestChunkFailedSignal    = capitan.NewSignal("nestor.ingest.chunk.failed", "Ingest chunk stage failed")
)

// Ingest embed stage signals.
var (
	ingestEmbedStartedSignal        = capitan.NewSignal("nestor.ingest.embed.started", "Ingest embed stage started")
	ingestEmbedCompletedSignal      = capitan.NewSignal("nestor.ingest.embed.completed", "Ingest embed stage completed")
	ingestEmbedFailedSignal         = capitan.NewSignal("nestor.ingest.embed.failed", "Ingest embed stage failed")
	ingestEmbedBatchCompletedSignal = capitan.NewSignal("nestor.ingest.embed.batch.completed", "Ingest embed batch completed")
)

// Ingest finalize stage signals.
var (
	ingestFinalizeStartedSignal   = capitan.NewSignal("nestor.ingest.finalize.started", "Ingest finalize stage started")
	ingestFinalizeCompletedSignal = capitan.NewSignal("nestor.ingest.finalize.completed", "Ingest finalize stage completed")
	ingestFinalizeFailedSignal    = capitan.NewSignal("nestor.ingest.finalize.failed", "Ingest finalize stage failed")
)

// Ingest contains all memory ingestion pipeline events.
var Ingest = struct {
	// Lifecycle events.
	Started   sum.Event[IngestStartedEvent]
	Completed sum.Event[IngestCompletedEvent]
	Failed    sum.Event[IngestFailedEvent]

	// Per-stage events.
	Validate struct {
		Started   sum.Event[StageEvent]
		Completed sum.Event[StageEvent]
		Failed    sum.Event[StageFailedEvent]
	}
	Store struct {
		Started   sum.Event[StageEvent]
		Completed sum.Event[StageEvent]
		Failed    sum.Event[StageFailedEvent]
	}
	Chunk struct {
		Started   sum.Event[StageEvent]
		Completed sum.Event[StageEvent]
		Failed    sum.Event[StageFailedEvent]
	}
	Embed struct {
		Started        sum.Event[StageEvent]
		Completed      sum.Event[EmbedCompletedEvent]
		Failed         sum.Event[StageFailedEvent]
		BatchCompleted sum.Event[EmbedBatchEvent]
	}
	Finalize struct {
		Started   sum.Event[StageEvent]
		Completed sum.Event[StageEvent]
		Failed    sum.Event[StageFailedEvent]
	}
}{
	Started:   sum.NewInfoEvent[IngestStartedEvent](ingestStartedSignal),
	Completed: sum.NewInfoEvent[IngestCompletedEvent](ingestCompletedSignal),
	Failed:    sum.NewErrorEvent[IngestFailedEvent](ingestFailedSignal),

	Validate: struct {
		Started   sum.Event[StageEvent]
		Completed sum.Event[StageEvent]
		Failed    sum.Event[StageFailedEvent]
	}{
		Started:   sum.NewInfoEvent[StageEvent](ingestValidateStartedSignal),
		Completed: sum.NewInfoEvent[StageEvent](ingestValidateCompletedSignal),
		Failed:    sum.NewErrorEvent[StageFailedEvent](ingestValidateFailedSignal),
	},

	Store: struct {
		Started   sum.Event[StageEvent]
		Completed sum.Event[StageEvent]
		Failed    sum.Event[StageFailedEvent]
	}{
		Started:   sum.NewInfoEvent[StageEvent](ingestStoreStartedSignal),
		Completed: sum.NewInfoEvent[StageEvent](ingestStoreCompletedSignal),
		Failed:    sum.NewErrorEvent[StageFailedEvent](ingestStoreFailedSignal),
	},

	Chunk: struct {
		Started   sum.Event[StageEvent]
		Completed sum.Event[StageEvent]
		Failed    sum.Event[StageFailedEvent]
	}{
		Started:   sum.NewInfoEvent[StageEvent](ingestChunkStartedSignal),
		Completed: sum.NewInfoEvent[StageEvent](ingestChunkCompletedSignal),
		Failed:    sum.NewErrorEvent[StageFailedEvent](ingestChunkFailedSignal),
	},

	Embed: struct {
		Started        sum.Event[StageEvent]
		Completed      sum.Event[EmbedCompletedEvent]
		Failed         sum.Event[StageFailedEvent]
		BatchCompleted sum.Event[EmbedBatchEvent]
	}{
		Started:        sum.NewInfoEvent[StageEvent](ingestEmbedStartedSignal),
		Completed:      sum.NewInfoEvent[EmbedCompletedEvent](ingestEmbedCompletedSignal),
		Failed:         sum.NewErrorEvent[StageFailedEvent](ingestEmbedFailedSignal),
		BatchCompleted: sum.NewInfoEvent[EmbedBatchEvent](ingestEmbedBatchCompletedSignal),
	},

	Finalize: struct {
		Started   sum.Event[StageEvent]
		Completed sum.Event[StageEvent]
		Failed    sum.Event[StageFailedEvent]
	}{
		Started:   sum.NewInfoEvent[StageEvent](ingestFinalizeStartedSignal),
		Completed: sum.NewInfoEvent[StageEvent](ingestFinalizeCompletedSignal),
		Failed:    sum.NewErrorEvent[StageFailedEvent](ingestFinalizeFailedSignal),
	},
}
