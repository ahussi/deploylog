package jenkins_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/deploylog/deploylog/internal/source/jenkins"
)

func TestClient_Name(t *testing.T) {
	c := jenkins.NewClient("http://localhost:8080", "my-job")
	if c.Name() != jenkins.SourceName {
		t.Errorf("expected %q, got %q", jenkins.SourceName, c.Name())
	}
}

func TestSourceName_Constant(t *testing.T) {
	if jenkins.SourceName != "jenkins" {
		t.Errorf("unexpected source name: %s", jenkins.SourceName)
	}
}

func TestClient_Fetch_Success(t *testing.T) {
	payload := map[string]interface{}{
		"builds": []map[string]interface{}{
			{
				"id":               "42",
				"fullDisplayName": "my-job #42",
				"result":          "SUCCESS",
				"timestamp":       int64(1700000000000),
				"duration":        int64(3200),
				"description":     "deploy to prod",
			},
		},
	}
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(payload)
	}))
	defer ts.Close()

	c := jenkins.NewClient(ts.URL, "my-job")
	events, err := c.Fetch()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(events))
	}
	if events[0].ID != "42" {
		t.Errorf("expected ID 42, got %s", events[0].ID)
	}
	if events[0].Source != "jenkins" {
		t.Errorf("expected source jenkins, got %s", events[0].Source)
	}
}

func TestClient_Fetch_HTTPError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ts.Close()

	c := jenkins.NewClient(ts.URL, "my-job")
	_, err := c.Fetch()
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestClient_Implements_Interface(t *testing.T) {
	// Compile-time check via assignment — if this builds, the interface is satisfied.
	c := jenkins.NewClient("http://localhost", "job")
	_ = c.Name
	_ = c.Fetch
}
