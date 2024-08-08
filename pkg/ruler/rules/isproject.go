package rules

import (
	"context"
	"errors"

	"github.com/qonto/standards-insights/pkg/project"
)

type IsProjectRule struct {
}

func NewIsProjectRule() *IsProjectRule {
	return &IsProjectRule{}
}

func (rule *IsProjectRule) Do(ctx context.Context, project project.Project) error {
	if project.SubProject != "" {
		return errors.New("this project is a subproject")
	}

	return nil
}