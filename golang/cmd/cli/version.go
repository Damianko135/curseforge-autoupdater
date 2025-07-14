package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

func newVersionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Show version info.",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("CurseForge Auto-Update CLI v0.1.0 (dev)")
		},
	}
}
