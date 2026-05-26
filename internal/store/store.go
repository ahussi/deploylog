// Package store provides persistent storage for deployment events.
package store

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"

	"github.com/deploylog/deploylog/internal/event"
)

// Store persists and retrieves deployment events from a JSON file.
type Store struct {
	mu   sync.RWMutex
	path string
}

// New creates a new Store backed by the file at path.
func New(path string) *Store {
	return &Store{path: path}
}

// Save writes the provided events to the store, overwriting any existing data.
func (s *Store) Save(events []event.Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	data, err := json.MarshalIndent(events, "", "  ")
	if err != nil {
		return fmt.Errorf("store: marshal events: %w", err)
	}

	if err := os.WriteFile(s.path, data, 0o644); err != nil {
		return fmt.Errorf("store: write file: %w", err)
	}

	return nil
}

// Load reads and returns all events from the store.
// Returns an empty slice if the file does not exist.
func (s *Store) Load() ([]event.Event, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	data, err := os.ReadFile(s.path)
	if os.IsNotExist(err) {
		return []event.Event{}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("store: read file: %w", err)
	}

	var events []event.Event
	if err := json.Unmarshal(data, &events); err != nil {
		return nil, fmt.Errorf("store: unmarshal events: %w", err)
	}

	return events, nil
}

// Append loads existing events, appends the new ones, and saves.
func (s *Store) Append(events []event.Event) error {
	existing, err := s.Load()
	if err != nil {
		return fmt.Errorf("store: append load: %w", err)
	}
	return s.Save(append(existing, events...))
}
