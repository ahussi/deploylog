package source_test

import (
	"context"
	"testing"

	"github.com/deploylog/internal/event"
	"github.com/deploylog/internal/source"
)

// stubCollector is a minimal Collector for testing.
type stubCollector struct {
	kind source.Kind
}

func (s *stubCollector) Kind() source.Kind { return s.kind }
func (s *stubCollector) Collect(_ context.Context, out chan<- event.Event) error {
	close(out)
	return nil
}

func TestRegistry_RegisterAndGet(t *testing.T) {
	reg := source.NewRegistry()
	col := &stubCollector{kind: source.KindGitHub}

	if err := reg.Register(col); err != nil {
		t.Fatalf("unexpected error registering collector: %v", err)
	}

	got, err := reg.Get(source.KindGitHub)
	if err != nil {
		t.Fatalf("expected collector, got error: %v", err)
	}
	if got.Kind() != source.KindGitHub {
		t.Errorf("expected kind %q, got %q", source.KindGitHub, got.Kind())
	}
}

func TestRegistry_DuplicateRegister(t *testing.T) {
	reg := source.NewRegistry()
	col := &stubCollector{kind: source.KindGitLab}

	_ = reg.Register(col)
	if err := reg.Register(col); err == nil {
		t.Error("expected error on duplicate registration, got nil")
	}
}

func TestRegistry_GetUnknown(t *testing.T) {
	reg := source.NewRegistry()

	_, err := reg.Get(source.KindCircleCI)
	if err == nil {
		t.Fatal("expected ErrUnsupportedSource, got nil")
	}

	unsupErr, ok := err.(*source.ErrUnsupportedSource)
	if !ok {
		t.Fatalf("expected *ErrUnsupportedSource, got %T", err)
	}
	if unsupErr.Kind != source.KindCircleCI {
		t.Errorf("expected kind %q in error, got %q", source.KindCircleCI, unsupErr.Kind)
	}
}

func TestRegistry_All(t *testing.T) {
	reg := source.NewRegistry()
	_ = reg.Register(&stubCollector{kind: source.KindGitHub})
	_ = reg.Register(&stubCollector{kind: source.KindGitLab})

	all := reg.All()
	if len(all) != 2 {
		t.Errorf("expected 2 collectors, got %d", len(all))
	}
}
