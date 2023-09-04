package argocd

import (
	"standards/config"
	"standards/rules/aggregate"
)

type Client struct {
	config *config.ArgoCDConfig
}

func New(config *config.ArgoCDConfig) *Client {
	return &Client{
		config: config,
	}
}

func (c *Client) FetchProjects() ([]*aggregate.Project, error) {
	// TODO: fixme
}
