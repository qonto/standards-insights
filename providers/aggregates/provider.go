package aggregates

type Provider interface {
	FetchProjects() ([]Project, error)
}

type Project struct {
	Name   string
	URL    string
	Branch string
	Path   string
}
