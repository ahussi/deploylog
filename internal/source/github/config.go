package github

import (
	"errors"
	"fmt"
	"strings"

	"github.com/yourorg/deploylog/internal/config"
)

// NewClientFromConfig constructs a GitHub Client from a source configuration
// entry. It validates required fields and returns a descriptive error if any
// are missing or malformed.
func NewClientFromConfig(cfg config.SourceConfig) (*Client, error) {
	var errs []string

	token, ok := cfg.Options["token"]
	if !ok || strings.TrimSpace(token) == "" {
		errs = append(errs, "missing required option: token")
	}

	repo, ok := cfg.Options["repo"]
	if !ok || strings.TrimSpace(repo) == "" {
		errs = append(errs, "missing required option: repo")
	}

	if len(errs) > 0 {
		return nil, fmt.Errorf("github source %q: %w", cfg.Name, joinErrors(errs))
	}

	baseURL := cfg.Options["base_url"] // optional; empty string is fine

	return NewClient(cfg.Name, token, repo, baseURL), nil
}

// joinErrors combines a slice of error strings into a single error.
func joinErrors(msgs []string) error {
	return errors.New(strings.Join(msgs, "; "))
}
