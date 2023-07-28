package rules

import (
	"context"
	"fmt"
	"standards/rules/aggregate"
)

func CheckRule(ctx context.Context, rule *aggregate.Rule) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}
	messages := []string{}
	success := true
	if len(rule.Files) != 0 {
		for _, fileRule := range rule.Files {
			err := fileRule.Verify(ctx)
			if err != nil {
				success = false
				summary := fileRule.Summary()
				message := fmt.Sprintf("%s, error: %s", summary, err.Error())
				messages = append(messages, message)
			}
		}
	}
	if !success {
		return &aggregate.RuleError{
			RuleName: rule.Name,
			Messages: messages,
		}
	}
	return nil
}

type Checker struct {
	Rules  map[string]aggregate.Rule
	Groups map[string]aggregate.Group
}

func NewChecker(rules []aggregate.Rule, groups []aggregate.Group) *Checker {
	checker := &Checker{
		Rules:  make(map[string]aggregate.Rule),
		Groups: make(map[string]aggregate.Group),
	}
	for _, rule := range rules {
		checker.Rules[rule.Name] = rule

	}
	for _, group := range groups {
		checker.Groups[group.Name] = group

	}
	return checker

}

func (c *Checker) CheckGroup(ctx context.Context, groupName string) (*aggregate.GroupResult, error) {
	group, ok := c.Groups[groupName]
	if !ok {
		return nil, fmt.Errorf("group %s not found", groupName)
	}
	result := &aggregate.GroupResult{
		GroupName: groupName,
	}

	if len(group.When) != 0 {
		for _, ruleName := range group.When {
			rule, ok := c.Rules[ruleName]
			if !ok {
				return nil, fmt.Errorf("rule %s not found in group %s when condition", ruleName, groupName)
			}
			err := CheckRule(ctx, &rule)
			if err != nil {
				result.Skipped = err
				// abort immediately
				return result, nil

			}
		}
	}
	for _, ruleName := range group.Rules {
		rule, ok := c.Rules[ruleName]
		if !ok {
			return nil, fmt.Errorf("rule %s not found in group %s", ruleName, groupName)
		}
		ruleResult := aggregate.RuleResult{
			RuleName: ruleName,
		}
		err := CheckRule(ctx, &rule)
		if err != nil {
			ruleResult.Err = err
		}
		result.RulesResults = append(result.RulesResults, ruleResult)
	}
	return result, nil
}
