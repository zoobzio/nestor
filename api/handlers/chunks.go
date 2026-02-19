package handlers

import (
	"time"

	"github.com/zoobzio/rocco"
	"github.com/zoobzio/sum"
	"github.com/zoobzio/nestor/api/contracts"
	"github.com/zoobzio/nestor/api/transformers"
	"github.com/zoobzio/nestor/api/wire"
	"github.com/zoobzio/nestor/models"
)

// ListChunks returns all chunks for a given memory, ordered by sequence.
var ListChunks = rocco.GET("/agents/{agent_id}/memories/{memory_id}/chunks", func(req *rocco.Request[rocco.NoBody]) (wire.ChunkListResponse, error) {
	chunks := sum.MustUse[contracts.Chunks](req.Context)

	memoryID := req.Params.Path["memory_id"]

	list, err := chunks.ListByMemoryID(req.Context, memoryID)
	if err != nil {
		return wire.ChunkListResponse{}, err
	}

	return transformers.ChunksToList(list), nil
}).WithPathParams("agent_id", "memory_id").
	WithSummary("List chunks").
	WithDescription("Returns all chunks belonging to the specified memory, ordered by sequence.").
	WithTags("Chunks").
	WithAuthentication()

// CreateChunk creates a new chunk for the specified memory.
var CreateChunk = rocco.POST("/agents/{agent_id}/memories/{memory_id}/chunks", func(req *rocco.Request[wire.CreateChunkRequest]) (wire.ChunkResponse, error) {
	chunks := sum.MustUse[contracts.Chunks](req.Context)

	chunk := &models.Chunk{
		ID:        newUUID(),
		MemoryID:  req.Params.Path["memory_id"],
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	transformers.ApplyCreateChunk(req.Body, chunk)

	if err := chunks.Set(req.Context, chunk.ID, chunk); err != nil {
		return wire.ChunkResponse{}, err
	}

	return transformers.ChunkToResponse(chunk), nil
}).WithPathParams("agent_id", "memory_id").
	WithSummary("Create chunk").
	WithDescription("Creates a new chunk for the specified memory.").
	WithTags("Chunks").
	WithAuthentication().
	WithSuccessStatus(201)

// GetChunk returns a single chunk by ID.
var GetChunk = rocco.GET("/agents/{agent_id}/memories/{memory_id}/chunks/{id}", func(req *rocco.Request[rocco.NoBody]) (wire.ChunkResponse, error) {
	chunks := sum.MustUse[contracts.Chunks](req.Context)

	id := req.Params.Path["id"]

	chunk, err := chunks.Get(req.Context, id)
	if err != nil {
		return wire.ChunkResponse{}, ErrChunkNotFound
	}

	return transformers.ChunkToResponse(chunk), nil
}).WithPathParams("agent_id", "memory_id", "id").
	WithSummary("Get chunk").
	WithDescription("Returns a chunk by ID.").
	WithTags("Chunks").
	WithErrors(ErrChunkNotFound).
	WithAuthentication()

// DeleteChunk removes a chunk by ID.
var DeleteChunk = rocco.DELETE("/agents/{agent_id}/memories/{memory_id}/chunks/{id}", func(req *rocco.Request[rocco.NoBody]) (rocco.NoBody, error) {
	chunks := sum.MustUse[contracts.Chunks](req.Context)

	id := req.Params.Path["id"]

	if err := chunks.Delete(req.Context, id); err != nil {
		return rocco.NoBody{}, ErrChunkNotFound
	}

	return rocco.NoBody{}, nil
}).WithPathParams("agent_id", "memory_id", "id").
	WithSummary("Delete chunk").
	WithDescription("Deletes a chunk by ID.").
	WithTags("Chunks").
	WithErrors(ErrChunkNotFound).
	WithAuthentication().
	WithSuccessStatus(204)

// SearchChunks performs vector similarity search across all chunks,
// returning results ordered by cosine distance to the query embedding.
var SearchChunks = rocco.POST("/chunks/search", func(req *rocco.Request[wire.ChunkSearchRequest]) (wire.ChunkSearchResponse, error) {
	chunks := sum.MustUse[contracts.Chunks](req.Context)

	results, err := chunks.SearchByEmbedding(req.Context, req.Body.Embedding, req.Body.Limit)
	if err != nil {
		return wire.ChunkSearchResponse{}, err
	}

	return transformers.ChunksToSearchResponse(results), nil
}).WithSummary("Search chunks by embedding").
	WithDescription("Performs vector similarity search across all chunks using cosine distance. Returns chunks ordered by similarity to the query embedding.").
	WithTags("Chunks").
	WithAuthentication()
