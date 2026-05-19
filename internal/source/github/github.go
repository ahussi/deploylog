// Package github implements a deploylog source for GitHub Actions deployment events.
package github

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/your-org/deploylog/internal/event"
)

const SourceName = "github"

// Client fetches deployment events from the GitHub Actions API.
type Client struct {
	owner   string
	repo    string
	token   string
	httpCli *http.Client
}

// NewClient creates a new GitHub source client.
func NewClient(owner, repo, token string) *Client {
	return &Client{
		owner:   owner,
		repo:    repo,
		token:   token,
		httpCli: &http.Client{Timeout: 15 * time.Second},
	}
}

// Name returns the source identifier.
func (c *Client) Name() string { return SourceName }

// Fetch retrieves recent deployment events from the GitHub Deployments API.
func (c *Client) Fetch(ctx context.Context) ([]event.Event, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/deployments", c.owner, c.repo)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("github: build request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+c.token)
	req.Header.Set("Accept", "application/vnd.github+json")

	resp, err := c.httpCli.Do(req)
	if err != nil {
		return nil, fmt.Errorf("github: http request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("github: unexpected status %d", resp.StatusCode)
	}

	var raw []struct {
		ID          int64     `json:"id"`
		Environment string    `json:"environment"`
		CreatedAt   time.Time `json:"created_at"`
		Description string    `json:"description"`
		Task        string    `json:"task"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
		return nil, fmt.Errorf("github: decode response: %w", err)
	}

	events := make([]event.Event, 0, len(raw))
	for _, d := range raw {
		ev := event.Event{
			ID:        fmt.Sprintf("github-%d", d.ID),
			Source:    SourceName,
			Service:   fmt.Sprintf("%s/%s", c.owner, c.repo),
			Status:    event.StatusRunning,
			Timestamp: d.CreatedAt,
			Metadata:  map[string]string{"environment": d.Environment, "task": d.Task},
		}
		if d.Description != "" {
			ev.Message = d.Description
		}
		events = append(events, ev)
	}
	return events, nil
}
