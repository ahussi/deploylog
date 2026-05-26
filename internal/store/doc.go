// Package store provides a simple file-backed JSON store for deployment events.
//
// It supports saving, loading, and appending events to a persistent file,
// making it suitable for caching fetched events between deploylog runs.
//
// Usage:
//
//	s := store.New("/var/lib/deploylog/events.json")
//
//	// Save a set of events
//	if err := s.Save(events); err != nil { ... }
//
//	// Load previously stored events
//	events, err := s.Load()
//
//	// Append new events to existing store
//	if err := s.Append(newEvents); err != nil { ... }
package store
