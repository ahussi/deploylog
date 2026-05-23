package event

import (
	"testing"
	"time"
)

func baseEvent() DeployEvent {
	return DeployEvent{
		ID:          "deploy-001",
		Source:      SourceGitHubActions,
		Environment: "production",
		Service:     "api-gateway",
		Version:     "v1.4.2",
		Status:      StatusSuccess,
		TriggeredBy: "alice",
		Timestamp:   time.Now(),
	}
}

func TestIsTerminal(t *testing.T) {
	tests := []struct {
		status   Status
		wantTerm bool
	}{
		{StatusSuccess, true},
		{StatusFailure, true},
		{StatusCanceled, true},
		{StatusPending, false},
	}

	for _, tc := range tests {
		e := baseEvent()
		e.Status = tc.status
		if got := e.IsTerminal(); got != tc.wantTerm {
			t.Errorf("IsTerminal() for status %q = %v, want %v", tc.status, got, tc.wantTerm)
		}
	}
}

func TestValidate_Valid(t *testing.T) {
	e := baseEvent()
	if err := e.Validate(); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestValidate_MissingFields(t *testing.T) {
	tests := []struct {
		name    string
		mutate  func(*DeployEvent)
		wantErr error
	}{
		{"missing id", func(e *DeployEvent) { e.ID = "" }, ErrMissingID},
		{"missing service", func(e *DeployEvent) { e.Service = "" }, ErrMissingService},
		{"missing environment", func(e *DeployEvent) { e.Environment = "" }, ErrMissingEnvironment},
		{"missing timestamp", func(e *DeployEvent) { e.Timestamp = time.Time{} }, ErrMissingTimestamp},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			e := baseEvent()
			tc.mutate(&e)
			if err := e.Validate(); err != tc.wantErr {
				t.Errorf("Validate() = %v, want %v", err, tc.wantErr)
			}
		})
	}
}

func TestValidate_InvalidStatus(t *testing.T) {
	e := baseEvent()
	e.Status = Status("unknown")
	if err := e.Validate(); err != ErrInvalidStatus {
		t.Errorf("Validate() = %v, want %v", err, ErrInvalidStatus)
	}
}
