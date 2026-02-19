package handlers

import (
	"crypto/rand"
	"fmt"
	"time"

	"github.com/zoobzio/rocco"
	"github.com/zoobzio/sum"
	"github.com/zoobzio/nestor/admin/contracts"
	"github.com/zoobzio/nestor/admin/transformers"
	"github.com/zoobzio/nestor/admin/wire"
	"github.com/zoobzio/nestor/models"
)

// newUUID generates a random UUID v4.
func newUUID() string {
	var b [16]byte
	_, _ = rand.Read(b[:])
	b[6] = (b[6] & 0x0f) | 0x40
	b[8] = (b[8] & 0x3f) | 0x80
	return fmt.Sprintf("%08x-%04x-%04x-%04x-%012x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
}

// ListUsers returns all users.
var ListUsers = rocco.GET("/users", func(req *rocco.Request[rocco.NoBody]) (wire.UserListResponse, error) {
	users := sum.MustUse[contracts.Users](req.Context)

	list, err := users.List(req.Context)
	if err != nil {
		return wire.UserListResponse{}, err
	}

	return transformers.UsersToList(list), nil
}).WithSummary("List users").
	WithDescription("Returns all users.").
	WithTags("Users").
	WithAuthentication()

// CreateUser creates a new user.
var CreateUser = rocco.POST("/users", func(req *rocco.Request[wire.CreateUserRequest]) (wire.UserResponse, error) {
	users := sum.MustUse[contracts.Users](req.Context)

	user := &models.User{
		ID:        newUUID(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	transformers.ApplyCreateUser(req.Body, user)

	if err := users.Set(req.Context, user.ID, user); err != nil {
		return wire.UserResponse{}, err
	}

	return transformers.UserToResponse(user), nil
}).WithSummary("Create user").
	WithDescription("Creates a new user.").
	WithTags("Users").
	WithAuthentication().
	WithSuccessStatus(201)

// GetUser returns a single user by ID.
var GetUser = rocco.GET("/users/{id}", func(req *rocco.Request[rocco.NoBody]) (wire.UserResponse, error) {
	users := sum.MustUse[contracts.Users](req.Context)

	id := req.Params.Path["id"]

	user, err := users.Get(req.Context, id)
	if err != nil {
		return wire.UserResponse{}, ErrUserNotFound
	}

	return transformers.UserToResponse(user), nil
}).WithPathParams("id").
	WithSummary("Get user").
	WithDescription("Returns a user by ID.").
	WithTags("Users").
	WithErrors(ErrUserNotFound).
	WithAuthentication()

// UpdateUser updates an existing user.
var UpdateUser = rocco.PATCH("/users/{id}", func(req *rocco.Request[wire.UpdateUserRequest]) (wire.UserResponse, error) {
	users := sum.MustUse[contracts.Users](req.Context)

	id := req.Params.Path["id"]

	user, err := users.Get(req.Context, id)
	if err != nil {
		return wire.UserResponse{}, ErrUserNotFound
	}

	transformers.ApplyUpdateUser(req.Body, user)
	user.UpdatedAt = time.Now()

	if err := users.Set(req.Context, id, user); err != nil {
		return wire.UserResponse{}, err
	}

	return transformers.UserToResponse(user), nil
}).WithPathParams("id").
	WithSummary("Update user").
	WithDescription("Updates a user's details.").
	WithTags("Users").
	WithErrors(ErrUserNotFound).
	WithAuthentication()

// DeleteUser removes a user by ID.
var DeleteUser = rocco.DELETE("/users/{id}", func(req *rocco.Request[rocco.NoBody]) (rocco.NoBody, error) {
	users := sum.MustUse[contracts.Users](req.Context)

	id := req.Params.Path["id"]

	if err := users.Delete(req.Context, id); err != nil {
		return rocco.NoBody{}, ErrUserNotFound
	}

	return rocco.NoBody{}, nil
}).WithPathParams("id").
	WithSummary("Delete user").
	WithDescription("Deletes a user by ID.").
	WithTags("Users").
	WithErrors(ErrUserNotFound).
	WithAuthentication().
	WithSuccessStatus(204)
