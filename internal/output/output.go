// Package output provides formatters for rendering the unified audit timeline
// to various output targets such as JSON, plain text, or structured logs.
package output

import (
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/yourorg/deploylog/internal/event"
)

// Formatter defines the interface for rendering deployment events.
type Formatter interface {
	Format(w io.Writer, events []event.Event) error
}

// JSONFormatter renders events as a JSON array.
type JSONFormatter struct {
	Indent bool
}

// Format writes events as JSON to w.
func (f *JSONFormatter) Format(w io.Writer, events []event.Event) error {
	var data []byte
	var err error
	if f.Indent {
		data, err = json.MarshalIndent(events, "", "  ")
	} else {
		data, err = json.Marshal(events)
	}
	if err != nil {
		return fmt.Errorf("output: json marshal: %w", err)
	}
	_, err = w.Write(data)
	return err
}

// TextFormatter renders events as human-readable lines.
type TextFormatter struct {
	// TimeFormat controls how timestamps are displayed. Defaults to RFC3339.
	TimeFormat string
}

// Format writes events as plain text lines to w.
func (f *TextFormatter) Format(w io.Writer, events []event.Event) error {
	tf := f.TimeFormat
	if tf == "" {
		tf = time.RFC3339
	}
	for _, e := range events {
		_, err := fmt.Fprintf(w, "[%s] source=%-12s status=%-10s %s\n",
			e.Timestamp.Format(tf),
			e.Source,
			e.Status,
			e.Description,
		)
		if err != nil {
			return fmt.Errorf("output: write: %w", err)
		}
	}
	return nil
}
