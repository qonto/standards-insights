package checker_test

import (
	"context"
	"log/slog"
	"testing"

	"github.com/qonto/standards-insights/config"
	"github.com/qonto/standards-insights/pkg/checker"
	"github.com/qonto/standards-insights/pkg/project"
	"github.com/qonto/standards-insights/pkg/ruler"
	"github.com/stretchr/testify/assert"
)

func TestRun(t *testing.T) {
	yes := true
	no := false
	logger := slog.Default()
	ruler := ruler.NewRuler(logger, []config.Rule{
		{
			Name:   "rule1",
			Simple: &yes,
		},
		{
			Name:   "rule2",
			Simple: &no,
		},
	})
	checks := []config.Check{
		{
			Name:  "check1",
			Rules: []string{"rule1"},
			Labels: map[string]string{
				"category": "cat1",
			},
		},
		{
			Name:  "check2",
			Rules: []string{"rule1", "rule2"},
			Labels: map[string]string{
				"category": "cat2",
			},
		},
	}
	groups := []config.Group{
		{
			Name:   "group1",
			Checks: []string{"check1", "check2"},
			When:   []string{"rule1"},
		},
	}

	checker := checker.NewChecker(logger, ruler, checks, groups)
	projects := []project.Project{
		{
			Name: "project1",
			Labels: map[string]string{
				"team": "sre",
			},
		},
	}
	results := checker.Run(context.Background(), projects)
	assert.Equal(t, 1, len(results))
	assert.Equal(t, results[0].Name, "project1")
	assert.Equal(t, 1, len(results[0].Labels))
	assert.Equal(t, "sre", results[0].Labels["team"])
	assert.Equal(t, 2, len(results[0].CheckResults))
	assert.Equal(t, "check1", results[0].CheckResults[0].Name)
	assert.Equal(t, "check2", results[0].CheckResults[1].Name)
	assert.Equal(t, "cat1", results[0].CheckResults[0].Labels["category"])
	assert.Equal(t, "cat2", results[0].CheckResults[1].Labels["category"])
	assert.True(t, results[0].CheckResults[0].Success)
	assert.False(t, results[0].CheckResults[1].Success)
}
