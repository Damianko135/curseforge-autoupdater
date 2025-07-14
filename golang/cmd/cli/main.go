package main

import (
	"fmt"
	"log"

	"github.com/damianko135/curseforge-autoupdate/golang/helper/env"
	"github.com/damianko135/curseforge-autoupdate/golang/internal/api"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type Config struct {
	APIKey string `mapstructure:"api_key"`
	ModID  int    `mapstructure:"mod_id"`
}

func main() {
	var configPath string
	var cfg Config

	rootCmd := &cobra.Command{
		Use:   "curseforge-autoupdate",
		Short: "A CLI tool to interact with CurseForge mods and configs.",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			// Skip config load for commands that opt-out
			if cmd.Annotations["skipConfig"] == "true" {
				return nil
			}

			if err := env.LoadDotenv(); err != nil {
				log.Printf("No .env loaded: %v", err)
			}
			if err := env.LoadTOMLConfig(configPath); err != nil {
				log.Printf("No config loaded from '%s': %v", configPath, err)
			}
			if err := viper.Unmarshal(&cfg); err != nil {
				return fmt.Errorf("failed to parse config: %w", err)
			}
			return nil
		},
	}

	rootCmd.PersistentFlags().StringVar(&configPath, "config", "config.toml", "Path to config file")

	rootCmd.AddCommand(
		newInitCmd(),
		newCheckCmd(&cfg),
		newUpdateCmd(),
		newBackupCmd(),
		newRestoreCmd(),
		newNotifyCmd(),
		newListCmd(),
		newVersionCmd(),
		newCreateConfigCmd(),
	)

	if err := rootCmd.Execute(); err != nil {
		log.Fatalf("Error: %v", err)
	}
}

func newCheckCmd(cfg *Config) *cobra.Command {
	return &cobra.Command{
		Use:   "check",
		Short: "Check if a mod exists using config/env variables.",
		Run: func(cmd *cobra.Command, args []string) {
			if cfg.APIKey == "" || cfg.ModID == 0 {
				log.Fatalf("Missing config: api_key='%s', mod_id='%d'. Hint: run `create-config` to scaffold one.", cfg.APIKey, cfg.ModID)
			}

			client := api.NewClient(cfg.APIKey)
			exists, err := client.CheckIfExists(cfg.ModID)
			if err != nil {
				log.Fatalf("Error checking mod: %v", err)
			}

			if exists {
				fmt.Printf("✅ Mod with ID %d found.\n", cfg.ModID)
			} else {
				fmt.Printf("❌ Mod with ID %d not found.\n", cfg.ModID)
			}
		},
	}
}

func newCreateConfigCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "create-config",
		Short: "Create a default config.toml file.",
		Annotations: map[string]string{
			"skipConfig": "true", // Skip loading config for this command
		},
		Run: func(cmd *cobra.Command, args []string) {
			if err := env.WriteTOMLTemplate("config.toml"); err != nil {
				log.Fatalf("Failed to create config.toml: %v", err)
			}
			fmt.Println("✅ config.toml created.")
		},
	}
}
