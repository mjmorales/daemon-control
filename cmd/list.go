package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/mjmorales/mac-daemon-control/internal/utils"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all available daemons",
	Long:  `List all available daemons in the ./daemons directory.`,
	Run: func(cmd *cobra.Command, args []string) {
		listDaemons()
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}

func listDaemons() {
	log.Info().Msg("Available daemons:")
	fmt.Println("-----------------")

	if _, err := os.Stat(utils.DaemonsDir); os.IsNotExist(err) {
		log.Error().Str("path", utils.DaemonsDir).Msg("Daemons directory not found")
		os.Exit(1)
	}

	files, err := filepath.Glob(filepath.Join(utils.DaemonsDir, "*.plist"))
	if err != nil {
		log.Error().Err(err).Msg("Failed to read daemons directory")
		os.Exit(1)
	}

	if len(files) == 0 {
		log.Warn().Str("path", utils.DaemonsDir).Msg("No daemons found")
		return
	}

	for _, file := range files {
		base := filepath.Base(file)
		daemonName := strings.TrimSuffix(base, ".plist")

		label, err := utils.GetDaemonLabel(file)
		if err != nil {
			label = "Unknown"
		}

		log.Info().Str("daemon", daemonName).Str("label", label).Msg("")
	}
}
