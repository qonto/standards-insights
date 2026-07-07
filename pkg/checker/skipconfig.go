package checker

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/qonto/standards-insights/pkg/project"
	"gopkg.in/yaml.v3"
)

// skipConfigFileName is the name of the optional file a repository can commit to
// its root to opt out of specific checks and/or groups (Rubocop-style).
const skipConfigFileName = ".standards-insights.yaml"

type skipConfig struct {
	Skip struct {
		Groups []string `yaml:"groups"`
		Checks []string `yaml:"checks"`
	} `yaml:"skip"`
}

// loadSkipConfig reads the per-repository skip configuration from the root of the
// given project path. It returns two sets (group names and check names to skip).
//
// It fails open: when the file is absent it returns empty sets and a nil error, and
// when the file is malformed it returns empty sets alongside the parse error so the
// caller can warn and skip nothing (a broken opt-out file must never silently disable
// enforcement).
func loadSkipConfig(projectPath string) (skippedGroups, skippedChecks map[string]bool, err error) {
	skippedGroups = map[string]bool{}
	skippedChecks = map[string]bool{}

	path := filepath.Join(projectPath, skipConfigFileName)
	content, err := os.ReadFile(path) //nolint
	if err != nil {
		if os.IsNotExist(err) {
			return skippedGroups, skippedChecks, nil
		}
		return skippedGroups, skippedChecks, fmt.Errorf("failed to read %s: %w", skipConfigFileName, err)
	}

	var config skipConfig
	if err := yaml.Unmarshal(content, &config); err != nil {
		return skippedGroups, skippedChecks, fmt.Errorf("failed to parse %s: %w", skipConfigFileName, err)
	}

	for _, group := range config.Skip.Groups {
		skippedGroups[group] = true
	}
	for _, check := range config.Skip.Checks {
		skippedChecks[check] = true
	}

	return skippedGroups, skippedChecks, nil
}

// warnUnknownSkips logs a warning for any skipped group or check name that does not
// match a configured group/check, to help repository owners catch typos in their
// per-repository skip file.
func (c *Checker) warnUnknownSkips(project project.Project, skippedGroups, skippedChecks map[string]bool) {
	knownGroups := make(map[string]bool, len(c.groups))
	for _, group := range c.groups {
		knownGroups[group.Name] = true
	}
	for group := range skippedGroups {
		if !knownGroups[group] {
			c.logger.Warn(fmt.Sprintf("%s for project %s skips unknown group %s", skipConfigFileName, project.Name, group))
		}
	}
	for check := range skippedChecks {
		if _, ok := c.checks[check]; !ok {
			c.logger.Warn(fmt.Sprintf("%s for project %s skips unknown check %s", skipConfigFileName, project.Name, check))
		}
	}
}
