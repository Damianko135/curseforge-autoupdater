package main

import (
	"fmt"
	"os"

	"github.com/damianko135/curseforge-autoupdate/golang/helper/env"
	"github.com/damianko135/curseforge-autoupdate/golang/internal/api"
	"github.com/damianko135/curseforge-autoupdate/golang/templates"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	embeddedTemplates = templates.EmbeddedTemplates
	verboseMode       bool
)

type Config struct {
	APIToken string `mapstructure:"api_key"`
	ModID    int    `mapstructure:"mod_id"`
}

// getConfigValue tries config, then env var, then default
func getConfigValue(key, defaultVal string) string {
	if val := viper.GetString(key); val != "" {
		return val
	}
	if envVal := os.Getenv(key); envVal != "" {
		return envVal
	}
	return defaultVal
}

func main() {
	var (
		configFilePath     string
		initTemplateFormat string
		userConfig         Config
	)

	rootCmd, err := setupRootCommand(&userConfig, &configFilePath, &initTemplateFormat)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error setting up root command: %v\n", err)
		return
	}

	// Let Cobra handle all CLI parsing, config, and command dispatching
	// All logic for --init, --config, --verbose, --version, etc. is now handled by the registered commands and PersistentPreRunE
	// This makes the CLI idiomatic and ensures all subcommands in cmd/cli are used

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		return
	}
}

// newRootCmd sets up the root command and all subcommands
func setupRootCommand(cfg *Config, configPath *string, initFormat *string) (*cobra.Command, error) {
	rootCmd := &cobra.Command{
		Use:   "curseforge-autoupdater",
		Short: "A CLI tool to interact with CurseForge mods and configs.",
	}

	rootCmd.PersistentFlags().StringVar(configPath, "config", "config.toml", "Path to config file")
	rootCmd.PersistentFlags().StringVar(initFormat, "init", "", "Initialize a new project with configuration templates (e.g. --init toml)")
	rootCmd.PersistentFlags().BoolVarP(&verboseMode, "verbose", "v", false, "Enable verbose output")

	// Register only essential top-level commands
	rootCmd.AddCommand(
		checkCmd(cfg),
		updateCmd(),
		backupCmd(),
		restoreCmd(),
		notifyCmd(),
		listCmd(),
		versionCmd(),
		initCmd(),
	)

	// Only load config for commands that need it
	rootCmd.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
		if cmd.Annotations["skipConfig"] == "true" {
			return nil
		}
		if *configPath == "" {
			*configPath = "config.toml"
		}
		if err := env.LoadConfig(*configPath); err != nil {
			if os.IsNotExist(err) {
				fmt.Printf("Config file '%s' not found. Would you like to create one? [Y/n]: ", *configPath)
				var resp string
				if _, err := fmt.Scanln(&resp); err != nil && err.Error() != "unexpected newline" {
					return fmt.Errorf("failed to read input: %w", err)
				}
				if resp == "" || resp == "y" || resp == "Y" {
					if err := env.WriteTOMLTemplate(*configPath); err != nil {
						return fmt.Errorf("failed to create config: %w", err)
					}
					fmt.Printf("Created %s. Please edit it and re-run.\n", *configPath)
					return nil
				} else {
					return fmt.Errorf("config file required: %s (user declined to create)", *configPath)
				}
			}
			return fmt.Errorf("failed to load config: %w", err)
		}
		if err := viper.Unmarshal(cfg); err != nil {
			return fmt.Errorf("failed to read values: %w", err)
		}
		if cfg.APIToken == "" {
			cfg.APIToken = getConfigValue("API_KEY", "")
		}
		if cfg.ModID == 0 {
			modIDStr := getConfigValue("MOD_ID", "0")
			if _, err := fmt.Sscanf(modIDStr, "%d", &cfg.ModID); err != nil {
				return fmt.Errorf("failed to parse MOD_ID: %w", err)
			}
		}
		return nil
	}
	return rootCmd, nil
}

func checkCmd(cfg *Config) *cobra.Command {
	return &cobra.Command{
		Use:     "check",
		Aliases: []string{"verify"},
		Short:   "Check if a mod exists using config/env variables.",
		Run: func(cmd *cobra.Command, args []string) {
			if cfg.APIToken == "" || cfg.ModID == 0 {
				fmt.Fprintf(os.Stderr, "Missing config: api_key='%s', mod_id='%d'. Hint: run `init` to scaffold one.\n", cfg.APIToken, cfg.ModID)
				return
			}

			client := api.NewClient(cfg.APIToken)
			exists, err := client.CheckIfExists(cfg.ModID)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error checking mod: %v\n", err)
				return
			}

			if exists {
				fmt.Printf("✅ Mod with ID %d found.\n", cfg.ModID)
			} else {
				fmt.Printf("❌ Mod with ID %d not found.\n", cfg.ModID)
			}
		},
	}
}
