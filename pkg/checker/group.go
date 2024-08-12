package checker

import (
	"context"

	"github.com/qonto/standards-insights/config"
	"github.com/qonto/standards-insights/pkg/checker/aggregates"
	"github.com/qonto/standards-insights/pkg/project"
)

func contains(slice []string, item string) bool {
    for _, s := range slice {
        if s == item {
            return true
        }
    }
    return false
}

func (c *Checker) shouldSkipGroup(ctx context.Context, group config.Group, project project.Project) bool {
	for _, rule := range group.When {
		ruleResult := c.ruler.Execute(ctx, rule, project)
		if !ruleResult.Success {
			return true
		}
	}
	return false
}

func (c *Checker) executeGroup(ctx context.Context, group config.Group, project project.Project) []aggregates.CheckResult {
	result := []aggregates.CheckResult{}
	for _, checkName := range group.Checks {
		// For now let's consider that we checked when the config is built
		// that checks always exist
		check := c.checks[checkName]
		checkResult := c.executeCheck(ctx, check, project)
		result = append(result, checkResult)
	}
	return result
}