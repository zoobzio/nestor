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

// chunkStage parses markdown into semantic chunks using chisel.
//
// Operations:
//  1. Parse markdown using Chunker contract (chisel markdown provider)
//  2. For each chunk, create Chunk model with MemoryID
//  3. Store chunks to database
//  4. Accumulate chunk IDs in job.ChunkIDs
//  5. Update job.ChunkCount
func chunkStage(ctx context.Context, job *models.IngestJob) (*models.IngestJob, error) {
	events.Ingest.Chunk.Started.Emit(ctx, events.StageEvent{JobID: job.ID})

	job.Stage = models.JobStageChunk

	if err := chunk(ctx, job); err != nil {
		events.Ingest.Chunk.Failed.Emit(ctx, events.StageFailedEvent{
			JobID: job.ID,
			Error: err.Error(),
		})
		return job, err
	}

	events.Ingest.Chunk.Completed.Emit(ctx, events.StageEvent{JobID: job.ID})
	return job, nil
}

func chunk(ctx context.Context, job *models.IngestJob) error {
	chunker := sum.MustUse[contracts.Chunker](ctx)
	chunks := sum.MustUse[contracts.Chunks](ctx)

	// Parse markdown into chunks.
	parsed, err := chunker.Chunk(ctx, []byte(job.Content))
	if err != nil {
		return err
	}

	// Reset chunk tracking.
	job.ChunkIDs = make([]string, 0, len(parsed))
	job.ChunkCount = 0

	now := time.Now()

	// Create and store each chunk.
	for i, c := range parsed {
		chunkID := uuid.New().String()

		var symbol *string
		if c.Symbol != "" {
			s := c.Symbol
			symbol = &s
		}

		var context []string
		if len(c.Context) > 0 {
			context = c.Context
		}

		chunk := &models.Chunk{
			ID:        chunkID,
			MemoryID:  job.MemoryID,
			Content:   c.Content,
			Sequence:  i,
			Kind:      models.ChunkKind(c.Kind),
			Symbol:    symbol,
			Context:   context,
			StartLine: c.StartLine,
			EndLine:   c.EndLine,
			CreatedAt: now,
			UpdatedAt: now,
		}

		if err := chunks.Set(ctx, chunkID, chunk); err != nil {
			return err
		}

		job.ChunkIDs = append(job.ChunkIDs, chunkID)
		job.ChunkCount++
	}

	return nil
}
