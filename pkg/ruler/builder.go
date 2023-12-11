package ruler

import (
	"context"

	"github.com/qonto/standards-insights/config"
	"github.com/qonto/standards-insights/pkg/project"
	"github.com/qonto/standards-insights/pkg/ruler/rules"
)

type ruleModule interface {
	Do(ctx context.Context, project project.Project) error
}

type rule struct {
	Name    string
	Modules []ruleModule
}

func newRule(config config.Rule) *rule {
	result := &rule{
		Name: config.Name,
	}
	modules := []ruleModule{}
	for _, fileConfig := range config.Files {
		modules = append(modules, rules.NewFileRule(fileConfig))
	}
	for _, grepConfig := range config.Grep {
		modules = append(modules, rules.NewGrepRule(grepConfig))
	}
	if config.Simple != nil {
		value := *config.Simple
		modules = append(modules, rules.NewSimpleRule(value))
	}
	result.Modules = modules
	return result
}
