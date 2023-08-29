package cmd

import (
	"fmt"
	"standards/checks"
	"standards/config"
	"standards/discovery"

	"github.com/spf13/cobra"
)

func batchCmd(configPath *string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "batch",
		Short: "Run checks on all projects",
		Run: func(cmd *cobra.Command, args []string) {
			config, err := config.New(*configPath)
			if err != nil {
				panic(err)
			}

			fmt.Println("Syncing all projects...")
			discovery := discovery.New(config)
			projects, err := discovery.SyncProjects()
			if err != nil {
				panic(err)
			}
			fmt.Println("Done!")

			processor := checks.NewProcessor(config)
			err = processor.Run(projects)
			if err != nil {
				panic(err)
			}
		},
	}

	return cmd
}
