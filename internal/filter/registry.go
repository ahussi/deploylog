package filter

import (
	"fmt"
	"sync"

	"github.com/yourorg/deploylog/internal/event"
)

// FilterFunc is a function that determines whether an event passes a filter.
type FilterFunc func(e event.Event) bool

// Registry holds named filter functions that can be composed and applied.
type Registry struct {
	mu      sync.RWMutex
	filters map[string]FilterFunc
}

// NewRegistry creates and returns an empty filter Registry.
func NewRegistry() *Registry {
	return &Registry{
		filters: make(map[string]FilterFunc),
	}
}

// Register adds a named FilterFunc to the registry.
// Returns an error if the name is already registered.
func (r *Registry) Register(name string, fn FilterFunc) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.filters[name]; exists {
		return fmt.Errorf("filter %q is already registered", name)
	}
	r.filters[name] = fn
	return nil
}

// Get retrieves a named FilterFunc from the registry.
// Returns an error if the name is not found.
func (r *Registry) Get(name string) (FilterFunc, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	fn, ok := r.filters[name]
	if !ok {
		return nil, fmt.Errorf("filter %q not found", name)
	}
	return fn, nil
}

// Names returns a sorted list of all registered filter names.
func (r *Registry) Names() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()

	names := make([]string, 0, len(r.filters))
	for name := range r.filters {
		names = append(names, name)
	}
	return names
}

// Apply runs all registered filters against the given events and returns
// only those that pass every registered filter.
func (r *Registry) Apply(events []event.Event) []event.Event {
	r.mu.RLock()
	defer r.mu.RUnlock()

	result := make([]event.Event, 0, len(events))
	for _, e := range events {
		pass := true
		for _, fn := range r.filters {
			if !fn(e) {
				pass = false
				break
			}
		}
		if pass {
			result = append(result, e)
		}
	}
	return result
}
