package cmd

import (
	"github.com/spf13/cobra"
)

func Run() error {
	configPath := ""
	rootCmd := &cobra.Command{
		Use:   "qstandards",
		Short: "Standard insights",
	}

	rootCmd.PersistentFlags().StringVarP(&configPath, "config", "c", "config.yaml", "Config path")

	rootCmd.AddCommand(batchCmd(&configPath))
	rootCmd.AddCommand(serverCmd(&configPath))
	rootCmd.AddCommand(runCmd(&configPath))

	return rootCmd.Execute()
}
