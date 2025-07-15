package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

func backupCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "backup",
		Short: "Manual backup operations.",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Fprintln(cmd.OutOrStdout(), "[backup] Manual backup not yet implemented.")
		},
	}
}
