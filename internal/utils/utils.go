package utils

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/rs/zerolog/log"

	"github.com/mjmorales/daemon-control/internal/core"
)

// GetDaemonsDir returns the daemons directory from config
func GetDaemonsDir() string {
	config := core.GetManager().GetConfig()
	if config != nil && config.DaemonsDir != "" {
		return config.DaemonsDir
	}
	return filepath.Join(getScriptDir(), "daemons")
}

// GetLaunchAgentsDir returns the LaunchAgents directory from config
func GetLaunchAgentsDir() string {
	config := core.GetManager().GetConfig()
	if config != nil && config.LaunchAgentsDir != "" {
		return config.LaunchAgentsDir
	}
	return filepath.Join(os.Getenv("HOME"), "Library", "LaunchAgents")
}

// Legacy path variables for backward compatibility
var (
	ScriptDir       = getScriptDir()
	DaemonsDir      = GetDaemonsDir()
	LaunchAgentsDir = GetLaunchAgentsDir()
)

func init() {
	// Let core config manager handle log configuration
	// It will be initialized when first accessed
}

func getScriptDir() string {
	ex, err := os.Executable()
	if err != nil {
		return "."
	}
	return filepath.Dir(ex)
}

// GetPlistPath returns the full path to a daemon's plist file
func GetPlistPath(daemonName string) string {
	return filepath.Join(DaemonsDir, daemonName+".plist")
}

// GetPlistValue reads a value from a plist file using defaults command
func GetPlistValue(plistPath, key string) (string, error) {
	ctx := context.Background()
	cmd := exec.CommandContext(ctx, "defaults", "read", plistPath, key)
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(output)), nil
}

// GetDaemonLabel extracts the Label value from a plist file
func GetDaemonLabel(plistPath string) (string, error) {
	return GetPlistValue(plistPath, "Label")
}

// GetWorkingDirectory extracts the WorkingDirectory value from a plist file
func GetWorkingDirectory(plistPath string) (string, error) {
	return GetPlistValue(plistPath, "WorkingDirectory")
}

// GetStdoutPath extracts the StandardOutPath value from a plist file
func GetStdoutPath(plistPath string) (string, error) {
	return GetPlistValue(plistPath, "StandardOutPath")
}

// GetStderrPath extracts the StandardErrorPath value from a plist file
func GetStderrPath(plistPath string) (string, error) {
	return GetPlistValue(plistPath, "StandardErrorPath")
}

// CheckPlistExists verifies that a daemon's plist file exists
func CheckPlistExists(daemonName string) error {
	plistPath := GetPlistPath(daemonName)
	if _, err := os.Stat(plistPath); os.IsNotExist(err) {
		log.Error().Str("daemon", daemonName).Msg("Daemon not found")
		log.Info().Msg("Available daemons:")

		files, _ := filepath.Glob(filepath.Join(DaemonsDir, "*.plist"))
		for _, file := range files {
			base := filepath.Base(file)
			name := strings.TrimSuffix(base, ".plist")
			log.Info().Str("daemon", name).Msg("")
		}
		return fmt.Errorf("daemon not found")
	}
	return nil
}

// IsInstalled checks if a daemon is installed in LaunchAgents
func IsInstalled(daemonName string) (bool, error) {
	plistPath := GetPlistPath(daemonName)
	label, err := GetDaemonLabel(plistPath)
	if err != nil {
		return false, err
	}

	installedPath := filepath.Join(LaunchAgentsDir, label+".plist")
	if _, err := os.Stat(installedPath); err == nil {
		return true, nil
	}
	return false, nil
}

// IsRunning checks if a daemon is currently running
func IsRunning(daemonName string) (bool, error) {
	plistPath := GetPlistPath(daemonName)
	label, err := GetDaemonLabel(plistPath)
	if err != nil {
		return false, err
	}

	ctx := context.Background()
	cmd := exec.CommandContext(ctx, "launchctl", "list")
	output, err := cmd.Output()
	if err != nil {
		return false, err
	}

	return strings.Contains(string(output), label), nil
}

// RunLaunchctl executes a launchctl command
func RunLaunchctl(args ...string) error {
	ctx := context.Background()
	cmd := exec.CommandContext(ctx, "launchctl", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// CopyFile copies a file from src to dst
func CopyFile(src, dst string) error {
	input, err := os.ReadFile(src)
	if err != nil {
		return err
	}

	err = os.WriteFile(dst, input, 0600)
	if err != nil {
		return err
	}

	return nil
}
