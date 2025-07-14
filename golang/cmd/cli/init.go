package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

func newInitCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "init",
		Short: "Initialize a new project with configuration templates.",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("[init] Project initialization not yet implemented.")
		},
	}
}
