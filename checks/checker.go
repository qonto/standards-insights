package checks

import (
	"context"

	"github.com/qonto/standards-insights/checks/aggregates"
	providerstypes "github.com/qonto/standards-insights/providers/aggregates"
	rulestypes "github.com/qonto/standards-insights/rules/aggregates"
)

type Ruler interface {
	Execute(ctx context.Context, ruleName string, project providerstypes.Project) rulestypes.RuleResult
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

func (c *Checker) Run(ctx context.Context, projects []providerstypes.Project) []aggregates.ProjectResult {
	projectResults := make([]aggregates.ProjectResult, len(projects))
	for i, project := range projects {
		projectResult := aggregates.ProjectResult{
			Labels:       project.Labels,
			Name:         project.Name,
			CheckResults: []aggregates.CheckResult{},
		}
		for _, group := range c.groups {
			if c.shouldSkipGroup(ctx, group, project) {
				continue
			}
			checkResults := c.executeGroup(ctx, group, project)
			projectResult.CheckResults = append(projectResult.CheckResults, checkResults...)
		}
		projectResults[i] = projectResult
	}
	return projectResults
}
