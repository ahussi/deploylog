package source

import (
	"context"
	"fmt"

	"github.com/deploylog/internal/event"
)

// Kind identifies the type of CI/CD source.
type Kind string

const (
	KindGitHub   Kind = "github"
	KindGitLab   Kind = "gitlab"
	KindCircleCI Kind = "circleci"
	KindUnknown  Kind = "unknown"
)

// ErrUnsupportedSource is returned when a source kind is not supported.
type ErrUnsupportedSource struct {
	Kind Kind
}

func (e *ErrUnsupportedSource) Error() string {
	return fmt.Sprintf("unsupported source kind: %q", e.Kind)
}

// Collector defines the interface for pulling deployment events from a CI/CD source.
type Collector interface {
	// Kind returns the source identifier.
	Kind() Kind

	// Collect fetches deployment events and sends them to the provided channel.
	// Implementations must close the channel when done or on context cancellation.
	Collect(ctx context.Context, out chan<- event.Event) error
}

// Registry holds registered Collector implementations keyed by Kind.
type Registry struct {
	collectors map[Kind]Collector
}

// NewRegistry creates an empty Registry.
func NewRegistry() *Registry {
	return &Registry{collectors: make(map[Kind]Collector)}
}

// Register adds a Collector to the registry.
// Returns an error if a collector for that Kind is already registered.
func (r *Registry) Register(c Collector) error {
	if _, exists := r.collectors[c.Kind()]; exists {
		return fmt.Errorf("collector already registered for kind %q", c.Kind())
	}
	r.collectors[c.Kind()] = c
	return nil
}

// Get retrieves a Collector by Kind.
func (r *Registry) Get(k Kind) (Collector, error) {
	c, ok := r.collectors[k]
	if !ok {
		return nil, &ErrUnsupportedSource{Kind: k}
	}
	return c, nil
}

// All returns every registered Collector.
func (r *Registry) All() []Collector {
	out := make([]Collector, 0, len(r.collectors))
	for _, c := range r.collectors {
		out = append(out, c)
	}
	return out
}
