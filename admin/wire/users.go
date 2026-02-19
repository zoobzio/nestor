// Package wire defines request and response types for the admin API surface.
package wire

import (
	"time"

	"github.com/zoobzio/check"
)

// UserResponse is the admin API response for user data.
type UserResponse struct {
	ID        string    `json:"id" description:"User UUID" example:"550e8400-e29b-41d4-a716-446655440000"`
	Name      string    `json:"name" description:"User name" example:"Jane Doe"`
	CreatedAt time.Time `json:"created_at" description:"Account creation time"`
	UpdatedAt time.Time `json:"updated_at" description:"Last update time"`
}

// Clone returns a shallow copy. All fields are value types.
func (u UserResponse) Clone() UserResponse {
	return u
}

// UserListResponse is the admin API response for listing users.
type UserListResponse struct {
	Users []UserResponse `json:"users" description:"List of users"`
}

// Clone returns a deep copy.
func (u UserListResponse) Clone() UserListResponse {
	r := u
	if u.Users != nil {
		r.Users = make([]UserResponse, len(u.Users))
		copy(r.Users, u.Users)
	}
	return r
}

// CreateUserRequest is the request body for creating a user.
type CreateUserRequest struct {
	Name string `json:"name" description:"User name" example:"Jane Doe"`
}

// Validate validates the request.
func (r *CreateUserRequest) Validate() error {
	return check.All(
		check.Str(r.Name, "name").Required().MaxLen(255).V(),
	).Err()
}

// Clone returns a shallow copy.
func (r CreateUserRequest) Clone() CreateUserRequest {
	return r
}

// UpdateUserRequest is the request body for updating a user.
type UpdateUserRequest struct {
	Name string `json:"name" description:"User name" example:"Jane Doe"`
}

// Validate validates the request.
func (r *UpdateUserRequest) Validate() error {
	return check.All(
		check.Str(r.Name, "name").Required().MaxLen(255).V(),
	).Err()
}

// Clone returns a shallow copy.
func (r UpdateUserRequest) Clone() UpdateUserRequest {
	return r
}
