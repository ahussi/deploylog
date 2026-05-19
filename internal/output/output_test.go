package output_test

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/yourorg/deploylog/internal/event"
	"github.com/yourorg/deploylog/internal/output"
)

var testTime = time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC)

func makeEvents() []event.Event {
	return []event.Event{
		{
			ID:          "evt-1",
			Source:      "github-actions",
			Status:      event.StatusSuccess,
			Description: "deploy to production",
			Timestamp:   testTime,
		},
		{
			ID:          "evt-2",
			Source:      "argocd",
			Status:      event.StatusFailed,
			Description: "sync failed",
			Timestamp:   testTime.Add(time.Minute),
		},
	}
}

func TestJSONFormatter_Format(t *testing.T) {
	events := makeEvents()
	var buf bytes.Buffer
	f := &output.JSONFormatter{}
	if err := f.Format(&buf, events); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var got []event.Event
	if err := json.Unmarshal(buf.Bytes(), &got); err != nil {
		t.Fatalf("invalid JSON output: %v", err)
	}
	if len(got) != len(events) {
		t.Errorf("expected %d events, got %d", len(events), len(got))
	}
}

func TestJSONFormatter_Indent(t *testing.T) {
	var buf bytes.Buffer
	f := &output.JSONFormatter{Indent: true}
	if err := f.Format(&buf, makeEvents()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "\n  ") {
		t.Error("expected indented JSON output")
	}
}

func TestTextFormatter_Format(t *testing.T) {
	events := makeEvents()
	var buf bytes.Buffer
	f := &output.TextFormatter{}
	if err := f.Format(&buf, events); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	if len(lines) != len(events) {
		t.Errorf("expected %d lines, got %d", len(events), len(lines))
	}
	if !strings.Contains(lines[0], "github-actions") {
		t.Errorf("expected source in output, got: %s", lines[0])
	}
	if !strings.Contains(lines[1], "sync failed") {
		t.Errorf("expected description in output, got: %s", lines[1])
	}
}

func TestTextFormatter_CustomTimeFormat(t *testing.T) {
	var buf bytes.Buffer
	f := &output.TextFormatter{TimeFormat: "2006-01-02"}
	if err := f.Format(&buf, makeEvents()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "2024-06-01") {
		t.Errorf("expected custom date format in output")
	}
}
