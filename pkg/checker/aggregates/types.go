package aggregates

import ruleraggregates "github.com/qonto/standards-insights/pkg/ruler/aggregates"

type CheckResult struct {
	Name    string
	Success bool
	Labels  map[string]string
	Results []ruleraggregates.RuleResult
}

type ProjectResult struct {
	Name         string
	CheckResults []CheckResult
	Labels       map[string]string
}
