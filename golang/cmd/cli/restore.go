package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

func restoreCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "restore",
		Short: "Restore from backup.",
		Run: func(cmd *cobra.Command, args []string) {
			// TODO: Replace with the restore logic
			fmt.Fprintln(cmd.OutOrStdout(), "[restore] Restore not yet implemented.")
		},
	}
}
