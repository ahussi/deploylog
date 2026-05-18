// Package timeline provides functionality for aggregating and ordering
// deployment events from multiple sources into a unified audit timeline.
package timeline

import (
	"sort"
	"sync"

	"github.com/deploylog/deploylog/internal/event"
)

// Timeline holds an ordered collection of deployment events.
type Timeline struct {
	mu     sync.RWMutex
	events []event.Event
}

// New creates and returns an empty Timeline.
func New() *Timeline {
	return &Timeline{
		events: make([]event.Event, 0),
	}
}

// Add appends a validated event to the timeline. Returns an error if the
// event fails validation.
func (t *Timeline) Add(e event.Event) error {
	if err := e.Validate(); err != nil {
		return err
	}

	t.mu.Lock()
	defer t.mu.Unlock()

	t.events = append(t.events, e)
	return nil
}

// Events returns a chronologically sorted copy of all events in the timeline.
func (t *Timeline) Events() []event.Event {
	t.mu.RLock()
	defer t.mu.RUnlock()

	copy := make([]event.Event, len(t.events))
	for i, e := range t.events {
		copy[i] = e
	}

	sort.Slice(copy, func(i, j int) bool {
		return copy[i].Timestamp.Before(copy[j].Timestamp)
	})

	return copy
}

// FilterBySource returns events emitted by the given source, in chronological order.
func (t *Timeline) FilterBySource(source string) []event.Event {
	all := t.Events()
	result := make([]event.Event, 0)
	for _, e := range all {
		if e.Source == source {
			result = append(result, e)
		}
	}
	return result
}

// Len returns the number of events currently in the timeline.
func (t *Timeline) Len() int {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return len(t.events)
}
