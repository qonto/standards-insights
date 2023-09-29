package aggregates

import (
	"context"

	"github.com/qonto/standards-insights/pkg/project"
)

type Provider interface {
	FetchProjects(ctx context.Context) ([]project.Project, error)
	Name() string
}
