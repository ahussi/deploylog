package gitlab_test

import (
	"testing"

	"github.com/deploylog/deploylog/internal/config"
	"github.com/deploylog/deploylog/internal/source/gitlab"
)

func TestNewClientFromConfig_Valid(t *testing.T) {
	cfg := config.SourceConfig{
		Name: "gl-prod",
		Type: "gitlab",
		Params: map[string]string{
			"token":      "glpat-abc123",
			"project_id": "42",
		},
	}
	client, err := gitlab.NewClientFromConfig(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if client.Name() != "gl-prod" {
		t.Errorf("expected name gl-prod, got %s", client.Name())
	}
}

func TestNewClientFromConfig_WithBaseURL(t *testing.T) {
	cfg := config.SourceConfig{
		Name: "gl-self",
		Type: "gitlab",
		Params: map[string]string{
			"token":      "glpat-xyz",
			"project_id": "7",
			"base_url":   "https://gitlab.example.com",
		},
	}
	client, err := gitlab.NewClientFromConfig(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if client == nil {
		t.Fatal("expected non-nil client")
	}
}

func TestNewClientFromConfig_MissingToken(t *testing.T) {
	cfg := config.SourceConfig{
		Name: "gl-bad",
		Type: "gitlab",
		Params: map[string]string{
			"project_id": "10",
		},
	}
	_, err := gitlab.NewClientFromConfig(cfg)
	if err == nil {
		t.Fatal("expected error for missing token")
	}
}

func TestNewClientFromConfig_MissingProjectID(t *testing.T) {
	cfg := config.SourceConfig{
		Name: "gl-bad",
		Type: "gitlab",
		Params: map[string]string{
			"token": "glpat-abc",
		},
	}
	_, err := gitlab.NewClientFromConfig(cfg)
	if err == nil {
		t.Fatal("expected error for missing project_id")
	}
}

func TestNewClientFromConfig_MissingBoth(t *testing.T) {
	cfg := config.SourceConfig{
		Name: "gl-empty",
		Type: "gitlab",
		Params: map[string]string{},
	}
	_, err := gitlab.NewClientFromConfig(cfg)
	if err == nil {
		t.Fatal("expected error for missing token and project_id")
	}
}
