package main

import (
	"fmt"

	"github.com/damianko135/curseforge-autoupdate/internal/config"
	"github.com/damianko135/curseforge-autoupdate/pkg/curseforge"
	"github.com/damianko135/curseforge-autoupdate/pkg/models"
	"github.com/muesli/coral"
	"github.com/sirupsen/logrus"
)

var (
	version = "1.0.0"
	logger  = logrus.New()
)

func main() {
	var rootCmd = &coral.Command{
		Use:     "curseforge-updater",
		Short:   "CurseForge mod auto-updater",
		Long:    "A CLI tool to automatically check and download the latest versions of CurseForge mods",
		Version: version,
		RunE:    runUpdate,
	}

	// Add flags
	rootCmd.PersistentFlags().String("api-key", "", "CurseForge API key")
	rootCmd.PersistentFlags().String("mod-id", "", "Mod ID to check for updates")
	rootCmd.PersistentFlags().String("download-path", "./downloads", "Path to download files")
	rootCmd.PersistentFlags().Int("game-id", 432, "Game ID (default: 432 for Minecraft)")
	rootCmd.PersistentFlags().String("log-level", "info", "Log level (debug, info, warn, error)")
	rootCmd.PersistentFlags().String("config", "", "Path to config file")

	// Add subcommands
	rootCmd.AddCommand(checkCmd())
	rootCmd.AddCommand(downloadCmd())
	rootCmd.AddCommand(infoCmd())

	if err := rootCmd.Execute(); err != nil {
		logger.Fatal(err)
	}
}

func setupLogger(level string) {
	switch level {
	case "debug":
		logger.SetLevel(logrus.DebugLevel)
	case "info":
		logger.SetLevel(logrus.InfoLevel)
	case "warn":
		logger.SetLevel(logrus.WarnLevel)
	case "error":
		logger.SetLevel(logrus.ErrorLevel)
	default:
		logger.SetLevel(logrus.InfoLevel)
	}

	logger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
		ForceColors:   true,
	})
}

func runUpdate(cmd *coral.Command, args []string) error {
	cfg, err := config.LoadConfig()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	setupLogger(cfg.LogLevel)
	logger.Info("Starting CurseForge Auto-Updater")

	client := curseforge.NewClient(cfg.APIKey, logger)

	// Get mod info
	modInfo, err := client.GetModInfo(cfg.ModID)
	if err != nil {
		return fmt.Errorf("failed to get mod info: %w", err)
	}

	logger.Infof("Checking updates for mod: %s", modInfo.Name)

	// Get mod files
	files, err := client.GetModFiles(cfg.ModID, cfg.GameID)
	if err != nil {
		return fmt.Errorf("failed to get mod files: %w", err)
	}

	if len(files) == 0 {
		logger.Warn("No files found for this mod")
		return nil
	}

	// Get latest file
	latestFile := client.GetLatestFile(files)
	if latestFile == nil {
		logger.Warn("No latest file found")
		return nil
	}

	logger.Infof("Latest file: %s (%s)", latestFile.FileName, latestFile.FileDate)

	// Load metadata
	metadata, err := curseforge.LoadDownloadMetadata(cfg.DownloadPath)
	if err != nil {
		return fmt.Errorf("failed to load download metadata: %w", err)
	}

	// Check if download is needed
	needsDownload, reason := curseforge.IsDownloadNeeded(latestFile, cfg.DownloadPath, metadata, logger)

	if needsDownload {
		logger.Infof("Download needed: %s", reason)

		if err := client.DownloadFile(latestFile, cfg.DownloadPath); err != nil {
			return fmt.Errorf("failed to download file: %w", err)
		}

		if err := curseforge.RecordDownload(latestFile, cfg.DownloadPath, metadata, logger); err != nil {
			return fmt.Errorf("failed to record download: %w", err)
		}

		logger.Info("Update completed successfully!")
	} else {
		logger.Infof("No download needed: %s", reason)
		logger.Info("Everything is up to date!")
	}

	return nil
}

func checkCmd() *coral.Command {
	return &coral.Command{
		Use:   "check",
		Short: "Check for updates without downloading",
		RunE: func(cmd *coral.Command, args []string) error {
			cfg, err := config.LoadConfig()
			if err != nil {
				return fmt.Errorf("failed to load config: %w", err)
			}

			setupLogger(cfg.LogLevel)
			client := curseforge.NewClient(cfg.APIKey, logger)

			modInfo, err := client.GetModInfo(cfg.ModID)
			if err != nil {
				return fmt.Errorf("failed to get mod info: %w", err)
			}

			files, err := client.GetModFiles(cfg.ModID, cfg.GameID)
			if err != nil {
				return fmt.Errorf("failed to get mod files: %w", err)
			}

			if len(files) == 0 {
				logger.Warn("No files found for this mod")
				return nil
			}

			latestFile := client.GetLatestFile(files)
			if latestFile == nil {
				logger.Warn("No latest file found")
				return nil
			}

			metadata, err := curseforge.LoadDownloadMetadata(cfg.DownloadPath)
			if err != nil {
				return fmt.Errorf("failed to load download metadata: %w", err)
			}

			needsDownload, reason := curseforge.IsDownloadNeeded(latestFile, cfg.DownloadPath, metadata, logger)

			fmt.Printf("Mod: %s\n", modInfo.Name)
			fmt.Printf("Latest File: %s\n", latestFile.FileName)
			fmt.Printf("File Date: %s\n", latestFile.FileDate)
			fmt.Printf("File Size: %d bytes\n", latestFile.FileLength)
			fmt.Printf("Update Needed: %t\n", needsDownload)
			if needsDownload {
				fmt.Printf("Reason: %s\n", reason)
			}

			return nil
		},
	}
}

func downloadCmd() *coral.Command {
	return &coral.Command{
		Use:   "download",
		Short: "Force download the latest file",
		RunE: func(cmd *coral.Command, args []string) error {
			cfg, err := config.LoadConfig()
			if err != nil {
				return fmt.Errorf("failed to load config: %w", err)
			}

			setupLogger(cfg.LogLevel)
			client := curseforge.NewClient(cfg.APIKey, logger)

			files, err := client.GetModFiles(cfg.ModID, cfg.GameID)
			if err != nil {
				return fmt.Errorf("failed to get mod files: %w", err)
			}

			if len(files) == 0 {
				logger.Warn("No files found for this mod")
				return nil
			}

			latestFile := client.GetLatestFile(files)
			if latestFile == nil {
				logger.Warn("No latest file found")
				return nil
			}

			if err := client.DownloadFile(latestFile, cfg.DownloadPath); err != nil {
				return fmt.Errorf("failed to download file: %w", err)
			}

			metadata, err := curseforge.LoadDownloadMetadata(cfg.DownloadPath)
			if err != nil {
				logger.Warnf("Failed to load metadata, creating new: %v", err)
				metadata = make(map[string]models.DownloadMetadata)
			}

			if err := curseforge.RecordDownload(latestFile, cfg.DownloadPath, metadata, logger); err != nil {
				return fmt.Errorf("failed to record download: %w", err)
			}

			logger.Info("Download completed successfully!")
			return nil
		},
	}
}

func infoCmd() *coral.Command {
	return &coral.Command{
		Use:   "info",
		Short: "Show mod information",
		RunE: func(cmd *coral.Command, args []string) error {
			cfg, err := config.LoadConfig()
			if err != nil {
				return fmt.Errorf("failed to load config: %w", err)
			}

			setupLogger(cfg.LogLevel)
			client := curseforge.NewClient(cfg.APIKey, logger)

			modInfo, err := client.GetModInfo(cfg.ModID)
			if err != nil {
				return fmt.Errorf("failed to get mod info: %w", err)
			}

			files, err := client.GetModFiles(cfg.ModID, cfg.GameID)
			if err != nil {
				return fmt.Errorf("failed to get mod files: %w", err)
			}

			fmt.Printf("Mod Information:\n")
			fmt.Printf("  Name: %s\n", modInfo.Name)
			fmt.Printf("  ID: %d\n", modInfo.ID)
			fmt.Printf("  Game ID: %d\n", modInfo.GameID)
			fmt.Printf("  Class ID: %d\n", modInfo.ClassID)
			fmt.Printf("  Authors: ")
			for i, author := range modInfo.Authors {
				if i > 0 {
					fmt.Printf(", ")
				}
				fmt.Printf("%s", author.Name)
			}
			fmt.Printf("\n")
			fmt.Printf("  Total Files: %d\n", len(files))

			if len(files) > 0 {
				latestFile := client.GetLatestFile(files)
				if latestFile != nil {
					fmt.Printf("  Latest File: %s (%s)\n", latestFile.FileName, latestFile.FileDate)
				}
			}

			return nil
		},
	}
}
