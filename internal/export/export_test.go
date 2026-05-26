package export_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/yourorg/deploylog/internal/event"
	"github.com/yourorg/deploylog/internal/export"
	"github.com/yourorg/deploylog/internal/output"
	"github.com/yourorg/deploylog/internal/timeline"
)

func makeTimeline(t *testing.T, n int) *timeline.Timeline {
	t.Helper()
	tl := timeline.New()
	for i := 0; i < n; i++ {
		e := event.Event{
			ID:        fmt.Sprintf("evt-%d", i),
			Source:    "github",
			Service:   "api",
			Status:    event.StatusSuccess,
			Timestamp: time.Now().Add(time.Duration(i) * time.Second),
		}
		_ = tl.Add(e)
	}
	return tl
}

func TestNew_NilFormatter(t *testing.T) {
	_, err := export.New(export.Options{Destination: export.DestinationStdout})
	if err == nil {
		t.Fatal("expected error for nil formatter")
	}
}

func TestNew_FileDestinationMissingPath(t *testing.T) {
	_, err := export.New(export.Options{
		Destination: export.DestinationFile,
		Formatter:   &output.JSONFormatter{},
	})
	if err == nil {
		t.Fatal("expected error for missing file path")
	}
}

func TestWrite_ToFile_JSON(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "out.json")

	ex, err := export.New(export.Options{
		Destination: export.DestinationFile,
		FilePath:    path,
		Formatter:   &output.JSONFormatter{},
	})
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	tl := timeline.New()
	ev := event.Event{
		ID: "e1", Source: "github", Service: "svc",
		Status: event.StatusSuccess, Timestamp: time.Now(),
	}
	_ = tl.Add(ev)

	if err := ex.Write(tl); err != nil {
		t.Fatalf("Write: %v", err)
	}

	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile: %v", err)
	}
	var events []map[string]interface{}
	if err := json.Unmarshal(raw, &events); err != nil {
		t.Fatalf("invalid JSON output: %v", err)
	}
	if len(events) != 1 {
		t.Errorf("expected 1 event, got %d", len(events))
	}
}

func TestWrite_InvalidFile(t *testing.T) {
	ex, _ := export.New(export.Options{
		Destination: export.DestinationFile,
		FilePath:    "/nonexistent/dir/out.json",
		Formatter:   &output.JSONFormatter{},
	})
	if err := ex.Write(timeline.New()); err == nil {
		t.Fatal("expected error writing to invalid path")
	}
}
