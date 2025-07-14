package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

func newRestoreCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "restore",
		Short: "Restore from backup.",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("[restore] Restore not yet implemented.")
		},
	}
}
