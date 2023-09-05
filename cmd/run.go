package cmd

import (
	"context"
	"os"
	"path/filepath"
	"standards/checks"
	"standards/config"
	"standards/providers"
	"standards/rules"

	"github.com/spf13/cobra"
)

func runCmd(configPath *string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "run",
		Short: "Run checks on current project",
		Run: func(cmd *cobra.Command, args []string) {
			config, err := config.New(*configPath)
			if err != nil {
				panic(err)
			}

			dir, err := os.Getwd()
			if err != nil {
				panic(err)
			}

			ruler := rules.NewRuler(config.Rules)

			checker := checks.NewChecker(ruler, config.Checks, config.Groups)
			projects := []providers.Project{
				{
					Path: ".",
					Name: filepath.Base(dir),
				},
			}
			err = checker.Run(context.Background(), projects)
			if err != nil {
				panic(err)
			}
		},
	}

	return cmd
}
