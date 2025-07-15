package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

func notifyCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "notify",
		Short: "Send notifications manually.",
		Run: func(cmd *cobra.Command, args []string) {
			// TODO: Replace with the notification logic
			fmt.Fprintln(cmd.OutOrStdout(), "[notify] Notification not yet implemented.")
		},
	}
}
