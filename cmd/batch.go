package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/qonto/standards-insights/config"
	"github.com/qonto/standards-insights/internal/build"
	"github.com/qonto/standards-insights/internal/providers"
	"github.com/qonto/standards-insights/pkg/checker"
	"github.com/qonto/standards-insights/pkg/ruler"

	"github.com/spf13/cobra"
)

func batchCmd(configPath, logLevel, logFormat *string) *cobra.Command {
	provider := ""
	cmd := &cobra.Command{
		Use:   "batch",
		Short: "Run checks on all projects",
		Run: func(cmd *cobra.Command, args []string) {
			logger := buildLogger(*logLevel, *logFormat)
			config, _, err := config.New(*configPath)
			exit(err)

			logger.Info(build.VersionMessage())

			providers, err := providers.NewProviders(logger, config.Providers)
			exit(err)

			ruler := ruler.NewRuler(logger, config.Rules)

			checker := checker.NewChecker(logger, ruler, config.Checks, config.Groups)

			for _, provider := range providers {
				ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
				projects, err := provider.FetchProjects(ctx)
				cancel()
				exit(err)
				fmt.Println("Done!")

				results := checker.Run(context.Background(), projects)
				stdoutResults(results)
			}
		},
	}

	// TODO: allow to repeat flag
	cmd.PersistentFlags().StringVarP(&provider, "provider", "p", "", "Filter providers")

	return cmd
}
