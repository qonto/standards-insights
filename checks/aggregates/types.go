package aggregates

import "github.com/qonto/standards-insights/rules/aggregates"

type Check struct {
	Name   string
	Labels map[string]string
	Rules  []string
}

type CheckResult struct {
	Name    string
	Success bool
	Labels  map[string]string
	Results []aggregates.RuleResult
}

type Group struct {
	Name   string
	Checks []string
	When   []string
}

type ProjectResult struct {
	Name         string
	CheckResults []CheckResult
}
