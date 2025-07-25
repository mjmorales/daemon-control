package core

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConfigDir(t *testing.T) {
	tests := []struct {
		name     string
		setup    func()
		teardown func()
		want     string
	}{
		{
			name: "returns home directory path",
			setup: func() {
				// No setup needed, uses actual home directory
			},
			teardown: func() {},
			want:     filepath.Join(os.Getenv("HOME"), ".daemon-control"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setup != nil {
				tt.setup()
			}
			if tt.teardown != nil {
				defer tt.teardown()
			}

			got := ConfigDir()
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestConfigPath(t *testing.T) {
	want := filepath.Join(ConfigDir(), "core.config.yaml")
	got := ConfigPath()
	assert.Equal(t, want, got)
}

func TestDefaultConfig(t *testing.T) {
	config := DefaultConfig()
	
	assert.Equal(t, "./daemons.yaml", config.DaemonConfigPath)
	assert.Equal(t, "./daemons", config.DaemonsDir)
	assert.Equal(t, "./out", config.OutputDir)
	assert.Equal(t, "./logs", config.LogsDir)
	assert.False(t, config.AutoGeneratePlists)
	assert.True(t, config.BackupOnGenerate)
	assert.True(t, config.ValidatePlists)
	assert.Equal(t, "info", config.LogLevel)
	assert.Equal(t, "console", config.LogFormat)
	assert.False(t, config.UseSystemLaunchd)
	assert.NotNil(t, config.CustomEnvVars)
	
	home, _ := os.UserHomeDir()
	assert.Equal(t, filepath.Join(home, "Library", "LaunchAgents"), config.LaunchAgentsDir)
}

func TestNewManager(t *testing.T) {
	manager := NewManager()
	
	assert.NotNil(t, manager)
	assert.Equal(t, ConfigPath(), manager.configPath)
	assert.NotNil(t, manager.viper)
}

func TestManager_Init(t *testing.T) {
	tests := []struct {
		name      string
		setup     func(string)
		wantError bool
		validate  func(*testing.T, *Manager)
	}{
		{
			name: "creates default config when not exists",
			setup: func(dir string) {
				// No config file exists
			},
			wantError: false,
			validate: func(t *testing.T, m *Manager) {
				assert.NotNil(t, m.config)
				assert.Equal(t, DefaultConfig().LogLevel, m.config.LogLevel)
			},
		},
		{
			name: "loads existing config",
			setup: func(dir string) {
				// Create a custom config file
				configPath := filepath.Join(dir, "core.config.yaml")
				content := `log_level: debug
log_format: json
daemon_config_path: /custom/path/daemons.yaml`
				err := os.WriteFile(configPath, []byte(content), 0600)
				require.NoError(t, err)
			},
			wantError: false,
			validate: func(t *testing.T, m *Manager) {
				assert.NotNil(t, m.config)
				assert.Equal(t, "debug", m.config.LogLevel)
				assert.Equal(t, "json", m.config.LogFormat)
				assert.Equal(t, "/custom/path/daemons.yaml", m.config.DaemonConfigPath)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create temp directory
			tempDir := t.TempDir()
			
			// Create manager with custom config path
			manager := &Manager{
				configPath: filepath.Join(tempDir, "core.config.yaml"),
				viper:      viper.New(),
			}
			
			// Since ConfigDir is a function, we can't override it directly
			// Instead, we'll use the manager with a custom config path
			
			if tt.setup != nil {
				tt.setup(tempDir)
			}
			
			err := manager.Init()
			
			if tt.wantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				if tt.validate != nil {
					tt.validate(t, manager)
				}
			}
		})
	}
}

func TestManager_Save(t *testing.T) {
	tests := []struct {
		name      string
		setup     func(*Manager)
		wantError bool
	}{
		{
			name: "saves loaded config successfully",
			setup: func(m *Manager) {
				m.config = &CoreConfig{
					LogLevel:  "debug",
					LogFormat: "json",
				}
			},
			wantError: false,
		},
		{
			name: "fails when no config loaded",
			setup: func(m *Manager) {
				m.config = nil
			},
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tempDir := t.TempDir()
			configPath := filepath.Join(tempDir, "core.config.yaml")
			
			manager := &Manager{
				configPath: configPath,
				viper:      viper.New(),
			}
			manager.viper.SetConfigFile(configPath)
			manager.viper.SetConfigType("yaml")
			
			if tt.setup != nil {
				tt.setup(manager)
			}
			
			err := manager.Save()
			
			if tt.wantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				// Verify file was created
				_, err := os.Stat(configPath)
				assert.NoError(t, err)
			}
		})
	}
}

func TestManager_GetAndSet(t *testing.T) {
	tests := []struct {
		name      string
		key       string
		setValue  interface{}
		wantValue interface{}
		wantError bool
	}{
		{
			name:      "get existing key",
			key:       "log_level",
			setValue:  "debug",
			wantValue: "debug",
			wantError: false,
		},
		{
			name:      "get non-existent key",
			key:       "nonexistent",
			setValue:  nil,
			wantValue: nil,
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tempDir := t.TempDir()
			manager := &Manager{
				configPath: filepath.Join(tempDir, "core.config.yaml"),
				viper:      viper.New(),
				config:     DefaultConfig(),
			}
			manager.viper.SetConfigFile(manager.configPath)
			manager.viper.SetConfigType("yaml")
			
			// Set value if needed
			if tt.setValue != nil {
				err := manager.Set(tt.key, tt.setValue)
				assert.NoError(t, err)
			}
			
			// Get value
			got, err := manager.Get(tt.key)
			
			if tt.wantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantValue, got)
			}
		})
	}
}

func TestManager_GetDaemonConfigPath(t *testing.T) {
	tests := []struct {
		name   string
		config *CoreConfig
		want   string
	}{
		{
			name: "returns configured path",
			config: &CoreConfig{
				DaemonConfigPath: "/custom/daemons.yaml",
			},
			want: "/custom/daemons.yaml",
		},
		{
			name: "expands home directory",
			config: &CoreConfig{
				DaemonConfigPath: "~/daemons.yaml",
			},
			want: filepath.Join(os.Getenv("HOME"), "daemons.yaml"),
		},
		{
			name:   "returns default when config is nil",
			config: nil,
			want:   DefaultConfig().DaemonConfigPath,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			manager := &Manager{
				config: tt.config,
			}
			
			got := manager.GetDaemonConfigPath()
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestValidKeys(t *testing.T) {
	keys := ValidKeys()
	
	expectedKeys := []string{
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
	
	assert.ElementsMatch(t, expectedKeys, keys)
}

func TestManager_applyLogSettings(t *testing.T) {
	tests := []struct {
		name     string
		logLevel string
		logFormat string
	}{
		{
			name:     "applies debug level with json format",
			logLevel: "debug",
			logFormat: "json",
		},
		{
			name:     "applies info level with console format",
			logLevel: "info",
			logFormat: "console",
		},
		{
			name:     "applies warn level",
			logLevel: "warn",
			logFormat: "console",
		},
		{
			name:     "applies error level",
			logLevel: "error",
			logFormat: "json",
		},
		{
			name:     "defaults to info for invalid level",
			logLevel: "invalid",
			logFormat: "console",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			manager := &Manager{
				config: &CoreConfig{
					LogLevel:  tt.logLevel,
					LogFormat: tt.logFormat,
				},
			}
			
			// This method sets global state, so we just verify it doesn't panic
			assert.NotPanics(t, func() {
				manager.applyLogSettings()
			})
		})
	}
}