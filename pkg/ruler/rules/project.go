package rules

import (
	"context"
	"fmt"

	"github.com/qonto/standards-insights/config"
	"github.com/qonto/standards-insights/pkg/project"
)

type ProjectRule struct {
	config config.ProjectRule
}

func NewProjectRule(config config.ProjectRule) *ProjectRule {
	return &ProjectRule{
		config: config,
	}
}

func (rule *ProjectRule) Do(ctx context.Context, project project.Project) error {
	match := true
	if rule.config.Match != nil {
		match = *rule.config.Match
	}

	if rule.config.Name != "" {
		if match && project.Name != rule.config.Name {
			return fmt.Errorf("project name %s is not %s", project.Name, rule.config.Name)
		}
		if !match && project.Name == rule.config.Name {
			return fmt.Errorf("project name %s is matching", project.Name)
		}
	}

	if len(rule.config.Names) > 0 {
		if match && !contains(project.Name, rule.config.Names) {
			return fmt.Errorf("project name %s is not in %v", project.Name, rule.config.Names)
		}
		if !match && contains(project.Name, rule.config.Names) {
			return fmt.Errorf("project name %s is matching one of %v", project.Name, rule.config.Names)
		}
	}

	if len(rule.config.Labels) > 0 {
		if match && !isSubset(rule.config.Labels, project.Labels) {
			return fmt.Errorf("project labels %v does not contain %v", project.Labels, rule.config.Labels)
		}
		if !match && isSubset(rule.config.Labels, project.Labels) {
			return fmt.Errorf("project labels %v contain %v", project.Labels, rule.config.Labels)
		}
	}

	return nil
}

func isSubset(a, b map[string]string) bool {
	if len(a) > len(b) {
		return false
	}
	for k, vsub := range a {
		if vm, found := b[k]; !found || vm != vsub {
			return false
		}
	}
	return true
}

func contains(s string, slice []string) bool {
	for _, v := range slice {
		if s == v {
			return true
		}
	}
	return false
}
