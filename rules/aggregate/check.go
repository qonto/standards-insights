package aggregate

type Check struct {
	Name   string
	Labels map[string]string
	Rules  []string
}

type CheckResult struct {
	Name    string
	Success bool
	Labels  map[string]string
	Results []RuleResult
}
