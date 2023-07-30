package checks

import (
	"fmt"
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

func (c *CheckProcessor) Run() error {
	fmt.Println("Syncing all projects...")
	err := c.discovery.SyncProjects()
	if err != nil {
		return err
	}
	fmt.Println("Done!")

	checkResults := []aggregate.CheckResults{}

	for true {
		if !c.discovery.HasNext() {
			break
		}

		project := c.discovery.GetNext()
		fmt.Printf("💡 Checking project '%s' against groups\n", project.Name)

		for _, group := range c.config.Groups {
			fmt.Printf("===== Group: %s =====\n", group.Name)
			if !group.IsMatchingRules() { // TODO: error handling (bool is not enough)
				fmt.Printf("[DEBUG]Group %s is not matching rules => ignore\n", group.Name)
				continue
			}

			for _, check := range group.Checks {
				fmt.Printf("Running check %s for project %s\n", check.Name, project.Name)

				checkResults = append(checkResults, aggregate.CheckResults{
					RepositoryName: project.Name,
					CheckName:      check.Name,
					Success:        check.IsMatchingRules(),
				})
			}
		}

		for _, result := range checkResults {
			if result.Success {
				fmt.Printf("✅ Repo: %s checked %s\n", result.RepositoryName, result.CheckName)
			} else {
				fmt.Printf("🚨 Repo: %s checked %s\n", result.RepositoryName, result.CheckName)
			}
		}
	}

	return nil
}
