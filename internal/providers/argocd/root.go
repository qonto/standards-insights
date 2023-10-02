package argocd

import (
	"context"
	"fmt"
	"os"
	"path"

	"github.com/argoproj/argo-cd/v2/pkg/apiclient"
	applicationpkg "github.com/argoproj/argo-cd/v2/pkg/apiclient/application"
	"github.com/qonto/standards-insights/config"
	"github.com/qonto/standards-insights/pkg/project"

	"log/slog"
)

type Client struct {
	config config.ArgoCDConfig
	client apiclient.Client
	logger *slog.Logger
}

func New(logger *slog.Logger, config config.ArgoCDConfig) (*Client, error) {
	clientConfig := apiclient.ClientOptions{
		ServerAddr: config.URL,
	}

	client, err := apiclient.NewClient(&clientConfig)
	if err != nil {
		return nil, err
	}
	err = os.Unsetenv("ARGOCD_AUTH_TOKEN")
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
	return "argocd"
}

func (c *Client) FetchProjects(ctx context.Context) ([]project.Project, error) {
	c.logger.Debug("fetching ArgoCD projects")
	conn, appClient, err := c.client.NewApplicationClient()
	if err != nil {
		return nil, err
	}
	defer conn.Close() //nolint

	query := &applicationpkg.ApplicationQuery{
		Projects: c.config.Projects,
	}
	if c.config.Selector != "" {
		query.Selector = &c.config.Selector
	}
	apps, err := appClient.List(ctx, query)
	if err != nil {
		return nil, err
	}
	c.logger.Debug(fmt.Sprintf("found %d ArgoCD projects", len(apps.Items)))
	result := make([]project.Project, len(apps.Items))
	for i, app := range apps.Items {
		result[i] = project.Project{
			Name:   app.Name,
			URL:    app.Spec.Source.RepoURL,
			Branch: app.Spec.Source.TargetRevision,
			Path:   path.Join(c.config.BasePath, app.Name),
			Labels: app.Labels,
		}
	}
	return result, nil
}
