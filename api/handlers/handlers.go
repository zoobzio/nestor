package handlers

import "github.com/zoobzio/rocco"

// All returns all public API handlers for registration with the router.
func All() []rocco.Endpoint {
	return []rocco.Endpoint{
		// Agents
		ListAgents,
		CreateAgent,
		GetAgent,
		UpdateAgent,
		DeleteAgent,

		// Memories
		ListMemories,
		CreateMemory,
		GetMemory,
		UpdateMemory,
		DeleteMemory,
		GetMemoryByCorrelationID,

		// Chunks
		ListChunks,
		CreateChunk,
		GetChunk,
		DeleteChunk,
		SearchChunks,

		// FilterGroups
		ListFilterGroups,
		CreateFilterGroup,
		GetFilterGroup,
		UpdateFilterGroup,
		DeleteFilterGroup,
	}
}
