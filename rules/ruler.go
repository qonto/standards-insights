package rules

import (
	"context"
	"fmt"
	"standards/rules/aggregate"
)

type Ruler struct {
	rules map[string]aggregate.Rule
}

func NewRuler(rules []aggregate.Rule) *Ruler {
	rulesMap := make(map[string]aggregate.Rule)

	for _, rule := range rules {
		rulesMap[rule.Name] = rule
	}

	return &Ruler{
		rules: rulesMap,
	}
}

func (r *Ruler) Execute(ctx context.Context, ruleName string) aggregate.RuleResult {
	rule, ok := r.rules[ruleName]
	if !ok {
		return aggregate.RuleResult{
			Success:  false,
			RuleName: ruleName,
			Messages: []string{fmt.Sprintf("rule %s not found in the rules configuration", ruleName)},
		}
	}
	errorMessages := []string{}
	err := executeFileRules(ctx, rule)
	if err != nil {
		errorMessages = append(errorMessages, err.Error())
	}
	if len(errorMessages) > 0 {
		return aggregate.RuleResult{
			Success:  false,
			RuleName: ruleName,
			Messages: errorMessages,
		}
	}
	return aggregate.RuleResult{
		Success:  true,
		RuleName: ruleName,
	}
}
