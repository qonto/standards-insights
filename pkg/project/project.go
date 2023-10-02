package project

type Project struct {
	Name   string
	URL    string
	Branch string
	Path   string
	Labels map[string]string
}
