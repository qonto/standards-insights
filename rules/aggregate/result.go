package aggregate

type RuleResult struct {
	RuleName string
	Err      error
}

type GroupResult struct {
	GroupName    string
	RulesResults []RuleResult
	Skipped      error
}
