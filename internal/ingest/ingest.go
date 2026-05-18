// Package ingest provides functionality for ingesting deployment events
// from multiple CI/CD sources into the unified timeline.
package ingest

import (
	"context"
	"fmt"
	"log"

	"github.com/deploylog/deploylog/internal/event"
	"github.com/deploylog/deploylog/internal/source"
	"github.com/deploylog/deploylog/internal/timeline"
)

// Fetcher is implemented by any CI/CD source adapter that can retrieve
// deployment events.
type Fetcher interface {
	Fetch(ctx context.Context) ([]event.Event, error)
}

// Processor ingests events from registered fetchers into a timeline.
type Processor struct {
	registry *source.Registry
	timeline *timeline.Timeline
	fetchers map[string]Fetcher
}

// NewProcessor creates a new Processor backed by the given registry and timeline.
func NewProcessor(reg *source.Registry, tl *timeline.Timeline) *Processor {
	return &Processor{
		registry: reg,
		timeline: tl,
		fetchers: make(map[string]Fetcher),
	}
}

// Register associates a Fetcher with a named source.
// The source name must already be registered in the Registry.
func (p *Processor) Register(name string, f Fetcher) error {
	if _, err := p.registry.Get(name); err != nil {
		return fmt.Errorf("ingest: unknown source %q: %w", name, err)
	}
	p.fetchers[name] = f
	return nil
}

// Run fetches events from all registered fetchers and adds them to the
// timeline. Errors from individual fetchers are logged but do not abort
// processing of other sources.
func (p *Processor) Run(ctx context.Context) error {
	for name, fetcher := range p.fetchers {
		events, err := fetcher.Fetch(ctx)
		if err != nil {
			log.Printf("ingest: fetch error from source %q: %v", name, err)
			continue
		}
		for _, ev := range events {
			if addErr := p.timeline.Add(ev); addErr != nil {
				log.Printf("ingest: skipping invalid event from %q: %v", name, addErr)
			}
		}
	}
	return nil
}
