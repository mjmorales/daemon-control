package plist

import (
	"encoding/xml"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mjmorales/mac-daemon-control/internal/config"
)

func TestNewGenerator(t *testing.T) {
	gen := NewGenerator("/tmp/output")
	assert.NotNil(t, gen)
	assert.Equal(t, "/tmp/output", gen.outputDir)
}

func TestGenerator_Generate(t *testing.T) {
	tests := []struct {
		name      string
		daemon    config.Daemon
		wantError bool
		validate  func(*testing.T, string)
	}{
		{
			name: "basic daemon",
			daemon: config.Daemon{
				Name:    "test-daemon",
				Label:   "com.example.test",
				Program: "/usr/bin/test",
			},
			wantError: false,
			validate: func(t *testing.T, content string) {
				assert.Contains(t, content, "<key>Label</key>")
				assert.Contains(t, content, "<string>com.example.test</string>")
				assert.Contains(t, content, "<key>ProgramArguments</key>")
				assert.Contains(t, content, "<string>/usr/bin/test</string>")
			},
		},
		{
			name: "daemon with program arguments",
			daemon: config.Daemon{
				Name:  "arg-daemon",
				Label: "com.example.args",
				ProgramArguments: []string{
					"/usr/bin/python3",
					"/path/to/script.py",
					"--verbose",
				},
			},
			wantError: false,
			validate: func(t *testing.T, content string) {
				assert.Contains(t, content, "<string>/usr/bin/python3</string>")
				assert.Contains(t, content, "<string>/path/to/script.py</string>")
				assert.Contains(t, content, "<string>--verbose</string>")
			},
		},
		{
			name: "daemon with environment variables",
			daemon: config.Daemon{
				Name:    "env-daemon",
				Label:   "com.example.env",
				Program: "/usr/bin/env-test",
				EnvironmentVariables: map[string]string{
					"PATH":     "/usr/local/bin:/usr/bin",
					"NODE_ENV": "production",
				},
			},
			wantError: false,
			validate: func(t *testing.T, content string) {
				assert.Contains(t, content, "<key>EnvironmentVariables</key>")
				assert.Contains(t, content, "<key>PATH</key>")
				assert.Contains(t, content, "<string>/usr/local/bin:/usr/bin</string>")
				assert.Contains(t, content, "<key>NODE_ENV</key>")
				assert.Contains(t, content, "<string>production</string>")
			},
		},
		{
			name: "daemon with keep alive",
			daemon: config.Daemon{
				Name:    "keepalive-daemon",
				Label:   "com.example.keepalive",
				Program: "/usr/bin/keepalive",
				KeepAlive: &config.KeepAlive{
					SuccessfulExit: boolPtr(false),
					NetworkState:   boolPtr(true),
					PathState: map[string]bool{
						"/var/run/app.pid": true,
					},
				},
			},
			wantError: false,
			validate: func(t *testing.T, content string) {
				assert.Contains(t, content, "<key>KeepAlive</key>")
				assert.Contains(t, content, "<key>SuccessfulExit</key>")
				assert.Contains(t, content, "<false></false>")
				assert.Contains(t, content, "<key>NetworkState</key>")
				assert.Contains(t, content, "<true></true>")
				assert.Contains(t, content, "<key>PathState</key>")
				assert.Contains(t, content, "<key>/var/run/app.pid</key>")
			},
		},
		{
			name: "daemon with resource limits",
			daemon: config.Daemon{
				Name:    "limited-daemon",
				Label:   "com.example.limited",
				Program: "/usr/bin/limited",
				ResourceLimits: &config.ResourceLimits{
					CPU:               intPtr(80),
					FileSize:          intPtr(1048576),
					NumberOfFiles:     intPtr(256),
					NumberOfProcesses: intPtr(64),
				},
			},
			wantError: false,
			validate: func(t *testing.T, content string) {
				assert.Contains(t, content, "<key>SoftResourceLimits</key>")
				assert.Contains(t, content, "<key>HardResourceLimits</key>")
				assert.Contains(t, content, "<key>CPU</key>")
				assert.Contains(t, content, "<integer>80</integer>")
				assert.Contains(t, content, "<key>FileSize</key>")
				assert.Contains(t, content, "<integer>1048576</integer>")
			},
		},
		{
			name: "daemon with single calendar interval",
			daemon: config.Daemon{
				Name:    "scheduled-daemon",
				Label:   "com.example.scheduled",
				Program: "/usr/bin/scheduled",
				StartCalendarInterval: []config.CalendarInterval{
					{
						Hour:   intPtr(9),
						Minute: intPtr(30),
					},
				},
			},
			wantError: false,
			validate: func(t *testing.T, content string) {
				assert.Contains(t, content, "<key>StartCalendarInterval</key>")
				assert.Contains(t, content, "<key>Hour</key>")
				assert.Contains(t, content, "<integer>9</integer>")
				assert.Contains(t, content, "<key>Minute</key>")
				assert.Contains(t, content, "<integer>30</integer>")
				// Should be a dict, not an array for single interval
				// Count dict tags after StartCalendarInterval - should be exactly one
				idx := strings.Index(content, "<key>StartCalendarInterval</key>")
				if idx > 0 {
					rest := content[idx:]
					assert.NotContains(t, rest[:100], "<array>")
				}
			},
		},
		{
			name: "daemon with multiple calendar intervals",
			daemon: config.Daemon{
				Name:    "multi-scheduled",
				Label:   "com.example.multi-scheduled",
				Program: "/usr/bin/scheduled",
				StartCalendarInterval: []config.CalendarInterval{
					{Hour: intPtr(9), Minute: intPtr(0)},
					{Hour: intPtr(17), Minute: intPtr(30)},
				},
			},
			wantError: false,
			validate: func(t *testing.T, content string) {
				assert.Contains(t, content, "<key>StartCalendarInterval</key>")
				assert.Contains(t, content, "<array>")
				// Count occurrences of Hour key - should be 2
				hourCount := strings.Count(content, "<key>Hour</key>")
				assert.Equal(t, 2, hourCount)
			},
		},
		{
			name: "daemon with sockets",
			daemon: config.Daemon{
				Name:    "socket-daemon",
				Label:   "com.example.socket",
				Program: "/usr/bin/socket-server",
				Sockets: map[string]config.Socket{
					"http": {
						SockType:        "stream",
						SockNodeName:    "localhost",
						SockServiceName: "8080",
						SockFamily:      "IPv4",
						SockProtocol:    "TCP",
					},
				},
			},
			wantError: false,
			validate: func(t *testing.T, content string) {
				assert.Contains(t, content, "<key>Sockets</key>")
				assert.Contains(t, content, "<key>http</key>")
				assert.Contains(t, content, "<key>SockType</key>")
				assert.Contains(t, content, "<string>stream</string>")
				assert.Contains(t, content, "<key>SockServiceName</key>")
				assert.Contains(t, content, "<string>8080</string>")
			},
		},
		{
			name: "daemon with all settings",
			daemon: config.Daemon{
				Name:                "full-daemon",
				Label:               "com.example.full",
				Description:         "A fully configured daemon",
				Program:             "/usr/bin/full",
				WorkingDirectory:    "/var/lib/daemon",
				StandardOutPath:     "/var/log/daemon.out",
				StandardErrorPath:   "/var/log/daemon.err",
				RunAtLoad:           true,
				StartInterval:       300,
				ThrottleInterval:    10,
				ProcessType:         "Background",
				Nice:                intPtr(10),
				InitGroups:          true,
				UserName:            "daemon",
				GroupName:           "daemon",
				RootDirectory:       "/chroot",
				WatchPaths:          []string{"/etc/config"},
				QueuePaths:          []string{"/var/spool"},
				EnableGlobbing:      true,
				EnableTransactions:  true,
				EnablePressuredExit: true,
				ExitTimeOut:         30,
			},
			wantError: false,
			validate: func(t *testing.T, content string) {
				// Check various settings
				assert.Contains(t, content, "<key>WorkingDirectory</key>")
				assert.Contains(t, content, "<string>/var/lib/daemon</string>")
				assert.Contains(t, content, "<key>RunAtLoad</key>")
				assert.Contains(t, content, "<true></true>")
				assert.Contains(t, content, "<key>StartInterval</key>")
				assert.Contains(t, content, "<integer>300</integer>")
				assert.Contains(t, content, "<key>ProcessType</key>")
				assert.Contains(t, content, "<string>Background</string>")
				assert.Contains(t, content, "<key>Nice</key>")
				assert.Contains(t, content, "<integer>10</integer>")
				assert.Contains(t, content, "<key>UserName</key>")
				assert.Contains(t, content, "<string>daemon</string>")
				assert.Contains(t, content, "<key>WatchPaths</key>")
				assert.Contains(t, content, "<key>QueueDirectories</key>")
				assert.Contains(t, content, "<key>EnableGlobbing</key>")
				assert.Contains(t, content, "<key>ExitTimeOut</key>")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create temp directory
			tempDir := t.TempDir()
			gen := NewGenerator(tempDir)

			// Generate plist
			err := gen.Generate(&tt.daemon)

			if tt.wantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)

				// Read generated file
				plistPath := filepath.Join(tempDir, tt.daemon.Name+".plist")
				content, err := os.ReadFile(plistPath)
				require.NoError(t, err)

				// Validate content
				contentStr := string(content)
				assert.Contains(t, contentStr, "<?xml version")
				assert.Contains(t, contentStr, "<!DOCTYPE plist")
				assert.Contains(t, contentStr, "<plist version=\"1.0\">")

				if tt.validate != nil {
					tt.validate(t, contentStr)
				}

				// Validate it's valid XML
				var plist Plist
				err = xml.Unmarshal(content, &plist)
				assert.NoError(t, err)
			}
		})
	}
}

func TestGenerator_GenerateAll(t *testing.T) {
	daemons := []config.Daemon{
		{
			Name:    "daemon1",
			Label:   "com.example.daemon1",
			Program: "/usr/bin/daemon1",
		},
		{
			Name:    "daemon2",
			Label:   "com.example.daemon2",
			Program: "/usr/bin/daemon2",
		},
		{
			Name:    "daemon3",
			Label:   "com.example.daemon3",
			Program: "/usr/bin/daemon3",
		},
	}

	// Create temp directory
	tempDir := t.TempDir()
	gen := NewGenerator(tempDir)

	// Generate all plists
	err := gen.GenerateAll(daemons)
	assert.NoError(t, err)

	// Check all files were created
	for _, daemon := range daemons {
		plistPath := filepath.Join(tempDir, daemon.Name+".plist")
		_, err := os.Stat(plistPath)
		assert.NoError(t, err, "plist file should exist for %s", daemon.Name)
	}
}

func TestGenerator_DirectoryCreation(t *testing.T) {
	// Use a nested path that doesn't exist
	tempDir := t.TempDir()
	outputDir := filepath.Join(tempDir, "nested", "output", "dir")

	gen := NewGenerator(outputDir)
	daemon := config.Daemon{
		Name:    "test",
		Label:   "com.example.test",
		Program: "/usr/bin/test",
	}

	// GenerateAll creates the directory, not Generate
	err := gen.GenerateAll([]config.Daemon{daemon})
	assert.NoError(t, err)

	// Check directory exists
	_, err = os.Stat(outputDir)
	assert.NoError(t, err)
}

func TestGenerator_FilePermissions(t *testing.T) {
	tempDir := t.TempDir()
	gen := NewGenerator(tempDir)

	daemon := config.Daemon{
		Name:    "test",
		Label:   "com.example.test",
		Program: "/usr/bin/test",
	}

	err := gen.Generate(&daemon)
	require.NoError(t, err)

	// Check file permissions
	plistPath := filepath.Join(tempDir, daemon.Name+".plist")
	info, err := os.Stat(plistPath)
	require.NoError(t, err)

	// Should be 0600
	assert.Equal(t, os.FileMode(0600), info.Mode().Perm())
}

func TestCalendarIntervalToDict(t *testing.T) {
	gen := NewGenerator("")

	tests := []struct {
		name     string
		interval config.CalendarInterval
		validate func(*testing.T, *Dict)
	}{
		{
			name: "all fields set",
			interval: config.CalendarInterval{
				Minute:  intPtr(30),
				Hour:    intPtr(14),
				Day:     intPtr(15),
				Weekday: intPtr(3),
				Month:   intPtr(6),
			},
			validate: func(t *testing.T, dict *Dict) {
				assert.Len(t, dict.Items, 10) // 5 key-value pairs
			},
		},
		{
			name:     "empty interval",
			interval: config.CalendarInterval{},
			validate: func(t *testing.T, dict *Dict) {
				assert.Len(t, dict.Items, 0)
			},
		},
		{
			name: "partial fields",
			interval: config.CalendarInterval{
				Hour:   intPtr(9),
				Minute: intPtr(0),
			},
			validate: func(t *testing.T, dict *Dict) {
				assert.Len(t, dict.Items, 4) // 2 key-value pairs
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dict := gen.calendarIntervalToDict(tt.interval)
			if tt.validate != nil {
				tt.validate(t, dict)
			}
		})
	}
}

// Helper functions
func intPtr(i int) *int {
	return &i
}

func boolPtr(b bool) *bool {
	return &b
}
