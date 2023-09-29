package rules_test

import (
	"context"
	"testing"

	"github.com/qonto/standards-insights/pkg/project"
	"github.com/qonto/standards-insights/pkg/ruler/rules"
	"github.com/qonto/standards-insights/types"
	"github.com/stretchr/testify/assert"
)

func TestFileRule(t *testing.T) {
	project := project.Project{
		Path: "./",
	}
	yes := true
	no := false
	r1, err := types.RegexpFromString("<3_ruby")
	assert.NoError(t, err)
	r2, err := types.RegexpFromString("azer.*")
	assert.NoError(t, err)
	r3, err := types.RegexpFromString("foobar")
	assert.NoError(t, err)
	cases := []struct {
		rule  rules.FileRule
		error string
	}{
		{
			rule: rules.FileRule{
				Path:   "_testdata/file1",
				Exists: &yes,
			},
		},
		{
			rule: rules.FileRule{
				Path:   "_testdata/file_ko",
				Exists: &no,
			},
		},
		{
			rule: rules.FileRule{
				Path:     "_testdata/file1",
				Contains: r1,
			},
		},
		{
			rule: rules.FileRule{
				Path:     "_testdata/file1",
				Contains: r2,
			},
		},
		{
			rule: rules.FileRule{
				Path:        "_testdata/file1",
				NotContains: r3,
			},
		},
		{
			rule: rules.FileRule{
				Path:        "_testdata/file1",
				NotContains: r1,
			},
			error: "pattern <3_ruby found in file",
		},
		{
			rule: rules.FileRule{
				Path:     "_testdata/file1",
				Contains: r3,
			},
			error: "pattern foobar not found in file",
		},
		{
			rule: rules.FileRule{
				Path:   "_testdata/file1",
				Exists: &no,
			},
			error: "file _testdata/file1 exists",
		},
		{
			rule: rules.FileRule{
				Path:   "_testdata/file_ko",
				Exists: &yes,
			},
			error: "does not exist",
		},
	}
	for _, c := range cases {
		err := c.rule.Do(context.Background(), project)
		if c.error == "" {
			assert.NoError(t, err)
		} else {
			assert.ErrorContains(t, err, c.error)
		}
	}
}
