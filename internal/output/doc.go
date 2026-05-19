// Package output provides pluggable formatters for rendering the deploylog
// unified audit timeline to different output targets.
//
// # Available Formatters
//
// JSONFormatter serialises the event slice as a JSON array. Set Indent to true
// for pretty-printed output suitable for human review.
//
// TextFormatter writes one human-readable line per event, including the
// timestamp, source, status, and description. The timestamp layout can be
// customised via the TimeFormat field (defaults to time.RFC3339).
//
// # Usage
//
//	var f output.Formatter = &output.JSONFormatter{Indent: true}
//	if err := f.Format(os.Stdout, events); err != nil {
//	    log.Fatal(err)
//	}
package output
