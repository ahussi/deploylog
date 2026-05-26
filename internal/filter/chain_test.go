package filter_test

import (
	"testing"
	"time"

	"github.com/yourorg/deploylog/internal/event"
	"github.com/yourorg/deploylog/internal/filter"
)

func chainEvent(source, status, service string) event.Event {
	return event.Event{
		ID:        "chain-test",
		Source:    source,
		Service:   service,
		Status:    event.Status(status),
		Timestamp: time.Now(),
	}
}

func TestChain_AllMatch(t *testing.T) {
	e := chainEvent("github", "success", "api")
	f := filter.Chain(filter.BySource("github"), filter.ByStatus("success"), filter.ByService("api"))
	if !f(e) {
		t.Error("expected Chain to match when all filters pass")
	}
}

func TestChain_OneFails(t *testing.T) {
	e := chainEvent("github", "failure", "api")
	f := filter.Chain(filter.BySource("github"), filter.ByStatus("success"))
	if f(e) {
		t.Error("expected Chain to reject when one filter fails")
	}
}

func TestChain_Empty(t *testing.T) {
	e := chainEvent("github", "success", "api")
	f := filter.Chain()
	if !f(e) {
		t.Error("expected empty Chain to match every event")
	}
}

func TestAnyOf_OneMatches(t *testing.T) {
	e := chainEvent("gitlab", "success", "api")
	f := filter.AnyOf(filter.BySource("github"), filter.BySource("gitlab"))
	if !f(e) {
		t.Error("expected AnyOf to match when at least one filter passes")
	}
}

func TestAnyOf_NoneMatch(t *testing.T) {
	e := chainEvent("jenkins", "success", "api")
	f := filter.AnyOf(filter.BySource("github"), filter.BySource("gitlab"))
	if f(e) {
		t.Error("expected AnyOf to reject when no filter passes")
	}
}

func TestAnyOf_Empty(t *testing.T) {
	e := chainEvent("github", "success", "api")
	f := filter.AnyOf()
	if f(e) {
		t.Error("expected empty AnyOf to reject every event")
	}
}

func TestNegate_InvertsTrue(t *testing.T) {
	e := chainEvent("github", "success", "api")
	f := filter.Negate(filter.BySource("github"))
	if f(e) {
		t.Error("expected Negate to invert a matching filter")
	}
}

func TestNegate_InvertsFalse(t *testing.T) {
	e := chainEvent("gitlab", "success", "api")
	f := filter.Negate(filter.BySource("github"))
	if !f(e) {
		t.Error("expected Negate to invert a non-matching filter")
	}
}
