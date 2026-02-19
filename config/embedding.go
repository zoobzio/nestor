// Package config holds static startup configuration types for the application.
package config

import "github.com/zoobzio/check"

// Embedding holds configuration for the embedding provider.
// Loaded once at startup; runtime changes require restart.
type Embedding struct {
	// Provider is the embedding backend to use.
	// Supported values: stub, openai, voyage, gemini, cohere.
	// Defaults to "stub" which returns zero-filled vectors without any network calls.
	Provider string `env:"NESTOR_EMBEDDING_PROVIDER" default:"stub"`

	// Model is the specific model name for the chosen provider.
	// Example: "text-embedding-ada-002" for OpenAI.
	Model string `env:"NESTOR_EMBEDDING_MODEL"`

	// APIKey is the authentication credential for the embedding provider.
	// Not required when Provider is "stub" or empty.
	APIKey string `env:"NESTOR_EMBEDDING_API_KEY"`

	// Dimensions is the vector dimensionality produced by the model.
	// Defaults to 1536, matching OpenAI's text-embedding-ada-002.
	Dimensions int `env:"NESTOR_EMBEDDING_DIMENSIONS" default:"1536"`
}

// Validate checks that the configuration is coherent.
func (c Embedding) Validate() error {
	return check.All(
		check.Int(c.Dimensions, "dimensions").Positive().V(),
	).Err()
}
