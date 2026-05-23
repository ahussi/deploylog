// Package github provides a deploylog source client for GitHub Actions.
//
// It fetches deployment events from the GitHub Deployments API and maps them
// to the unified deploylog event model. Authentication is performed via a
// personal access token supplied through source configuration options.
//
// # Configuration Options
//
//   - token (required): GitHub personal access token with repo scope.
//   - repo (required): Repository in "owner/name" format, e.g. "acme/api".
//   - base_url (optional): Override the GitHub API base URL, useful for
//     GitHub Enterprise Server instances.
//
// # Usage
//
// Construct a client directly with NewClient or from a config.SourceConfig
// entry using NewClientFromConfig. The latter is the preferred approach when
// building clients from a YAML configuration file.
package github
