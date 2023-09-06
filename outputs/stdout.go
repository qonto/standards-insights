package outputs

import (
	"fmt"

	"github.com/qonto/standards-insights/checks/aggregates"
)

func Stdout(results []aggregates.ProjectResult) {
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
