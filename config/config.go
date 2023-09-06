package config

import (
	checkstypes "standards/checks/aggregates"
	rulestypes "standards/rules/aggregates"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Providers ProvidersConfig
	Groups    []checkstypes.Group
	Checks    []checkstypes.Check
	Rules     []rulestypes.Rule
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
