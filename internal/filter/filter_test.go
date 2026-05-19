package filter_test

import (
	"testing"
	"time"

	"github.com/yourorg/deploylog/internal/event"
	"github.com/yourorg/deploylog/internal/filter"
)

var baseTime = time.Date(2024, 1, 15, 12, 0, 0, 0, time.UTC)

func makeEvent(source, service string, status event.Status, ts time.Time) event.Event {
	return event.Event{
		ID:        "evt-1",
		Source:    source,
		Service:   service,
		Status:    status,
		Timestamp: ts,
	}
}

func TestBySource(t *testing.T) {
	events := []event.Event{
		makeEvent("github", "api", event.StatusSuccess, baseTime),
		makeEvent("gitlab", "api", event.StatusSuccess, baseTime),
		makeEvent("github", "worker", event.StatusFailed, baseTime),
	}
	got := filter.Apply(events, filter.BySource("github"))
	if len(got) != 2 {
		t.Fatalf("expected 2 events, got %d", len(got))
	}
}

func TestByStatus(t *testing.T) {
	events := []event.Event{
		makeEvent("github", "api", event.StatusSuccess, baseTime),
		makeEvent("github", "api", event.StatusFailed, baseTime),
		makeEvent("github", "api", event.StatusRunning, baseTime),
	}
	got := filter.Apply(events, filter.ByStatus(event.StatusFailed))
	if len(got) != 1 {
		t.Fatalf("expected 1 event, got %d", len(got))
	}
}

func TestByTimeRange(t *testing.T) {
	events := []event.Event{
		makeEvent("github", "api", event.StatusSuccess, baseTime.Add(-2*time.Hour)),
		makeEvent("github", "api", event.StatusSuccess, baseTime),
		makeEvent("github", "api", event.StatusSuccess, baseTime.Add(2*time.Hour)),
	}
	from := baseTime.Add(-1 * time.Hour)
	to := baseTime.Add(1 * time.Hour)
	got := filter.Apply(events, filter.ByTimeRange(from, to))
	if len(got) != 1 {
		t.Fatalf("expected 1 event in range, got %d", len(got))
	}
}

func TestByService(t *testing.T) {
	events := []event.Event{
		makeEvent("github", "api", event.StatusSuccess, baseTime),
		makeEvent("github", "worker", event.StatusSuccess, baseTime),
	}
	got := filter.Apply(events, filter.ByService("worker"))
	if len(got) != 1 || got[0].Service != "worker" {
		t.Fatalf("expected 1 worker event, got %d", len(got))
	}
}

func TestAll(t *testing.T) {
	events := []event.Event{
		makeEvent("github", "api", event.StatusSuccess, baseTime),
		makeEvent("github", "api", event.StatusFailed, baseTime),
		makeEvent("gitlab", "api", event.StatusSuccess, baseTime),
	}
	pred := filter.All(filter.BySource("github"), filter.ByStatus(event.StatusSuccess))
	got := filter.Apply(events, pred)
	if len(got) != 1 {
		t.Fatalf("expected 1 event matching all predicates, got %d", len(got))
	}
}

func TestAny(t *testing.T) {
	events := []event.Event{
		makeEvent("github", "api", event.StatusSuccess, baseTime),
		makeEvent("gitlab", "api", event.StatusFailed, baseTime),
		makeEvent("jenkins", "api", event.StatusRunning, baseTime),
	}
	pred := filter.Any(filter.BySource("github"), filter.BySource("gitlab"))
	got := filter.Apply(events, pred)
	if len(got) != 2 {
		t.Fatalf("expected 2 events matching any predicate, got %d", len(got))
	}
}
