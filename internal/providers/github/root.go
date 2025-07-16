package github

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"

	"github.com/bradleyfalzon/ghinstallation/v2"
	"github.com/google/go-github/v71/github"
	"github.com/qonto/standards-insights/config"
	"github.com/qonto/standards-insights/pkg/daemon"
	"github.com/qonto/standards-insights/pkg/project"
)

type Client struct {
	config config.GithubConfig
	client *github.Client
	logger *slog.Logger
}

func New(logger *slog.Logger, config config.GithubConfig) (*Client, error) {
	appID := os.Getenv("GITHUB_APP_ID")
	if appID != "" {
		appID, err := strconv.ParseInt(appID, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("could not parse github app id: %w", err)
		}
		config.AppID = appID
		err = os.Unsetenv("GITHUB_APP_ID")
		if err != nil {
			return nil, err
		}
	}

	installationID := os.Getenv("GITHUB_INSTALLATION_ID")
	if installationID != "" {
		installationID, err := strconv.ParseInt(installationID, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("could not parse github installation id: %w", err)
		}
		config.InstallationID = installationID
		err = os.Unsetenv("GITHUB_INSTALLATION_ID")
		if err != nil {
			return nil, err
		}
	}

	privateKey := os.Getenv("GITHUB_PRIVATE_KEY")
	if privateKey != "" {
		config.PrivateKey = privateKey
		err := os.Unsetenv("GITHUB_PRIVATE_KEY")
		if err != nil {
			return nil, err
		}
	}

	itr, err := ghinstallation.New(http.DefaultTransport, config.AppID, config.InstallationID, []byte(config.PrivateKey))
	if err != nil {
		return nil, err
	}

	client := github.NewClient(&http.Client{Transport: itr})

	return &Client{
		client: client,
		config: config,
		logger: logger,
	}, nil
}

func (c *Client) Name() string {
	return "github"
}

func (c *Client) FetchProjects(ctx context.Context) ([]project.Project, error) {
	options := &github.SearchOptions{
		TextMatch: true,
		ListOptions: github.ListOptions{
			PerPage: 100,
		},
	}

	c.logger.Debug("fetching GitHub repositories")
	projects, err := c.listAllRepositories(options)
	if err != nil {
		return nil, err
	}

	c.logger.Debug(fmt.Sprintf("found %d GitHub repositories", len(projects)))
	return c.makeProjects(projects), nil
}

func (c *Client) listAllRepositories(opts *github.SearchOptions) ([]*github.Repository, error) {
	result := []*github.Repository{}
	var queryParts []string

	if len(c.config.Topics) > 0 {
		topics := make([]string, len(c.config.Topics))
		for i, topic := range c.config.Topics {
			topics[i] = fmt.Sprintf("topic:%s", topic)
		}
		queryParts = append(queryParts, strings.Join(topics, " "))
	}

	if len(c.config.Organizations) > 0 {
		orgs := make([]string, len(c.config.Organizations))
		for i, org := range c.config.Organizations {
			orgs[i] = fmt.Sprintf("org:%s", org)
		}
		queryParts = append(queryParts, strings.Join(orgs, " "))
	}

	query := strings.Join(queryParts, " ")

	for {
		searchResult, resp, err := c.client.Search.Repositories(context.Background(), query, opts)
		if err != nil {
			return nil, err
		}

		result = append(result, searchResult.Repositories...)

		if resp.NextPage == 0 {
			break
		}
		opts.Page = resp.NextPage
	}
	return result, nil
}

func (c *Client) makeProjects(repos []*github.Repository) []project.Project {
	result := make([]project.Project, len(repos))
	for i, repo := range repos {
		labels := make(map[string]string)
		for _, topic := range repo.Topics {
			labels[topic] = "true"
		}

		result[i] = project.Project{
			Name:   repo.GetName(),
			URL:    repo.GetCloneURL(),
			Branch: repo.GetDefaultBranch(),
			Path:   path.Join(c.config.BasePath, repo.GetFullName()),
			Labels: labels,
		}
	}
	return result
}

func (c *Client) ConfigureGit(g daemon.Git) error {
	itr, err := ghinstallation.New(http.DefaultTransport, c.config.AppID, c.config.InstallationID, []byte(c.config.PrivateKey))
	if err != nil {
		return err
	}

	token, err := itr.Token(context.Background())
	if err != nil {
		c.logger.Error("failed to generate GitHub installation token", "error", err)
		return fmt.Errorf("failed to generate GitHub installation token: %w", err)
	}

	g.SetToken(token)
	return nil
}
