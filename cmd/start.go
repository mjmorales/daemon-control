package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	"github.com/mjmorales/mac-daemon-control/internal/utils"
)

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start <daemon-name>",
	Short: "Start a daemon",
	Long:  `Start a daemon that has been installed.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		daemonName := args[0]
		if err := startDaemon(daemonName); err != nil {
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(startCmd)
}

func startDaemon(daemonName string) error {
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
		log.Error().Str("daemon", daemonName).Msg("Daemon not installed. Run 'daemon-control install' first")
		return fmt.Errorf("daemon not installed")
	}

	running, err := utils.IsRunning(daemonName)
	if err != nil {
		log.Error().Err(err).Msg("Failed to check running status")
		return err
	}

	if running {
		log.Warn().Str("daemon", daemonName).Msg("Daemon already running")
		return nil
	}

	log.Info().Str("daemon", daemonName).Msg("Starting daemon")

	if err := utils.RunLaunchctl("start", label); err != nil {
		log.Error().Err(err).Msg("Failed to start daemon")
		return err
	}

	// Wait a moment and check status
	time.Sleep(2 * time.Second)

	running, err = utils.IsRunning(daemonName)
	if err != nil {
		log.Error().Err(err).Msg("Failed to check running status after start")
		return err
	}

	if running {
		log.Info().Str("daemon", daemonName).Msg("Daemon started successfully")
	} else {
		log.Error().Str("daemon", daemonName).Msg("Failed to start daemon. Check logs with 'daemon-control logs'")
		return fmt.Errorf("daemon failed to start")
	}

	return nil
}
