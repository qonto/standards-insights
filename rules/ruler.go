package rules

import (
	"context"
	"fmt"

	providerstypes "github.com/qonto/standards-insights/providers/aggregates"
	"github.com/qonto/standards-insights/rules/aggregates"
)

type Ruler struct {
	rules map[string]aggregates.Rule
}

func NewRuler(rules []aggregates.Rule) *Ruler {
	rulesMap := make(map[string]aggregates.Rule)

	for _, rule := range rules {
		rulesMap[rule.Name] = rule
	}

	return &Ruler{
		rules: rulesMap,
	}
}

func (r *Ruler) Execute(ctx context.Context, ruleName string, project providerstypes.Project) aggregates.RuleResult {
	rule, ok := r.rules[ruleName]
	if !ok {
		return aggregates.RuleResult{
			Success:  false,
			RuleName: ruleName,
			Messages: []string{fmt.Sprintf("rule %s not found in the rules configuration", ruleName)},
		}
	}
	errorMessages := []string{}
	err := executeFileRules(ctx, rule, project)
	if err != nil {
		errorMessages = append(errorMessages, err.Error())
	}
	if len(errorMessages) > 0 {
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
