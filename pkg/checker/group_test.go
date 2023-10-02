package checker

import (
	"context"
	"testing"

	"log/slog"

	"github.com/qonto/standards-insights/config"
	"github.com/qonto/standards-insights/pkg/project"
	"github.com/qonto/standards-insights/pkg/ruler"
	"github.com/stretchr/testify/assert"
)

func TestShouldSkipGroup(t *testing.T) {
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
	checker := NewChecker(logger, ruler, nil, nil)
	group := config.Group{When: []string{"rule1"}}
	project := project.Project{}
	result := checker.shouldSkipGroup(context.Background(), group, project)
	assert.False(t, result)

	group = config.Group{When: []string{"rule1", "rule2"}}
	result = checker.shouldSkipGroup(context.Background(), group, project)
	assert.True(t, result)

	group = config.Group{When: []string{"ruledoesnotexist"}}
	result = checker.shouldSkipGroup(context.Background(), group, project)
	assert.True(t, result)
}
