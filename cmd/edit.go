package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/mjmorales/mac-daemon-control/internal/core"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var (
	editCore bool
)

// editCmd represents the edit command
var editCmd = &cobra.Command{
	Use:   "edit",
	Short: "Open configuration file in editor",
	Long: `Open the daemon configuration file in your default editor.
	
The editor is determined by the EDITOR environment variable,
or falls back to common editors based on your platform.`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := runEdit(); err != nil {
			log.Error().Err(err).Msg("Failed to open editor")
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(editCmd)
	
	editCmd.Flags().BoolVar(&editCore, "core", false, "Edit core configuration instead of daemon configuration")
}

func runEdit() error {
	// Determine which config file to edit
	var configPath string
	
	if editCore {
		// Edit core config
		configPath = core.ConfigPath()
		log.Info().Str("path", configPath).Msg("Opening core configuration")
	} else {
		// Edit daemon config
		manager := core.GetManager()
		configPath = manager.GetDaemonConfigPath()
		
		// Check if file exists, create if not
		if _, err := os.Stat(configPath); os.IsNotExist(err) {
			log.Warn().Str("path", configPath).Msg("Daemon config not found, creating example")
			
			// Create a minimal example config
			exampleConfig := `# Daemon configuration file
# Define your daemons here

daemons:
  # Example daemon
  # - name: my-daemon
  #   label: com.example.my-daemon
  #   program: /usr/local/bin/my-program
  #   working_directory: /path/to/working/dir
  #   run_at_load: false
`
			if err := os.MkdirAll(filepath.Dir(configPath), 0755); err != nil {
				return fmt.Errorf("failed to create config directory: %w", err)
			}
			
			if err := os.WriteFile(configPath, []byte(exampleConfig), 0644); err != nil {
				return fmt.Errorf("failed to create example config: %w", err)
			}
		}
		
		log.Info().Str("path", configPath).Msg("Opening daemon configuration")
	}
	
	// Find editor
	editor := findEditor()
	if editor == "" {
		return fmt.Errorf("no editor found. Please set EDITOR environment variable")
	}
	
	log.Debug().Str("editor", editor).Str("file", configPath).Msg("Launching editor")
	
	// Launch editor
	cmd := exec.Command(editor, configPath)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("editor failed: %w", err)
	}
	
	log.Info().Msg("Editor closed")
	
	// If editing daemon config, offer to generate plists
	if !editCore {
		fmt.Println("\nDaemon configuration edited.")
		fmt.Println("Run 'daemon-control generate' to generate/update plist files.")
	}
	
	return nil
}

// findEditor finds an appropriate editor
func findEditor() string {
	// Check EDITOR environment variable first
	if editor := os.Getenv("EDITOR"); editor != "" {
		return editor
	}
	
	// Check VISUAL environment variable
	if visual := os.Getenv("VISUAL"); visual != "" {
		return visual
	}
	
	// Try common editors
	editors := []string{
		"nano",     // Simple and universally available
		"vim",      // Vi improved
		"vi",       // Basic vi
		"emacs",    // Emacs
		"code",     // VS Code
		"subl",     // Sublime Text
		"atom",     // Atom
		"mate",     // TextMate
		"bbedit",   // BBEdit
		"nova",     // Nova
	}
	
	// On macOS, also try 'open' which uses the default app
	if runtime.GOOS == "darwin" {
		editors = append([]string{"open", "-t"}, editors...)
	}
	
	// Find first available editor
	for _, editor := range editors {
		if _, err := exec.LookPath(editor); err == nil {
			// Special handling for 'open -t' on macOS
			if editor == "open" {
				return editor
			}
			return editor
		}
	}
	
	return ""
}

// Additional helper command for quick access
var editDaemonCmd = &cobra.Command{
	Use:   "edit-daemon",
	Short: "Open daemon configuration in editor (alias for 'edit')",
	Long:  `Open the daemon configuration file in your default editor.`,
	Run: func(cmd *cobra.Command, args []string) {
		editCore = false
		if err := runEdit(); err != nil {
			log.Error().Err(err).Msg("Failed to open editor")
			os.Exit(1)
		}
	},
}

var editCoreCmd = &cobra.Command{
	Use:   "edit-core",  
	Short: "Open core configuration in editor (alias for 'edit --core')",
	Long:  `Open the core configuration file in your default editor.`,
	Run: func(cmd *cobra.Command, args []string) {
		editCore = true
		if err := runEdit(); err != nil {
			log.Error().Err(err).Msg("Failed to open editor")
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(editDaemonCmd)
	rootCmd.AddCommand(editCoreCmd)
}