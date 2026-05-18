// Package source defines the Collector interface and Registry for integrating
// multiple CI/CD providers (GitHub Actions, GitLab CI, CircleCI, etc.) into
// deploylog's unified audit timeline.
//
// # Implementing a new source
//
// Create a struct that satisfies the Collector interface:
//
//	type MyCollector struct{}
//
//	func (m *MyCollector) Kind() source.Kind { return "myprovider" }
//
//	func (m *MyCollector) Collect(ctx context.Context, out chan<- event.Event) error {
//		defer close(out)
//		// fetch and emit events ...
//		return nil
//	}
//
// Then register it with a Registry before starting collection:
//
//	reg := source.NewRegistry()
//	_ = reg.Register(&MyCollector{})
package source
