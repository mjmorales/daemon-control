package cmd

import (
	"os"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	"github.com/mjmorales/mac-daemon-control/internal/utils"
)

// stopCmd represents the stop command
var stopCmd = &cobra.Command{
	Use:   "stop <daemon-name>",
	Short: "Stop a daemon",
	Long:  `Stop a running daemon.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		daemonName := args[0]
		if err := stopDaemon(daemonName); err != nil {
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(stopCmd)
}

func stopDaemon(daemonName string) error {
	if err := utils.CheckPlistExists(daemonName); err != nil {
		return err
	}

	plistPath := utils.GetPlistPath(daemonName)
	label, err := utils.GetDaemonLabel(plistPath)
	if err != nil {
		log.Error().Err(err).Msg("Failed to read daemon label")
		return err
	}

	running, err := utils.IsRunning(daemonName)
	if err != nil {
		log.Error().Err(err).Msg("Failed to check running status")
		return err
	}

	if !running {
		log.Warn().Str("daemon", daemonName).Msg("Daemon not running")
		return nil
	}

	log.Info().Str("daemon", daemonName).Msg("Stopping daemon")

	if err := utils.RunLaunchctl("stop", label); err != nil {
		log.Error().Err(err).Msg("Failed to stop daemon")
		return err
	}

	log.Info().Str("daemon", daemonName).Msg("Daemon stopped")
	return nil
}
