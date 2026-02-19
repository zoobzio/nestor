// Package transformers provides pure functions for mapping between models and wire types
// on the admin API surface.
package transformers

import (
	"github.com/zoobzio/nestor/admin/wire"
	"github.com/zoobzio/nestor/models"
)

// UserToResponse transforms a User model to an admin API response.
func UserToResponse(u *models.User) wire.UserResponse {
	return wire.UserResponse{
		ID:        u.ID,
		Name:      u.Name,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}

// UsersToList transforms a slice of User models to a list response.
func UsersToList(users []*models.User) wire.UserListResponse {
	resp := wire.UserListResponse{
		Users: make([]wire.UserResponse, len(users)),
	}
	for i, u := range users {
		resp.Users[i] = UserToResponse(u)
	}
	return resp
}

// ApplyCreateUser applies a CreateUserRequest to a User model.
func ApplyCreateUser(req wire.CreateUserRequest, u *models.User) {
	u.Name = req.Name
}

// ApplyUpdateUser applies an UpdateUserRequest to a User model.
func ApplyUpdateUser(req wire.UpdateUserRequest, u *models.User) {
	u.Name = req.Name
}
