package cmd

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/qonto/standards-insights/pkg/checker/aggregates"
	"github.com/spf13/cobra"
)

func exit(err error) {
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(2)
	}
}

func Run() error {
	configPath := ""
	logLevel := ""
	logFormat := ""
	rootCmd := &cobra.Command{
		Use:   "qstandards",
		Short: "Standard insights",
	}

	rootCmd.PersistentFlags().StringVarP(&configPath, "config", "c", "config.yaml", "Config path")
	rootCmd.PersistentFlags().StringVarP(&logLevel, "log-level", "v", "info", "Logger log level (debug, info, warn, error)")
	rootCmd.PersistentFlags().StringVar(&logFormat, "log-format", "text", "Logger logs format (text, json)")

	rootCmd.AddCommand(batchCmd(&configPath, &logLevel, &logFormat))
	rootCmd.AddCommand(serverCmd(&configPath, &logLevel, &logFormat))
	rootCmd.AddCommand(runCmd(&configPath, &logLevel, &logFormat))

	return rootCmd.Execute()
}

func stdoutResults(results []aggregates.ProjectResult) {
	for _, project := range results {
		fmt.Printf("== Project %s\n", project.Name)
		for _, result := range project.CheckResults {
			if result.Success {
				fmt.Printf("✅ Check %s PASS (labels: %s)\n", result.Name, result.Labels)
			} else {
				fmt.Printf("🚨 Check %s FAILED (labels: %s)\n", result.Name, result.Labels)
				fmt.Printf("🚨 %+v\n", result)
			}
		}
	}
}

func buildLogger(level string, format string) *slog.Logger {
	var programLevel = new(slog.LevelVar)
	switch level {
	case "debug":
		programLevel.Set(slog.LevelDebug)
	case "info":
		programLevel.Set(slog.LevelInfo)
	case "warn":
		programLevel.Set(slog.LevelWarn)
	case "error":
		programLevel.Set(slog.LevelError)
	default:
		programLevel.Set(slog.LevelInfo)
	}

	options := &slog.HandlerOptions{Level: programLevel}
	switch format {
	case "text":
		return slog.New(slog.NewTextHandler(os.Stdout, options))
	case "json":
		return slog.New(slog.NewJSONHandler(os.Stdout, options))
	default:
		return slog.New(slog.NewTextHandler(os.Stdout, options))
	}
}
