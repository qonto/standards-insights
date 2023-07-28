package config

import "standards/rules/aggregate"

type Config struct {
	Rules  []aggregate.Rule
	Groups []aggregate.Group
}
