// Package filter provides predicates for selecting deployment events
// from a timeline based on various criteria.
package filter

import (
	"time"

	"github.com/yourorg/deploylog/internal/event"
)

// Predicate is a function that returns true if the event should be included.
type Predicate func(e event.Event) bool

// BySource returns a Predicate that matches events from the given source.
func BySource(source string) Predicate {
	return func(e event.Event) bool {
		return e.Source == source
	}
}

// ByStatus returns a Predicate that matches events with the given status.
func ByStatus(status event.Status) Predicate {
	return func(e event.Event) bool {
		return e.Status == status
	}
}

// ByTimeRange returns a Predicate that matches events whose timestamp
// falls within [from, to] inclusive. A zero value for either bound is
// treated as unbounded.
func ByTimeRange(from, to time.Time) Predicate {
	return func(e event.Event) bool {
		if !from.IsZero() && e.Timestamp.Before(from) {
			return false
		}
		if !to.IsZero() && e.Timestamp.After(to) {
			return false
		}
		return true
	}
}

// ByService returns a Predicate that matches events for the given service.
func ByService(service string) Predicate {
	return func(e event.Event) bool {
		return e.Service == service
	}
}

// All combines multiple predicates with logical AND; an event must satisfy
// every predicate to be included.
func All(predicates ...Predicate) Predicate {
	return func(e event.Event) bool {
		for _, p := range predicates {
			if !p(e) {
				return false
			}
		}
		return true
	}
}

// Any combines multiple predicates with logical OR; an event must satisfy
// at least one predicate to be included.
func Any(predicates ...Predicate) Predicate {
	return func(e event.Event) bool {
		for _, p := range predicates {
			if p(e) {
				return true
			}
		}
		return false
	}
}

// Apply filters a slice of events, returning only those that satisfy pred.
func Apply(events []event.Event, pred Predicate) []event.Event {
	out := make([]event.Event, 0, len(events))
	for _, e := range events {
		if pred(e) {
			out = append(out, e)
		}
	}
	return out
}
