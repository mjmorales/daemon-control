package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "daemon-control",
	Short: "Manage macOS LaunchAgent daemons",
	Long: `A generic daemon control tool for managing macOS LaunchAgent daemons.
	
This tool allows you to install, uninstall, start, stop, and monitor
daemons defined as plist files in the ./daemons directory.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
}


