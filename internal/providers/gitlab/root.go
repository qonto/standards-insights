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
	options := makeGitlabListProjectsOptions(c.config)

	c.logger.Debug("fetching Gitlab projects")
	projects, err := c.listAllProject(options)
	if err != nil {
		return nil, err
	}

	c.logger.Debug(fmt.Sprintf("found %d Gitlab projects", len(projects)))

	return c.makeProjects(projects), nil
}

func makeGitlabListProjectsOptions(config config.GitlabConfig) *gitlab.ListProjectsOptions {
	options := &gitlab.ListProjectsOptions{
		ListOptions: gitlab.ListOptions{
			PerPage: 25,
			Sort:    "asc",
		},
	}

	if len(config.Topics) > 0 {
		topic := strings.Join(config.Topics, ",")
		options.Topic = &topic
	}

	if config.Search != "" {
		options.Search = &config.Search
	}

	return options
}

func (c *Client) listAllProject(opts *gitlab.ListProjectsOptions) ([]*gitlab.Project, error) {
	result := []*gitlab.Project{}

	for {
		ps, resp, err := c.client.Projects.ListProjects(opts)
		if err != nil {
			return nil, err
		}

		result = append(result, ps...)

		if resp.NextPage == 0 {
			break
		}

		opts.Page = resp.NextPage
	}
	return result, nil
}

func (c *Client) makeProjects(projects []*gitlab.Project) []project.Project {
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
	return result
}
