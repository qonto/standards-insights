package providers

type Provider interface {
	FetchProjects() ([]Project, error)
}

type Project struct {
	Path string
	Name string
}
