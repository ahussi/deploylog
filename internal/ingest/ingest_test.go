package ingest_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/deploylog/deploylog/internal/event"
	"github.com/deploylog/deploylog/internal/ingest"
	"github.com/deploylog/deploylog/internal/source"
	"github.com/deploylog/deploylog/internal/timeline"
)

type stubFetcher struct {
	events []event.Event
	err    error
}

func (s *stubFetcher) Fetch(_ context.Context) ([]event.Event, error) {
	return s.events, s.err
}

func makeTestEvent(src, id string) event.Event {
	return event.Event{
		ID:        id,
		Source:    src,
		Status:    event.StatusSuccess,
		Timestamp: time.Now(),
		Service:   "api",
	}
}

func TestProcessor_RegisterUnknownSource(t *testing.T) {
	reg := source.NewRegistry()
	tl := timeline.New()
	p := ingest.NewProcessor(reg, tl)

	if err := p.Register("unknown", &stubFetcher{}); err == nil {
		t.Fatal("expected error for unregistered source, got nil")
	}
}

func TestProcessor_Run_AddsEvents(t *testing.T) {
	reg := source.NewRegistry()
	_ = reg.Register("github-actions", source.Meta{DisplayName: "GitHub Actions"})
	tl := timeline.New()
	p := ingest.NewProcessor(reg, tl)

	events := []event.Event{
		makeTestEvent("github-actions", "evt-1"),
		makeTestEvent("github-actions", "evt-2"),
	}
	_ = p.Register("github-actions", &stubFetcher{events: events})

	if err := p.Run(context.Background()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := tl.Len(); got != 2 {
		t.Errorf("expected 2 events in timeline, got %d", got)
	}
}

func TestProcessor_Run_FetchError_ContinuesOtherSources(t *testing.T) {
	reg := source.NewRegistry()
	_ = reg.Register("gitlab-ci", source.Meta{DisplayName: "GitLab CI"})
	_ = reg.Register("circleci", source.Meta{DisplayName: "CircleCI"})
	tl := timeline.New()
	p := ingest.NewProcessor(reg, tl)

	_ = p.Register("gitlab-ci", &stubFetcher{err: errors.New("connection refused")})
	_ = p.Register("circleci", &stubFetcher{events: []event.Event{makeTestEvent("circleci", "evt-3")}})

	if err := p.Run(context.Background()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := tl.Len(); got != 1 {
		t.Errorf("expected 1 event in timeline, got %d", got)
	}
}
