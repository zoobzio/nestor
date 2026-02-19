// Package config holds static startup configuration types for the application.
package config

import "github.com/zoobzio/check"

// App holds server configuration for the public API.
// Loaded once at startup; runtime changes require restart.
type App struct {
	// Port is the HTTP port the server listens on.
	Port int `env:"NESTOR_PORT" default:"8080"`
}

// Validate checks App configuration for required values.
func (c App) Validate() error {
	return check.Int(c.Port, "port").Positive().V().Err()
}
