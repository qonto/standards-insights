package cmd

import (
	"context"
	"os"
	"path/filepath"

	"github.com/qonto/standards-insights/config"
	"github.com/qonto/standards-insights/internal/build"
	"github.com/qonto/standards-insights/pkg/checker"
	"github.com/qonto/standards-insights/pkg/project"
	"github.com/qonto/standards-insights/pkg/ruler"

	"github.com/spf13/cobra"
)

func runCmd(configPath *string, logLevel, logFormat *string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "run",
		Short: "Run checks on current project",
		Run: func(cmd *cobra.Command, args []string) {
			RunLocalCheck(configPath, logLevel, logFormat)
		},
	}
	return cmd
}

func RunLocalCheck(configPath, logLevel, logFormat *string) {
	logger := buildLogger(*logLevel, *logFormat)
	config, _, err := config.New(*configPath)
	exit(err)

	logger.Info(build.VersionMessage())

	dir, err := os.Getwd()
	exit(err)

	ruler := ruler.NewRuler(logger, config.Rules)

	checker := checker.NewChecker(logger, ruler, config.Checks, config.Groups)
	projects := []project.Project{
		{
			Path: ".",
			Name: filepath.Base(dir),
		},
	}
	results := checker.Run(context.Background(), projects)
	stdoutResults(results)
}
