package config

import (
	"standards/rules/aggregate"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Providers ProvidersConfig
	Groups    []aggregate.Group
	Checks    []aggregate.Check
	Rules     []aggregate.Rule
}

type ProvidersConfig struct {
	ArgoCD ArgoCDConfig
	Git    GitConfig
}

type GitConfig []string

type ArgoCDConfig struct {
	URL string
}

func New(path string) (*Config, error) {
	var config Config

	content, err := getConfigYaml(path)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(content, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
