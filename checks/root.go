package checks

import (
	"context"
	"standards/rules/aggregate"
)

type Ruler interface {
	Execute(ctx context.Context, ruleName string) aggregate.RuleResult
}

type Checker struct {
	ruler  Ruler
	checks map[string]aggregate.Check
	groups []aggregate.Group
}

func NewChecker(ruler Ruler, checks []aggregate.Check, groups []aggregate.Group) *Checker {

	checksMap := make(map[string]aggregate.Check)
	for _, check := range checks {
		checksMap[check.Name] = check
	}
	return &Checker{
		ruler:  ruler,
		checks: checksMap,
		groups: groups,
	}
}
