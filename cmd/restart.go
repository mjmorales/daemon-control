package cmd

import (
	"os"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

// restartCmd represents the restart command
var restartCmd = &cobra.Command{
	Use:   "restart <daemon-name>",
	Short: "Restart a daemon",
	Long:  `Stop and then start a daemon.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		daemonName := args[0]
		if err := restartDaemon(daemonName); err != nil {
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(restartCmd)
}

func restartDaemon(daemonName string) error {
	log.Info().Str("daemon", daemonName).Msg("Restarting daemon")
	
	// Stop the daemon
	if err := stopDaemon(daemonName); err != nil {
		return err
	}
	
	// Wait a moment
	time.Sleep(2 * time.Second)
	
	// Start the daemon
	if err := startDaemon(daemonName); err != nil {
		return err
	}
	
	return nil
}