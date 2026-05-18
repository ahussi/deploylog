package timeline_test

import (
	"testing"
	"time"

	"github.com/deploylog/deploylog/internal/event"
	"github.com/deploylog/deploylog/internal/timeline"
)

func makeEvent(id, source, status string, ts time.Time) event.Event {
	return event.Event{
		ID:        id,
		Source:    source,
		Service:   "api",
		Status:    status,
		Timestamp: ts,
	}
}

func TestTimeline_AddAndLen(t *testing.T) {
	tl := timeline.New()

	if tl.Len() != 0 {
		t.Fatalf("expected empty timeline, got %d events", tl.Len())
	}

	e := makeEvent("evt-1", "github-actions", "success", time.Now())
	if err := tl.Add(e); err != nil {
		t.Fatalf("unexpected error adding event: %v", err)
	}

	if tl.Len() != 1 {
		t.Fatalf("expected 1 event, got %d", tl.Len())
	}
}

func TestTimeline_Add_InvalidEvent(t *testing.T) {
	tl := timeline.New()

	invalid := event.Event{} // missing required fields
	if err := tl.Add(invalid); err == nil {
		t.Fatal("expected validation error, got nil")
	}

	if tl.Len() != 0 {
		t.Fatalf("expected timeline to remain empty, got %d events", tl.Len())
	}
}

func TestTimeline_Events_ChronologicalOrder(t *testing.T) {
	tl := timeline.New()
	now := time.Now()

	_ = tl.Add(makeEvent("evt-3", "circleci", "running", now.Add(2*time.Minute)))
	_ = tl.Add(makeEvent("evt-1", "github-actions", "success", now))
	_ = tl.Add(makeEvent("evt-2", "gitlab-ci", "failed", now.Add(time.Minute)))

	events := tl.Events()
	if len(events) != 3 {
		t.Fatalf("expected 3 events, got %d", len(events))
	}

	for i := 1; i < len(events); i++ {
		if events[i].Timestamp.Before(events[i-1].Timestamp) {
			t.Errorf("events not in chronological order at index %d", i)
		}
	}
}

func TestTimeline_FilterBySource(t *testing.T) {
	tl := timeline.New()
	now := time.Now()

	_ = tl.Add(makeEvent("evt-1", "github-actions", "success", now))
	_ = tl.Add(makeEvent("evt-2", "circleci", "failed", now.Add(time.Minute)))
	_ = tl.Add(makeEvent("evt-3", "github-actions", "running", now.Add(2*time.Minute)))

	ghaEvents := tl.FilterBySource("github-actions")
	if len(ghaEvents) != 2 {
		t.Fatalf("expected 2 github-actions events, got %d", len(ghaEvents))
	}

	for _, e := range ghaEvents {
		if e.Source != "github-actions" {
			t.Errorf("unexpected source %q in filtered results", e.Source)
		}
	}
}
