package rules

import (
	"context"
	"errors"

	"github.com/qonto/standards-insights/pkg/project"
)

type SimpleRule struct {
	result bool
}

func NewSimpleRule(result bool) *SimpleRule {
	return &SimpleRule{
		result: result,
	}
}

func (rule *SimpleRule) Do(ctx context.Context, project project.Project) error {
	if rule.result {
		return nil
	}
	return errors.New("this rule always fail")
}
