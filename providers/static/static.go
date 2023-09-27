package static

import (
	"context"

	"log/slog"

	"github.com/qonto/standards-insights/providers/aggregates"
)

type Static struct {
	logger   *slog.Logger
	projects []aggregates.Project
}

func New(logger *slog.Logger, projects []aggregates.Project) *Static {
	return &Static{
		logger:   logger,
		projects: projects,
	}
}

func (c *Static) FetchProjects(_ context.Context) ([]aggregates.Project, error) {
	c.logger.Debug("fetching static projects")
	return c.projects, nil
}
