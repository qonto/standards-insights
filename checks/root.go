package checks

import (
	"context"

	"standards/checks/aggregates"
	rulestypes "standards/rules/aggregates"
)

type Ruler interface {
	Execute(ctx context.Context, ruleName string) rulestypes.RuleResult
}

type Checker struct {
	ruler  Ruler
	checks map[string]aggregates.Check
	groups []aggregates.Group
}

func NewChecker(ruler Ruler, checks []aggregates.Check, groups []aggregates.Group) *Checker {
	checksMap := make(map[string]aggregates.Check)
	for _, check := range checks {
		checksMap[check.Name] = check
	}
	return &Checker{
		ruler:  ruler,
		checks: checksMap,
		groups: groups,
	}
}
