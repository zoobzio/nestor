//go:build testing

// Package testing provides test infrastructure: fixtures, mocks, and helpers.
package testing

import (
	"testing"
	"time"

	"github.com/zoobzio/nestor/models"
)

// NewAgent returns an Agent fixture with sensible defaults.
func NewAgent(t *testing.T) *models.Agent {
	t.Helper()
	now := time.Now()
	return &models.Agent{
		ID:        "550e8400-e29b-41d4-a716-446655440001",
		UserID:    "550e8400-e29b-41d4-a716-446655440000",
		Name:      "Research Assistant",
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// NewAgents returns n Agent fixtures with distinct IDs.
func NewAgents(t *testing.T, n int) []*models.Agent {
	t.Helper()
	agents := make([]*models.Agent, n)
	for i := range agents {
		a := NewAgent(t)
		// Vary the last byte of the ID to produce distinct values.
		a.ID = "550e8400-e29b-41d4-a716-4466554400" + twoDigit(i+1)
		a.Name = "Agent " + twoDigit(i+1)
		agents[i] = a
	}
	return agents
}

// NewMemory returns a Memory fixture with sensible defaults.
func NewMemory(t *testing.T) *models.Memory {
	t.Helper()
	now := time.Now()
	return &models.Memory{
		ID:        "550e8400-e29b-41d4-a716-446655440002",
		AgentID:   "550e8400-e29b-41d4-a716-446655440001",
		Content:   "The user prefers concise responses.",
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// NewMemories returns n Memory fixtures with distinct IDs.
func NewMemories(t *testing.T, n int) []*models.Memory {
	t.Helper()
	memories := make([]*models.Memory, n)
	for i := range memories {
		m := NewMemory(t)
		m.ID = "550e8400-e29b-41d4-a716-4466554400" + twoDigit(i+2)
		m.Content = "Memory " + twoDigit(i+1)
		memories[i] = m
	}
	return memories
}

// NewChunk returns a Chunk fixture with sensible defaults.
func NewChunk(t *testing.T) *models.Chunk {
	t.Helper()
	now := time.Now()
	return &models.Chunk{
		ID:        "550e8400-e29b-41d4-a716-446655440003",
		MemoryID:  "550e8400-e29b-41d4-a716-446655440002",
		Content:   "The user prefers concise responses.",
		Sequence:  0,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// NewChunks returns n Chunk fixtures with distinct IDs and sequential Sequence values.
func NewChunks(t *testing.T, n int) []*models.Chunk {
	t.Helper()
	chunks := make([]*models.Chunk, n)
	for i := range chunks {
		c := NewChunk(t)
		c.ID = "550e8400-e29b-41d4-a716-4466554400" + twoDigit(i+3)
		c.Content = "Chunk " + twoDigit(i+1)
		c.Sequence = i
		chunks[i] = c
	}
	return chunks
}

// NewIngestJob returns an IngestJob fixture with sensible defaults.
func NewIngestJob(t *testing.T) *models.IngestJob {
	t.Helper()
	now := time.Now()
	return &models.IngestJob{
		ID:         "550e8400-e29b-41d4-a716-446655440010",
		CustomerID: "550e8400-e29b-41d4-a716-446655440000",
		AgentID:    "550e8400-e29b-41d4-a716-446655440001",
		MemoryID:   "550e8400-e29b-41d4-a716-446655440002",
		Content:    "Ingest this content for testing.",
		Stage:      models.JobStageValidate,
		Status:     models.JobStatusPending,
		CreatedAt:  now,
	}
}

// NewIngestJobs returns n IngestJob fixtures with distinct IDs.
func NewIngestJobs(t *testing.T, n int) []*models.IngestJob {
	t.Helper()
	jobs := make([]*models.IngestJob, n)
	for i := range jobs {
		j := NewIngestJob(t)
		j.ID = "550e8400-e29b-41d4-a716-4466554401" + twoDigit(i+1)
		j.Content = "Ingest job " + twoDigit(i+1)
		jobs[i] = j
	}
	return jobs
}

// twoDigit formats n as a zero-padded two-character decimal string for n < 100.
func twoDigit(n int) string {
	if n < 10 {
		return "0" + string(rune('0'+n))
	}
	tens := n / 10
	ones := n % 10
	return string(rune('0'+tens)) + string(rune('0'+ones))
}
