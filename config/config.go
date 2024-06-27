package config

import (
	"github.com/qonto/standards-insights/pkg/project"
	"github.com/qonto/standards-insights/types"

	"gopkg.in/yaml.v3"
)

type Config struct {
	HTTP      HTTPConfig      `validate:"omitempty"`
	Providers ProvidersConfig `validate:"omitempty"`
	Groups    []Group         `validate:"dive"`
	Checks    []Check         `validate:"dive"`
	Rules     []Rule          `validate:"dive"`
	Labels    []string        `validate:"dive"`
	Interval  int
	Git       GitConfig
}

type ProvidersConfig struct {
	ArgoCD ArgoCDConfig      `validate:"omitempty"`
	Static []project.Project `validate:"omitempty,dive"`
	Gitlab GitlabConfig      `validate:"omitempty"`
}

type ArgoCDConfig struct {
	URL      string `validate:"required"`
	Projects []string
	Selector string
	BasePath string `yaml:"base-path"`
}

type GitlabConfig struct {
	URL      string `validate:"required"`
	Token    string
	Topics   []string
	Search   string
	BasePath string `yaml:"base-path"`
}

type Rule struct {
	Name    string        `validate:"required"`
	Files   []FileRule    `validate:"dive"`
	Grep    []GrepRule    `validate:"dive"`
	Project []ProjectRule `validate:"dive"`
	Simple  *bool
}

type FileRule struct {
	Path        string `validate:"required"`
	Contains    *types.Regexp
	NotContains *types.Regexp `yaml:"not-contains"`
	Exists      *bool
}

type GrepRule struct {
	Path           string `validate:"required"`
	Pattern        string `validate:"required"`
	Include        string
	ExtendedRegexp bool `yaml:"extended-regexp"`
	Recursive      bool
	Match          bool
	SkipNotFound   bool `yaml:"skip-not-found"`
}

type ProjectRule struct {
	Names  []string
	Labels map[string]string
	Match  *bool
}

type Check struct {
	Name     string `validate:"required"`
	Labels   map[string]string
	Operator string   `validate:"oneof='and' 'or' ''"`
	Rules    []string `validate:"required,min=1"`
}

func (c Check) IsAND() bool {
	return c.Operator == "" || c.Operator == "and"
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
