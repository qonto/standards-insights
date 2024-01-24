package rules

import (
	"context"
	"errors"
	"testing"

	"github.com/qonto/standards-insights/config"
	"github.com/qonto/standards-insights/pkg/project"
	"github.com/stretchr/testify/assert"
)

func TestProjectRule_Do(t *testing.T) {
	f := false
	tests := []struct {
		name    string
		config  config.ProjectRule
		project project.Project
		wantErr error
	}{
		{
			name:    "empty rule should pass",
			config:  config.ProjectRule{},
			project: project.Project{Name: "a", Labels: map[string]string{"a": "b"}},
			wantErr: nil,
		},
		{
			name: "project name match",
			config: config.ProjectRule{
				Name: "a",
			},
			project: project.Project{Name: "a"},
			wantErr: nil,
		},
		{
			name: "project name does not match",
			config: config.ProjectRule{
				Name: "a",
			},
			project: project.Project{Name: "c"},
			wantErr: errors.New("project name c is not a"),
		},
		{
			name: "project name matches and should not",
			config: config.ProjectRule{
				Name:  "a",
				Match: &f,
			},
			project: project.Project{Name: "a"},
			wantErr: errors.New("project name a is matching"),
		},
		{
			name: "project name does not match and should not",
			config: config.ProjectRule{
				Name:  "a",
				Match: &f,
			},
			project: project.Project{Name: "c"},
			wantErr: nil,
		},

		{
			name: "project name list match",
			config: config.ProjectRule{
				Names: []string{"a", "b", "c"},
			},
			project: project.Project{Name: "b"},
			wantErr: nil,
		},
		{
			name: "project name does not match list",
			config: config.ProjectRule{
				Names: []string{"a", "b", "c"},
			},
			project: project.Project{Name: "f"},
			wantErr: errors.New("project name f is not in [a b c]"),
		},
		{
			name: "project name matches list and should not",
			config: config.ProjectRule{
				Names: []string{"a", "b", "c"},
				Match: &f,
			},
			project: project.Project{Name: "a"},
			wantErr: errors.New("project name a is matching one of [a b c]"),
		},
		{
			name: "project name does not match list and should not",
			config: config.ProjectRule{
				Names: []string{"a", "b", "c"},
				Match: &f,
			},
			project: project.Project{Name: "e"},
			wantErr: nil,
		},

		{
			name: "project labels full match",
			config: config.ProjectRule{
				Labels: map[string]string{
					"a": "b",
					"e": "f",
				},
			},
			project: project.Project{Labels: map[string]string{
				"a": "b",
				"c": "d",
				"e": "f",
			}},
			wantErr: nil,
		},
		{
			name: "project labels not matching",
			config: config.ProjectRule{
				Labels: map[string]string{
					"a": "b",
					"e": "f",
				},
			},
			project: project.Project{Labels: map[string]string{
				"a": "b",
				"c": "d",
				"g": "h",
			}},
			wantErr: errors.New("project labels map[a:b c:d g:h] does not contain map[a:b e:f]"),
		},
		{
			name: "project labels should not match",
			config: config.ProjectRule{
				Labels: map[string]string{
					"a": "b",
					"e": "f",
				},
				Match: &f,
			},
			project: project.Project{Labels: map[string]string{
				"a": "b",
				"c": "d",
				"e": "f",
			}},
			wantErr: errors.New("project labels map[a:b c:d e:f] contain map[a:b e:f]"),
		},
		{
			name: "project labels should not match and don't match",
			config: config.ProjectRule{
				Labels: map[string]string{
					"a": "b",
					"e": "h",
				},
				Match: &f,
			},
			project: project.Project{Labels: map[string]string{
				"a": "b",
				"c": "d",
				"e": "f",
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rule := NewProjectRule(tt.config)

			err := rule.Do(context.Background(), tt.project)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}
