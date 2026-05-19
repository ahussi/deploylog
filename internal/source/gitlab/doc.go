// Package gitlab provides a deploylog source client for GitLab CI/CD.
//
// It fetches deployment events from the GitLab Deployments API
// (https://docs.gitlab.com/ee/api/deployments.html) and converts them
// into the unified event.Event format used across the deploylog system.
//
// # Usage
//
//	client := gitlab.NewClient(
//		"<project-id>",
//		"<private-token>",
//		"", // empty string uses the default https://gitlab.com/api/v4 base URL
//	)
//	events, err := client.Fetch()
//
// # Status Mapping
//
// GitLab deployment statuses are mapped as follows:
//
//	"success" → event.StatusSuccess
//	"failed"  → event.StatusFailure
//	"running" → event.StatusRunning
//	(other)   → event.StatusPending
package gitlab
