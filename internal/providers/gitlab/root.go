package gitlab

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"path"
	"strings"

	"github.com/qonto/standards-insights/config"
	"github.com/qonto/standards-insights/pkg/project"
	"github.com/xanzy/go-gitlab"
)

type Client struct {
	config config.GitlabConfig
	client *gitlab.Client
	logger *slog.Logger
}

func New(logger *slog.Logger, config config.GitlabConfig) (*Client, error) {
	token := os.Getenv("GITLAB_TOKEN")
	if token != "" {
		config.Token = token
		err := os.Unsetenv("GITLAB_TOKEN")
		if err != nil {
			return nil, err
		}
	}

	client, err := gitlab.NewClient(config.Token, gitlab.WithBaseURL(config.URL))
	if err != nil {
		return nil, err
	}

	return &Client{
		client: client,
		config: config,
		logger: logger,
	}, nil
}

func (c *Client) Name() string {
	return "gitlab"
}

func (c *Client) FetchProjects(ctx context.Context) ([]project.Project, error) {
	options := &gitlab.ListProjectsOptions{
		ListOptions: gitlab.ListOptions{PerPage: c.config.Limit},
	}

	if len(c.config.Topics) > 0 {
		topic := strings.Join(c.config.Topics, ",")
		options.Topic = &topic
	}

	if c.config.Search != "" {
		options.Search = &c.config.Search
	}

	c.logger.Debug("fetching Gitlab projects")
	projects, _, err := c.client.Projects.ListProjects(
		options,
	)
	if err != nil {
		return nil, err
	}

	c.logger.Debug(fmt.Sprintf("found %d Gitlab projects", len(projects)))
	result := make([]project.Project, len(projects))
	for i, proj := range projects {
		labels := make(map[string]string)
		for _, topic := range proj.Topics {
			if split := strings.Split(topic, ": "); len(split) > 1 {
				labels[split[0]] = split[1]
			}
		}

		result[i] = project.Project{
			Name:   proj.Name,
			URL:    proj.SSHURLToRepo,
			Branch: proj.DefaultBranch,
			Path:   path.Join(c.config.BasePath, proj.PathWithNamespace),
			Labels: labels,
		}
	}
	return result, nil
}
