package aggregate

type Rule struct {
	Name  string
	Files []FileRule
}

type FileRule struct {
	Path     string
	Contains *Regexp
	Exists   *bool
}

type RuleResult struct {
	RuleName string
	Messages []string
	Success  bool
}
