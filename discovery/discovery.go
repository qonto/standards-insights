package discovery

import (
	"standards/rules/aggregate"
)

func (d *Discovery) SyncProjects() ([]*aggregate.Project, error) {
	// FIXME: get projects & clone them in temporary folder
	project := &aggregate.Project{
		Path: ".",
		Name: "standards-poc",
	}
	projects := []*aggregate.Project{project}
	return projects, nil
}
