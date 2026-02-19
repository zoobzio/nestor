package ingest

import (
	"context"
	"sync"
	"sync/atomic"
	"time"

	"github.com/zoobzio/nestor/api/contracts"
	"github.com/zoobzio/nestor/capacitors"
	"github.com/zoobzio/nestor/events"
	"github.com/zoobzio/nestor/models"
	"github.com/zoobzio/pipz"
	"github.com/zoobzio/sum"
)

// embedWork is the carrier for the embed batch pipeline.
// It carries a batch of chunks through the embedding process.
type embedWork struct {
	JobID    string
	MemoryID string
	BatchIdx int
	Chunks   []*models.Chunk
}

// Clone implements pipz.Cloner for parallel execution.
func (w *embedWork) Clone() *embedWork {
	c := *w
	if w.Chunks != nil {
		c.Chunks = make([]*models.Chunk, len(w.Chunks))
		for i, chunk := range w.Chunks {
			cloned := chunk.Clone()
			c.Chunks[i] = &cloned
		}
	}
	return &c
}

// Embed batch stage identity.
var EmbedBatchStageID = pipz.NewIdentity("ingest.embed.batch", "Embeds a batch of chunks")

// embedStage generates vector embeddings for all chunks.
//
// Operations:
//  1. Load chunks from database using job.ChunkIDs
//  2. Batch chunks according to capacitor config
//  3. Process batches in parallel using nested WorkerPool
//  4. Update chunks with vectors
//  5. Update job.EmbeddedCount
func embedStage(ctx context.Context, job *models.IngestJob) (*models.IngestJob, error) {
	events.Ingest.Embed.Started.Emit(ctx, events.StageEvent{JobID: job.ID})

	job.Stage = models.JobStageEmbed

	if err := embed(ctx, job); err != nil {
		events.Ingest.Embed.Failed.Emit(ctx, events.StageFailedEvent{
			JobID: job.ID,
			Error: err.Error(),
		})
		return job, err
	}

	events.Ingest.Embed.Completed.Emit(ctx, events.EmbedCompletedEvent{
		JobID:         job.ID,
		EmbeddedCount: job.EmbeddedCount,
	})
	return job, nil
}

func embed(ctx context.Context, job *models.IngestJob) error {
	if len(job.ChunkIDs) == 0 {
		// No chunks to embed.
		return nil
	}

	chunks := sum.MustUse[contracts.Chunks](ctx)

	// Load all chunks.
	allChunks := make([]*models.Chunk, 0, len(job.ChunkIDs))
	for _, id := range job.ChunkIDs {
		chunk, err := chunks.Get(ctx, id)
		if err != nil {
			return err
		}
		if chunk != nil {
			allChunks = append(allChunks, chunk)
		}
	}

	if len(allChunks) == 0 {
		return nil
	}

	// Get configuration.
	cfg := capacitors.DefaultEmbedConfig()

	// Create batches.
	batches := batchChunks(allChunks, cfg.BatchSize)

	// Create nested pipeline for batch processing.
	batchStage := pipz.Apply(EmbedBatchStageID, embedBatchStage)
	batchWithTimeout := pipz.NewTimeout(
		pipz.NewIdentity("ingest.embed.batch.timeout", "Batch timeout"),
		batchStage,
		cfg.Timeout,
	)

	pool := pipz.NewWorkerPool(EmbedPoolID, cfg.Workers, batchWithTimeout)
	defer pool.Close()

	// Track embedded count atomically.
	var embeddedCount int64

	// Process batches concurrently.
	var wg sync.WaitGroup
	errCh := make(chan error, len(batches))

	for i, batch := range batches {
		wg.Add(1)
		go func(idx int, chunks []*models.Chunk) {
			defer wg.Done()

			work := &embedWork{
				JobID:    job.ID,
				MemoryID: job.MemoryID,
				BatchIdx: idx,
				Chunks:   chunks,
			}

			result, err := pool.Process(ctx, work)
			if err != nil {
				errCh <- err
				return
			}

			atomic.AddInt64(&embeddedCount, int64(len(result.Chunks)))

			events.Ingest.Embed.BatchCompleted.Emit(ctx, events.EmbedBatchEvent{
				JobID:    job.ID,
				BatchIdx: idx,
				Count:    len(result.Chunks),
			})
		}(i, batch)
	}

	wg.Wait()
	close(errCh)

	// Check for errors.
	for err := range errCh {
		if err != nil {
			return err
		}
	}

	job.EmbeddedCount = int(embeddedCount)
	return nil
}

// embedBatchStage processes a single batch of chunks.
func embedBatchStage(ctx context.Context, work *embedWork) (*embedWork, error) {
	embedder := sum.MustUse[contracts.Embedder](ctx)
	chunks := sum.MustUse[contracts.Chunks](ctx)

	if len(work.Chunks) == 0 {
		return work, nil
	}

	// Extract text content for embedding.
	texts := make([]string, len(work.Chunks))
	for i, c := range work.Chunks {
		texts[i] = c.Content
	}

	// Generate embeddings.
	embeddings, err := embedder.Embed(ctx, texts)
	if err != nil {
		return work, err
	}

	// Update chunks with vectors.
	now := time.Now()
	for i, chunk := range work.Chunks {
		if i < len(embeddings) {
			chunk.Embedding = embeddings[i]
			chunk.UpdatedAt = now

			if err := chunks.Set(ctx, chunk.ID, chunk); err != nil {
				return work, err
			}
		}
	}

	return work, nil
}

// batchChunks splits chunks into batches of the specified size.
func batchChunks(chunks []*models.Chunk, batchSize int) [][]*models.Chunk {
	if batchSize <= 0 {
		batchSize = 128
	}

	var batches [][]*models.Chunk
	for i := 0; i < len(chunks); i += batchSize {
		end := i + batchSize
		if end > len(chunks) {
			end = len(chunks)
		}
		batches = append(batches, chunks[i:end])
	}
	return batches
}
