package rules

import (
	"context"
	"errors"
	"fmt"

	"github.com/qonto/standards-insights/pkg/project"
)

type IsSubProjectRule struct {
}

func NewIsSubProjectRule() *IsSubProjectRule {
	return &IsSubProjectRule{}
}

func (rule *IsSubProjectRule) Do(ctx context.Context, project project.Project) error {
	if project.SubProject == "" {
		return errors.New("this project is not a subproject")
	}

	return nil
}