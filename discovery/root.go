package discovery

import "standards/config"

type Discovery struct {
	discovery config.ConfigDiscovery
}

func New(config *config.Config) *Discovery {
	return &Discovery{
		discovery: config.Discovery,
	}
}
