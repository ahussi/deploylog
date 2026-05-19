package output

import (
	"fmt"
	"sort"
	"strings"
)

// FormatterFactory is a constructor function for a Formatter.
type FormatterFactory func() Formatter

// FormatterRegistry maps format names to their factory functions.
type FormatterRegistry struct {
	factories map[string]FormatterFactory
}

// NewFormatterRegistry returns a FormatterRegistry pre-populated with the
// built-in formatters ("json" and "text").
func NewFormatterRegistry() *FormatterRegistry {
	r := &FormatterRegistry{factories: make(map[string]FormatterFactory)}
	r.Register("json", func() Formatter { return &JSONFormatter{Indent: true} })
	r.Register("text", func() Formatter { return &TextFormatter{} })
	return r
}

// Register adds a named formatter factory to the registry.
// It returns an error if the name is already taken.
func (r *FormatterRegistry) Register(name string, factory FormatterFactory) error {
	name = strings.ToLower(strings.TrimSpace(name))
	if _, exists := r.factories[name]; exists {
		return fmt.Errorf("output: formatter %q already registered", name)
	}
	r.factories[name] = factory
	return nil
}

// Get returns a new Formatter for the given name, or an error if unknown.
func (r *FormatterRegistry) Get(name string) (Formatter, error) {
	name = strings.ToLower(strings.TrimSpace(name))
	factory, ok := r.factories[name]
	if !ok {
		return nil, fmt.Errorf("output: unknown formatter %q", name)
	}
	return factory(), nil
}

// Names returns the sorted list of registered formatter names.
func (r *FormatterRegistry) Names() []string {
	names := make([]string, 0, len(r.factories))
	for n := range r.factories {
		names = append(names, n)
	}
	sort.Strings(names)
	return names
}
