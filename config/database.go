// Package config holds static startup configuration types for the application.
package config

import "github.com/zoobzio/check"

// Database holds configuration for the PostgreSQL database connection.
// Loaded once at startup; runtime changes require restart.
type Database struct {
	// DSN is the full PostgreSQL connection string.
	// Example: "host=localhost port=5432 user=nestor password=secret dbname=nestor sslmode=disable"
	DSN string `env:"NESTOR_DATABASE_DSN"`
}

// Validate checks that the configuration is coherent.
func (c Database) Validate() error {
	return check.All(
		check.Str(c.DSN, "dsn").Required().V(),
	).Err()
}
