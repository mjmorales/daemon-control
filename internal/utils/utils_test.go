package utils

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetDaemonsDir(t *testing.T) {
	// Test default behavior
	dir := GetDaemonsDir()
	assert.NotEmpty(t, dir)
}

func TestGetLaunchAgentsDir(t *testing.T) {
	// Test default behavior
	dir := GetLaunchAgentsDir()
	assert.NotEmpty(t, dir)
	assert.Contains(t, dir, "LaunchAgents")
}

func TestGetScriptDir(t *testing.T) {
	dir := getScriptDir()
	assert.NotEmpty(t, dir)
	// Should return directory of executable or "." if error
}

func TestGetPlistPath(t *testing.T) {
	tests := []struct {
		name       string
		daemonName string
		want       string
	}{
		{
			name:       "simple daemon name",
			daemonName: "test-daemon",
			want:       "test-daemon.plist",
		},
		{
			name:       "daemon with dots",
			daemonName: "com.example.daemon",
			want:       "com.example.daemon.plist",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := GetPlistPath(tt.daemonName)
			assert.True(t, strings.HasSuffix(path, tt.want))
		})
	}
}

func TestGetPlistValue(t *testing.T) {
	// Skip if 'defaults' command is not available
	if _, err := exec.LookPath("defaults"); err != nil {
		t.Skip("defaults command not available")
	}

	// Create a test plist file
	tempDir := t.TempDir()
	plistPath := filepath.Join(tempDir, "test.plist")
	
	plistContent := `<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
	<key>Label</key>
	<string>com.example.test</string>
	<key>WorkingDirectory</key>
	<string>/usr/local</string>
</dict>
</plist>`
	
	err := os.WriteFile(plistPath, []byte(plistContent), 0644)
	require.NoError(t, err)

	tests := []struct {
		name      string
		key       string
		want      string
		wantError bool
	}{
		{
			name:      "existing key",
			key:       "Label",
			want:      "com.example.test",
			wantError: false,
		},
		{
			name:      "non-existent key",
			key:       "NonExistent",
			want:      "",
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Remove .plist extension as defaults command doesn't want it
			plistPathWithoutExt := strings.TrimSuffix(plistPath, ".plist")
			value, err := GetPlistValue(plistPathWithoutExt, tt.key)
			
			if tt.wantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, value)
			}
		})
	}
}

func TestGetDaemonLabel(t *testing.T) {
	// Skip if 'defaults' command is not available
	if _, err := exec.LookPath("defaults"); err != nil {
		t.Skip("defaults command not available")
	}

	// Create a test plist file
	tempDir := t.TempDir()
	plistPath := filepath.Join(tempDir, "test.plist")
	
	plistContent := `<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
	<key>Label</key>
	<string>com.example.test</string>
</dict>
</plist>`
	
	err := os.WriteFile(plistPath, []byte(plistContent), 0644)
	require.NoError(t, err)

	plistPathWithoutExt := strings.TrimSuffix(plistPath, ".plist")
	label, err := GetDaemonLabel(plistPathWithoutExt)
	assert.NoError(t, err)
	assert.Equal(t, "com.example.test", label)
}

func TestCheckPlistExists(t *testing.T) {
	// Create temp directory for daemons
	tempDir := t.TempDir()
	
	// Override DaemonsDir for testing
	oldDaemonsDir := DaemonsDir
	DaemonsDir = tempDir
	defer func() { DaemonsDir = oldDaemonsDir }()

	// Create a test plist
	plistPath := filepath.Join(tempDir, "existing-daemon.plist")
	err := os.WriteFile(plistPath, []byte("test"), 0644)
	require.NoError(t, err)

	tests := []struct {
		name       string
		daemonName string
		wantError  bool
	}{
		{
			name:       "existing daemon",
			daemonName: "existing-daemon",
			wantError:  false,
		},
		{
			name:       "non-existent daemon",
			daemonName: "non-existent",
			wantError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := CheckPlistExists(tt.daemonName)
			
			if tt.wantError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "daemon not found")
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestIsInstalled(t *testing.T) {
	// Skip if 'defaults' command is not available
	if _, err := exec.LookPath("defaults"); err != nil {
		t.Skip("defaults command not available")
	}

	// Create temp directories
	tempDaemonsDir := t.TempDir()
	tempLaunchAgentsDir := t.TempDir()
	
	// Override directories for testing
	oldDaemonsDir := DaemonsDir
	oldLaunchAgentsDir := LaunchAgentsDir
	DaemonsDir = tempDaemonsDir
	LaunchAgentsDir = tempLaunchAgentsDir
	defer func() {
		DaemonsDir = oldDaemonsDir
		LaunchAgentsDir = oldLaunchAgentsDir
	}()

	// Create a daemon plist
	daemonPlistPath := filepath.Join(tempDaemonsDir, "test-daemon.plist")
	plistContent := `<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
	<key>Label</key>
	<string>com.example.test</string>
</dict>
</plist>`
	err := os.WriteFile(daemonPlistPath, []byte(plistContent), 0644)
	require.NoError(t, err)

	// Test not installed
	installed, err := IsInstalled("test-daemon")
	assert.NoError(t, err)
	assert.False(t, installed)

	// Create installed plist
	installedPlistPath := filepath.Join(tempLaunchAgentsDir, "com.example.test.plist")
	err = os.WriteFile(installedPlistPath, []byte(plistContent), 0644)
	require.NoError(t, err)

	// Test installed
	installed, err = IsInstalled("test-daemon")
	assert.NoError(t, err)
	assert.True(t, installed)
}

func TestIsRunning(t *testing.T) {
	// Skip if not on macOS or if launchctl is not available
	if _, err := exec.LookPath("launchctl"); err != nil {
		t.Skip("launchctl command not available")
	}
	if _, err := exec.LookPath("defaults"); err != nil {
		t.Skip("defaults command not available")
	}

	// Create temp directory
	tempDir := t.TempDir()
	
	// Override DaemonsDir for testing
	oldDaemonsDir := DaemonsDir
	DaemonsDir = tempDir
	defer func() { DaemonsDir = oldDaemonsDir }()

	// Create a daemon plist
	daemonPlistPath := filepath.Join(tempDir, "test-daemon.plist")
	plistContent := `<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
	<key>Label</key>
	<string>com.example.test-not-running</string>
</dict>
</plist>`
	err := os.WriteFile(daemonPlistPath, []byte(plistContent), 0644)
	require.NoError(t, err)

	// Test with a daemon that's definitely not running
	running, err := IsRunning("test-daemon")
	assert.NoError(t, err)
	assert.False(t, running)
}

func TestCopyFile(t *testing.T) {
	// Create temp directory
	tempDir := t.TempDir()
	
	// Create source file
	srcPath := filepath.Join(tempDir, "source.txt")
	content := []byte("test content")
	err := os.WriteFile(srcPath, content, 0644)
	require.NoError(t, err)

	// Test successful copy
	dstPath := filepath.Join(tempDir, "destination.txt")
	err = CopyFile(srcPath, dstPath)
	assert.NoError(t, err)

	// Verify content
	copiedContent, err := os.ReadFile(dstPath)
	require.NoError(t, err)
	assert.Equal(t, content, copiedContent)

	// Verify permissions
	info, err := os.Stat(dstPath)
	require.NoError(t, err)
	assert.Equal(t, os.FileMode(0600), info.Mode().Perm())

	// Test error cases
	err = CopyFile("/non/existent/file", dstPath)
	assert.Error(t, err)

	err = CopyFile(srcPath, "/invalid/path/file")
	assert.Error(t, err)
}

func TestRunLaunchctl(t *testing.T) {
	// Skip if launchctl is not available
	if _, err := exec.LookPath("launchctl"); err != nil {
		t.Skip("launchctl command not available")
	}

	// Test with a safe command
	err := RunLaunchctl("help")
	assert.NoError(t, err)
}

func TestGetWorkingDirectory(t *testing.T) {
	// Skip if 'defaults' command is not available
	if _, err := exec.LookPath("defaults"); err != nil {
		t.Skip("defaults command not available")
	}

	// Create a test plist file
	tempDir := t.TempDir()
	plistPath := filepath.Join(tempDir, "test.plist")
	
	plistContent := `<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
	<key>WorkingDirectory</key>
	<string>/usr/local/test</string>
</dict>
</plist>`
	
	err := os.WriteFile(plistPath, []byte(plistContent), 0644)
	require.NoError(t, err)

	plistPathWithoutExt := strings.TrimSuffix(plistPath, ".plist")
	dir, err := GetWorkingDirectory(plistPathWithoutExt)
	assert.NoError(t, err)
	assert.Equal(t, "/usr/local/test", dir)
}

func TestGetStdoutPath(t *testing.T) {
	// Skip if 'defaults' command is not available
	if _, err := exec.LookPath("defaults"); err != nil {
		t.Skip("defaults command not available")
	}

	// Create a test plist file
	tempDir := t.TempDir()
	plistPath := filepath.Join(tempDir, "test.plist")
	
	plistContent := `<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
	<key>StandardOutPath</key>
	<string>/var/log/test.out</string>
</dict>
</plist>`
	
	err := os.WriteFile(plistPath, []byte(plistContent), 0644)
	require.NoError(t, err)

	plistPathWithoutExt := strings.TrimSuffix(plistPath, ".plist")
	path, err := GetStdoutPath(plistPathWithoutExt)
	assert.NoError(t, err)
	assert.Equal(t, "/var/log/test.out", path)
}

func TestGetStderrPath(t *testing.T) {
	// Skip if 'defaults' command is not available
	if _, err := exec.LookPath("defaults"); err != nil {
		t.Skip("defaults command not available")
	}

	// Create a test plist file
	tempDir := t.TempDir()
	plistPath := filepath.Join(tempDir, "test.plist")
	
	plistContent := `<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
	<key>StandardErrorPath</key>
	<string>/var/log/test.err</string>
</dict>
</plist>`
	
	err := os.WriteFile(plistPath, []byte(plistContent), 0644)
	require.NoError(t, err)

	plistPathWithoutExt := strings.TrimSuffix(plistPath, ".plist")
	path, err := GetStderrPath(plistPathWithoutExt)
	assert.NoError(t, err)
	assert.Equal(t, "/var/log/test.err", path)
}