package cmd

import (
	"os"
	"path/filepath"

	"github.com/mjmorales/mac-daemon-control/internal/utils"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

// uninstallCmd represents the uninstall command
var uninstallCmd = &cobra.Command{
	Use:   "uninstall <daemon-name>",
	Short: "Uninstall a daemon",
	Long:  `Uninstall a daemon by removing its plist file from LaunchAgents directory.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		daemonName := args[0]
		if err := uninstallDaemon(daemonName); err != nil {
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(uninstallCmd)
}

func uninstallDaemon(daemonName string) error {
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
	
	if !installed {
		log.Warn().Str("daemon", daemonName).Msg("Daemon not installed")
		return nil
	}
	
	log.Info().Str("daemon", daemonName).Msg("Uninstalling daemon")
	
	// Stop if running
	running, err := utils.IsRunning(daemonName)
	if err != nil {
		log.Error().Err(err).Msg("Failed to check running status")
		return err
	}
	
	installedPath := filepath.Join(utils.LaunchAgentsDir, label+".plist")
	
	if running {
		if err := utils.RunLaunchctl("unload", installedPath); err != nil {
			log.Error().Err(err).Msg("Failed to unload daemon")
			return err
		}
	}
	
	// Remove plist
	if err := os.Remove(installedPath); err != nil {
		log.Error().Err(err).Msg("Failed to remove plist file")
		return err
	}
	
	log.Info().Str("daemon", daemonName).Msg("Daemon uninstalled successfully")
	return nil
}