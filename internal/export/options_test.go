package export_test

import (
	"fmt"
	"testing"

	"github.com/yourorg/deploylog/internal/export"
	"github.com/yourorg/deploylog/internal/output"
)

func TestNewWithOptions_Stdout(t *testing.T) {
	ex, err := export.NewWithOptions(
		export.WithStdout(),
		export.WithFormatter(&output.JSONFormatter{}),
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ex == nil {
		t.Fatal("expected non-nil exporter")
	}
}

func TestNewWithOptions_File(t *testing.T) {
	ex, err := export.NewWithOptions(
		export.WithFile("/tmp/test-deploylog.json"),
		export.WithFormatter(&output.JSONFormatter{}),
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ex == nil {
		t.Fatal("expected non-nil exporter")
	}
}

func TestNewWithOptions_MissingFormatter(t *testing.T) {
	_, err := export.NewWithOptions(export.WithStdout())
	if err == nil {
		t.Fatal("expected error when formatter is nil")
	}
}

func TestNewWithOptions_FileMissingPath(t *testing.T) {
	_, err := export.NewWithOptions(
		export.WithFormatter(&output.TextFormatter{}),
		// WithFile not called, destination defaults to stdout so this should succeed
	)
	if err != nil {
		t.Fatalf("unexpected error for stdout destination: %v", err)
	}
}

func TestNewWithOptions_CombinedApply(t *testing.T) {
	path := fmt.Sprintf("/tmp/deploylog-combined-%d.json", 42)
	ex, err := export.NewWithOptions(
		export.WithStdout(),
		export.WithFile(path),
		export.WithFormatter(&output.JSONFormatter{Indent: true}),
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ex == nil {
		t.Fatal("expected non-nil exporter")
	}
}
