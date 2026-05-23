package gitlab_test

import (
	"testing"

	"github.com/deploylog/deploylog/internal/source/gitlab"
)

func TestWithToken(t *testing.T) {
	client := gitlab.NewClient("test", gitlab.WithToken("mytoken"))
	if client == nil {
		t.Fatal("expected non-nil client")
	}
}

func TestWithProjectID(t *testing.T) {
	client := gitlab.NewClient("test", gitlab.WithProjectID("99"))
	if client == nil {
		t.Fatal("expected non-nil client")
	}
}

func TestWithBaseURL(t *testing.T) {
	client := gitlab.NewClient("test",
		gitlab.WithToken("tok"),
		gitlab.WithProjectID("1"),
		gitlab.WithBaseURL("https://gitlab.mycompany.io"),
	)
	if client == nil {
		t.Fatal("expected non-nil client")
	}
	if client.Name() != "test" {
		t.Errorf("expected name 'test', got %s", client.Name())
	}
}

func TestWithOptions_CombinedApply(t *testing.T) {
	client := gitlab.NewClient("combined",
		gitlab.WithToken("tok123"),
		gitlab.WithProjectID("55"),
	)
	if client.Name() != "combined" {
		t.Errorf("expected name 'combined', got %s", client.Name())
	}
}
