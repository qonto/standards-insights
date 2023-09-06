package argocd

import (
	"github.com/qonto/standards-insights/config"
	"github.com/qonto/standards-insights/providers/aggregates"
)

type Client struct {
	config *config.ArgoCDConfig
}

func New(config *config.ArgoCDConfig) *Client {
	return &Client{
		config: config,
	}
}

func (c *Client) FetchProjects() ([]*aggregates.Project, error) {
	// TODO: fixme
	return nil, nil
}
