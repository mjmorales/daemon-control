package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

// Loader handles configuration loading
type Loader struct {
	configPath string
	config     *Config
}

// NewLoader creates a new configuration loader
func NewLoader(configPath string) *Loader {
	return &Loader{
		configPath: configPath,
	}
}

// Load reads and parses the configuration file
func (l *Loader) Load() (*Config, error) {
	v := viper.New()
	
	// Set config name and paths
	if l.configPath != "" {
		v.SetConfigFile(l.configPath)
	} else {
		// Default config locations
		v.SetConfigName("daemons")
		v.SetConfigType("yaml")
		v.AddConfigPath(".")
		v.AddConfigPath("./config")
		v.AddConfigPath("$HOME/.daemon-control")
		v.AddConfigPath("/etc/daemon-control")
	}
	
	// Read config file
	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Warn().Msg("No config file found, using defaults")
			return &Config{}, nil
		}
		return nil, fmt.Errorf("error reading config: %w", err)
	}
	
	log.Info().Str("config", v.ConfigFileUsed()).Msg("Using config file")
	
	// Unmarshal config
	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("error unmarshaling config: %w", err)
	}
	
	// Validate config
	if err := l.validateConfig(&cfg); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}
	
	l.config = &cfg
	return &cfg, nil
}

// validateConfig validates the configuration
func (l *Loader) validateConfig(cfg *Config) error {
	// Check for duplicate names
	names := make(map[string]bool)
	labels := make(map[string]bool)
	
	for i, daemon := range cfg.Daemons {
		// Validate required fields
		if daemon.Name == "" {
			return fmt.Errorf("daemon[%d]: name is required", i)
		}
		
		if daemon.Label == "" {
			return fmt.Errorf("daemon[%d]: label is required", i)
		}
		
		if daemon.Program == "" && len(daemon.ProgramArguments) == 0 {
			return fmt.Errorf("daemon[%d]: program or program_arguments is required", i)
		}
		
		// Check for duplicates
		if names[daemon.Name] {
			return fmt.Errorf("duplicate daemon name: %s", daemon.Name)
		}
		names[daemon.Name] = true
		
		if labels[daemon.Label] {
			return fmt.Errorf("duplicate daemon label: %s", daemon.Label)
		}
		labels[daemon.Label] = true
		
		// Validate paths
		if daemon.WorkingDirectory != "" {
			if !filepath.IsAbs(daemon.WorkingDirectory) {
				return fmt.Errorf("daemon[%s]: working_directory must be absolute path", daemon.Name)
			}
		}
		
		// Validate process type
		if daemon.ProcessType != "" {
			validTypes := map[string]bool{
				"Background":  true,
				"Standard":    true,
				"Adaptive":    true,
				"Interactive": true,
			}
			if !validTypes[daemon.ProcessType] {
				return fmt.Errorf("daemon[%s]: invalid process_type: %s", daemon.Name, daemon.ProcessType)
			}
		}
		
		// Validate calendar intervals
		for j, interval := range daemon.StartCalendarInterval {
			if err := validateCalendarInterval(interval); err != nil {
				return fmt.Errorf("daemon[%s].start_calendar_interval[%d]: %w", daemon.Name, j, err)
			}
		}
	}
	
	return nil
}

// validateCalendarInterval validates a calendar interval
func validateCalendarInterval(interval CalendarInterval) error {
	if interval.Minute != nil && (*interval.Minute < 0 || *interval.Minute > 59) {
		return fmt.Errorf("minute must be between 0 and 59")
	}
	
	if interval.Hour != nil && (*interval.Hour < 0 || *interval.Hour > 23) {
		return fmt.Errorf("hour must be between 0 and 23")
	}
	
	if interval.Day != nil && (*interval.Day < 1 || *interval.Day > 31) {
		return fmt.Errorf("day must be between 1 and 31")
	}
	
	if interval.Weekday != nil && (*interval.Weekday < 0 || *interval.Weekday > 7) {
		return fmt.Errorf("weekday must be between 0 and 7")
	}
	
	if interval.Month != nil && (*interval.Month < 1 || *interval.Month > 12) {
		return fmt.Errorf("month must be between 1 and 12")
	}
	
	return nil
}

// GetDaemon returns a daemon by name
func (l *Loader) GetDaemon(name string) (*Daemon, error) {
	if l.config == nil {
		return nil, fmt.Errorf("config not loaded")
	}
	
	for _, daemon := range l.config.Daemons {
		if daemon.Name == name {
			return &daemon, nil
		}
	}
	
	return nil, fmt.Errorf("daemon not found: %s", name)
}

// GetAllDaemons returns all configured daemons
func (l *Loader) GetAllDaemons() []Daemon {
	if l.config == nil {
		return []Daemon{}
	}
	return l.config.Daemons
}

// ConfigExists checks if a config file exists
func ConfigExists(path string) bool {
	if path != "" {
		_, err := os.Stat(path)
		return err == nil
	}
	
	// Check default locations
	locations := []string{
		"daemons.yaml",
		"daemons.yml",
		"config/daemons.yaml",
		"config/daemons.yml",
	}
	
	for _, loc := range locations {
		if _, err := os.Stat(loc); err == nil {
			return true
		}
	}
	
	return false
}