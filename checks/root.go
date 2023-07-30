package checks

import (
	"standards/config"
	"standards/rules/aggregate"
)

type CheckProcessor struct {
	config    *config.Config
	discovery Discovery
}

type Discovery interface {
	SyncProjects() ([]*aggregate.Project, error)
}

func NewProcessor(config *config.Config, discovery Discovery) *CheckProcessor {
	return &CheckProcessor{
		config:    config,
		discovery: discovery,
	}
}
