package argocd

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path"

	"github.com/qonto/standards-insights/config"
	"github.com/qonto/standards-insights/pkg/project"

	"log/slog"
)

type Client struct {
	config config.ArgoCDConfig
	token  string
	client *http.Client
	logger *slog.Logger
}

type argocdSource struct {
	RepoURL        string
	TargetRevision string
}

type argocdSpec struct {
	Source argocdSource
}

type argocdMetadata struct {
	Labels map[string]string
	Name   string
}

type argocdApplication struct {
	Metadata argocdMetadata
	Spec     argocdSpec
}

type argocdListResponse struct {
	Items []argocdApplication
}

func New(logger *slog.Logger, config config.ArgoCDConfig) (*Client, error) {
	client := &http.Client{}
	token := os.Getenv("ARGOCD_AUTH_TOKEN")
	if token == "" {
		return nil, errors.New("You need to set the ARGOCD_AUTH_TOKEN environment variable")
	}
	err := os.Unsetenv("ARGOCD_AUTH_TOKEN")
	if err != nil {
		return nil, err
	}

	return &Client{
		client: client,
		config: config,
		logger: logger,
		token:  token,
	}, nil
}

func (c *Client) Name() string {
	return "argocd"
}

func (c *Client) FetchProjects(ctx context.Context) ([]project.Project, error) {
	c.logger.Debug("fetching ArgoCD applications")

	url, err := url.JoinPath(c.config.URL, "/api/v1/applications")
	if err != nil {
		return nil, fmt.Errorf("fail to build argocd url: %w", err)
	}
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", c.token))

	q := req.URL.Query()
	for _, project := range c.config.Projects {
		q.Add("projects", project)
	}
	if c.config.Selector != "" {
		q.Add("selector", c.config.Selector)
	}
	req.URL.RawQuery = q.Encode()
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close() //nolint
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("fail to list ArgoCD applications: status %d", resp.StatusCode)
	}

	var apps argocdListResponse
	err = json.NewDecoder(resp.Body).Decode(&apps)
	if err != nil {
		return nil, err
	}

	c.logger.Debug(fmt.Sprintf("found %d ArgoCD applications", len(apps.Items)))
	result := make([]project.Project, len(apps.Items))
	for i, app := range apps.Items {
		result[i] = project.Project{
			Name:   app.Metadata.Name,
			URL:    app.Spec.Source.RepoURL,
			Branch: app.Spec.Source.TargetRevision,
			Path:   path.Join(c.config.BasePath, app.Metadata.Name),
			Labels: app.Metadata.Labels,
		}
	}
	return result, nil
}
