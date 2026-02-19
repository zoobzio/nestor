//go:build testing

// Package testing provides test infrastructure: fixtures, mocks, and helpers.
package testing

import (
	"context"

	"github.com/zoobzio/nestor/api/contracts"
	"github.com/zoobzio/nestor/models"
)

// MockAgents implements contracts.Agents using function fields.
// Set each On* field to control behaviour per test.
type MockAgents struct {
	OnGet           func(ctx context.Context, key string) (*models.Agent, error)
	OnSet           func(ctx context.Context, key string, agent *models.Agent) error
	OnDelete        func(ctx context.Context, key string) error
	OnListByUserID  func(ctx context.Context, userID string) ([]*models.Agent, error)
}

// Get implements contracts.Agents.
func (m *MockAgents) Get(ctx context.Context, key string) (*models.Agent, error) {
	if m.OnGet != nil {
		return m.OnGet(ctx, key)
	}
	return &models.Agent{}, nil
}

// Set implements contracts.Agents.
func (m *MockAgents) Set(ctx context.Context, key string, agent *models.Agent) error {
	if m.OnSet != nil {
		return m.OnSet(ctx, key, agent)
	}
	return nil
}

// Delete implements contracts.Agents.
func (m *MockAgents) Delete(ctx context.Context, key string) error {
	if m.OnDelete != nil {
		return m.OnDelete(ctx, key)
	}
	return nil
}

// ListByUserID implements contracts.Agents.
func (m *MockAgents) ListByUserID(ctx context.Context, userID string) ([]*models.Agent, error) {
	if m.OnListByUserID != nil {
		return m.OnListByUserID(ctx, userID)
	}
	return []*models.Agent{}, nil
}

// MockMemories implements contracts.Memories using function fields.
// Set each On* field to control behaviour per test.
type MockMemories struct {
	OnGet                func(ctx context.Context, key string) (*models.Memory, error)
	OnSet                func(ctx context.Context, key string, memory *models.Memory) error
	OnDelete             func(ctx context.Context, key string) error
	OnListByAgentID      func(ctx context.Context, agentID string) ([]*models.Memory, error)
	OnGetByCorrelationID func(ctx context.Context, correlationID string) (*models.Memory, error)
}

// Get implements contracts.Memories.
func (m *MockMemories) Get(ctx context.Context, key string) (*models.Memory, error) {
	if m.OnGet != nil {
		return m.OnGet(ctx, key)
	}
	return &models.Memory{}, nil
}

// Set implements contracts.Memories.
func (m *MockMemories) Set(ctx context.Context, key string, memory *models.Memory) error {
	if m.OnSet != nil {
		return m.OnSet(ctx, key, memory)
	}
	return nil
}

// Delete implements contracts.Memories.
func (m *MockMemories) Delete(ctx context.Context, key string) error {
	if m.OnDelete != nil {
		return m.OnDelete(ctx, key)
	}
	return nil
}

// ListByAgentID implements contracts.Memories.
func (m *MockMemories) ListByAgentID(ctx context.Context, agentID string) ([]*models.Memory, error) {
	if m.OnListByAgentID != nil {
		return m.OnListByAgentID(ctx, agentID)
	}
	return []*models.Memory{}, nil
}

// GetByCorrelationID implements contracts.Memories.
func (m *MockMemories) GetByCorrelationID(ctx context.Context, correlationID string) (*models.Memory, error) {
	if m.OnGetByCorrelationID != nil {
		return m.OnGetByCorrelationID(ctx, correlationID)
	}
	return &models.Memory{}, nil
}

// MockChunks implements contracts.Chunks using function fields.
// Set each On* field to control behaviour per test.
type MockChunks struct {
	OnGet                func(ctx context.Context, key string) (*models.Chunk, error)
	OnSet                func(ctx context.Context, key string, chunk *models.Chunk) error
	OnDelete             func(ctx context.Context, key string) error
	OnListByMemoryID     func(ctx context.Context, memoryID string) ([]*models.Chunk, error)
	OnSearchByEmbedding  func(ctx context.Context, embedding []float32, limit int) ([]*models.Chunk, error)
}

// Get implements contracts.Chunks.
func (m *MockChunks) Get(ctx context.Context, key string) (*models.Chunk, error) {
	if m.OnGet != nil {
		return m.OnGet(ctx, key)
	}
	return &models.Chunk{}, nil
}

// Set implements contracts.Chunks.
func (m *MockChunks) Set(ctx context.Context, key string, chunk *models.Chunk) error {
	if m.OnSet != nil {
		return m.OnSet(ctx, key, chunk)
	}
	return nil
}

// Delete implements contracts.Chunks.
func (m *MockChunks) Delete(ctx context.Context, key string) error {
	if m.OnDelete != nil {
		return m.OnDelete(ctx, key)
	}
	return nil
}

// ListByMemoryID implements contracts.Chunks.
func (m *MockChunks) ListByMemoryID(ctx context.Context, memoryID string) ([]*models.Chunk, error) {
	if m.OnListByMemoryID != nil {
		return m.OnListByMemoryID(ctx, memoryID)
	}
	return []*models.Chunk{}, nil
}

// SearchByEmbedding implements contracts.Chunks.
func (m *MockChunks) SearchByEmbedding(ctx context.Context, embedding []float32, limit int) ([]*models.Chunk, error) {
	if m.OnSearchByEmbedding != nil {
		return m.OnSearchByEmbedding(ctx, embedding, limit)
	}
	return []*models.Chunk{}, nil
}

// MockChunker implements contracts.Chunker using a function field.
// Set OnChunk to control behaviour per test.
type MockChunker struct {
	OnChunk func(ctx context.Context, content []byte) ([]contracts.Chunk, error)
}

// Chunk implements contracts.Chunker.
func (m *MockChunker) Chunk(ctx context.Context, content []byte) ([]contracts.Chunk, error) {
	if m.OnChunk != nil {
		return m.OnChunk(ctx, content)
	}
	return []contracts.Chunk{}, nil
}

// MockEmbedder implements contracts.Embedder using a function field.
// Set OnEmbed to control behaviour per test.
type MockEmbedder struct {
	OnEmbed func(ctx context.Context, texts []string) ([][]float32, error)
}

// Embed implements contracts.Embedder.
func (m *MockEmbedder) Embed(ctx context.Context, texts []string) ([][]float32, error) {
	if m.OnEmbed != nil {
		return m.OnEmbed(ctx, texts)
	}
	result := make([][]float32, len(texts))
	for i := range result {
		result[i] = []float32{0.1, 0.2, 0.3}
	}
	return result, nil
}

// MockIngestJobs implements contracts.IngestJobs using function fields.
// Set each On* field to control behaviour per test.
type MockIngestJobs struct {
	OnGet            func(ctx context.Context, id string) (*models.IngestJob, error)
	OnSave           func(ctx context.Context, job *models.IngestJob) error
	OnListByStatus   func(ctx context.Context, status models.JobStatus) ([]*models.IngestJob, error)
}

// Get implements contracts.IngestJobs.
func (m *MockIngestJobs) Get(ctx context.Context, id string) (*models.IngestJob, error) {
	if m.OnGet != nil {
		return m.OnGet(ctx, id)
	}
	return &models.IngestJob{}, nil
}

// Save implements contracts.IngestJobs.
func (m *MockIngestJobs) Save(ctx context.Context, job *models.IngestJob) error {
	if m.OnSave != nil {
		return m.OnSave(ctx, job)
	}
	return nil
}

// ListByStatus implements contracts.IngestJobs.
func (m *MockIngestJobs) ListByStatus(ctx context.Context, status models.JobStatus) ([]*models.IngestJob, error) {
	if m.OnListByStatus != nil {
		return m.OnListByStatus(ctx, status)
	}
	return []*models.IngestJob{}, nil
}
