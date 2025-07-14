package main

import (
	"fmt"
	"log"
	"os"

	"github.com/damianko135/curseforge-autoupdate/golang/helper/env"
	"github.com/damianko135/curseforge-autoupdate/golang/internal/api"
	"github.com/damianko135/curseforge-autoupdate/golang/templates"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var embeddedTemplates = templates.EmbeddedTemplates

var verbose bool

type Config struct {
	APIKey string `mapstructure:"api_key"`
	ModID  int    `mapstructure:"mod_id"`
}

// getConfigValue tries config, then env var, then default
func getConfigValue(key string, defaultVal string) string {
	val := viper.GetString(key)
	if val != "" {
		return val
	}
	envVal := os.Getenv(key)
	if envVal != "" {
		return envVal
	}
	return defaultVal
}

func main() {
	var configPath string
	var initFormat string
	var cfg Config

	rootCmd := newRootCmd(&cfg, &configPath, &initFormat)

	// Check for --init flag before executing rootCmd
	for i, arg := range os.Args {
		if arg == "--init" || arg == "-init" || arg == "-i" || arg == "--initialize" {
			if i+1 < len(os.Args) {
				initFormat = os.Args[i+1]
			}
			if initFormat != "" {
				if err := runInitDirect(initFormat); err != nil {
					fmt.Fprintf(os.Stderr, "[init] %v\n", err)
					os.Exit(1)
				}
				os.Exit(0)
			}
		}
	}

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

// newRootCmd sets up the root command and all subcommands
func newRootCmd(cfg *Config, configPath *string, initFormat *string) *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "curseforge-autoupdater",
		Short: "A CLI tool to interact with CurseForge mods and configs.",
	}

	rootCmd.PersistentFlags().StringVar(configPath, "config", "config.toml", "Path to config file")
	rootCmd.PersistentFlags().StringVar(initFormat, "init", "", "Initialize a new project with configuration templates (e.g. --init toml)")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Enable verbose output")

	// Register commands
	rootCmd.AddCommand(
		newInitCmd(),
		newCheckCmd(cfg),
		newUpdateCmd(),
		newBackupCmd(),
		newRestoreCmd(),
		newNotifyCmd(),
		newListCmd(),
		newVersionCmd(),
		newCreateConfigCmd(),
		newHelpCmd(rootCmd),
		newCompletionCmd(rootCmd),
	)

	// Only load config for commands that need it
	rootCmd.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
		// Skip config load for commands that opt-out
		if cmd.Annotations["skipConfig"] == "true" {
			return nil
		}

		if *configPath == "" {
			*configPath = "config.toml" // Default config path
		}

		if err := env.LoadConfig(*configPath); err != nil {
			// Fallback: prompt to create config
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
					os.Exit(0)
				} else {
					return fmt.Errorf("config file required, exiting")
				}
			}
			return fmt.Errorf("failed to load config: %w", err)
		}

		if err := viper.Unmarshal(cfg); err != nil {
			return fmt.Errorf("failed to unmarshal config: %w", err)
		}

		// Fallback to env vars if missing
		if cfg.APIKey == "" {
			cfg.APIKey = getConfigValue("API_KEY", "")
		}
		if cfg.ModID == 0 {
			modIDStr := getConfigValue("MOD_ID", "0")
			if _, err := fmt.Sscanf(modIDStr, "%d", &cfg.ModID); err != nil {
				return fmt.Errorf("failed to parse MOD_ID: %w", err)
			}
		}

		return nil
	}

	return rootCmd
}

// runInitDirect runs the init logic for --init flag (cross-platform)
func runInitDirect(format string) error {

	switch format {
	case "toml", "yaml", "json", "yml":
		filename := "config." + format
		if _, err := os.Stat(filename); err == nil {
			return fmt.Errorf("%s already exists", filename)
		}
		templateName := "templates/template." + format
		contentBytes, err := embeddedTemplates.ReadFile(templateName)
		if err != nil {
			return fmt.Errorf("failed to read embedded template: %w", err)
		}
		if err := os.WriteFile(filename, contentBytes, 0600); err != nil {
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
		Use:     "check",
		Aliases: []string{"verify"},
		Short:   "Check if a mod exists using config/env variables.",
		Run: func(cmd *cobra.Command, args []string) {
			if cfg.APIKey == "" || cfg.ModID == 0 {
				fmt.Fprintf(os.Stderr, "Missing config: api_key='%s', mod_id='%d'. Hint: run `create-config` to scaffold one.\n", cfg.APIKey, cfg.ModID)
				os.Exit(1)
			}

			client := api.NewClient(cfg.APIKey)
			exists, err := client.CheckIfExists(cfg.ModID)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error checking mod: %v\n", err)
				os.Exit(1)
			}

			if exists {
				fmt.Printf("✅ Mod with ID %d found.\n", cfg.ModID)
			} else {
				fmt.Printf("❌ Mod with ID %d not found.\n", cfg.ModID)
			}
		},
	}
}

// newHelpCmd adds a help command for better UX
func newHelpCmd(root *cobra.Command) *cobra.Command {
	return &cobra.Command{
		Use:   "help",
		Short: "Show help for any command",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) > 0 {
				c, _, err := root.Find(args)
				if err == nil && c != nil {
					if err := c.Help(); err != nil {
						fmt.Fprintf(os.Stderr, "Help error: %v\n", err)
					}
					return
				}
			}
			if err := root.Help(); err != nil {
				fmt.Fprintf(os.Stderr, "Help error: %v\n", err)
			}
		},
	}
}

// newCompletionCmd adds shell completion generation
func newCompletionCmd(root *cobra.Command) *cobra.Command {
	return &cobra.Command{
		Use:   "completion",
		Short: "Generate shell completion scripts",
		Long:  "Generate shell completion scripts for bash, zsh, fish, or PowerShell.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			switch args[0] {
			case "bash":
				return root.GenBashCompletion(os.Stdout)
			case "zsh":
				return root.GenZshCompletion(os.Stdout)
			case "fish":
				return root.GenFishCompletion(os.Stdout, true)
			case "powershell":
				return root.GenPowerShellCompletionWithDesc(os.Stdout)
			default:
				return fmt.Errorf("unsupported shell: %s", args[0])
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
