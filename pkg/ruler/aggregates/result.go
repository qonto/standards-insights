package aggregates

type RuleResult struct {
	RuleName string
	Messages []string
	Success  bool
}
