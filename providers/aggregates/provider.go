package aggregates

import "context"

type Provider interface {
	FetchProjects(ctx context.Context) ([]Project, error)
}

type Project struct {
	Name   string
	URL    string
	Branch string
	Path   string
	Labels map[string]string
}
