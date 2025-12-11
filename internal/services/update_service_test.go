package services

import (
	"encoding/json"
	"freelance-flow/internal/update"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestUpdateService_CheckForUpdate(t *testing.T) {
	// Mock GitHub API
	mockRelease := struct {
		TagName string `json:"tag_name"`
	}{
		TagName: "v1.0.1",
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/repos/royzhu/freelance-flow/releases/latest" {
			if err := json.NewEncoder(w).Encode(mockRelease); err != nil {
				t.Fatalf("failed to encode mock release: %v", err)
			}
			return
		}
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	// Override base URL
	oldURL := update.GitHubAPIBaseURL
	update.GitHubAPIBaseURL = server.URL
	defer func() { update.GitHubAPIBaseURL = oldURL }()

	service := NewUpdateService()

	// Test CheckForUpdate
	err := service.CheckForUpdate()
	if err != nil {
		t.Errorf("CheckForUpdate() error = %v", err)
	}

	state := service.GetUpdateState()
	if state.Status != update.StatusAvailable {
		t.Errorf("Expected status Available, got %s. Error: %s", state.Status, state.Error)
	}
	if state.LatestVersion != "1.0.1" {
		t.Errorf("Expected latest version 1.0.1, got %s", state.LatestVersion)
	}
}
