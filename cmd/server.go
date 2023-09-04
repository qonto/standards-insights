package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func serverCmd(configPath *string) *cobra.Command {
	provider := ""
	cmd := &cobra.Command{
		Use:   "server",
		Short: "Run server",
		Run: func(cmd *cobra.Command, args []string) {
			// TODO: fix me
			fmt.Println("Not implemented yet")
		},
	}

	// TODO: allow to repeat flag
	cmd.PersistentFlags().StringVarP(&provider, "provider", "p", "", "Filter providers")

	return cmd
}
