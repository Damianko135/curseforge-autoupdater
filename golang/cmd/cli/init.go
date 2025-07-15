package main

import (
	"fmt"
	"os"

	"github.com/damianko135/curseforge-autoupdate/golang/helper/filesystem"
	"github.com/spf13/cobra"
)

func initCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "init [format]",
		Short: "Initialize a new project with configuration templates.",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			format := "toml"
			if len(args) > 0 {
				format = args[0]
			}
			switch format {
			case "toml", "yaml", "json", "yml":
				filename := "config." + format
				if filesystem.FileExists(filename) {
					return fmt.Errorf("%s already exists", filename)
				}
				templateName := "template." + format
				contentBytes, err := embeddedTemplates.ReadFile(templateName)
				if err != nil {
					return fmt.Errorf("failed to read embedded template: %w", err)
				}
				if err := os.WriteFile(filename, contentBytes, 0600); err != nil {
					return fmt.Errorf("failed to write %s: %w", filename, err)
				}
				fmt.Fprintf(cmd.OutOrStdout(), "✅ %s created.\n", filename)
				return nil
			case "":
				return fmt.Errorf("no config format specified (got empty string), please use one of: toml, yaml, json, yml, dotenv")
			default:
				return fmt.Errorf("unsupported format: %s (supported: toml, yaml, yml, json)", format)
			}
		},
	}
}
