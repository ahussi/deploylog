// Package export provides functionality for writing formatted deployment
// timelines to various output destinations such as files or stdout.
package export

import (
	"fmt"
	"io"
	"os"

	"github.com/yourorg/deploylog/internal/output"
	"github.com/yourorg/deploylog/internal/timeline"
)

// Destination represents a write target for exported timeline data.
type Destination int

const (
	// DestinationStdout writes output to os.Stdout.
	DestinationStdout Destination = iota
	// DestinationFile writes output to a file path.
	DestinationFile
)

// Options configures an export operation.
type Options struct {
	Destination Destination
	FilePath    string
	Formatter   output.Formatter
}

// Exporter writes a timeline using a configured formatter and destination.
type Exporter struct {
	opts Options
}

// New creates an Exporter with the provided options.
func New(opts Options) (*Exporter, error) {
	if opts.Formatter == nil {
		return nil, fmt.Errorf("export: formatter must not be nil")
	}
	if opts.Destination == DestinationFile && opts.FilePath == "" {
		return nil, fmt.Errorf("export: file path must be set when destination is file")
	}
	return &Exporter{opts: opts}, nil
}

// Write formats all events from tl and writes them to the configured destination.
func (e *Exporter) Write(tl *timeline.Timeline) error {
	data, err := e.opts.Formatter.Format(tl.Events())
	if err != nil {
		return fmt.Errorf("export: format error: %w", err)
	}

	writer, closer, err := e.openWriter()
	if err != nil {
		return err
	}
	if closer != nil {
		defer closer()
	}

	if _, err := writer.Write(data); err != nil {
		return fmt.Errorf("export: write error: %w", err)
	}
	return nil
}

func (e *Exporter) openWriter() (io.Writer, func(), error) {
	if e.opts.Destination == DestinationStdout {
		return os.Stdout, nil, nil
	}
	f, err := os.Create(e.opts.FilePath)
	if err != nil {
		return nil, nil, fmt.Errorf("export: cannot open file %q: %w", e.opts.FilePath, err)
	}
	return f, func() { f.Close() }, nil
}
