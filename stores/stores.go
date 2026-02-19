// Package stores provides shared data access layer implementations.
package stores

import (
	"github.com/jmoiron/sqlx"
	"github.com/zoobzio/astql"
)

// Stores aggregates all store implementations.
// Stores are shared across all API surfaces.
type Stores struct {
	Users        *Users
	Agents       *Agents
	Memories     *Memories
	Chunks       *Chunks
	IngestJobs   *IngestJobs
	FilterGroups *FilterGroups
}

// New creates and returns all stores.
func New(db *sqlx.DB, renderer astql.Renderer) (*Stores, error) {
	users, err := NewUsers(db, renderer)
	if err != nil {
		return nil, err
	}

	agents, err := NewAgents(db, renderer)
	if err != nil {
		return nil, err
	}

	memories, err := NewMemories(db, renderer)
	if err != nil {
		return nil, err
	}

	chunks, err := NewChunks(db, renderer)
	if err != nil {
		return nil, err
	}

	ingestJobs, err := NewIngestJobs(db, renderer)
	if err != nil {
		return nil, err
	}

	filterGroups, err := NewFilterGroups(db, renderer)
	if err != nil {
		return nil, err
	}

	return &Stores{
		Users:        users,
		Agents:       agents,
		Memories:     memories,
		Chunks:       chunks,
		IngestJobs:   ingestJobs,
		FilterGroups: filterGroups,
	}, nil
}
