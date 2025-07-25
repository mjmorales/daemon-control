package core

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

// CoreConfig represents the core daemon-control configuration
type CoreConfig struct {
	// Daemon configuration file path
	DaemonConfigPath string `mapstructure:"daemon_config_path" yaml:"daemon_config_path" json:"daemon_config_path"`

	// Default paths
	DaemonsDir string `mapstructure:"daemons_dir" yaml:"daemons_dir" json:"daemons_dir"`
	OutputDir  string `mapstructure:"output_dir" yaml:"output_dir" json:"output_dir"`
	LogsDir    string `mapstructure:"logs_dir" yaml:"logs_dir" json:"logs_dir"`

	// Behavior settings
	AutoGeneratePlists bool `mapstructure:"auto_generate_plists" yaml:"auto_generate_plists" json:"auto_generate_plists"`
	BackupOnGenerate   bool `mapstructure:"backup_on_generate" yaml:"backup_on_generate" json:"backup_on_generate"`
	ValidatePlists     bool `mapstructure:"validate_plists" yaml:"validate_plists" json:"validate_plists"`

	// Logging settings
	LogLevel  string `mapstructure:"log_level" yaml:"log_level" json:"log_level"`
	LogFormat string `mapstructure:"log_format" yaml:"log_format" json:"log_format"` // json or console

	// LaunchAgent settings
	LaunchAgentsDir string `mapstructure:"launch_agents_dir" yaml:"launch_agents_dir" json:"launch_agents_dir"`

	// Advanced settings
	UseSystemLaunchd bool              `mapstructure:"use_system_launchd" yaml:"use_system_launchd" json:"use_system_launchd"`
	CustomEnvVars    map[string]string `mapstructure:"custom_env_vars" yaml:"custom_env_vars" json:"custom_env_vars"`
}

// Manager handles core configuration operations
type Manager struct {
	configPath string
	config     *CoreConfig
	viper      *viper.Viper
}

// ConfigDir returns the daemon-control config directory path
func ConfigDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ".daemon-control"
	}
	return filepath.Join(home, ".daemon-control")
}

// ConfigPath returns the core config file path
func ConfigPath() string {
	return filepath.Join(ConfigDir(), "core.config.yaml")
}

// DefaultConfig returns the default core configuration
func DefaultConfig() *CoreConfig {
	home, _ := os.UserHomeDir()
	return &CoreConfig{
		DaemonConfigPath:   "./daemons.yaml",
		DaemonsDir:         "./daemons",
		OutputDir:          "./out",
		LogsDir:            "./logs",
		AutoGeneratePlists: false,
		BackupOnGenerate:   true,
		ValidatePlists:     true,
		LogLevel:           "info",
		LogFormat:          "console",
		LaunchAgentsDir:    filepath.Join(home, "Library", "LaunchAgents"),
		UseSystemLaunchd:   false,
		CustomEnvVars:      make(map[string]string),
	}
}

// NewManager creates a new core config manager
func NewManager() *Manager {
	return &Manager{
		configPath: ConfigPath(),
		viper:      viper.New(),
	}
}

// Init initializes the core configuration
func (m *Manager) Init() error {
	// Create config directory if it doesn't exist
	configDir := ConfigDir()
	if err := os.MkdirAll(configDir, 0700); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	// Check if config file exists
	if _, err := os.Stat(m.configPath); os.IsNotExist(err) {
		// Create default config
		if err := m.createDefaultConfig(); err != nil {
			return fmt.Errorf("failed to create default config: %w", err)
		}
		log.Info().Str("path", m.configPath).Msg("Created default core configuration")
	}

	return m.Load()
}

// Load loads the core configuration
func (m *Manager) Load() error {
	m.viper.SetConfigFile(m.configPath)
	m.viper.SetConfigType("yaml")

	// Set defaults
	defaults := DefaultConfig()
	m.viper.SetDefault("daemon_config_path", defaults.DaemonConfigPath)
	m.viper.SetDefault("daemons_dir", defaults.DaemonsDir)
	m.viper.SetDefault("output_dir", defaults.OutputDir)
	m.viper.SetDefault("logs_dir", defaults.LogsDir)
	m.viper.SetDefault("auto_generate_plists", defaults.AutoGeneratePlists)
	m.viper.SetDefault("backup_on_generate", defaults.BackupOnGenerate)
	m.viper.SetDefault("validate_plists", defaults.ValidatePlists)
	m.viper.SetDefault("log_level", defaults.LogLevel)
	m.viper.SetDefault("log_format", defaults.LogFormat)
	m.viper.SetDefault("launch_agents_dir", defaults.LaunchAgentsDir)
	m.viper.SetDefault("use_system_launchd", defaults.UseSystemLaunchd)

	// Read config
	if err := m.viper.ReadInConfig(); err != nil {
		return fmt.Errorf("failed to read config: %w", err)
	}

	// Unmarshal config
	var config CoreConfig
	if err := m.viper.Unmarshal(&config); err != nil {
		return fmt.Errorf("failed to unmarshal config: %w", err)
	}

	m.config = &config

	// Apply log settings
	m.applyLogSettings()

	return nil
}

// Save saves the current configuration
func (m *Manager) Save() error {
	if m.config == nil {
		return fmt.Errorf("no config loaded")
	}

	// Set all values in viper
	m.viper.Set("daemon_config_path", m.config.DaemonConfigPath)
	m.viper.Set("daemons_dir", m.config.DaemonsDir)
	m.viper.Set("output_dir", m.config.OutputDir)
	m.viper.Set("logs_dir", m.config.LogsDir)
	m.viper.Set("auto_generate_plists", m.config.AutoGeneratePlists)
	m.viper.Set("backup_on_generate", m.config.BackupOnGenerate)
	m.viper.Set("validate_plists", m.config.ValidatePlists)
	m.viper.Set("log_level", m.config.LogLevel)
	m.viper.Set("log_format", m.config.LogFormat)
	m.viper.Set("launch_agents_dir", m.config.LaunchAgentsDir)
	m.viper.Set("use_system_launchd", m.config.UseSystemLaunchd)
	m.viper.Set("custom_env_vars", m.config.CustomEnvVars)

	// Write config
	if err := m.viper.WriteConfig(); err != nil {
		return fmt.Errorf("failed to write config: %w", err)
	}

	return nil
}

// Get returns a config value by key
func (m *Manager) Get(key string) (interface{}, error) {
	if m.viper == nil {
		return nil, fmt.Errorf("config not loaded")
	}

	if !m.viper.IsSet(key) {
		return nil, fmt.Errorf("key not found: %s", key)
	}

	return m.viper.Get(key), nil
}

// Set sets a config value by key
func (m *Manager) Set(key string, value interface{}) error {
	if m.viper == nil {
		return fmt.Errorf("config not loaded")
	}

	m.viper.Set(key, value)

	// Reload config to update struct
	var config CoreConfig
	if err := m.viper.Unmarshal(&config); err != nil {
		return fmt.Errorf("failed to update config: %w", err)
	}
	m.config = &config

	// Apply log settings if changed
	if key == "log_level" || key == "log_format" {
		m.applyLogSettings()
	}

	return m.Save()
}

// GetConfig returns the loaded configuration
func (m *Manager) GetConfig() *CoreConfig {
	return m.config
}

// GetDaemonConfigPath returns the resolved daemon config path
func (m *Manager) GetDaemonConfigPath() string {
	if m.config == nil {
		return DefaultConfig().DaemonConfigPath
	}

	// Expand home directory if needed
	path := m.config.DaemonConfigPath
	if strings.HasPrefix(path, "~/") {
		home, _ := os.UserHomeDir()
		path = filepath.Join(home, path[2:])
	}

	return path
}

// createDefaultConfig creates the default config file
func (m *Manager) createDefaultConfig() error {
	config := DefaultConfig()

	// Create viper instance with defaults
	v := viper.New()
	v.SetConfigFile(m.configPath)
	v.SetConfigType("yaml")

	v.Set("daemon_config_path", config.DaemonConfigPath)
	v.Set("daemons_dir", config.DaemonsDir)
	v.Set("output_dir", config.OutputDir)
	v.Set("logs_dir", config.LogsDir)
	v.Set("auto_generate_plists", config.AutoGeneratePlists)
	v.Set("backup_on_generate", config.BackupOnGenerate)
	v.Set("validate_plists", config.ValidatePlists)
	v.Set("log_level", config.LogLevel)
	v.Set("log_format", config.LogFormat)
	v.Set("launch_agents_dir", config.LaunchAgentsDir)
	v.Set("use_system_launchd", config.UseSystemLaunchd)
	v.Set("custom_env_vars", config.CustomEnvVars)

	// Add header comment
	v.Set("_comment", "daemon-control core configuration file")

	return v.WriteConfig()
}

// applyLogSettings applies the log level and format settings
func (m *Manager) applyLogSettings() {
	if m.config == nil {
		return
	}

	// Set log level
	switch m.config.LogLevel {
	case "debug":
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	case "info":
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	case "warn":
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
	case "error":
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	default:
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}

	// Set log format
	if m.config.LogFormat == "json" {
		log.Logger = zerolog.New(os.Stderr).With().Timestamp().Logger()
	} else {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}
}

// ValidKeys returns all valid configuration keys
func ValidKeys() []string {
	return []string{
		"daemon_config_path",
		"daemons_dir",
		"output_dir",
		"logs_dir",
		"auto_generate_plists",
		"backup_on_generate",
		"validate_plists",
		"log_level",
		"log_format",
		"launch_agents_dir",
		"use_system_launchd",
		"custom_env_vars",
	}
}
