package filter

import "github.com/yourorg/deploylog/internal/event"

// Chain combines multiple FilterFuncs into a single FilterFunc that returns
// true only when all provided filters match (logical AND).
//
// An empty chain always returns true, matching every event.
func Chain(filters ...FilterFunc) FilterFunc {
	return func(e event.Event) bool {
		for _, f := range filters {
			if !f(e) {
				return false
			}
		}
		return true
	}
}

// AnyOf combines multiple FilterFuncs into a single FilterFunc that returns
// true when at least one of the provided filters matches (logical OR).
//
// An empty AnyOf always returns false.
func AnyOf(filters ...FilterFunc) FilterFunc {
	return func(e event.Event) bool {
		for _, f := range filters {
			if f(e) {
				return true
			}
		}
		return false
	}
}

// Negate inverts the result of the given FilterFunc.
func Negate(f FilterFunc) FilterFunc {
	return func(e event.Event) bool {
		return !f(e)
	}
}
