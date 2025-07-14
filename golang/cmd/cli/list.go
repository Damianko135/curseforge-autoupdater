package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

func newListCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List available commands/info.",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("[list] List not yet implemented.")
		},
	}
}
