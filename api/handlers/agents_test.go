//go:build testing

package handlers_test

import (
	"context"
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/zoobzio/rocco"
	rtesting "github.com/zoobzio/rocco/testing"
	"github.com/zoobzio/nestor/api/handlers"
	"github.com/zoobzio/nestor/models"
	nesttest "github.com/zoobzio/nestor/testing"
)

const testUserID = "550e8400-e29b-41d4-a716-446655440000"
const testAgentID = "550e8400-e29b-41d4-a716-446655440001"

// newTestEngine builds a rocco Engine with authentication and the given handlers registered.
func newTestEngine(t *testing.T, endpoints ...rocco.Endpoint) *rocco.Engine {
	t.Helper()
	identity := rtesting.NewMockIdentity(testUserID)
	engine := rtesting.TestEngineWithAuth(func(_ context.Context, _ *http.Request) (rocco.Identity, error) {
		return identity, nil
	})
	engine.WithHandlers(endpoints...)
	return engine
}

// newTestAgent returns a consistent agent for handler tests.
func newTestAgent(t *testing.T) *models.Agent {
	t.Helper()
	now := time.Now()
	return &models.Agent{
		ID:        testAgentID,
		UserID:    testUserID,
		Name:      "Research Assistant",
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// TestListAgents_Success verifies a 200 response with an agents list.
func TestListAgents_Success(t *testing.T) {
	agent := newTestAgent(t)
	mock := &nesttest.MockAgents{
		OnListByUserID: func(_ context.Context, userID string) ([]*models.Agent, error) {
			return []*models.Agent{agent}, nil
		},
	}
	nesttest.SetupRegistry(t, nesttest.WithAgents(mock))

	engine := newTestEngine(t, handlers.ListAgents)
	capture := rtesting.ServeRequest(engine, http.MethodGet, "/agents", nil)

	rtesting.AssertStatus(t, capture, http.StatusOK)

	var resp struct {
		Agents []struct {
			ID   string `json:"id"`
			Name string `json:"name"`
		} `json:"agents"`
	}
	if err := capture.DecodeJSON(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if len(resp.Agents) != 1 {
		t.Fatalf("expected 1 agent, got %d", len(resp.Agents))
	}
	if resp.Agents[0].ID != testAgentID {
		t.Errorf("agent ID %q, want %q", resp.Agents[0].ID, testAgentID)
	}
	if resp.Agents[0].Name != "Research Assistant" {
		t.Errorf("agent name %q, want %q", resp.Agents[0].Name, "Research Assistant")
	}
}

// TestListAgents_Empty verifies a 200 response with an empty list.
func TestListAgents_Empty(t *testing.T) {
	mock := &nesttest.MockAgents{
		OnListByUserID: func(_ context.Context, _ string) ([]*models.Agent, error) {
			return []*models.Agent{}, nil
		},
	}
	nesttest.SetupRegistry(t, nesttest.WithAgents(mock))

	engine := newTestEngine(t, handlers.ListAgents)
	capture := rtesting.ServeRequest(engine, http.MethodGet, "/agents", nil)

	rtesting.AssertStatus(t, capture, http.StatusOK)

	var resp struct {
		Agents []any `json:"agents"`
	}
	if err := capture.DecodeJSON(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if len(resp.Agents) != 0 {
		t.Errorf("expected empty agents list, got %d", len(resp.Agents))
	}
}

// TestListAgents_StoreError verifies a 500 response when the store errors.
func TestListAgents_StoreError(t *testing.T) {
	mock := &nesttest.MockAgents{
		OnListByUserID: func(_ context.Context, _ string) ([]*models.Agent, error) {
			return nil, errors.New("database unavailable")
		},
	}
	nesttest.SetupRegistry(t, nesttest.WithAgents(mock))

	engine := newTestEngine(t, handlers.ListAgents)
	capture := rtesting.ServeRequest(engine, http.MethodGet, "/agents", nil)

	rtesting.AssertStatus(t, capture, http.StatusInternalServerError)
}

// TestCreateAgent_Success verifies a 201 response with the created agent.
func TestCreateAgent_Success(t *testing.T) {
	var storedAgent *models.Agent
	mock := &nesttest.MockAgents{
		OnSet: func(_ context.Context, _ string, agent *models.Agent) error {
			storedAgent = agent
			return nil
		},
	}
	nesttest.SetupRegistry(t, nesttest.WithAgents(mock))

	engine := newTestEngine(t, handlers.CreateAgent)

	body := map[string]string{"name": "New Agent"}
	capture := rtesting.ServeRequest(engine, http.MethodPost, "/agents", body)

	rtesting.AssertStatus(t, capture, http.StatusCreated)

	if storedAgent == nil {
		t.Fatal("expected agent to be stored, got nil")
	}
	if storedAgent.Name != "New Agent" {
		t.Errorf("stored agent name %q, want %q", storedAgent.Name, "New Agent")
	}
	if storedAgent.UserID != testUserID {
		t.Errorf("stored agent UserID %q, want %q", storedAgent.UserID, testUserID)
	}

	var resp struct {
		Name string `json:"name"`
	}
	if err := capture.DecodeJSON(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if resp.Name != "New Agent" {
		t.Errorf("response name %q, want %q", resp.Name, "New Agent")
	}
}

// TestCreateAgent_StoreError verifies a 500 response when Set fails.
func TestCreateAgent_StoreError(t *testing.T) {
	mock := &nesttest.MockAgents{
		OnSet: func(_ context.Context, _ string, _ *models.Agent) error {
			return errors.New("storage failure")
		},
	}
	nesttest.SetupRegistry(t, nesttest.WithAgents(mock))

	engine := newTestEngine(t, handlers.CreateAgent)

	body := map[string]string{"name": "New Agent"}
	capture := rtesting.ServeRequest(engine, http.MethodPost, "/agents", body)

	rtesting.AssertStatus(t, capture, http.StatusInternalServerError)
}

// TestGetAgent_Success verifies a 200 response with agent data.
func TestGetAgent_Success(t *testing.T) {
	agent := newTestAgent(t)
	mock := &nesttest.MockAgents{
		OnGet: func(_ context.Context, key string) (*models.Agent, error) {
			if key != agent.ID {
				return nil, errors.New("not found")
			}
			return agent, nil
		},
	}
	nesttest.SetupRegistry(t, nesttest.WithAgents(mock))

	engine := newTestEngine(t, handlers.GetAgent)
	capture := rtesting.ServeRequest(engine, http.MethodGet, "/agents/"+testAgentID, nil)

	rtesting.AssertStatus(t, capture, http.StatusOK)

	var resp struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	}
	if err := capture.DecodeJSON(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if resp.ID != testAgentID {
		t.Errorf("response ID %q, want %q", resp.ID, testAgentID)
	}
	if resp.Name != "Research Assistant" {
		t.Errorf("response name %q, want %q", resp.Name, "Research Assistant")
	}
}

// TestGetAgent_NotFound verifies a 404 response when the agent does not exist.
func TestGetAgent_NotFound(t *testing.T) {
	mock := &nesttest.MockAgents{
		OnGet: func(_ context.Context, _ string) (*models.Agent, error) {
			return nil, errors.New("not found")
		},
	}
	nesttest.SetupRegistry(t, nesttest.WithAgents(mock))

	engine := newTestEngine(t, handlers.GetAgent)
	capture := rtesting.ServeRequest(engine, http.MethodGet, "/agents/does-not-exist", nil)

	rtesting.AssertStatus(t, capture, http.StatusNotFound)
	rtesting.AssertErrorCode(t, capture, "NOT_FOUND")
}

// TestUpdateAgent_Success verifies a 200 response with updated agent data.
func TestUpdateAgent_Success(t *testing.T) {
	agent := newTestAgent(t)
	mock := &nesttest.MockAgents{
		OnGet: func(_ context.Context, key string) (*models.Agent, error) {
			if key != agent.ID {
				return nil, errors.New("not found")
			}
			return agent, nil
		},
		OnSet: func(_ context.Context, _ string, updated *models.Agent) error {
			agent = updated
			return nil
		},
	}
	nesttest.SetupRegistry(t, nesttest.WithAgents(mock))

	engine := newTestEngine(t, handlers.UpdateAgent)

	body := map[string]string{"name": "Updated Name"}
	capture := rtesting.ServeRequest(engine, http.MethodPatch, "/agents/"+testAgentID, body)

	rtesting.AssertStatus(t, capture, http.StatusOK)

	var resp struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	}
	if err := capture.DecodeJSON(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if resp.Name != "Updated Name" {
		t.Errorf("response name %q, want %q", resp.Name, "Updated Name")
	}
}

// TestUpdateAgent_NotFound verifies a 404 when the agent does not exist.
func TestUpdateAgent_NotFound(t *testing.T) {
	mock := &nesttest.MockAgents{
		OnGet: func(_ context.Context, _ string) (*models.Agent, error) {
			return nil, errors.New("not found")
		},
	}
	nesttest.SetupRegistry(t, nesttest.WithAgents(mock))

	engine := newTestEngine(t, handlers.UpdateAgent)

	body := map[string]string{"name": "Updated Name"}
	capture := rtesting.ServeRequest(engine, http.MethodPatch, "/agents/does-not-exist", body)

	rtesting.AssertStatus(t, capture, http.StatusNotFound)
	rtesting.AssertErrorCode(t, capture, "NOT_FOUND")
}

// TestUpdateAgent_SetError verifies a 500 when Set fails after a successful Get.
func TestUpdateAgent_SetError(t *testing.T) {
	agent := newTestAgent(t)
	mock := &nesttest.MockAgents{
		OnGet: func(_ context.Context, key string) (*models.Agent, error) {
			if key != agent.ID {
				return nil, errors.New("not found")
			}
			return agent, nil
		},
		OnSet: func(_ context.Context, _ string, _ *models.Agent) error {
			return errors.New("storage failure")
		},
	}
	nesttest.SetupRegistry(t, nesttest.WithAgents(mock))

	engine := newTestEngine(t, handlers.UpdateAgent)

	body := map[string]string{"name": "Updated Name"}
	capture := rtesting.ServeRequest(engine, http.MethodPatch, "/agents/"+testAgentID, body)

	rtesting.AssertStatus(t, capture, http.StatusInternalServerError)
}

// TestDeleteAgent_Success verifies a 204 response on successful deletion.
func TestDeleteAgent_Success(t *testing.T) {
	var deletedID string
	mock := &nesttest.MockAgents{
		OnDelete: func(_ context.Context, key string) error {
			deletedID = key
			return nil
		},
	}
	nesttest.SetupRegistry(t, nesttest.WithAgents(mock))

	engine := newTestEngine(t, handlers.DeleteAgent)
	capture := rtesting.ServeRequest(engine, http.MethodDelete, "/agents/"+testAgentID, nil)

	rtesting.AssertStatus(t, capture, http.StatusNoContent)

	if deletedID != testAgentID {
		t.Errorf("deleted ID %q, want %q", deletedID, testAgentID)
	}
}

// TestDeleteAgent_NotFound verifies a 404 when the agent does not exist.
func TestDeleteAgent_NotFound(t *testing.T) {
	mock := &nesttest.MockAgents{
		OnDelete: func(_ context.Context, _ string) error {
			return errors.New("not found")
		},
	}
	nesttest.SetupRegistry(t, nesttest.WithAgents(mock))

	engine := newTestEngine(t, handlers.DeleteAgent)
	capture := rtesting.ServeRequest(engine, http.MethodDelete, "/agents/does-not-exist", nil)

	rtesting.AssertStatus(t, capture, http.StatusNotFound)
	rtesting.AssertErrorCode(t, capture, "NOT_FOUND")
}
