package cmd

import (
	"os"
	"path/filepath"

	"github.com/mjmorales/mac-daemon-control/internal/config"
	"github.com/mjmorales/mac-daemon-control/internal/plist"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
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
	
	generateCmd.Flags().StringVarP(&configFile, "config", "c", "", "Configuration file path (default: ./daemons.yaml)")
	generateCmd.Flags().StringVarP(&outputDir, "output", "o", "out", "Output directory for generated plist files")
}

func runGenerate() error {
	// Load configuration
	loader := config.NewLoader(configFile)
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
		Str("output", outputDir).
		Msg("Generating plist files")
	
	// Create generator
	generator := plist.NewGenerator(outputDir)
	
	// Generate plist files
	if err := generator.GenerateAll(cfg.Daemons); err != nil {
		return err
	}
	
	log.Info().
		Int("count", len(cfg.Daemons)).
		Str("directory", outputDir).
		Msg("Successfully generated plist files")
	
	// Also update the daemons directory if it exists
	if _, err := os.Stat("daemons"); err == nil {
		log.Info().Msg("Updating daemons directory...")
		for _, daemon := range cfg.Daemons {
			src := filepath.Join(outputDir, daemon.Name+".plist")
			dst := filepath.Join("daemons", daemon.Name+".plist")
			
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
			if err := os.WriteFile(dst, data, 0644); err != nil {
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
	
	return nil
}