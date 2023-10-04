package config

import (
	"github.com/qonto/standards-insights/pkg/project"
	"github.com/qonto/standards-insights/types"

	"gopkg.in/yaml.v3"
)

type Config struct {
	HTTP      HTTPConfig      `validate:"omitempty"`
	Providers ProvidersConfig `validate:"omitempty"`
	Groups    []Group
	Checks    []Check
	Rules     []Rule
	Labels    []string
	Interval  int
	Git       GitConfig
}

type ProvidersConfig struct {
	ArgoCD ArgoCDConfig      `validate:"omitempty"`
	Static []project.Project `validate:"omitempty"`
}

type ArgoCDConfig struct {
	URL      string `validate:"required"`
	Projects []string
	Selector string
	BasePath string `yaml:"base-path"`
}

type Rule struct {
	Name   string `validate:"required"`
	Files  []FileRule
	Simple *bool
}

type FileRule struct {
	Path        string `validate:"required"`
	Contains    *types.Regexp
	NotContains *types.Regexp `yaml:"not-contains"`
	Exists      *bool
}

type Check struct {
	Name   string `validate:"required"`
	Labels map[string]string
	Rules  []string `validate:"required,min=1"`
}

type Group struct {
	Name   string   `validate:"required"`
	Checks []string `validate:"required,min=1"`
	When   []string
}

func New(path string) (*Config, []byte, error) {
	var config Config

	content, err := getConfigYaml(path)
	if err != nil {
		return nil, nil, err
	}

	err = yaml.Unmarshal(content, &config)
	if err != nil {
		return nil, nil, err
	}
	err = validate(config)
	if err != nil {
		return nil, nil, err
	}
	return &config, content, nil
}
