package cmd

import (
	"os"
	"path/filepath"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	"github.com/mjmorales/mac-daemon-control/internal/config"
	"github.com/mjmorales/mac-daemon-control/internal/core"
	"github.com/mjmorales/mac-daemon-control/internal/plist"
	"github.com/mjmorales/mac-daemon-control/internal/utils"
)

var (
	configFile string
	outputDir  string
)

// generateCmd represents the generate command
var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate plist files from configuration",
	Long: `Generate macOS plist files from a daemon configuration file.
	
This command reads a YAML configuration file containing daemon definitions
and generates corresponding plist files that can be used with launchd.`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := runGenerate(); err != nil {
			log.Error().Err(err).Msg("Failed to generate plist files")
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(generateCmd)

	generateCmd.Flags().StringVarP(&configFile, "config", "c", "", "Configuration file path (default: from core config)")
	generateCmd.Flags().StringVarP(&outputDir, "output", "o", "", "Output directory for generated plist files (default: from core config)")
}

func runGenerate() error {
	// Get core config
	coreManager := core.GetManager()
	coreConfig := coreManager.GetConfig()

	// Determine config file path
	configPath := configFile
	if configPath == "" {
		configPath = coreManager.GetDaemonConfigPath()
	}

	// Determine output directory
	outDir := outputDir
	if outDir == "" && coreConfig != nil {
		outDir = coreConfig.OutputDir
	}
	if outDir == "" {
		outDir = "out"
	}

	// Load daemon configuration
	loader := config.NewLoader(configPath)
	cfg, err := loader.Load()
	if err != nil {
		return err
	}

	if len(cfg.Daemons) == 0 {
		log.Warn().Msg("No daemons defined in configuration")
		return nil
	}

	log.Info().
		Int("count", len(cfg.Daemons)).
		Str("output", outDir).
		Msg("Generating plist files")

	// Create generator
	generator := plist.NewGenerator(outDir)

	// Generate plist files
	if err := generator.GenerateAll(cfg.Daemons); err != nil {
		return err
	}

	log.Info().
		Int("count", len(cfg.Daemons)).
		Str("directory", outDir).
		Msg("Successfully generated plist files")

	// Also update the daemons directory if configured
	daemonsDir := utils.GetDaemonsDir()
	if coreConfig != nil && coreConfig.AutoGeneratePlists {
		if _, err := os.Stat(daemonsDir); err == nil {
			log.Info().Msg("Updating daemons directory...")

			// Backup existing plists if configured
			if coreConfig.BackupOnGenerate {
				backupDir := filepath.Join(daemonsDir, ".backup")
				if err := os.MkdirAll(backupDir, 0750); err == nil {
					// Copy existing plists to backup
					files, _ := filepath.Glob(filepath.Join(daemonsDir, "*.plist"))
					for _, file := range files {
						base := filepath.Base(file)
						dst := filepath.Join(backupDir, base+".bak")
						data, err := os.ReadFile(file)
						if err != nil {
							log.Warn().Err(err).Str("file", file).Msg("Failed to read file for backup")
							continue
						}
						if err := os.WriteFile(dst, data, 0600); err != nil {
							log.Warn().Err(err).Str("dst", dst).Msg("Failed to write backup file")
						}
					}
					log.Info().Str("dir", backupDir).Msg("Backed up existing plists")
				}
			}

			// Copy generated plists
			for _, daemon := range cfg.Daemons {
				src := filepath.Join(outDir, daemon.Name+".plist")
				dst := filepath.Join(daemonsDir, daemon.Name+".plist")

				// Read the generated file
				data, err := os.ReadFile(src)
				if err != nil {
					log.Error().
						Err(err).
						Str("daemon", daemon.Name).
						Msg("Failed to read generated plist")
					continue
				}

				// Write to daemons directory
				if err := os.WriteFile(dst, data, 0600); err != nil {
					log.Error().
						Err(err).
						Str("daemon", daemon.Name).
						Msg("Failed to copy plist to daemons directory")
					continue
				}

				log.Info().
					Str("daemon", daemon.Name).
					Str("path", dst).
					Msg("Copied plist to daemons directory")
			}
		}
	}

	return nil
}
