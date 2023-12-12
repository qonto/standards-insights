package rules_test

import (
	"context"
	"testing"

	"github.com/qonto/standards-insights/pkg/project"
	"github.com/qonto/standards-insights/pkg/ruler/rules"
	"github.com/stretchr/testify/assert"
)

func TestGrepRule(t *testing.T) {
	project := project.Project{
		Path: "./",
	}

	cases := []struct {
		rule  rules.GrepRule
		error string
	}{
		{
			rule: rules.GrepRule{
				Path:      "_testdata",
				Recursive: true,
				Pattern:   "abcdefg",
				Match:     true,
			},
		},
		{
			rule: rules.GrepRule{
				Path:         "_testdata",
				Recursive:    true,
				Pattern:      "abcdefg",
				Match:        true,
				SkipNotFound: true,
			},
		},
		{
			rule: rules.GrepRule{
				Path:      "_testdata",
				Recursive: true,
				Pattern:   "abc*",
				Match:     true,
			},
		},
		{
			rule: rules.GrepRule{
				Path:    "_testdata/file1",
				Pattern: "abcdefg",
				Match:   true,
			},
		},
		{
			rule: rules.GrepRule{
				Path:      "_testdata",
				Recursive: true,
				Pattern:   "aIOJij89Yaa",
				Match:     false,
			},
		},
		{
			rule: rules.GrepRule{
				Path:    "_testdata/file1",
				Pattern: "aIOJij89Yaa",
				Match:   false,
			},
		},
		{
			rule: rules.GrepRule{
				Path:    "_testdata",
				Pattern: "abcdefg",
				Match:   true,
			},
			error: "Is a directory",
		},
		{
			rule: rules.GrepRule{
				Path:      "_testdata",
				Recursive: true,
				Pattern:   "abcdefg",
				Match:     false,
			},
			error: "match found for pattern",
		},
		{
			rule: rules.GrepRule{
				Path:    "_testdata/file1",
				Pattern: "abcdefg",
				Match:   false,
			},
			error: "match found for pattern",
		},
		{
			rule: rules.GrepRule{
				Path:    "_testdata/file1",
				Pattern: "abc*",
				Match:   false,
			},
			error: "match found for pattern",
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
