package cmd

import (
	"standards/checks"
	"standards/config"
	"standards/discovery"
)

func Run() {
	config, err := config.New()
	if err != nil {
		panic(err)
	}

	// Debug config parsing
	// fmt.Printf("Config: %s\n", config)

	discovery := discovery.New(config)
	processor := checks.NewProcessor(config, discovery) // TODO: find better name than processor
	err = processor.Run()
	if err != nil {
		panic(err)
	}
}
