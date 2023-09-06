package providers

import (
	"github.com/qonto/standards-insights/config"
	"github.com/qonto/standards-insights/providers/aggregates"
)

type Client struct{}

func NewProviders(config *config.ProvidersConfig, filters []string) ([]aggregates.Provider, error) {
	// TODO: fixme
	return []aggregates.Provider{}, nil
}
