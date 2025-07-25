package cmd

import (
	"os"
	"os/exec"
	"strings"

	"github.com/mjmorales/mac-daemon-control/internal/utils"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

// statusCmd represents the status command
var statusCmd = &cobra.Command{
	Use:   "status <daemon-name>",
	Short: "Check daemon status",
	Long:  `Check the status of a daemon including installation and running state.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		daemonName := args[0]
		if err := checkStatus(daemonName); err != nil {
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(statusCmd)
}

func checkStatus(daemonName string) error {
	if err := utils.CheckPlistExists(daemonName); err != nil {
		return err
	}
	
	plistPath := utils.GetPlistPath(daemonName)
	label, err := utils.GetDaemonLabel(plistPath)
	if err != nil {
		log.Error().Err(err).Msg("Failed to read daemon label")
		return err
	}
	
	log.Info().Str("daemon", daemonName).Msg("Daemon status")
	log.Info().Str("label", label).Msg("Label")
	
	installed, err := utils.IsInstalled(daemonName)
	if err != nil {
		log.Error().Err(err).Msg("Failed to check installation status")
		return err
	}
	
	if installed {
		log.Info().Bool("installed", true).Msg("Installation status")
	} else {
		log.Warn().Bool("installed", false).Msg("Installation status")
	}
	
	running, err := utils.IsRunning(daemonName)
	if err != nil {
		log.Error().Err(err).Msg("Failed to check running status")
		return err
	}
	
	if running {
		log.Info().Bool("running", true).Msg("Running status")
		
		// Get process info
		cmd := exec.Command("launchctl", "list")
		output, err := cmd.Output()
		if err == nil {
			lines := strings.Split(string(output), "\n")
			for _, line := range lines {
				if strings.Contains(line, label) {
					log.Info().Str("process_info", strings.TrimSpace(line)).Msg("Process details")
					break
				}
			}
		}
	} else {
		log.Warn().Bool("running", false).Msg("Running status")
	}
	
	// Show additional info from plist
	workingDir, err := utils.GetWorkingDirectory(plistPath)
	if err == nil && workingDir != "" {
		log.Info().Str("working_directory", workingDir).Msg("Working directory")
	}
	
	return nil
}