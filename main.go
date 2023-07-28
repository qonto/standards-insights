package main

import (
	"context"
	"fmt"
	"os"
	"standards/config"
	"standards/rules"

	"gopkg.in/yaml.v3"
)

func main() {
	var config config.Config

	file, err := os.ReadFile("config.yaml")
	if err != nil {
		panic(err)
	}
	err = yaml.Unmarshal(file, &config)
	if err != nil {
		panic(err)
	}
	checker := rules.NewChecker(config.Rules, config.Groups)

	result, err := checker.CheckGroup(context.Background(), "golang")
	if err != nil {
		fmt.Println(err.Error())
	} else {
		if result.Skipped != nil {
			fmt.Printf("group %s skipped: %s\n", result.GroupName, result.Skipped.Error())
		} else {
			for _, ruleResult := range result.RulesResults {
				success := ruleResult.Err == nil
				if success {
					fmt.Printf("rule %s successful\n", ruleResult.RuleName)
				} else {
					fmt.Printf("rule %s failed: %s\n", ruleResult.RuleName, ruleResult.Err.Error())
				}

			}
		}

	}
}
