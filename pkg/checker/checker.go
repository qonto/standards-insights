package checker

import (
	"context"
	"fmt"

	"log/slog"

	"github.com/qonto/standards-insights/config"
	"github.com/qonto/standards-insights/pkg/checker/aggregates"
	"github.com/qonto/standards-insights/pkg/project"
	ruleraggregates "github.com/qonto/standards-insights/pkg/ruler/aggregates"
)

type Ruler interface {
	Execute(ctx context.Context, ruleName string, project project.Project) ruleraggregates.RuleResult
}

type Checker struct {
	logger *slog.Logger
	ruler  Ruler
	checks map[string]config.Check
	groups []config.Group
}

func NewChecker(logger *slog.Logger, ruler Ruler, checks []config.Check, groups []config.Group) *Checker {
	checksMap := make(map[string]config.Check)
	for _, check := range checks {
		checksMap[check.Name] = check
	}
	return &Checker{
		logger: logger,
		ruler:  ruler,
		checks: checksMap,
		groups: groups,
	}
}

func (c *Checker) Run(ctx context.Context, projects []project.Project) []aggregates.ProjectResult {
	projectResults := make([]aggregates.ProjectResult, len(projects))
	for _, project := range projects {
		c.logger.Info("checking project " + project.Name)
		projectResult := aggregates.ProjectResult{
			Labels:       project.Labels,
			Name:         project.Name,
			CheckResults: []aggregates.CheckResult{},
		}
		for _, group := range c.groups {
			fmt.Printf("group %+v", group)
			if c.shouldSkipGroup(ctx, group, project) {
				c.logger.Info(fmt.Sprintf("skipping group %s for project %s", group.Name, project.Name))
				continue
			}
			checkResults := c.executeGroup(ctx, group, project)
			projectResult.CheckResults = append(projectResult.CheckResults, checkResults...)

			if group.ApplyOnSubProjects {
				// print group name, project name and apply on sub projects
				c.logger.Info(fmt.Sprintf("applying group %s for project %s and sub projects", group.Name, project.Name))
				for _, subProject := range project.SubProjects {
					if c.shouldSkipGroup(ctx, group, subProject) {
						c.logger.Info(fmt.Sprintf("skipping group %s for subproject %s", group.Name, subProject.SubProject))
						continue
					}
					subProjectResult := aggregates.ProjectResult{
						Labels:       subProject.Labels,
						Name:         subProject.Name,
						Subproject:   subProject.SubProject,
						CheckResults: []aggregates.CheckResult{},
					}
					subProjectCheckResults := c.executeGroup(ctx, group, subProject)
					subProjectResult.CheckResults = append(subProjectResult.CheckResults, subProjectCheckResults...)
					projectResults = append(projectResults, subProjectResult)
				}
			}
		}
		projectResults = append(projectResults, projectResult)
	}
	return projectResults
}
