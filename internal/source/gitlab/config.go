package gitlab

import (
	"errors"
	"fmt"

	"github.com/deploylog/deploylog/internal/config"
)

// NewClientFromConfig constructs a GitLab Client from a source config entry.
// Required params: token, project_id
// Optional params: base_url (defaults to https://gitlab.com)
func NewClientFromConfig(cfg config.SourceConfig) (*Client, error) {
	var errs []error

	token, ok := cfg.Params["token"]
	if !ok || token == "" {
		errs = append(errs, errors.New("missing required param: token"))
	}

	projectID, ok := cfg.Params["project_id"]
	if !ok || projectID == "" {
		errs = append(errs, errors.New("missing required param: project_id"))
	}

	if len(errs) > 0 {
		return nil, joinErrors(errs)
	}

	opts := []func(*Client){
		WithToken(token),
		WithProjectID(projectID),
	}

	if baseURL, ok := cfg.Params["base_url"]; ok && baseURL != "" {
		opts = append(opts, WithBaseURL(baseURL))
	}

	return NewClient(cfg.Name, opts...), nil
}

func joinErrors(errs []error) error {
	if len(errs) == 0 {
		return nil
	}
	msg := ""
	for i, e := range errs {
		if i > 0 {
			msg += "; "
		}
		msg += e.Error()
	}
	return fmt.Errorf("%s", msg)
}
