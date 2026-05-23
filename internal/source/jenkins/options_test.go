package jenkins

import (
	"testing"
)

func TestWithBaseURL(t *testing.T) {
	c := &Client{}
	opt := WithBaseURL("http://jenkins.example.com")
	opt(c)
	if c.baseURL != "http://jenkins.example.com" {
		t.Errorf("expected baseURL %q, got %q", "http://jenkins.example.com", c.baseURL)
	}
}

func TestWithJob(t *testing.T) {
	c := &Client{}
	opt := WithJob("my-pipeline")
	opt(c)
	if c.job != "my-pipeline" {
		t.Errorf("expected job %q, got %q", "my-pipeline", c.job)
	}
}

func TestWithCredentials(t *testing.T) {
	c := &Client{}
	opt := WithCredentials("admin", "secret-token")
	opt(c)
	if c.username != "admin" {
		t.Errorf("expected username %q, got %q", "admin", c.username)
	}
	if c.token != "secret-token" {
		t.Errorf("expected token %q, got %q", "secret-token", c.token)
	}
}

func TestWithOptions_CombinedApply(t *testing.T) {
	c := &Client{}
	opts := []Option{
		WithBaseURL("http://ci.local"),
		WithJob("deploy-job"),
		WithCredentials("user", "tok"),
	}
	for _, o := range opts {
		o(c)
	}
	if c.baseURL != "http://ci.local" {
		t.Errorf("unexpected baseURL: %q", c.baseURL)
	}
	if c.job != "deploy-job" {
		t.Errorf("unexpected job: %q", c.job)
	}
	if c.username != "user" || c.token != "tok" {
		t.Errorf("unexpected credentials: %q / %q", c.username, c.token)
	}
}
