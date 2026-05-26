// Package export writes formatted deployment timeline data to a configured
// output destination.
//
// Supported destinations are stdout and file paths. Formatting is delegated
// to any output.Formatter implementation, decoupling serialisation from
// the write target.
//
// Basic usage:
//
//	ex, err := export.New(export.Options{
//		Destination: export.DestinationFile,
//		FilePath:    "timeline.json",
//		Formatter:   &output.JSONFormatter{},
//	})
//	if err != nil { ... }
//	if err := ex.Write(tl); err != nil { ... }
package export
