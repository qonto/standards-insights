package checks

import (
	"standards/config"
)

type CheckProcessor struct {
	config *config.Config
}

func NewProcessor(config *config.Config) *CheckProcessor {
	return &CheckProcessor{
		config: config,
	}
}
