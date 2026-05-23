package github

// Option is a functional option for configuring a GitHub Client.
type Option func(*Client)

// WithToken sets the personal access token used for GitHub API authentication.
func WithToken(token string) Option {
	return func(c *Client) {
		c.token = token
	}
}

// WithRepo sets the repository in "owner/repo" format to fetch workflow runs from.
func WithRepo(repo string) Option {
	return func(c *Client) {
		c.repo = repo
	}
}

// WithBaseURL overrides the default GitHub API base URL (useful for GitHub Enterprise).
func WithBaseURL(baseURL string) Option {
	return func(c *Client) {
		c.baseURL = baseURL
	}
}
