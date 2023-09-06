package aggregates

import "github.com/qonto/standards-insights/types"

type Rule struct {
	Name  string
	Files []FileRule
}

type FileRule struct {
	Path     string
	Contains *types.Regexp
	Exists   *bool
}

type RuleResult struct {
	RuleName string
	Messages []string
	Success  bool
}
