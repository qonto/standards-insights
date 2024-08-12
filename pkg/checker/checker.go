package checker

import (
	"context"
	"fmt"
	"regexp"
	"time"

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
			// measure time to skip group
			start := time.Now()
			if c.shouldSkipGroup(ctx, group, project) {
				c.logger.Info(fmt.Sprintf("skipping group %s for project %s", group.Name, project.Name))
				continue
			}
			c.logger.Info(fmt.Sprintf("time to skip group %s for project %s: %s", group.Name, project.Name, time.Since(start)))
			// measure time to execute group
			start = time.Now()
			checkResults := c.executeGroup(ctx, group, project)
			c.logger.Info(fmt.Sprintf("time to execute group %s for project %s: %s", group.Name, project.Name, time.Since(start)))
			projectResult.CheckResults = append(projectResult.CheckResults, checkResults...)

			if group.ApplyOnSubProjects {
				// use group pattern to filter sub projects
				filteredSubProjects := filterSubProjects(project.SubProjects, group.SubprojectsPattern)
				// print group name, project name and apply on sub projects
				c.logger.Info(fmt.Sprintf("applying group %s for project %s and sub projects", group.Name, project.Name))
				totalSubProjectExecutionTime := time.Duration(0)
				totalSkipTime := time.Duration(0) // Initialize total skip time
				for _, subProject := range filteredSubProjects {
					// measure time to execute group on sub project
					start := time.Now()
					if c.shouldSkipGroup(ctx, group, subProject) {
						skipTime := time.Since(start)
						totalSkipTime += skipTime // Accumulate skip time
						c.logger.Info(fmt.Sprintf("skipping group %s for subproject %s", group.Name, subProject.SubProject))
						continue
					}
					skipTime := time.Since(start)
					totalSkipTime += skipTime // Accumulate skip time
					//c.logger.Info(fmt.Sprintf("time to skip group %s for subproject %s: %s", group.Name, subProject.SubProject, skipTime))
					subProjectResult := aggregates.ProjectResult{
						Labels:       subProject.Labels,
						Name:         subProject.Name,
						Subproject:   subProject.SubProject,
						CheckResults: []aggregates.CheckResult{},
					}
					// measure time to execute group on sub project
					start = time.Now()
					subProjectCheckResults := c.executeGroup(ctx, group, subProject)
					executionTime := time.Since(start)
					totalSubProjectExecutionTime += executionTime
					//c.logger.Info(fmt.Sprintf("time to execute group %s for subproject %s: %s", group.Name, subProject.SubProject, executionTime))
					subProjectResult.CheckResults = append(subProjectResult.CheckResults, subProjectCheckResults...)
					projectResults = append(projectResults, subProjectResult)
				}
				c.logger.Info(fmt.Sprintf("total time to skip group %s across all subprojects for project %s: %s", group.Name, project.Name, totalSkipTime))
				c.logger.Info(fmt.Sprintf("total time to execute group %s across all subprojects for project %s: %s", group.Name, project.Name, totalSubProjectExecutionTime))
			}
		}
		projectResults = append(projectResults, projectResult)
	}
	return projectResults
}

func filterSubProjects(projects []project.Project, pattern string) []project.Project {
	var filteredProjects []project.Project
	re := regexp.MustCompile(pattern)
	for _, proj := range projects {
		if re.MatchString(proj.SubProject) {
			filteredProjects = append(filteredProjects, proj)
		}
	}
	return filteredProjects
}