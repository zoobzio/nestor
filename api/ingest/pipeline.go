// Package ingest implements the memory ingestion pipeline.
//
// The pipeline transforms markdown content from agents into searchable,
// semantically-indexed memories through five stages:
//
//   - Validate: Guard clause ensuring valid customer, agent, and content
//   - Store: Persist raw memory before further processing
//   - Chunk: Parse markdown into semantic chunks using chisel
//   - Embed: Generate vector embeddings for semantic search
//   - Finalize: Mark memory as ready for search
//
// Each stage is wrapped with appropriate reliability (timeout, retry, backoff)
// and emits events for observability.
package ingest

import (
	"time"

	"github.com/zoobzio/nestor/models"
	"github.com/zoobzio/pipz"
)

// Pipeline identities for observability and debugging.
var (
	PipelineID = pipz.NewIdentity("ingest.pipeline", "Memory ingestion pipeline")

	// Stage identities.
	ValidateStageID  = pipz.NewIdentity("ingest.validate", "Validates job inputs")
	StoreStageID     = pipz.NewIdentity("ingest.store", "Stores raw memory")
	ChunkStageID     = pipz.NewIdentity("ingest.chunk", "Parses markdown into chunks")
	EmbedStageID     = pipz.NewIdentity("ingest.embed", "Generates vector embeddings")
	FinalizeStageID  = pipz.NewIdentity("ingest.finalize", "Marks memory as ready")

	// Timeout wrapper identities.
	StoreTimeoutID    = pipz.NewIdentity("ingest.store.timeout", "Store stage timeout")
	ChunkTimeoutID    = pipz.NewIdentity("ingest.chunk.timeout", "Chunk stage timeout")
	EmbedTimeoutID    = pipz.NewIdentity("ingest.embed.timeout", "Embed stage timeout")
	FinalizeTimeoutID = pipz.NewIdentity("ingest.finalize.timeout", "Finalize stage timeout")

	// Retry wrapper identities.
	StoreRetryID    = pipz.NewIdentity("ingest.store.retry", "Store stage retry")
	ChunkRetryID    = pipz.NewIdentity("ingest.chunk.retry", "Chunk stage retry")
	EmbedBackoffID  = pipz.NewIdentity("ingest.embed.backoff", "Embed stage backoff")
	FinalizeRetryID = pipz.NewIdentity("ingest.finalize.retry", "Finalize stage retry")

	// Embed nested pool identity.
	EmbedPoolID = pipz.NewIdentity("ingest.embed.pool", "Embedding batch worker pool")
)

// NewPipeline constructs the memory ingestion pipeline.
//
// The pipeline processes IngestJob carriers through five stages:
//
//	validate -> store -> chunk -> embed -> finalize
//
// Reliability wrapping:
//   - Validate: No retry (fast, local validation)
//   - Store: Timeout 30s, retry 3x
//   - Chunk: Timeout 5m, retry 3x
//   - Embed: Timeout 30m, backoff 3x with 1s base
//   - Finalize: Timeout 10s, retry 3x
func NewPipeline() *pipz.Sequence[*models.IngestJob] {
	// Create base stages.
	validate := pipz.Apply(ValidateStageID, validateStage)
	store := pipz.Apply(StoreStageID, storeStage)
	chunk := pipz.Apply(ChunkStageID, chunkStage)
	embed := pipz.Apply(EmbedStageID, embedStage)
	finalize := pipz.Apply(FinalizeStageID, finalizeStage)

	// Wrap with timeout.
	storeWithTimeout := pipz.NewTimeout(StoreTimeoutID, store, 30*time.Second)
	chunkWithTimeout := pipz.NewTimeout(ChunkTimeoutID, chunk, 5*time.Minute)
	embedWithTimeout := pipz.NewTimeout(EmbedTimeoutID, embed, 30*time.Minute)
	finalizeWithTimeout := pipz.NewTimeout(FinalizeTimeoutID, finalize, 10*time.Second)

	// Wrap with retry/backoff.
	storeReliable := pipz.NewRetry(StoreRetryID, storeWithTimeout, 3)
	chunkReliable := pipz.NewRetry(ChunkRetryID, chunkWithTimeout, 3)
	embedReliable := pipz.NewBackoff(EmbedBackoffID, embedWithTimeout, 3, time.Second)
	finalizeReliable := pipz.NewRetry(FinalizeRetryID, finalizeWithTimeout, 3)

	// Build sequence.
	return pipz.NewSequence(PipelineID,
		validate, // No retry - fast validation
		storeReliable,
		chunkReliable,
		embedReliable,
		finalizeReliable,
	)
}
