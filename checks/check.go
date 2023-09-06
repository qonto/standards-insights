package checks

import (
	"context"

	providerstypes "github.com/qonto/standards-insights/providers/aggregates"
	rulestypes "github.com/qonto/standards-insights/rules/aggregates"

	"github.com/qonto/standards-insights/checks/aggregates"
)

func (c *Checker) executeCheck(ctx context.Context, check aggregates.Check, project providerstypes.Project) aggregates.CheckResult {
	success := true
	checkResult := aggregates.CheckResult{
		Name:    check.Name,
		Labels:  check.Labels,
		Results: []rulestypes.RuleResult{},
	}
	for _, rule := range check.Rules {
		ruleResult := c.ruler.Execute(ctx, rule, project)
		if !ruleResult.Success {
			success = false
		}
		checkResult.Results = append(checkResult.Results, ruleResult)
	}
	checkResult.Success = success
	return checkResult
}
