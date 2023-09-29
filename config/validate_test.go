package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidate(t *testing.T) {
	yes := true
	cases := []struct {
		config Config
		error  string
	}{
		{
			config: Config{
				Groups: []Group{
					{
						Name:   "gr1",
						Checks: []string{"check1"},
						When:   []string{"rule2"},
					},
				},
				Checks: []Check{
					{
						Name:  "check1",
						Rules: []string{"rule1"},
					},
				},
				Rules: []Rule{
					{
						Name:   "rule1",
						Simple: &yes,
					},
					{
						Name:   "rule2",
						Simple: &yes,
					},
				},
			},
		},
		{
			config: Config{
				Groups: []Group{
					{
						Name:   "gr1",
						Checks: []string{"check1"},
						When:   []string{"rule3"},
					},
				},
				Checks: []Check{
					{
						Name:  "check1",
						Rules: []string{"rule1"},
					},
				},
				Rules: []Rule{
					{
						Name:   "rule1",
						Simple: &yes,
					},
					{
						Name:   "rule2",
						Simple: &yes,
					},
				},
			},
			error: "rule3 not defined in group gr1",
		},
		{
			config: Config{
				Groups: []Group{
					{
						Name:   "gr1",
						Checks: []string{"check1"},
						When:   []string{"rule2"},
					},
				},
				Checks: []Check{
					{
						Name:  "check1",
						Rules: []string{"rule3"},
					},
				},
				Rules: []Rule{
					{
						Name:   "rule1",
						Simple: &yes,
					},
					{
						Name:   "rule2",
						Simple: &yes,
					},
				},
			},
			error: "rule3 not defined in check check1",
		},
		{
			config: Config{
				Groups: []Group{
					{
						Name:   "gr1",
						Checks: []string{"check2"},
						When:   []string{"rule2"},
					},
				},
				Checks: []Check{
					{
						Name:  "check1",
						Rules: []string{"rule1"},
					},
				},
				Rules: []Rule{
					{
						Name:   "rule1",
						Simple: &yes,
					},
					{
						Name:   "rule2",
						Simple: &yes,
					},
				},
			},
			error: "check2 not defined in group gr1",
		},
	}
	for _, c := range cases {
		err := validate(c.config)
		if c.error == "" {
			assert.NoError(t, err)
		} else {
			assert.ErrorContains(t, err, c.error)
		}
	}
}
