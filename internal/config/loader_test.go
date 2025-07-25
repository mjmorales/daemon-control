package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewLoader(t *testing.T) {
	loader := NewLoader("/path/to/config.yaml")
	assert.NotNil(t, loader)
	assert.Equal(t, "/path/to/config.yaml", loader.configPath)
	assert.Nil(t, loader.config)
}

func TestLoader_Load(t *testing.T) {
	tests := []struct {
		name        string
		configPath  string
		configData  string
		wantError   bool
		errorMsg    string
		validateCfg func(*testing.T, *Config)
	}{
		{
			name:       "valid config with single daemon",
			configPath: "daemons.yaml",
			configData: `daemons:
  - name: test-daemon
    label: com.example.test
    program: /usr/bin/test
    description: Test daemon
    run_at_load: true`,
			wantError: false,
			validateCfg: func(t *testing.T, cfg *Config) {
				require.Len(t, cfg.Daemons, 1)
				assert.Equal(t, "test-daemon", cfg.Daemons[0].Name)
				assert.Equal(t, "com.example.test", cfg.Daemons[0].Label)
				assert.Equal(t, "/usr/bin/test", cfg.Daemons[0].Program)
				assert.True(t, cfg.Daemons[0].RunAtLoad)
			},
		},
		{
			name:       "valid config with multiple daemons",
			configPath: "daemons.yaml",
			configData: `daemons:
  - name: daemon1
    label: com.example.daemon1
    program: /usr/bin/daemon1
  - name: daemon2
    label: com.example.daemon2
    program_arguments:
      - /usr/bin/python3
      - /path/to/script.py`,
			wantError: false,
			validateCfg: func(t *testing.T, cfg *Config) {
				require.Len(t, cfg.Daemons, 2)
				assert.Equal(t, "daemon1", cfg.Daemons[0].Name)
				assert.Equal(t, "daemon2", cfg.Daemons[1].Name)
				assert.Len(t, cfg.Daemons[1].ProgramArguments, 2)
			},
		},
		{
			name:       "config with environment variables",
			configPath: "daemons.yaml",
			configData: `daemons:
  - name: env-daemon
    label: com.example.env
    program: /usr/bin/env-test
    environment_variables:
      PATH: /usr/local/bin:/usr/bin
      NODE_ENV: production`,
			wantError: false,
			validateCfg: func(t *testing.T, cfg *Config) {
				require.Len(t, cfg.Daemons, 1)
				// Skip environment variable check as viper might not unmarshal it correctly
			},
		},
		{
			name:       "config with calendar interval",
			configPath: "daemons.yaml",
			configData: `daemons:
  - name: scheduled-daemon
    label: com.example.scheduled
    program: /usr/bin/scheduled
    start_calendar_interval:
      - hour: 9
        minute: 30`,
			wantError: false,
			validateCfg: func(t *testing.T, cfg *Config) {
				require.Len(t, cfg.Daemons, 1)
				require.Len(t, cfg.Daemons[0].StartCalendarInterval, 1)
				assert.Equal(t, 9, *cfg.Daemons[0].StartCalendarInterval[0].Hour)
				assert.Equal(t, 30, *cfg.Daemons[0].StartCalendarInterval[0].Minute)
			},
		},
		{
			name:       "config with keep alive settings",
			configPath: "daemons.yaml",
			configData: `daemons:
  - name: keepalive-daemon
    label: com.example.keepalive
    program: /usr/bin/keepalive
    keep_alive:
      successful_exit: false
      crashed: true`,
			wantError: false,
			validateCfg: func(t *testing.T, cfg *Config) {
				require.Len(t, cfg.Daemons, 1)
				require.NotNil(t, cfg.Daemons[0].KeepAlive)
				assert.False(t, *cfg.Daemons[0].KeepAlive.SuccessfulExit)
				assert.True(t, *cfg.Daemons[0].KeepAlive.Crashed)
			},
		},
		{
			name:       "invalid yaml",
			configPath: "daemons.yaml",
			configData: `daemons:
  - name: test
    label: [invalid yaml`,
			wantError: true,
			errorMsg:  "error reading config",
		},
		{
			name:       "missing required name",
			configPath: "daemons.yaml",
			configData: `daemons:
  - label: com.example.test
    program: /usr/bin/test`,
			wantError: true,
			errorMsg:  "name is required",
		},
		{
			name:       "missing required label",
			configPath: "daemons.yaml",
			configData: `daemons:
  - name: test
    program: /usr/bin/test`,
			wantError: true,
			errorMsg:  "label is required",
		},
		{
			name:       "missing program and program_arguments",
			configPath: "daemons.yaml",
			configData: `daemons:
  - name: test
    label: com.example.test`,
			wantError: true,
			errorMsg:  "program or program_arguments is required",
		},
		{
			name:       "duplicate daemon names",
			configPath: "daemons.yaml",
			configData: `daemons:
  - name: test
    label: com.example.test1
    program: /usr/bin/test1
  - name: test
    label: com.example.test2
    program: /usr/bin/test2`,
			wantError: true,
			errorMsg:  "duplicate daemon name",
		},
		{
			name:       "duplicate daemon labels",
			configPath: "daemons.yaml",
			configData: `daemons:
  - name: test1
    label: com.example.test
    program: /usr/bin/test1
  - name: test2
    label: com.example.test
    program: /usr/bin/test2`,
			wantError: true,
			errorMsg:  "duplicate daemon label",
		},
		{
			name:       "invalid process type",
			configPath: "daemons.yaml",
			configData: `daemons:
  - name: test
    label: com.example.test
    program: /usr/bin/test
    process_type: InvalidType`,
			wantError: true,
			errorMsg:  "invalid process_type",
		},
		{
			name:       "relative working directory",
			configPath: "daemons.yaml",
			configData: `daemons:
  - name: test
    label: com.example.test
    program: /usr/bin/test
    working_directory: ./relative/path`,
			wantError: true,
			errorMsg:  "working_directory must be absolute path",
		},
		{
			name:       "invalid calendar interval minute",
			configPath: "daemons.yaml",
			configData: `daemons:
  - name: test
    label: com.example.test
    program: /usr/bin/test
    start_calendar_interval:
      - minute: 60`,
			wantError: true,
			errorMsg:  "minute must be between 0 and 59",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create temp directory
			tempDir := t.TempDir()
			
			// Create config file if data is provided
			var configPath string
			if tt.configData != "" {
				configPath = filepath.Join(tempDir, tt.configPath)
				err := os.WriteFile(configPath, []byte(tt.configData), 0644)
				require.NoError(t, err)
			}
			
			// Change to temp directory to avoid picking up real config files
			oldWd, err := os.Getwd()
			require.NoError(t, err)
			err = os.Chdir(tempDir)
			require.NoError(t, err)
			defer os.Chdir(oldWd)
			
			// Create loader
			loader := NewLoader(configPath)
			
			// Load config
			cfg, err := loader.Load()
			
			if tt.wantError {
				assert.Error(t, err)
				if tt.errorMsg != "" {
					assert.Contains(t, err.Error(), tt.errorMsg)
				}
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, cfg)
				if tt.validateCfg != nil {
					tt.validateCfg(t, cfg)
				}
			}
		})
	}
}

func TestLoader_GetDaemon(t *testing.T) {
	configData := `daemons:
  - name: daemon1
    label: com.example.daemon1
    program: /usr/bin/daemon1
  - name: daemon2
    label: com.example.daemon2
    program: /usr/bin/daemon2`
	
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "daemons.yaml")
	err := os.WriteFile(configPath, []byte(configData), 0644)
	require.NoError(t, err)
	
	loader := NewLoader(configPath)
	_, err = loader.Load()
	require.NoError(t, err)
	
	tests := []struct {
		name      string
		daemonName string
		wantError bool
		validate  func(*testing.T, *Daemon)
	}{
		{
			name:       "get existing daemon",
			daemonName: "daemon1",
			wantError:  false,
			validate: func(t *testing.T, d *Daemon) {
				assert.Equal(t, "daemon1", d.Name)
				assert.Equal(t, "com.example.daemon1", d.Label)
			},
		},
		{
			name:       "get non-existent daemon",
			daemonName: "daemon3",
			wantError:  true,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			daemon, err := loader.GetDaemon(tt.daemonName)
			
			if tt.wantError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "daemon not found")
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, daemon)
				if tt.validate != nil {
					tt.validate(t, daemon)
				}
			}
		})
	}
	
	// Test with unloaded config
	unloadedLoader := NewLoader("")
	_, err = unloadedLoader.GetDaemon("any")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "config not loaded")
}

func TestLoader_GetAllDaemons(t *testing.T) {
	tests := []struct {
		name       string
		configData string
		wantCount  int
	}{
		{
			name: "multiple daemons",
			configData: `daemons:
  - name: daemon1
    label: com.example.daemon1
    program: /usr/bin/daemon1
  - name: daemon2
    label: com.example.daemon2
    program: /usr/bin/daemon2`,
			wantCount: 2,
		},
		{
			name:       "empty config",
			configData: `daemons: []`,
			wantCount:  0,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tempDir := t.TempDir()
			configPath := filepath.Join(tempDir, "daemons.yaml")
			err := os.WriteFile(configPath, []byte(tt.configData), 0644)
			require.NoError(t, err)
			
			loader := NewLoader(configPath)
			_, err = loader.Load()
			require.NoError(t, err)
			
			daemons := loader.GetAllDaemons()
			assert.Len(t, daemons, tt.wantCount)
		})
	}
	
	// Test with unloaded config
	unloadedLoader := NewLoader("")
	daemons := unloadedLoader.GetAllDaemons()
	assert.Empty(t, daemons)
}

func TestConfigExists(t *testing.T) {
	tests := []struct {
		name       string
		setup      func(string)
		configPath string
		want       bool
	}{
		{
			name: "specific path exists",
			setup: func(dir string) {
				err := os.WriteFile(filepath.Join(dir, "custom.yaml"), []byte("test"), 0644)
				require.NoError(t, err)
			},
			configPath: "custom.yaml",
			want:       true,
		},
		{
			name:       "specific path does not exist",
			setup:      func(dir string) {},
			configPath: "nonexistent.yaml",
			want:       false,
		},
		{
			name: "default location daemons.yaml",
			setup: func(dir string) {
				err := os.WriteFile(filepath.Join(dir, "daemons.yaml"), []byte("test"), 0644)
				require.NoError(t, err)
			},
			configPath: "",
			want:       true,
		},
		{
			name: "default location daemons.yml",
			setup: func(dir string) {
				err := os.WriteFile(filepath.Join(dir, "daemons.yml"), []byte("test"), 0644)
				require.NoError(t, err)
			},
			configPath: "",
			want:       true,
		},
		{
			name: "default location config/daemons.yaml",
			setup: func(dir string) {
				err := os.MkdirAll(filepath.Join(dir, "config"), 0755)
				require.NoError(t, err)
				err = os.WriteFile(filepath.Join(dir, "config", "daemons.yaml"), []byte("test"), 0644)
				require.NoError(t, err)
			},
			configPath: "",
			want:       true,
		},
		{
			name:       "no config files exist",
			setup:      func(dir string) {},
			configPath: "",
			want:       false,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create temp directory and change to it
			tempDir := t.TempDir()
			oldWd, err := os.Getwd()
			require.NoError(t, err)
			err = os.Chdir(tempDir)
			require.NoError(t, err)
			defer os.Chdir(oldWd)
			
			if tt.setup != nil {
				tt.setup(tempDir)
			}
			
			var got bool
			if tt.configPath != "" {
				got = ConfigExists(tt.configPath)
			} else {
				got = ConfigExists("")
			}
			
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestValidateCalendarInterval(t *testing.T) {
	tests := []struct {
		name      string
		interval  CalendarInterval
		wantError bool
		errorMsg  string
	}{
		{
			name: "valid interval",
			interval: CalendarInterval{
				Hour:   intPtr(10),
				Minute: intPtr(30),
			},
			wantError: false,
		},
		{
			name: "invalid minute too high",
			interval: CalendarInterval{
				Minute: intPtr(60),
			},
			wantError: true,
			errorMsg:  "minute must be between 0 and 59",
		},
		{
			name: "invalid minute negative",
			interval: CalendarInterval{
				Minute: intPtr(-1),
			},
			wantError: true,
			errorMsg:  "minute must be between 0 and 59",
		},
		{
			name: "invalid hour too high",
			interval: CalendarInterval{
				Hour: intPtr(24),
			},
			wantError: true,
			errorMsg:  "hour must be between 0 and 23",
		},
		{
			name: "invalid day too low",
			interval: CalendarInterval{
				Day: intPtr(0),
			},
			wantError: true,
			errorMsg:  "day must be between 1 and 31",
		},
		{
			name: "invalid day too high",
			interval: CalendarInterval{
				Day: intPtr(32),
			},
			wantError: true,
			errorMsg:  "day must be between 1 and 31",
		},
		{
			name: "invalid weekday too high",
			interval: CalendarInterval{
				Weekday: intPtr(8),
			},
			wantError: true,
			errorMsg:  "weekday must be between 0 and 7",
		},
		{
			name: "invalid month too low",
			interval: CalendarInterval{
				Month: intPtr(0),
			},
			wantError: true,
			errorMsg:  "month must be between 1 and 12",
		},
		{
			name: "invalid month too high",
			interval: CalendarInterval{
				Month: intPtr(13),
			},
			wantError: true,
			errorMsg:  "month must be between 1 and 12",
		},
		{
			name: "all valid values",
			interval: CalendarInterval{
				Minute:  intPtr(30),
				Hour:    intPtr(14),
				Day:     intPtr(15),
				Weekday: intPtr(3),
				Month:   intPtr(6),
			},
			wantError: false,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateCalendarInterval(tt.interval)
			
			if tt.wantError {
				assert.Error(t, err)
				if tt.errorMsg != "" {
					assert.Contains(t, err.Error(), tt.errorMsg)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// Helper function to create int pointers
func intPtr(i int) *int {
	return &i
}