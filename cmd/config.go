package cmd

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/mjmorales/mac-daemon-control/internal/core"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage daemon-control configuration",
	Long:  `Manage the core daemon-control configuration settings.`,
}

// configInitCmd represents the config init command
var configInitCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize configuration",
	Long:  `Initialize the daemon-control configuration with default values.`,
	Run: func(cmd *cobra.Command, args []string) {
		manager := core.NewManager()
		if err := manager.Init(); err != nil {
			log.Error().Err(err).Msg("Failed to initialize configuration")
			os.Exit(1)
		}
		log.Info().Str("path", core.ConfigPath()).Msg("Configuration initialized")
	},
}

// configShowCmd represents the config show command
var configShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Show current configuration",
	Long:  `Display the current daemon-control configuration.`,
	Run: func(cmd *cobra.Command, args []string) {
		manager := core.NewManager()
		if err := manager.Init(); err != nil {
			log.Error().Err(err).Msg("Failed to load configuration")
			os.Exit(1)
		}
		
		config := manager.GetConfig()
		if config == nil {
			log.Error().Msg("No configuration loaded")
			os.Exit(1)
		}
		
		// Marshal to YAML for display
		data, err := yaml.Marshal(config)
		if err != nil {
			log.Error().Err(err).Msg("Failed to format configuration")
			os.Exit(1)
		}
		
		fmt.Println("Current configuration:")
		fmt.Println("---------------------")
		fmt.Print(string(data))
	},
}

// configGetCmd represents the config get command
var configGetCmd = &cobra.Command{
	Use:   "get <key>",
	Short: "Get a configuration value",
	Long:  `Get a specific configuration value by key.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		key := args[0]
		
		manager := core.NewManager()
		if err := manager.Init(); err != nil {
			log.Error().Err(err).Msg("Failed to load configuration")
			os.Exit(1)
		}
		
		value, err := manager.Get(key)
		if err != nil {
			log.Error().Err(err).Str("key", key).Msg("Failed to get configuration value")
			os.Exit(1)
		}
		
		// Format output based on type
		switch v := value.(type) {
		case string:
			fmt.Println(v)
		case bool:
			fmt.Println(v)
		case int, int64, float64:
			fmt.Println(v)
		case map[string]interface{}:
			data, _ := yaml.Marshal(v)
			fmt.Print(string(data))
		default:
			fmt.Printf("%v\n", v)
		}
	},
}

// configSetCmd represents the config set command
var configSetCmd = &cobra.Command{
	Use:   "set <key> <value>",
	Short: "Set a configuration value",
	Long: `Set a configuration value by key.
	
For boolean values, use: true, false, yes, no, on, off
For map values (like custom_env_vars), use key.subkey format`,
	Args: cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		key := args[0]
		value := args[1]
		
		manager := core.NewManager()
		if err := manager.Init(); err != nil {
			log.Error().Err(err).Msg("Failed to load configuration")
			os.Exit(1)
		}
		
		// Parse value based on key type
		var parsedValue interface{}
		
		// Check if it's a boolean field
		boolFields := []string{"auto_generate_plists", "backup_on_generate", "validate_plists", "use_system_launchd"}
		isBoolField := false
		for _, field := range boolFields {
			if key == field {
				isBoolField = true
				break
			}
		}
		
		if isBoolField {
			// Parse boolean
			switch strings.ToLower(value) {
			case "true", "yes", "on", "1":
				parsedValue = true
			case "false", "no", "off", "0":
				parsedValue = false
			default:
				log.Error().Str("value", value).Msg("Invalid boolean value")
				os.Exit(1)
			}
		} else if strings.HasPrefix(key, "custom_env_vars.") {
			// Handle map values
			parts := strings.SplitN(key, ".", 2)
			if len(parts) == 2 {
				// Get current map
				currentVal, _ := manager.Get("custom_env_vars")
				envVars, ok := currentVal.(map[string]interface{})
				if !ok {
					envVars = make(map[string]interface{})
				}
				
				// Set the specific key
				envVars[parts[1]] = value
				key = "custom_env_vars"
				parsedValue = envVars
			}
		} else {
			// Use string value as-is
			parsedValue = value
		}
		
		if err := manager.Set(key, parsedValue); err != nil {
			log.Error().Err(err).Msg("Failed to set configuration value")
			os.Exit(1)
		}
		
		log.Info().Str("key", key).Interface("value", parsedValue).Msg("Configuration updated")
	},
}

// configListCmd represents the config list command
var configListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all configuration keys",
	Long:  `List all available configuration keys.`,
	Run: func(cmd *cobra.Command, args []string) {
		keys := core.ValidKeys()
		sort.Strings(keys)
		
		fmt.Println("Available configuration keys:")
		fmt.Println("----------------------------")
		for _, key := range keys {
			fmt.Printf("  %s\n", key)
		}
		
		fmt.Println("\nFor map values like custom_env_vars, use dot notation:")
		fmt.Println("  custom_env_vars.KEY_NAME")
	},
}

// configPathCmd represents the config path command
var configPathCmd = &cobra.Command{
	Use:   "path",
	Short: "Show configuration file path",
	Long:  `Display the path to the core configuration file.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(core.ConfigPath())
	},
}

func init() {
	rootCmd.AddCommand(configCmd)
	
	// Add subcommands
	configCmd.AddCommand(configInitCmd)
	configCmd.AddCommand(configShowCmd)
	configCmd.AddCommand(configGetCmd)
	configCmd.AddCommand(configSetCmd)
	configCmd.AddCommand(configListCmd)
	configCmd.AddCommand(configPathCmd)
}