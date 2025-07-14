package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

func newBackupCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "backup",
		Short: "Manual backup operations.",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("[backup] Manual backup not yet implemented.")
		},
	}
}
