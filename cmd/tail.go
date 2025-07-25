package cmd

import (
	"context"
	"os"
	"os/exec"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	"github.com/mjmorales/daemon-control/internal/utils"
)

// tailCmd represents the tail command
var tailCmd = &cobra.Command{
	Use:   "tail <daemon-name>",
	Short: "Tail daemon logs",
	Long:  `Tail daemon logs in real-time.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		daemonName := args[0]
		if err := tailLogs(daemonName); err != nil {
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(tailCmd)
}

func tailLogs(daemonName string) error {
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

	log.Info().Str("daemon", daemonName).Msg("Tailing logs (Ctrl+C to stop)...")

	// Build tail command args
	var files []string
	if stdoutPath != "" {
		if _, err := os.Stat(stdoutPath); err == nil {
			files = append(files, stdoutPath)
		}
	}
	if stderrPath != "" {
		if _, err := os.Stat(stderrPath); err == nil {
			files = append(files, stderrPath)
		}
	}

	if len(files) == 0 {
		log.Error().Msg("No log files found")
		return nil
	}

	// Use tail command with context
	ctx := context.Background()
	args := append([]string{"-f"}, files...)
	cmd := exec.CommandContext(ctx, "tail", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	return cmd.Run()
}
