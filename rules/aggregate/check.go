package aggregate

import "context"

type Check struct {
	Name     string
	Category string
	Level    string
	Exclude  []string
	Include  []string
	Rules    []*Rule
}

func (c *Check) IsMatchingRules() bool {
	for _, rule := range c.Rules {
		err := rule.Check(context.Background())
		if err != nil {
			return false
		}
	}
	return true
}
