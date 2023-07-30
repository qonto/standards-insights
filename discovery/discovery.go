package discovery

import (
	"fmt"
	"standards/rules/aggregate"
)

type iterator struct {
	index int
}

func (d *Discovery) SyncProjects() error {
	// FIXME: get projects & clone them in temporary folder
	fmt.Println("FIXME")
	return nil
}

func (d *Discovery) GetNext() aggregate.Project {
	// FIXME: Use iterator
	d.iterator.index++

	return aggregate.Project{
		Path: ".",
		Name: "standards-poc",
	}
}

func (d *Discovery) HasNext() bool {
	// FIXME: Use iterator
	return d.iterator.index == -1
}
