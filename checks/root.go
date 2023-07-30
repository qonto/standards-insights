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
	SyncProjects() error
	HasNext() bool
	GetNext() aggregate.Project
}

func NewProcessor(config *config.Config, discovery Discovery) *CheckProcessor {
	return &CheckProcessor{
		config:    config,
		discovery: discovery,
	}
}
