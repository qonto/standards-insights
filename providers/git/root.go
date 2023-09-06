package git

import (
	"standards/config"
	"standards/providers/aggregates"
)

type Client struct {
	config *config.GitConfig
}

func New(config *config.GitConfig) *Client {
	return &Client{
		config: config,
	}
}

func (c *Client) FetchProjects() ([]*aggregates.Project, error) {
	// TODO: fixme
	return nil, nil
}
