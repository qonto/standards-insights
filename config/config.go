package config

import (
	checkstypes "github.com/qonto/standards-insights/checks/aggregates"
	"github.com/qonto/standards-insights/git"
	"github.com/qonto/standards-insights/http"
	providerstypes "github.com/qonto/standards-insights/providers/aggregates"
	rulestypes "github.com/qonto/standards-insights/rules/aggregates"

	"gopkg.in/yaml.v3"
)

type Config struct {
	HTTP      http.Config
	Providers ProvidersConfig
	Groups    []checkstypes.Group
	Checks    []checkstypes.Check
	Rules     []rulestypes.Rule
	Labels    []string
	Interval  int
	Git       git.Config
}

type ProvidersConfig struct {
	ArgoCD ArgoCDConfig
	Static []providerstypes.Project
}

type GitConfig []string

type ArgoCDConfig struct {
	URL      string
	Projects []string
	Selector string
	BasePath string `yaml:"base-path"`
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
