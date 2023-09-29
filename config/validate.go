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
	}
	validate := validator.New()
	return validate.Struct(config)
}
