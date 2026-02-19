package handlers

import "github.com/zoobzio/rocco"

// All returns all admin API handlers for registration with the router.
func All() []rocco.Endpoint {
	return []rocco.Endpoint{
		// Users
		ListUsers,
		CreateUser,
		GetUser,
		UpdateUser,
		DeleteUser,
	}
}
