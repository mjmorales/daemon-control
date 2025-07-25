package cmd

import (
	"os"
	"path/filepath"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	"github.com/mjmorales/mac-daemon-control/internal/utils"
)

// installCmd represents the install command
var installCmd = &cobra.Command{
	Use:   "install <daemon-name>",
	Short: "Install a daemon",
	Long:  `Install a daemon by copying its plist file to LaunchAgents directory.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		daemonName := args[0]
		if err := installDaemon(daemonName); err != nil {
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(installCmd)
}

func installDaemon(daemonName string) error {
	if err := utils.CheckPlistExists(daemonName); err != nil {
		return err
	}

	plistPath := utils.GetPlistPath(daemonName)
	label, err := utils.GetDaemonLabel(plistPath)
	if err != nil {
		log.Error().Err(err).Msg("Failed to read daemon label")
		return err
	}

	installed, err := utils.IsInstalled(daemonName)
	if err != nil {
		log.Error().Err(err).Msg("Failed to check installation status")
		return err
	}

	if installed {
		log.Warn().Str("daemon", daemonName).Msg("Daemon already installed")
		return nil
	}

	log.Info().Str("daemon", daemonName).Msg("Installing daemon")

	// Create LaunchAgents directory if it doesn't exist
	if err := os.MkdirAll(utils.LaunchAgentsDir, 0750); err != nil {
		log.Error().Err(err).Msg("Failed to create LaunchAgents directory")
		return err
	}

	// Copy plist file
	destPath := filepath.Join(utils.LaunchAgentsDir, label+".plist")
	if err := utils.CopyFile(plistPath, destPath); err != nil {
		log.Error().Err(err).Msg("Failed to copy plist file")
		return err
	}

	// Load the daemon
	if err := utils.RunLaunchctl("load", destPath); err != nil {
		log.Error().Err(err).Msg("Failed to load daemon")
		return err
	}

	log.Info().Str("daemon", daemonName).Msg("Daemon installed successfully")
	return nil
}
