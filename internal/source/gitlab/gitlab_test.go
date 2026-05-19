package gitlab_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/deploylog/deploylog/internal/source/gitlab"
)

func TestClient_Name(t *testing.T) {
	c := gitlab.NewClient("123", "token", "")
	if c.Name() != gitlab.SourceName {
		t.Errorf("expected %q, got %q", gitlab.SourceName, c.Name())
	}
}

func TestSourceName_Constant(t *testing.T) {
	if gitlab.SourceName != "gitlab" {
		t.Errorf("unexpected source name: %q", gitlab.SourceName)
	}
}

func TestClient_Fetch_Success(t *testing.T) {
	payload := []map[string]interface{}{
		{
			"id":         1,
			"status":     "success",
			"ref":        "main",
			"created_at": "2024-01-15T10:00:00Z",
			"environment": map[string]string{"name": "production"},
		},
		{
			"id":         2,
			"status":     "failed",
			"ref":        "feature-x",
			"created_at": "2024-01-15T09:00:00Z",
			"environment": map[string]string{"name": "staging"},
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("PRIVATE-TOKEN") == "" {
			t.Error("expected PRIVATE-TOKEN header")
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(payload)
	}))
	defer server.Close()

	client := gitlab.NewClient("42", "secret", server.URL)
	events, err := client.Fetch()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(events) != 2 {
		t.Fatalf("expected 2 events, got %d", len(events))
	}
	if events[0].Service != "production" {
		t.Errorf("expected service 'production', got %q", events[0].Service)
	}
	if events[1].Service != "staging" {
		t.Errorf("expected service 'staging', got %q", events[1].Service)
	}
}

func TestClient_Fetch_HTTPError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer server.Close()

	client := gitlab.NewClient("42", "bad-token", server.URL)
	_, err := client.Fetch()
	if err == nil {
		t.Fatal("expected error for non-200 response")
	}
}

func TestClient_Implements_Interface(t *testing.T) {
	// Compile-time check via assignment — if this builds, the interface is satisfied.
	client := gitlab.NewClient("1", "tok", "")
	_ = client.Name
	_ = client.Fetch
}
