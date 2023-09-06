package cmd

import (
	"context"
	"os"
	"path/filepath"

	"github.com/qonto/standards-insights/checks"
	"github.com/qonto/standards-insights/config"
	"github.com/qonto/standards-insights/providers/aggregates"
	"github.com/qonto/standards-insights/rules"

	"github.com/spf13/cobra"
)

func runCmd(configPath *string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "run",
		Short: "Run checks on current project",
		Run: func(cmd *cobra.Command, args []string) {
			RunLocalCheck(*configPath)
		},
	}

	return cmd
}

func RunLocalCheck(configPath string) {
	config, err := config.New(configPath)
	if err != nil {
		panic(err)
	}

	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	ruler := rules.NewRuler(config.Rules)

	checker := checks.NewChecker(ruler, config.Checks, config.Groups)
	projects := []aggregates.Project{
		{
			Path: ".",
			Name: filepath.Base(dir),
		},
	}
	err = checker.Run(context.Background(), projects)
	if err != nil {
		panic(err)
	}
}
