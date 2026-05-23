package jenkins

import (
	"testing"

	"github.com/user/deploylog/internal/config"
)

func TestNewClientFromConfig_Valid(t *testing.T) {
	cfg := config.SourceConfig{
		Name: "jenkins-prod",
		Type: "jenkins",
		Options: map[string]string{
			"base_url": "http://jenkins.example.com",
			"job":      "deploy-app",
			"username": "admin",
			"password": "secret",
		},
	}
	client, err := NewClientFromConfig(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if client == nil {
		t.Fatal("expected non-nil client")
	}
	if client.Name() != SourceName {
		t.Errorf("expected source name %q, got %q", SourceName, client.Name())
	}
}

func TestNewClientFromConfig_WithoutCredentials(t *testing.T) {
	cfg := config.SourceConfig{
		Name: "jenkins-ci",
		Type: "jenkins",
		Options: map[string]string{
			"base_url": "http://jenkins.example.com",
			"job":      "build-service",
		},
	}
	client, err := NewClientFromConfig(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if client == nil {
		t.Fatal("expected non-nil client")
	}
}

func TestNewClientFromConfig_MissingBaseURL(t *testing.T) {
	cfg := config.SourceConfig{
		Name: "jenkins-ci",
		Type: "jenkins",
		Options: map[string]string{
			"job": "deploy-app",
		},
	}
	_, err := NewClientFromConfig(cfg)
	if err == nil {
		t.Fatal("expected error for missing base_url")
	}
}

func TestNewClientFromConfig_MissingJob(t *testing.T) {
	cfg := config.SourceConfig{
		Name: "jenkins-ci",
		Type: "jenkins",
		Options: map[string]string{
			"base_url": "http://jenkins.example.com",
		},
	}
	_, err := NewClientFromConfig(cfg)
	if err == nil {
		t.Fatal("expected error for missing job")
	}
}

func TestNewClientFromConfig_MissingBoth(t *testing.T) {
	cfg := config.SourceConfig{
		Name: "jenkins-ci",
		Type: "jenkins",
		Options: map[string]string{},
	}
	_, err := NewClientFromConfig(cfg)
	if err == nil {
		t.Fatal("expected error for missing base_url and job")
	}
}
