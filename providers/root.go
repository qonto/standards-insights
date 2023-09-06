package providers

import (
	"standards/config"
	"standards/providers/aggregates"
)

type Client struct{}

func NewProviders(config *config.ProvidersConfig, filters []string) ([]aggregates.Provider, error) {
	// TODO: fixme
	return []aggregates.Provider{}, nil
}
