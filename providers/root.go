package providers

import "standards/config"

type Client struct {
}

func NewProviders(config *config.ProvidersConfig, filters []string) ([]Provider, error) {
	// TODO: fixme
	return []Provider{}, nil
}
