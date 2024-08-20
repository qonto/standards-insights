package project

type Project struct {
	Name        string
	URL         string
	Branch      string
	Path        string
	FilePath    string    // a path to a file within the project
	SubProjects []Project // a list of subprojects within the project
	Labels      map[string]string
}
