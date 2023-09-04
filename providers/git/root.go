package git

import "standards/rules/aggregate"

type Client struct {
	config *config.GitConfig
}

func New(config *config.GitConfig) *Client {
	return &Client{
		config: config,
	}
}

func (c *Client) FetchProjects() ([]*aggregate.Project, error) {
	// TODO: fixme
}
