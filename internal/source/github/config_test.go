package github_test

import (
	"testing"

	"github.com/yourorg/deploylog/internal/config"
	"github.com/yourorg/deploylog/internal/source/github"
)

func TestNewClientFromConfig_Valid(t *testing.T) {
	cfg := config.SourceConfig{
		Name: "gh-prod",
		Type: "github",
		Options: map[string]string{
			"token": "secret",
			"repo":  "org/repo",
		},
	}

	client, err := github.NewClientFromConfig(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if client == nil {
		t.Fatal("expected non-nil client")
	}
	if client.Name() != "gh-prod" {
		t.Errorf("expected name %q, got %q", "gh-prod", client.Name())
	}
}

func TestNewClientFromConfig_WithBaseURL(t *testing.T) {
	cfg := config.SourceConfig{
		Name: "gh-enterprise",
		Type: "github",
		Options: map[string]string{
			"token":    "secret",
			"repo":     "org/repo",
			"base_url": "https://github.example.com/api/v3",
		},
	}

	client, err := github.NewClientFromConfig(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if client == nil {
		t.Fatal("expected non-nil client")
	}
}

func TestNewClientFromConfig_MissingToken(t *testing.T) {
	cfg := config.SourceConfig{
		Name: "gh-prod",
		Type: "github",
		Options: map[string]string{
			"repo": "org/repo",
		},
	}

	_, err := github.NewClientFromConfig(cfg)
	if err == nil {
		t.Fatal("expected error for missing token")
	}
}

func TestNewClientFromConfig_MissingRepo(t *testing.T) {
	cfg := config.SourceConfig{
		Name: "gh-prod",
		Type: "github",
		Options: map[string]string{
			"token": "secret",
		},
	}

	_, err := github.NewClientFromConfig(cfg)
	if err == nil {
		t.Fatal("expected error for missing repo")
	}
}

func TestNewClientFromConfig_MissingBoth(t *testing.T) {
	cfg := config.SourceConfig{
		Name:    "gh-prod",
		Type:    "github",
		Options: map[string]string{},
	}

	_, err := github.NewClientFromConfig(cfg)
	if err == nil {
		t.Fatal("expected error when both token and repo are missing")
	}
}
