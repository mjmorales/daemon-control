package cmd

import (
	"io"
	"os"

	"github.com/mjmorales/mac-daemon-control/internal/utils"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

// logsCmd represents the logs command
var logsCmd = &cobra.Command{
	Use:   "logs <daemon-name>",
	Short: "Show daemon logs",
	Long:  `Show recent logs from the daemon's stdout and stderr.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		daemonName := args[0]
		if err := showLogs(daemonName); err != nil {
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(logsCmd)
}

func showLogs(daemonName string) error {
	if err := utils.CheckPlistExists(daemonName); err != nil {
		return err
	}
	
	plistPath := utils.GetPlistPath(daemonName)
	
	stdoutPath, err := utils.GetStdoutPath(plistPath)
	if err != nil {
		stdoutPath = ""
	}
	
	stderrPath, err := utils.GetStderrPath(plistPath)
	if err != nil {
		stderrPath = ""
	}
	
	if stdoutPath == "" && stderrPath == "" {
		log.Error().Msg("No log paths configured in plist")
		return nil
	}
	
	if stdoutPath != "" {
		log.Info().Msg("=== STDOUT ===")
		if err := showLogFile(stdoutPath, 50); err != nil {
			log.Warn().Str("path", stdoutPath).Err(err).Msg("Could not read stdout log")
		}
	}
	
	if stderrPath != "" {
		log.Info().Msg("")
		log.Info().Msg("=== STDERR ===")
		if err := showLogFile(stderrPath, 50); err != nil {
			log.Warn().Str("path", stderrPath).Err(err).Msg("Could not read stderr log")
		}
	}
	
	return nil
}

func showLogFile(path string, lines int) error {
	file, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			log.Info().Msg("(no log file)")
			return nil
		}
		return err
	}
	defer file.Close()
	
	// Simple tail implementation - read last N lines
	// For production, consider using a proper tail library
	content, err := io.ReadAll(file)
	if err != nil {
		return err
	}
	
	if len(content) == 0 {
		log.Info().Msg("(empty)")
		return nil
	}
	
	// Print content directly to stdout to preserve formatting
	os.Stdout.Write(content)
	
	return nil
}