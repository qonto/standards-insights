package config

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

func validate(config Config) error {
	rules := make(map[string]bool)
	for _, rule := range config.Rules {
		rules[rule.Name] = true
	}
	checks := make(map[string]bool)
	for _, check := range config.Checks {
		checks[check.Name] = true
		for _, rule := range check.Rules {
			_, ok := rules[rule]
			if !ok {
				return fmt.Errorf("rule %s not defined in check %s", rule, check.Name)
			}
		}
	}
	for _, group := range config.Groups {
		for _, check := range group.Checks {
			_, ok := checks[check]
			if !ok {
				return fmt.Errorf("check %s not defined in group %s", check, group.Name)
			}
			for _, rule := range group.When {
				_, ok := rules[rule]
				if !ok {
					return fmt.Errorf("rule (when clause) %s not defined in group %s", rule, group.Name)
				}
			}
		}
		if group.Files.ApplyToFiles {
			for _, groupCheck := range group.Checks {
				check, ok := findCheckByName(config.Checks, groupCheck)
				if !ok {
					continue
				}
				for _, checkRule := range check.Rules {
					rule, ok := findRuleByName(config.Rules, checkRule)
					if !ok {
						continue
					}
					for _, grepRule := range rule.Grep {
						if grepRule.Path != "" {
							return fmt.Errorf("grep rule %s in check %s has path defined but this is not allowed when the group applies to files", rule.Name, check.Name)
						}
					}
					if rule.Files != nil {
						return fmt.Errorf("file rule %s in check %s is not allowed when the group applies to files", rule.Name, check.Name)
					}
				}
			}

			for _, groupWhen := range group.When {
				rule, ok := findRuleByName(config.Rules, groupWhen)
				if !ok {
					return fmt.Errorf("when clause rule %s not defined in group %s", groupWhen, group.Name)
				}
				for _, grepRule := range rule.Grep {
					if grepRule.Path != "" {
						return fmt.Errorf("grep rule %s in when clause rule %s has path defined but this is not allowed when the group applies to files", rule.Name, groupWhen)
					}
				}
			}
		}
	}
	validate := validator.New()
	return validate.Struct(config)
}

// Helper functions to find check and rule by name
func findCheckByName(checks []Check, name string) (Check, bool) {
	for _, check := range checks {
		if check.Name == name {
			return check, true
		}
	}
	return Check{}, false
}

func findRuleByName(rules []Rule, name string) (Rule, bool) {
	for _, rule := range rules {
		if rule.Name == name {
			return rule, true
		}
	}
	return Rule{}, false
}
