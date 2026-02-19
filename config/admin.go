// Package config holds static startup configuration types for the application.
package config

import "github.com/zoobzio/check"

// Admin holds configuration for the admin API server.
// Loaded once at startup; runtime changes require restart.
type Admin struct {
	// Port is the TCP port the admin HTTP server listens on.
	Port int `env:"NESTOR_ADMIN_PORT" default:"8081"`
}

// Validate checks that the configuration is coherent.
func (c Admin) Validate() error {
	return check.All(
		check.Int(c.Port, "port").Positive().Max(65535).V(),
	).Err()
}
