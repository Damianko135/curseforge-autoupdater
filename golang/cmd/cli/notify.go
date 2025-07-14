package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

func newNotifyCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "notify",
		Short: "Send notifications manually.",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("[notify] Notification not yet implemented.")
		},
	}
}
