package gitlab

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/deploylog/deploylog/internal/event"
)

const SourceName = "gitlab"

// Client fetches deployment events from the GitLab Deployments API.
type Client struct {
	projectID string
	token     string
	baseURL   string
	httpClient *http.Client
}

type gitlabDeployment struct {
	ID          int    `json:"id"`
	Status      string `json:"status"`
	Ref         string `json:"ref"`
	Environment struct {
		Name string `json:"name"`
	} `json:"environment"`
	CreatedAt string `json:"created_at"`
}

// NewClient creates a new GitLab source client.
func NewClient(projectID, token, baseURL string) *Client {
	if baseURL == "" {
		baseURL = "https://gitlab.com/api/v4"
	}
	return &Client{
		projectID:  projectID,
		token:      token,
		baseURL:    baseURL,
		httpClient: &http.Client{Timeout: 10 * time.Second},
	}
}

// Name returns the source identifier.
func (c *Client) Name() string {
	return SourceName
}

// Fetch retrieves deployment events from GitLab.
func (c *Client) Fetch() ([]event.Event, error) {
	url := fmt.Sprintf("%s/projects/%s/deployments?order_by=created_at&sort=desc", c.baseURL, c.projectID)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("gitlab: build request: %w", err)
	}
	req.Header.Set("PRIVATE-TOKEN", c.token)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("gitlab: http request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("gitlab: unexpected status %d", resp.StatusCode)
	}

	var deployments []gitlabDeployment
	if err := json.NewDecoder(resp.Body).Decode(&deployments); err != nil {
		return nil, fmt.Errorf("gitlab: decode response: %w", err)
	}

	events := make([]event.Event, 0, len(deployments))
	for _, d := range deployments {
		t, err := time.Parse(time.RFC3339, d.CreatedAt)
		if err != nil {
			continue
		}
		events = append(events, event.Event{
			Source:    SourceName,
			Service:   d.Environment.Name,
			Status:    mapStatus(d.Status),
			Timestamp: t,
			Meta:      map[string]string{"ref": d.Ref, "deployment_id": fmt.Sprintf("%d", d.ID)},
		})
	}
	return events, nil
}

func mapStatus(s string) event.Status {
	switch s {
	case "success":
		return event.StatusSuccess
	case "failed":
		return event.StatusFailure
	case "running":
		return event.StatusRunning
	default:
		return event.StatusPending
	}
}
