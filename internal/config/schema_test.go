package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
)

func TestDaemonYAMLMarshalUnmarshal(t *testing.T) {
	tests := []struct {
		name   string
		daemon Daemon
	}{
		{
			name: "basic daemon",
			daemon: Daemon{
				Name:    "test-daemon",
				Label:   "com.example.test",
				Program: "/usr/bin/test",
			},
		},
		{
			name: "daemon with all fields",
			daemon: Daemon{
				Name:             "full-daemon",
				Label:            "com.example.full",
				Description:      "A fully configured daemon",
				Program:          "/usr/bin/full",
				ProgramArguments: []string{"--arg1", "value1", "--arg2"},
				WorkingDirectory: "/var/lib/daemon",
				EnvironmentVariables: map[string]string{
					"PATH": "/usr/local/bin:/usr/bin",
					"ENV":  "production",
				},
				StandardOutPath:   "/var/log/daemon.log",
				StandardErrorPath: "/var/log/daemon.err",
				RunAtLoad:         true,
				StartInterval:     300,
				KeepAlive: &KeepAlive{
					SuccessfulExit: boolPtr(false),
					NetworkState:   boolPtr(true),
					PathState: map[string]bool{
						"/var/run/daemon.pid": true,
					},
				},
				ThrottleInterval: 10,
				ResourceLimits: &ResourceLimits{
					CPU:               intPtr(80),
					FileSize:          intPtr(1048576),
					NumberOfFiles:     intPtr(256),
					Core:              intPtr(0),
					Data:              intPtr(268435456),
					MemoryLock:        intPtr(0),
					NumberOfProcesses: intPtr(64),
					ResidentSetSize:   intPtr(536870912),
					Stack:             intPtr(8388608),
				},
				ProcessType:   "Background",
				Nice:          intPtr(10),
				InitGroups:    true,
				UserName:      "daemon",
				GroupName:     "daemon",
				RootDirectory: "/",
				Sockets: map[string]Socket{
					"main": {
						SockType:        "stream",
						SockPassive:     boolPtr(true),
						SockNodeName:    "localhost",
						SockServiceName: "8080",
						SockFamily:      "IPv4",
						SockProtocol:    "TCP",
					},
				},
				StartCalendarInterval: []CalendarInterval{
					{
						Hour:   intPtr(2),
						Minute: intPtr(30),
					},
				},
				WatchPaths:          []string{"/etc/config.conf"},
				QueuePaths:          []string{"/var/spool/daemon"},
				EnableGlobbing:      true,
				EnableTransactions:  true,
				EnablePressuredExit: false,
				ExitTimeOut:         30,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Marshal to YAML
			data, err := yaml.Marshal(tt.daemon)
			assert.NoError(t, err)
			assert.NotEmpty(t, data)

			// Unmarshal back
			var unmarshaled Daemon
			err = yaml.Unmarshal(data, &unmarshaled)
			assert.NoError(t, err)

			// Compare
			assert.Equal(t, tt.daemon, unmarshaled)
		})
	}
}

func TestConfigYAMLMarshalUnmarshal(t *testing.T) {
	config := Config{
		Daemons: []Daemon{
			{
				Name:    "daemon1",
				Label:   "com.example.daemon1",
				Program: "/usr/bin/daemon1",
			},
			{
				Name:  "daemon2",
				Label: "com.example.daemon2",
				ProgramArguments: []string{
					"/usr/bin/python3",
					"/path/to/script.py",
				},
			},
		},
	}

	// Marshal to YAML
	data, err := yaml.Marshal(config)
	assert.NoError(t, err)
	assert.NotEmpty(t, data)

	// Unmarshal back
	var unmarshaled Config
	err = yaml.Unmarshal(data, &unmarshaled)
	assert.NoError(t, err)

	// Compare
	assert.Equal(t, config, unmarshaled)
	assert.Len(t, unmarshaled.Daemons, 2)
}

func TestKeepAliveYAMLMarshalUnmarshal(t *testing.T) {
	keepAlive := KeepAlive{
		SuccessfulExit: boolPtr(false),
		NetworkState:   boolPtr(true),
		PathState: map[string]bool{
			"/var/run/app.pid": true,
			"/etc/config":      false,
		},
		OtherJobEnabled: map[string]bool{
			"com.example.other": true,
		},
		Crashed:            boolPtr(true),
		AfterInitialDemand: boolPtr(false),
	}

	// Marshal to YAML
	data, err := yaml.Marshal(keepAlive)
	assert.NoError(t, err)
	assert.NotEmpty(t, data)

	// Unmarshal back
	var unmarshaled KeepAlive
	err = yaml.Unmarshal(data, &unmarshaled)
	assert.NoError(t, err)

	// Compare
	assert.Equal(t, keepAlive, unmarshaled)
}

func TestResourceLimitsYAMLMarshalUnmarshal(t *testing.T) {
	limits := ResourceLimits{
		CPU:               intPtr(75),
		FileSize:          intPtr(524288),
		NumberOfFiles:     intPtr(1024),
		Core:              intPtr(0),
		Data:              intPtr(134217728),
		MemoryLock:        intPtr(65536),
		NumberOfProcesses: intPtr(32),
		ResidentSetSize:   intPtr(268435456),
		Stack:             intPtr(4194304),
	}

	// Marshal to YAML
	data, err := yaml.Marshal(limits)
	assert.NoError(t, err)
	assert.NotEmpty(t, data)

	// Unmarshal back
	var unmarshaled ResourceLimits
	err = yaml.Unmarshal(data, &unmarshaled)
	assert.NoError(t, err)

	// Compare
	assert.Equal(t, limits, unmarshaled)
}

func TestSocketYAMLMarshalUnmarshal(t *testing.T) {
	socket := Socket{
		SockType:        "stream",
		SockPassive:     boolPtr(true),
		SockNodeName:    "0.0.0.0",
		SockServiceName: "8080",
		SockFamily:      "IPv4",
		SockProtocol:    "TCP",
		SockPathName:    "/var/run/daemon.sock",
		SockPathMode:    intPtr(0666),
		Bonjour:         boolPtr(true),
		BonjourMultiple: []string{"_http._tcp", "_custom._tcp"},
	}

	// Marshal to YAML
	data, err := yaml.Marshal(socket)
	assert.NoError(t, err)
	assert.NotEmpty(t, data)

	// Unmarshal back
	var unmarshaled Socket
	err = yaml.Unmarshal(data, &unmarshaled)
	assert.NoError(t, err)

	// Compare
	assert.Equal(t, socket, unmarshaled)
}

func TestCalendarIntervalYAMLMarshalUnmarshal(t *testing.T) {
	tests := []struct {
		name     string
		interval CalendarInterval
	}{
		{
			name: "daily at specific time",
			interval: CalendarInterval{
				Hour:   intPtr(14),
				Minute: intPtr(30),
			},
		},
		{
			name: "weekly on Monday",
			interval: CalendarInterval{
				Weekday: intPtr(1),
				Hour:    intPtr(9),
				Minute:  intPtr(0),
			},
		},
		{
			name: "monthly on 15th",
			interval: CalendarInterval{
				Day:    intPtr(15),
				Hour:   intPtr(3),
				Minute: intPtr(0),
			},
		},
		{
			name: "yearly on Jan 1st",
			interval: CalendarInterval{
				Month:  intPtr(1),
				Day:    intPtr(1),
				Hour:   intPtr(0),
				Minute: intPtr(0),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Marshal to YAML
			data, err := yaml.Marshal(tt.interval)
			assert.NoError(t, err)
			assert.NotEmpty(t, data)

			// Unmarshal back
			var unmarshaled CalendarInterval
			err = yaml.Unmarshal(data, &unmarshaled)
			assert.NoError(t, err)

			// Compare
			assert.Equal(t, tt.interval, unmarshaled)
		})
	}
}

// Helper function to create bool pointers
func boolPtr(b bool) *bool {
	return &b
}
