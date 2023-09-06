package checks

import (
	"context"
	"fmt"

	"standards/checks/aggregates"
	providerstypes "standards/providers/aggregates"
	rulestypes "standards/rules/aggregates"
)

func (c *Checker) executeCheck(ctx context.Context, check aggregates.Check) aggregates.CheckResult {
	success := true
	checkResult := aggregates.CheckResult{
		Name:    check.Name,
		Labels:  check.Labels,
		Results: []rulestypes.RuleResult{},
	}
	for _, rule := range check.Rules {
		ruleResult := c.ruler.Execute(ctx, rule)
		if !ruleResult.Success {
			success = false
		}
		checkResult.Results = append(checkResult.Results, ruleResult)
	}
	checkResult.Success = success
	return checkResult
}

func (c *Checker) shouldSkipGroup(ctx context.Context, group aggregates.Group) bool {
	for _, rule := range group.When {
		ruleResult := c.ruler.Execute(ctx, rule)
		if !ruleResult.Success {
			return true
		}
	}
	return false
}

func (c *Checker) executeGroup(ctx context.Context, group aggregates.Group) []aggregates.CheckResult {
	result := []aggregates.CheckResult{}
	for _, checkName := range group.Checks {
		// TODO check exists
		check := c.checks[checkName]
		checkResult := c.executeCheck(ctx, check)
		result = append(result, checkResult)
	}
	return result
}

func (c *Checker) Run(ctx context.Context, projects []providerstypes.Project) error {
	projectResults := make([]aggregates.ProjectResult, len(projects))
	for i, project := range projects {
		// TODO
		// the project should be passed to every layers because the rules should be executed for
		// each project
		projectResult := aggregates.ProjectResult{
			Name:         project.Name,
			CheckResults: []aggregates.CheckResult{},
		}
		for _, group := range c.groups {
			if c.shouldSkipGroup(ctx, group) {
				continue
			}
			checkResults := c.executeGroup(ctx, group)
			projectResult.CheckResults = append(projectResult.CheckResults, checkResults...)
		}
		projectResults[i] = projectResult
	}
	for _, project := range projectResults {
		fmt.Printf("== Project %s\n", project.Name)
		for _, result := range project.CheckResults {
			if result.Success {
				fmt.Printf("✅ Check %s PASS (labels: %s)\n", result.Name, result.Labels)
			} else {
				fmt.Printf("🚨 Check %s FAILED (labels: %s)\n", result.Name, result.Labels)
				fmt.Printf("🚨 %+v\n", result)
			}
		}
	}
	return nil
}
