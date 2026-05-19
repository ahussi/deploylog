// Package config provides configuration loading and validation for deploylog.
package config

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// SourceConfig holds configuration for a single CI/CD source.
type SourceConfig struct {
	Name    string            `json:"name"`
	Type    string            `json:"type"`
	Options map[string]string `json:"options,omitempty"`
}

// OutputConfig holds configuration for a single output formatter.
type OutputConfig struct {
	Format     string `json:"format"`
	Indent     bool   `json:"indent,omitempty"`
	TimeFormat string `json:"time_format,omitempty"`
}

// Config is the top-level deploylog configuration.
type Config struct {
	Sources      []SourceConfig `json:"sources"`
	Output       OutputConfig   `json:"output"`
	PollInterval Duration       `json:"poll_interval,omitempty"`
}

// Duration is a wrapper around time.Duration for JSON unmarshalling.
type Duration struct {
	time.Duration
}

func (d *Duration) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	parsed, err := time.ParseDuration(s)
	if err != nil {
		return fmt.Errorf("invalid duration %q: %w", s, err)
	}
	d.Duration = parsed
	return nil
}

// Load reads and parses a JSON config file from the given path.
func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading config file: %w", err)
	}
	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parsing config file: %w", err)
	}
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}
	return &cfg, nil
}

// Validate checks that the configuration is semantically valid.
func (c *Config) Validate() error {
	if len(c.Sources) == 0 {
		return fmt.Errorf("at least one source must be configured")
	}
	seen := make(map[string]bool)
	for i, s := range c.Sources {
		if s.Name == "" {
			return fmt.Errorf("source[%d]: name is required", i)
		}
		if s.Type == "" {
			return fmt.Errorf("source[%d] %q: type is required", i, s.Name)
		}
		if seen[s.Name] {
			return fmt.Errorf("duplicate source name %q", s.Name)
		}
		seen[s.Name] = true
	}
	if c.Output.Format == "" {
		c.Output.Format = "text"
	}
	return nil
}
