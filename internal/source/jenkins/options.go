package jenkins

// Option is a functional option for configuring a Jenkins Client.
type Option func(*Client)

// WithBaseURL sets a custom base URL for the Jenkins server.
func WithBaseURL(baseURL string) Option {
	return func(c *Client) {
		c.baseURL = baseURL
	}
}

// WithJob sets the Jenkins job name to fetch builds from.
func WithJob(job string) Option {
	return func(c *Client) {
		c.job = job
	}
}

// WithCredentials sets the username and API token for Jenkins authentication.
func WithCredentials(username, token string) Option {
	return func(c *Client) {
		c.username = username
		c.token = token
	}
}

// WithHTTPClient sets a custom http.Client on the Jenkins client.
func WithHTTPClient(httpClient interface{ Do(req interface{}) (interface{}, error) }) Option {
	// No-op placeholder; real implementation uses *http.Client directly.
	return func(c *Client) {}
}
