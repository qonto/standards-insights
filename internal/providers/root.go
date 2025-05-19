package providers

import (
	"log/slog"

	"github.com/qonto/standards-insights/config"
	"github.com/qonto/standards-insights/internal/providers/aggregates"
	"github.com/qonto/standards-insights/internal/providers/argocd"
	"github.com/qonto/standards-insights/internal/providers/github"
	"github.com/qonto/standards-insights/internal/providers/gitlab"
	"github.com/qonto/standards-insights/internal/providers/static"
)

func NewProviders(logger *slog.Logger, config config.ProvidersConfig) ([]aggregates.Provider, error) {
	result := []aggregates.Provider{}
	if config.ArgoCD.URL != "" {
		argoProvider, err := argocd.New(logger, config.ArgoCD)
		if err != nil {
			return nil, err
		}
		result = append(result, argoProvider)
	}
	if config.Github.URL != "" {
		githubProvider, err := github.New(logger, config.Github)
		if err != nil {
			return nil, err
		}
		result = append(result, githubProvider)
	}
	if config.Gitlab.URL != "" {
		gitlabProvider, err := gitlab.New(logger, config.Gitlab)
		if err != nil {
			return nil, err
		}
		result = append(result, gitlabProvider)
	}
	if len(config.Static) != 0 {
		staticProvider := static.New(logger, config.Static)
		result = append(result, staticProvider)
	}
	return result, nil
}
