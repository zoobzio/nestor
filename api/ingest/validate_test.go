//go:build testing

package ingest

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/zoobzio/nestor/models"
	nesttest "github.com/zoobzio/nestor/testing"
)

const (
	testAgentID    = "550e8400-e29b-41d4-a716-446655440001"
	testCustomerID = "550e8400-e29b-41d4-a716-446655440000"
	testJobID      = "550e8400-e29b-41d4-a716-446655440010"
)

// newValidateJob returns a minimal IngestJob ready for the validate stage.
func newValidateJob() *models.IngestJob {
	return &models.IngestJob{
		ID:         testJobID,
		CustomerID: testCustomerID,
		AgentID:    testAgentID,
		Content:    "# Hello\n\nSome markdown content.",
		Stage:      models.JobStageValidate,
		Status:     models.JobStatusPending,
	}
}

// TestValidateStage_Success verifies a valid job passes without error.
func TestValidateStage_Success(t *testing.T) {
	agent := &models.Agent{ID: testAgentID}
	mock := &nesttest.MockAgents{
		OnGet: func(_ context.Context, key string) (*models.Agent, error) {
			if key == testAgentID {
				return agent, nil
			}
			return nil, nil
		},
	}
	nesttest.SetupRegistry(t, nesttest.WithAgents(mock))

	job := newValidateJob()
	result, err := validateStage(context.Background(), job)

	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if result.Stage != models.JobStageValidate {
		t.Errorf("stage = %q, want %q", result.Stage, models.JobStageValidate)
	}
}

// TestValidateStage_AgentNotFound verifies ErrAgentNotFound when Get returns nil.
func TestValidateStage_AgentNotFound(t *testing.T) {
	mock := &nesttest.MockAgents{
		OnGet: func(_ context.Context, _ string) (*models.Agent, error) {
			return nil, nil
		},
	}
	nesttest.SetupRegistry(t, nesttest.WithAgents(mock))

	job := newValidateJob()
	_, err := validateStage(context.Background(), job)

	if !errors.Is(err, ErrAgentNotFound) {
		t.Errorf("error = %v, want ErrAgentNotFound", err)
	}
}

// TestValidateStage_AgentGetError verifies the store error propagates.
func TestValidateStage_AgentGetError(t *testing.T) {
	storeErr := errors.New("database unavailable")
	mock := &nesttest.MockAgents{
		OnGet: func(_ context.Context, _ string) (*models.Agent, error) {
			return nil, storeErr
		},
	}
	nesttest.SetupRegistry(t, nesttest.WithAgents(mock))

	job := newValidateJob()
	_, err := validateStage(context.Background(), job)

	if !errors.Is(err, storeErr) {
		t.Errorf("error = %v, want %v", err, storeErr)
	}
}

// TestValidateStage_EmptyContent verifies ErrContentEmpty for a zero-length content field.
func TestValidateStage_EmptyContent(t *testing.T) {
	agent := &models.Agent{ID: testAgentID}
	mock := &nesttest.MockAgents{
		OnGet: func(_ context.Context, _ string) (*models.Agent, error) {
			return agent, nil
		},
	}
	nesttest.SetupRegistry(t, nesttest.WithAgents(mock))

	job := newValidateJob()
	job.Content = ""

	_, err := validateStage(context.Background(), job)

	if !errors.Is(err, ErrContentEmpty) {
		t.Errorf("error = %v, want ErrContentEmpty", err)
	}
}

// TestValidateStage_ContentTooLarge verifies ErrContentTooLarge when content exceeds the limit.
func TestValidateStage_ContentTooLarge(t *testing.T) {
	agent := &models.Agent{ID: testAgentID}
	mock := &nesttest.MockAgents{
		OnGet: func(_ context.Context, _ string) (*models.Agent, error) {
			return agent, nil
		},
	}
	nesttest.SetupRegistry(t, nesttest.WithAgents(mock))

	// Default limit is 10MB; build content just over that.
	job := newValidateJob()
	job.Content = strings.Repeat("a", 10*1024*1024+1)

	_, err := validateStage(context.Background(), job)

	if !errors.Is(err, ErrContentTooLarge) {
		t.Errorf("error = %v, want ErrContentTooLarge", err)
	}
}
