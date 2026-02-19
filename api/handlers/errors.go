// Package handlers provides HTTP handlers for the public API surface.
package handlers

import "github.com/zoobzio/rocco"

var (
	ErrAgentNotFound       = rocco.ErrNotFound.WithMessage("agent not found")
	ErrMemoryNotFound      = rocco.ErrNotFound.WithMessage("memory not found")
	ErrChunkNotFound       = rocco.ErrNotFound.WithMessage("chunk not found")
	ErrFilterGroupNotFound = rocco.ErrNotFound.WithMessage("filter group not found")
)
