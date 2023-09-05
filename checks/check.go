package checks

import (
	"context"
	"fmt"
	"standards/providers"
	"standards/rules/aggregate"
)

type ProjectResult struct {
	Name        string
	CheckErrors []CheckError
}

type CheckError struct {
	Name    string
	Message string
}

func (c *Checker) executeCheck(ctx context.Context, check aggregate.Check) aggregate.CheckResult {
	success := true
	checkResult := aggregate.CheckResult{
		Name:    check.Name,
		Labels:  check.Labels,
		Results: []aggregate.RuleResult{},
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

func (c *Checker) skipGroup(ctx context.Context, group aggregate.Group) bool {
	for _, rule := range group.Rules {
		ruleResult := c.ruler.Execute(ctx, rule)
		if !ruleResult.Success {
			return true
		}
	}
	return false
}

func (c *Checker) executeGroup(ctx context.Context, group aggregate.Group) []aggregate.CheckResult {
	result := []aggregate.CheckResult{}
	for _, checkName := range group.Checks {
		// TODO check exists
		check := c.checks[checkName]
		checkResult := c.executeCheck(ctx, check)
		result = append(result, checkResult)
	}
	return result
}

func (c *Checker) Run(ctx context.Context, projects []providers.Project) error {
	projectResults := []aggregate.ProjectResult{}
	for _, project := range projects {
		// TODO
		// the project should be passed to every layers because the rules should be executed for
		// each project
		projectResult := aggregate.ProjectResult{
			Name:         project.Name,
			CheckResults: []aggregate.CheckResult{},
		}
		for _, group := range c.groups {
			if c.skipGroup(ctx, group) {
				continue
			}
			checkResults := c.executeGroup(ctx, group)
			projectResult.CheckResults = append(projectResult.CheckResults, checkResults...)
		}
		projectResults = append(projectResults, projectResult)
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
