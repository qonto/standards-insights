package cmd

import (
	"fmt"
	"strings"

	"github.com/qonto/standards-insights/pkg/checker/aggregates"
)

func stdoutResults(results []aggregates.ProjectResult) {
	for _, project := range results {
		fmt.Printf("== Project %s\n", project.Name)
		for _, result := range project.CheckResults {
			if result.Success {
				fmt.Printf("✅ Check %s PASS\n", result.Name)
				for key, value := range result.Labels {
					fmt.Printf("\t%s: %s\n", key, value)
				}
			} else {
				fmt.Printf("🚨 Check %s FAILED\n", result.Name)
				for key, value := range result.Labels {
					fmt.Printf("\t%s: %s\n", key, value)
				}
				for _, ruleResult := range result.Results {
					fmt.Printf("\tMessage: %s\n", strings.Join(ruleResult.Messages, ","))
				}
			}
		}
	}
}
