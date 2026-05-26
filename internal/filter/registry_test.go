package filter_test

import (
	"testing"
	"time"

	"github.com/yourorg/deploylog/internal/event"
	"github.com/yourorg/deploylog/internal/filter"
)

func makeRegistryEvent(source, status string) event.Event {
	return event.Event{
		ID:        "test-id",
		Source:    source,
		Service:   "svc",
		Status:    event.Status(status),
		Timestamp: time.Now(),
	}
}

func TestRegistry_RegisterAndGet(t *testing.T) {
	reg := filter.NewRegistry()

	fn := func(e event.Event) bool { return e.Source == "github" }
	if err := reg.Register("by-github", fn); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got, err := reg.Get("by-github")
	if err != nil {
		t.Fatalf("expected to find filter: %v", err)
	}
	if got == nil {
		t.Fatal("expected non-nil filter func")
	}
}

func TestRegistry_DuplicateRegister(t *testing.T) {
	reg := filter.NewRegistry()
	fn := func(e event.Event) bool { return true }

	if err := reg.Register("dup", fn); err != nil {
		t.Fatalf("first register failed: %v", err)
	}
	if err := reg.Register("dup", fn); err == nil {
		t.Fatal("expected error on duplicate registration")
	}
}

func TestRegistry_GetUnknown(t *testing.T) {
	reg := filter.NewRegistry()
	if _, err := reg.Get("nonexistent"); err == nil {
		t.Fatal("expected error for unknown filter")
	}
}

func TestRegistry_Names(t *testing.T) {
	reg := filter.NewRegistry()
	fn := func(e event.Event) bool { return true }

	_ = reg.Register("alpha", fn)
	_ = reg.Register("beta", fn)

	names := reg.Names()
	if len(names) != 2 {
		t.Fatalf("expected 2 names, got %d", len(names))
	}
}

func TestRegistry_Apply_AllPass(t *testing.T) {
	reg := filter.NewRegistry()
	_ = reg.Register("github-only", func(e event.Event) bool { return e.Source == "github" })

	events := []event.Event{
		makeRegistryEvent("github", "success"),
		makeRegistryEvent("gitlab", "success"),
		makeRegistryEvent("github", "failed"),
	}

	result := reg.Apply(events)
	if len(result) != 2 {
		t.Fatalf("expected 2 events, got %d", len(result))
	}
	for _, e := range result {
		if e.Source != "github" {
			t.Errorf("expected source github, got %s", e.Source)
		}
	}
}

func TestRegistry_Apply_MultipleFilters(t *testing.T) {
	reg := filter.NewRegistry()
	_ = reg.Register("github-only", func(e event.Event) bool { return e.Source == "github" })
	_ = reg.Register("success-only", func(e event.Event) bool { return e.Status == event.StatusSuccess })

	events := []event.Event{
		makeRegistryEvent("github", "success"),
		makeRegistryEvent("github", "failed"),
		makeRegistryEvent("gitlab", "success"),
	}

	result := reg.Apply(events)
	if len(result) != 1 {
		t.Fatalf("expected 1 event, got %d", len(result))
	}
	if result[0].Source != "github" || result[0].Status != event.StatusSuccess {
		t.Errorf("unexpected event: %+v", result[0])
	}
}
