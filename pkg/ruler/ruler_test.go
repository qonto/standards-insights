package ruler_test

import (
	"context"
	"testing"

	"log/slog"

	"github.com/qonto/standards-insights/config"
	"github.com/qonto/standards-insights/pkg/project"
	"github.com/qonto/standards-insights/pkg/ruler"
	"github.com/stretchr/testify/assert"
)

func TestRuler(t *testing.T) {
	yes := true
	ruler := ruler.NewRuler(slog.Default(), []config.Rule{
		{
			Name: "rule1",
			Files: []config.FileRule{
				{
					Path:   "_testdata/file1",
					Exists: &yes,
				},
			},
		},
		{
			Name: "rule2",
			Files: []config.FileRule{
				{
					Path:   "_testdata/file_ko",
					Exists: &yes,
				},
			},
		},
	})

	project := project.Project{
		Path: "./",
	}
	result := ruler.Execute(context.Background(), "rule1", project)
	assert.True(t, result.Success)
	assert.Equal(t, 0, len(result.Messages))
	assert.Equal(t, "rule1", result.RuleName)

	result = ruler.Execute(context.Background(), "rule2", project)
	assert.False(t, result.Success)
	assert.Equal(t, 1, len(result.Messages))
	assert.Equal(t, "rule2", result.RuleName)
	assert.Contains(t, result.Messages[0], "does not exist")

	result = ruler.Execute(context.Background(), "rule3", project)
	assert.False(t, result.Success)
	assert.Equal(t, 1, len(result.Messages))
	assert.Equal(t, "rule3", result.RuleName)
	assert.Contains(t, result.Messages[0], "not found in the rules configuration")
}
