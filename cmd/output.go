package cmd

import (
	"fmt"

	"github.com/qonto/standards-insights/pkg/checker/aggregates"
)

func stdoutResults(results []aggregates.ProjectResult) {
	for _, project := range results {
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
}
