package github

import (
	"net/http"
	"testing"
)

func TestWithToken(t *testing.T) {
	c := &Client{}
	WithToken("tok123")(c)
	if c.token != "tok123" {
		t.Errorf("expected token %q, got %q", "tok123", c.token)
	}
}

func TestWithRepo(t *testing.T) {
	c := &Client{}
	WithRepo("owner/repo")(c)
	if c.repo != "owner/repo" {
		t.Errorf("expected repo %q, got %q", "owner/repo", c.repo)
	}
}

func TestWithBaseURL(t *testing.T) {
	c := &Client{}
	WithBaseURL("https://github.example.com")(c)
	if c.baseURL != "https://github.example.com" {
		t.Errorf("expected baseURL %q, got %q", "https://github.example.com", c.baseURL)
	}
}

func TestWithOptions_CombinedApply(t *testing.T) {
	httpClient := &http.Client{}
	c := &Client{}
	opts := []Option{
		WithToken("mytoken"),
		WithRepo("acme/app"),
		WithBaseURL("https://api.github.com"),
	}
	for _, o := range opts {
		o(c)
	}
	_ = httpClient
	if c.token != "mytoken" {
		t.Errorf("token mismatch")
	}
	if c.repo != "acme/app" {
		t.Errorf("repo mismatch")
	}
	if c.baseURL != "https://api.github.com" {
		t.Errorf("baseURL mismatch")
	}
}
