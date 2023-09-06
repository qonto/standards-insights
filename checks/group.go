package checks

import (
	"context"

	"github.com/qonto/standards-insights/checks/aggregates"
	providerstypes "github.com/qonto/standards-insights/providers/aggregates"
)

func (c *Checker) shouldSkipGroup(ctx context.Context, group aggregates.Group, project providerstypes.Project) bool {
	for _, rule := range group.When {
		ruleResult := c.ruler.Execute(ctx, rule, project)
		if !ruleResult.Success {
			return true
		}
	}
	return false
}

func (c *Checker) executeGroup(ctx context.Context, group aggregates.Group, project providerstypes.Project) []aggregates.CheckResult {
	result := []aggregates.CheckResult{}
	for _, checkName := range group.Checks {
		// TODO check exists
		check := c.checks[checkName]
		checkResult := c.executeCheck(ctx, check, project)
		result = append(result, checkResult)
	}
	return result
}
