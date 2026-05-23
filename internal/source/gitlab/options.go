package gitlab

// WithToken sets the personal access token used for GitLab API authentication.
func WithToken(token string) func(*Client) {
	return func(c *Client) {
		c.token = token
	}
}

// WithProjectID sets the GitLab project ID (numeric or namespace/project path).
func WithProjectID(projectID string) func(*Client) {
	return func(c *Client) {
		c.projectID = projectID
	}
}

// WithBaseURL overrides the default GitLab instance URL.
// Useful for self-hosted GitLab installations.
func WithBaseURL(baseURL string) func(*Client) {
	return func(c *Client) {
		c.baseURL = baseURL
	}
}
