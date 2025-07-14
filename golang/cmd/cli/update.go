package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

func newUpdateCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "update",
		Short: "Perform the full update process.",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("[update] Update process not yet implemented.")
		},
	}
}
