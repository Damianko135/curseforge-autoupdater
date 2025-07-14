package main

import (
	"embed"
	"fmt"
	"log"
	"os"

	"github.com/damianko135/curseforge-autoupdate/golang/helper/env"
	"github.com/damianko135/curseforge-autoupdate/golang/internal/api"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

//go:embed templates/init.toml templates/init.yaml templates/init.json templates/init.yml templates/init.env
var embeddedTemplates embed.FS

type Config struct {
	APIKey string `mapstructure:"api_key"`
	ModID  int    `mapstructure:"mod_id"`
}

func main() {
	var configPath string
	var initFormat string
	var cfg Config

	rootCmd := &cobra.Command{
		Use:   "curseforge-autoupdater",
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
	rootCmd.PersistentFlags().StringVar(&initFormat, "init", "", "Initialize a new project with configuration templates (e.g. --init toml)")

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

	// Check for --init flag before executing rootCmd
	for i, arg := range os.Args {
		if arg == "--init" || arg == "-init" || arg == "-i" || arg == "--initialize" {
			if i+1 < len(os.Args) {
				initFormat = os.Args[i+1]
			}
			if initFormat != "" {
				if err := runInitDirect(initFormat); err != nil {
					log.Fatalf("[init] %v", err)
				}
				os.Exit(0)
			}
		}
	}

	if err := rootCmd.Execute(); err != nil {
		log.Fatalf("Error: %v", err)
	}
}

// runInitDirect runs the init logic for --init flag (cross-platform)
func runInitDirect(format string) error {

	switch format {
	case "toml", "yaml", "json", "yml":
		filename := "config." + format
		if _, err := os.Stat(filename); err == nil {
			return fmt.Errorf("%s already exists", filename)
		}
		templateName := "templates/init." + format
		contentBytes, err := embeddedTemplates.ReadFile(templateName)
		if err != nil {
			return fmt.Errorf("failed to read embedded template: %w", err)
		}
		if err := os.WriteFile(filename, contentBytes, 0644); err != nil {
			return fmt.Errorf("failed to write %s: %w", filename, err)
		}
		fmt.Printf("✅ %s created.\n", filename)
		return nil
	case "":
		return fmt.Errorf("no format specified for --init, please use one of: toml, yaml, json, yml, dotenv")
	default:
		return fmt.Errorf("unsupported format: %s (supported: toml, yaml, yml, json)", format)
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
