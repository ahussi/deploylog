package jenkins

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/deploylog/deploylog/internal/event"
)

// SourceName is the canonical identifier for the Jenkins source.
const SourceName = "jenkins"

type buildResponse struct {
	Builds []build `json:"builds"`
}

type build struct {
	ID          string  `json:"id"`
	FullName    string  `json:"fullDisplayName"`
	Result      string  `json:"result"`
	Timestamp   int64   `json:"timestamp"`
	DurationMS  int64   `json:"duration"`
	Description string  `json:"description"`
}

// Client fetches deployment events from a Jenkins instance.
type Client struct {
	baseURL    string
	jobName    string
	httpClient *http.Client
}

// NewClient creates a new Jenkins source client.
func NewClient(baseURL, jobName string) *Client {
	return &Client{
		baseURL:    baseURL,
		jobName:    jobName,
		httpClient: &http.Client{Timeout: 10 * time.Second},
	}
}

// Name returns the source identifier.
func (c *Client) Name() string {
	return SourceName
}

// Fetch retrieves recent builds from Jenkins and maps them to events.
func (c *Client) Fetch() ([]event.Event, error) {
	url := fmt.Sprintf("%s/job/%s/api/json?tree=builds[id,fullDisplayName,result,timestamp,duration,description]",
		c.baseURL, c.jobName)

	resp, err := c.httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("jenkins fetch: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("jenkins fetch: unexpected status %d", resp.StatusCode)
	}

	var br buildResponse
	if err := json.NewDecoder(resp.Body).Decode(&br); err != nil {
		return nil, fmt.Errorf("jenkins decode: %w", err)
	}

	events := make([]event.Event, 0, len(br.Builds))
	for _, b := range br.Builds {
		e := event.Event{
			ID:        b.ID,
			Source:    SourceName,
			Service:   b.FullName,
			Status:    mapStatus(b.Result),
			Timestamp: time.UnixMilli(b.Timestamp).UTC(),
			Meta: map[string]string{
				"duration_ms":  fmt.Sprintf("%d", b.DurationMS),
				"description": b.Description,
			},
		}
		if err := e.Validate(); err == nil {
			events = append(events, e)
		}
	}
	return events, nil
}

func mapStatus(result string) event.Status {
	switch result {
	case "SUCCESS":
		return event.StatusSuccess
	case "FAILURE":
		return event.StatusFailed
	case "ABORTED":
		return event.StatusCancelled
	case "":
		return event.StatusRunning
	default:
		return event.StatusUnknown
	}
}
