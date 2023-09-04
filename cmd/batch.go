package cmd

import (
	"fmt"
	"standards/checks"
	"standards/config"
	"standards/providers"

	"github.com/spf13/cobra"
)

func batchCmd(configPath *string) *cobra.Command {
	provider := ""
	cmd := &cobra.Command{
		Use:   "batch",
		Short: "Run checks on all projects",
		Run: func(cmd *cobra.Command, args []string) {
			config, err := config.New(*configPath)
			if err != nil {
				panic(err)
			}

			providers, err := providers.NewProviders(&config.Providers, []string{"argocd"})
			if err != nil {
				panic(err)
			}

			processor := checks.NewProcessor(config)

			for _, provider := range providers {
				projects, err := provider.FetchProjects()
				if err != nil {
					panic(err)
				}
				fmt.Println("Done!")

				err = processor.Run(projects)
				if err != nil {
					panic(err)
				}
			}
		},
	}

	// TODO: allow to repeat flag
	cmd.PersistentFlags().StringVarP(&provider, "provider", "p", "", "Filter providers")

	return cmd
}
