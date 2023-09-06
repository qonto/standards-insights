package cmd

import (
	"context"
	"fmt"

	"github.com/qonto/standards-insights/checks"
	"github.com/qonto/standards-insights/config"
	"github.com/qonto/standards-insights/outputs"
	"github.com/qonto/standards-insights/providers"
	"github.com/qonto/standards-insights/rules"

	"github.com/spf13/cobra"
)

func batchCmd(configPath *string) *cobra.Command {
	provider := ""
	cmd := &cobra.Command{
		Use:   "batch",
		Short: "Run checks on all projects",
		Run: func(cmd *cobra.Command, args []string) {
			config, err := config.New(*configPath)
			exit(err)

			providers, err := providers.NewProviders(&config.Providers, []string{"argocd"})
			exit(err)

			ruler := rules.NewRuler(config.Rules)

			checker := checks.NewChecker(ruler, config.Checks, config.Groups)

			for _, provider := range providers {
				projects, err := provider.FetchProjects()
				exit(err)
				fmt.Println("Done!")

				results := checker.Run(context.Background(), projects)
				outputs.Stdout(results)
			}
		},
	}

	// TODO: allow to repeat flag
	cmd.PersistentFlags().StringVarP(&provider, "provider", "p", "", "Filter providers")

	return cmd
}
