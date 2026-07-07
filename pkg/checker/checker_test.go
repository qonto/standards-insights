package checker_test

import (
	"context"
	"log/slog"
	"os"
	"path/filepath"
	"testing"

	"github.com/qonto/standards-insights/config"
	"github.com/qonto/standards-insights/pkg/checker"
	"github.com/qonto/standards-insights/pkg/project"
	"github.com/qonto/standards-insights/pkg/ruler"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
		{
			Name:     "check3",
			Rules:    []string{"rule1", "rule2"},
			Operator: "or",
			Labels: map[string]string{
				"category": "cat3",
			},
		},
	}
	groups := []config.Group{
		{
			Name:   "group1",
			Checks: []string{"check1", "check2", "check3"},
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
	assert.Equal(t, 3, len(results[0].CheckResults))
	assert.Equal(t, "check1", results[0].CheckResults[0].Name)
	assert.Equal(t, "check2", results[0].CheckResults[1].Name)
	assert.Equal(t, "check3", results[0].CheckResults[2].Name)
	assert.Equal(t, "cat1", results[0].CheckResults[0].Labels["category"])
	assert.Equal(t, "cat2", results[0].CheckResults[1].Labels["category"])
	assert.Equal(t, "cat3", results[0].CheckResults[2].Labels["category"])
	assert.True(t, results[0].CheckResults[0].Success)
	assert.False(t, results[0].CheckResults[1].Success)
	assert.True(t, results[0].CheckResults[2].Success)
}

func TestRunWithSkipConfig(t *testing.T) {
	yes := true
	logger := slog.Default()
	r := ruler.NewRuler(logger, []config.Rule{
		{Name: "rule1", Simple: &yes},
	})
	checks := []config.Check{
		{Name: "check1", Rules: []string{"rule1"}},
		{Name: "check2", Rules: []string{"rule1"}},
		{Name: "check3", Rules: []string{"rule1"}},
	}
	groups := []config.Group{
		{Name: "group1", Checks: []string{"check1", "check2"}},
		{Name: "group2", Checks: []string{"check3"}},
	}
	c := checker.NewChecker(logger, r, checks, groups)

	// Repo opts out of check1 (individual) and group2 (whole group).
	dir := t.TempDir()
	err := os.WriteFile(filepath.Join(dir, ".standards-insights.yaml"), []byte(`
skip:
  groups:
    - group2
  checks:
    - check1
`), 0o600)
	require.NoError(t, err)

	projects := []project.Project{{Name: "project1", Path: dir}}
	results := c.Run(context.Background(), projects)

	require.Equal(t, 1, len(results))
	// check1 (skipped) and check3 (group2 skipped) are gone; only check2 remains.
	require.Equal(t, 1, len(results[0].CheckResults))
	assert.Equal(t, "check2", results[0].CheckResults[0].Name)
}
