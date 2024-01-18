package checker

import (
	"context"
	"fmt"

	"github.com/qonto/standards-insights/config"

	"github.com/qonto/standards-insights/pkg/checker/aggregates"
	"github.com/qonto/standards-insights/pkg/project"
	ruleraggregates "github.com/qonto/standards-insights/pkg/ruler/aggregates"
)

func (c *Checker) executeCheck(ctx context.Context, check config.Check, project project.Project) aggregates.CheckResult {
	if check.IsAND() {
		return c.executeANDCheck(ctx, check, project)
	} else {
		return c.executeORCheck(ctx, check, project)
	}
}

func (c *Checker) executeANDCheck(ctx context.Context, check config.Check, project project.Project) aggregates.CheckResult {
	success := true
	checkResult := aggregates.CheckResult{
		Name:    check.Name,
		Labels:  check.Labels,
		Results: []ruleraggregates.RuleResult{},
	}

	for _, rule := range check.Rules {
		ruleResult := c.ruler.Execute(ctx, rule, project)
		if !ruleResult.Success {
			c.logger.Debug(fmt.Sprintf("rule %s failed on project %s for check %s", rule, project.Name, check.Name))
			success = false
		} else {
			c.logger.Debug(fmt.Sprintf("rule %s successful on project %s for check %s", rule, project.Name, check.Name))
		}
		checkResult.Results = append(checkResult.Results, ruleResult)
	}
	checkResult.Success = success
	return checkResult
}

func (c *Checker) executeORCheck(ctx context.Context, check config.Check, project project.Project) aggregates.CheckResult {
	success := false
	checkResult := aggregates.CheckResult{
		Name:    check.Name,
		Labels:  check.Labels,
		Results: []ruleraggregates.RuleResult{},
	}

	for _, rule := range check.Rules {
		ruleResult := c.ruler.Execute(ctx, rule, project)
		if ruleResult.Success {
			success = true
			c.logger.Debug(fmt.Sprintf("rule %s successful on project %s for check %s", rule, project.Name, check.Name))
		} else {
			c.logger.Debug(fmt.Sprintf("rule %s failed on project %s for check %s", rule, project.Name, check.Name))
		}
		checkResult.Results = append(checkResult.Results, ruleResult)
	}
	checkResult.Success = success
	return checkResult
}
