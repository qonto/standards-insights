package discovery

import "standards/config"

type Discovery struct {
	discovery config.ConfigDiscovery
	iterator  *iterator
}

func New(config *config.Config) *Discovery {
	return &Discovery{
		discovery: config.Discovery,
		iterator: &iterator{
			index: -1,
		},
	}
}
