package config_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/deploylog/deploylog/internal/config"
)

func writeConfig(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, "config.json")
	if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
		t.Fatalf("writing temp config: %v", err)
	}
	return p
}

func TestLoad_Valid(t *testing.T) {
	path := writeConfig(t, `{
		"sources": [{"name": "ci", "type": "github-actions"}],
		"output": {"format": "json", "indent": true},
		"poll_interval": "30s"
	}`)
	cfg, err := config.Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(cfg.Sources) != 1 {
		t.Errorf("expected 1 source, got %d", len(cfg.Sources))
	}
	if cfg.Sources[0].Name != "ci" {
		t.Errorf("expected source name 'ci', got %q", cfg.Sources[0].Name)
	}
	if cfg.Output.Format != "json" {
		t.Errorf("expected format 'json', got %q", cfg.Output.Format)
	}
	if cfg.PollInterval.Seconds() != 30 {
		t.Errorf("expected 30s poll interval, got %v", cfg.PollInterval)
	}
}

func TestLoad_DefaultOutputFormat(t *testing.T) {
	path := writeConfig(t, `{"sources": [{"name": "ci", "type": "jenkins"}], "output": {}}`)
	cfg, err := config.Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Output.Format != "text" {
		t.Errorf("expected default format 'text', got %q", cfg.Output.Format)
	}
}

func TestLoad_MissingSources(t *testing.T) {
	path := writeConfig(t, `{"sources": [], "output": {"format": "text"}}`)
	_, err := config.Load(path)
	if err == nil {
		t.Fatal("expected error for empty sources, got nil")
	}
}

func TestLoad_DuplicateSourceName(t *testing.T) {
	path := writeConfig(t, `{
		"sources": [
			{"name": "ci", "type": "github-actions"},
			{"name": "ci", "type": "jenkins"}
		],
		"output": {"format": "text"}
	}`)
	_, err := config.Load(path)
	if err == nil {
		t.Fatal("expected error for duplicate source name, got nil")
	}
}

func TestLoad_MissingFile(t *testing.T) {
	_, err := config.Load("/nonexistent/path/config.json")
	if err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
}

func TestLoad_InvalidDuration(t *testing.T) {
	path := writeConfig(t, `{
		"sources": [{"name": "ci", "type": "github-actions"}],
		"output": {"format": "text"},
		"poll_interval": "not-a-duration"
	}`)
	_, err := config.Load(path)
	if err == nil {
		t.Fatal("expected error for invalid duration, got nil")
	}
}
