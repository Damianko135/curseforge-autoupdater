package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

func updateCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "update",
		Short: "Perform the full update process.",
		Run: func(cmd *cobra.Command, args []string) {
			// TODO: Replace with the update logic
			fmt.Fprintln(cmd.OutOrStdout(), "[update] Update process not yet implemented.")
		},
	}
}
