package jenkins

import (
	"errors"
	"fmt"
)

// Config holds configuration required to instantiate a Jenkins client.
type Config struct {
	// BaseURL is the root URL of the Jenkins instance, e.g. "https://ci.example.com".
	BaseURL string `yaml:"base_url"`
	// JobName is the Jenkins job (project) name to query for builds.
	JobName string `yaml:"job_name"`
}

// Validate checks that all required Config fields are present.
func (c Config) Validate() error {
	var errs []error
	if c.BaseURL == "" {
		errs = append(errs, errors.New("jenkins: base_url is required"))
	}
	if c.JobName == "" {
		errs = append(errs, errors.New("jenkins: job_name is required"))
	}
	if len(errs) > 0 {
		return joinErrors(errs)
	}
	return nil
}

// NewClientFromConfig creates a Jenkins Client from a validated Config.
func NewClientFromConfig(cfg Config) (*Client, error) {
	if err := cfg.Validate(); err != nil {
		return nil, err
	}
	return NewClient(cfg.BaseURL, cfg.JobName), nil
}

func joinErrors(errs []error) error {
	if len(errs) == 1 {
		return errs[0]
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
