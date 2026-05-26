package export

import "github.com/yourorg/deploylog/internal/output"

// Option is a functional option for configuring an Exporter.
type Option func(*Options)

// WithFile sets the destination to a file at the given path.
func WithFile(path string) Option {
	return func(o *Options) {
		o.Destination = DestinationFile
		o.FilePath = path
	}
}

// WithStdout sets the destination to standard output.
func WithStdout() Option {
	return func(o *Options) {
		o.Destination = DestinationStdout
		o.FilePath = ""
	}
}

// WithFormatter sets the formatter used to serialise events.
func WithFormatter(f output.Formatter) Option {
	return func(o *Options) {
		o.Formatter = f
	}
}

// NewWithOptions creates an Exporter by applying the provided functional options.
func NewWithOptions(opts ...Option) (*Exporter, error) {
	cfg := Options{Destination: DestinationStdout}
	for _, o := range opts {
		o(&cfg)
	}
	return New(cfg)
}
