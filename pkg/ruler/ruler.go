package ruler

import (
	"context"
	"fmt"

	"log/slog"

	"github.com/qonto/standards-insights/config"
	"github.com/qonto/standards-insights/pkg/project"
	aggregates "github.com/qonto/standards-insights/pkg/ruler/aggregates"
)

type Ruler struct {
	logger *slog.Logger
	rules  map[string]*rule
}

func NewRuler(logger *slog.Logger, rulesConfig []config.Rule) *Ruler {
	rulesMap := make(map[string]*rule)
	for i := range rulesConfig {
		ruleConfig := rulesConfig[i]
		rulesMap[ruleConfig.Name] = newRule(ruleConfig)
	}
	return &Ruler{
		logger: logger,
		rules:  rulesMap,
	}
}

func (r *Ruler) Execute(ctx context.Context, ruleName string, project project.Project) aggregates.RuleResult {
	r.logger.Debug(fmt.Sprintf("executing rule %s on project %s", ruleName, project.Name))
	rule, ok := r.rules[ruleName]
	if !ok {
		return aggregates.RuleResult{
			Success:  false,
			RuleName: ruleName,
			Messages: []string{fmt.Sprintf("rule %s not found in the rules configuration", ruleName)},
		}
	}
	errorMessages := []string{}

	for _, module := range rule.Modules {
		err := module.Do(ctx, project)
		if err != nil {
			errorMessages = append(errorMessages, err.Error())
		}
	}
	if len(errorMessages) > 0 {
		for _, errorMessage := range errorMessages {
			r.logger.Error(fmt.Sprintf("error on rule %s for project %s subproject %s: %s", ruleName, project.Name, project.SubProject, errorMessage))
		}
		return aggregates.RuleResult{
			Success:  false,
			RuleName: ruleName,
			Messages: errorMessages,
		}
	}
	return aggregates.RuleResult{
		Success:  true,
		RuleName: ruleName,
	}
}
