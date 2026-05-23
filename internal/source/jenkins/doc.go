// Package jenkins provides a deploylog source client for Jenkins CI.
//
// It connects to a Jenkins instance via its JSON API and retrieves
// build records for a configured job, mapping each build to a
// unified deploylog event.
//
// # Usage
//
//	client := jenkins.NewClient("https://ci.example.com", "deploy-prod")
//	events, err := client.Fetch()
//
// # Status Mapping
//
// Jenkins build results are mapped as follows:
//
//	SUCCESS  → event.StatusSuccess
//	FAILURE  → event.StatusFailed
//	ABORTED  → event.StatusCancelled
//	""       → event.StatusRunning  (build still in progress)
//	other    → event.StatusUnknown
package jenkins
