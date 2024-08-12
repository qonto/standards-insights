package checker

import (
	"context"
	"fmt"
	"regexp"

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
			if c.shouldSkipGroup(ctx, group, project) {
				c.logger.Debug(fmt.Sprintf("skipping group %s for project %s", group.Name, project.Name))
				continue
			}
			checkResults := c.executeGroup(ctx, group, project)
			projectResult.CheckResults = append(projectResult.CheckResults, checkResults...)

			if group.Files.ApplyToFiles {
				// use group pattern to filter sub projects
				filteredSubProjects := filterSubProjects(project.SubProjects, group.Files.FilesPattern)
				// print group name, project name and apply on sub projects
				c.logger.Debug(fmt.Sprintf("applying group %s for project %s and sub projects", group.Name, project.Name))
				for _, subProject := range filteredSubProjects {
					if c.shouldSkipGroup(ctx, group, subProject) {
						c.logger.Debug(fmt.Sprintf("skipping group %s for file %s", group.Name, subProject.FilePath))
						continue
					}
					subProjectResult := aggregates.ProjectResult{
						Labels:       subProject.Labels,
						Name:         subProject.Name,
						FilePath:     subProject.FilePath,
						CheckResults: []aggregates.CheckResult{},
					}
					subProjectCheckResults := c.executeGroup(ctx, group, subProject)
					subProjectResult.CheckResults = append(subProjectResult.CheckResults, subProjectCheckResults...)
					projectResults = append(projectResults, subProjectResult)
				}
			}
		}
		projectResults[0] = projectResult
	}
	return projectResults
}

func filterSubProjects(projects []project.Project, pattern string) []project.Project {
	var filteredProjects []project.Project
	re := regexp.MustCompile(pattern)
	for _, proj := range projects {
		if re.MatchString(proj.FilePath) {
			filteredProjects = append(filteredProjects, proj)
		}
	}
	return filteredProjects
}
