package static

import (
	"context"

	"log/slog"

	"github.com/qonto/standards-insights/pkg/project"
)

type Static struct {
	logger   *slog.Logger
	projects []project.Project
}

func New(logger *slog.Logger, projects []project.Project) *Static {
	return &Static{
		logger:   logger,
		projects: projects,
	}
}

func (c *Static) Name() string {
	return "static"
}

func (c *Static) FetchProjects(_ context.Context) ([]project.Project, error) {
	c.logger.Debug("fetching static projects")
	return c.projects, nil
}
