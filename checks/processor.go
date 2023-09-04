package checks

import (
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

func (c *CheckProcessor) Run(projects []providers.Project) error {
	projectResults := make([]*aggregate.ProjectResult, len(projects))
	for i, project := range projects {
		projectResult := &aggregate.ProjectResult{
			Name:         project.Name,
			CheckResults: []aggregate.CheckResult{},
		}
		projectResults[i] = projectResult
		fmt.Printf("💡 Checking project '%s' against groups\n", project.Name)

		for _, group := range c.config.Groups {
			fmt.Printf("===== Group: %s =====\n", group.Name)
			if !group.IsMatchingRules() { // TODO: error handling (bool is not enough)
				fmt.Printf("[DEBUG]Group %s is not matching rules => ignore\n", group.Name)
				continue
			}

			for _, check := range group.Checks {
				fmt.Printf("Running check %s for project %s\n", check.Name, project.Name)

				projectResult.CheckResults = append(projectResult.CheckResults, aggregate.CheckResult{
					Check:   check,
					Success: check.IsMatchingRules(),
				})
			}
		}
	}

	fmt.Println("\n\nResults:")
	for _, project := range projectResults {
		fmt.Printf("== Project %s\n", project.Name)
		for _, result := range project.CheckResults {
			if result.Success {
				fmt.Printf("✅ Check %s PASS (level: %s, category: %s)\n", result.Check.Name, result.Check.Level, result.Check.Category)
			} else {
				fmt.Printf("🚨 Check %s FAILED (level: %s, category: %s)\n", result.Check.Name, result.Check.Level, result.Check.Category)
			}
		}
	}

	return nil
}
