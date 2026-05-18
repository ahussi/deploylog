package event

import "time"

// Source represents the CI/CD platform that produced the deployment event.
type Source string

const (
	SourceGitHubActions Source = "github_actions"
	SourceGitLabCI      Source = "gitlab_ci"
	SourceCircleCI      Source = "circleci"
	SourceJenkins       Source = "jenkins"
	SourceUnknown       Source = "unknown"
)

// Status represents the outcome of a deployment.
type Status string

const (
	StatusSuccess  Status = "success"
	StatusFailure  Status = "failure"
	StatusPending  Status = "pending"
	StatusCanceled Status = "canceled"
)

// DeployEvent is the unified representation of a deployment event
// regardless of the originating CI/CD source.
type DeployEvent struct {
	ID          string            `json:"id"`
	Source      Source            `json:"source"`
	Environment string            `json:"environment"`
	Service     string            `json:"service"`
	Version     string            `json:"version"`
	Status      Status            `json:"status"`
	TriggeredBy string            `json:"triggered_by"`
	Timestamp   time.Time         `json:"timestamp"`
	Duration    time.Duration     `json:"duration_ms"`
	Metadata    map[string]string `json:"metadata,omitempty"`
}

// IsTerminal returns true if the event status represents a final state.
func (e *DeployEvent) IsTerminal() bool {
	return e.Status == StatusSuccess ||
		e.Status == StatusFailure ||
		e.Status == StatusCanceled
}

// Validate checks that the required fields of a DeployEvent are populated.
func (e *DeployEvent) Validate() error {
	if e.ID == "" {
		return ErrMissingID
	}
	if e.Service == "" {
		return ErrMissingService
	}
	if e.Environment == "" {
		return ErrMissingEnvironment
	}
	if e.Timestamp.IsZero() {
		return ErrMissingTimestamp
	}
	return nil
}
