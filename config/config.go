package config

import (
	"encoding/json"
	"fmt"
	"os"
	"standards/rules/aggregate"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Discovery ConfigDiscovery
	Groups    []*aggregate.Group

	// Should stay private because only groups is needed in public
	rules  []*aggregate.Rule
	checks []*aggregate.Check
}

type rawConfig struct {
	Discovery ConfigDiscovery
	Groups    []Group
	Checks    []Check
	Rules     []Rule
}

type Group struct {
	Name   string
	Checks []string
	Rules  []string
}

type Check struct {
	Name     string
	Category string
	Level    string
	Owner    string
	Exclude  []string
	Include  []string
	Rules    []string
}

type Rule struct {
	Name  string
	Files []aggregate.FileRule
}

type ConfigCheck struct {
	Name string
}

type ConfigDiscovery struct {
	ArgoCD ArgoCDConfig
}

type ArgoCDConfig struct {
	URL string
}

func New(path string) (*Config, error) {
	var raw rawConfig

	file, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("Could not find config file: %v", err)
	}
	err = yaml.Unmarshal(file, &raw)
	if err != nil {
		return nil, err
	}

	// Init new struct
	config := &Config{
		Discovery: raw.Discovery,
		Groups:    make([]*aggregate.Group, len(raw.Groups)),
		rules:     make([]*aggregate.Rule, len(raw.Rules)),
		checks:    make([]*aggregate.Check, len(raw.Checks)),
	}

	// Hydrate rules to use them as pointers in checks & groups
	for i, rule := range raw.Rules {
		config.rules[i] = &aggregate.Rule{
			Name:  rule.Name,
			Files: rule.Files,
		}
	}

	// Hydrate checks with rules pointers to use it in groups
	for i, rawCheck := range raw.Checks {
		rules := make([]*aggregate.Rule, len(rawCheck.Rules))
		for y, ruleName := range rawCheck.Rules {
			rule, err := config.findRule(ruleName)
			if err != nil {
				return config, err
			}
			rules[y] = rule
		}

		config.checks[i] = &aggregate.Check{
			Name:     rawCheck.Name,
			Category: rawCheck.Category,
			Level:    rawCheck.Level,
			Exclude:  rawCheck.Exclude,
			Include:  rawCheck.Include,
			Rules:    rules,
		}
	}

	// Hydrate public groups with checks & rules pointers
	for i, rawGroup := range raw.Groups {

		// Re-use rules
		rules := make([]*aggregate.Rule, len(rawGroup.Rules))
		for y, ruleName := range rawGroup.Rules {
			rule, err := config.findRule(ruleName)
			if err != nil {
				return config, err
			}
			rules[y] = rule
		}

		// Re-use checks
		checks := make([]*aggregate.Check, len(rawGroup.Checks))
		for y, checkName := range rawGroup.Checks {
			check, err := config.findCheck(checkName)
			if err != nil {
				return config, err
			}
			checks[y] = check
		}

		config.Groups[i] = &aggregate.Group{
			Name:   rawGroup.Name,
			Checks: checks,
			Rules:  rules,
		}
	}

	return config, nil
}

func (c *Config) findRule(name string) (*aggregate.Rule, error) {
	for _, rule := range c.rules {
		if rule.Name == name {
			return rule, nil
		}
	}

	return nil, fmt.Errorf("Unable to find rule %s", name)
}

func (c *Config) findCheck(name string) (*aggregate.Check, error) {
	for _, check := range c.checks {
		if check.Name == name {
			return check, nil
		}
	}

	return nil, fmt.Errorf("Unable to find check %s", name)
}

func (c *Config) String() string {
	b, err := json.MarshalIndent(c, "", "    ")
	if err != nil {
		return "Could json marshal"
	}
	return string(b)
}
