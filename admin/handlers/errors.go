// Package handlers provides HTTP handlers for the admin API surface.
package handlers

import "github.com/zoobzio/rocco"

var (
	// ErrUserNotFound is returned when a user cannot be located by ID.
	ErrUserNotFound = rocco.ErrNotFound.WithMessage("user not found")
)
