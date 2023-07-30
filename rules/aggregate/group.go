package aggregate

import (
	"context"
)

type Group struct {
	Name   string
	Checks []*Check
	Rules  []*Rule
}

func (g *Group) IsMatchingRules() bool {
	for _, rule := range g.Rules {
		err := rule.Check(context.Background())
		if err != nil {
			return false
		}
	}

	return true
}
