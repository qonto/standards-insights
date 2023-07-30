package rules

/*
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
*/
