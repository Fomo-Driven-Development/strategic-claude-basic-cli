package config

import (
	"strings"
	"testing"
	"time"
)

func TestGetFrameworkDirectories(t *testing.T) {
	dirs := GetFrameworkDirectories()

	expectedDirs := []string{CoreDir, TemplatesDir}

	if len(dirs) != len(expectedDirs) {
		t.Errorf("Expected %d directories, got %d", len(expectedDirs), len(dirs))
	}

	for i, expected := range expectedDirs {
		if i >= len(dirs) || dirs[i] != expected {
			t.Errorf("Expected directory %s at index %d, got %s", expected, i, dirs[i])
		}
	}

	// Test that all expected directories are present
	dirMap := make(map[string]bool)
	for _, dir := range dirs {
		dirMap[dir] = true
	}

	for _, expected := range expectedDirs {
		if !dirMap[expected] {
			t.Errorf("Expected directory %s not found in result", expected)
		}
	}
}

func TestGetCoreDirectories(t *testing.T) {
	dirs := GetCoreDirectories()

	expectedDirs := []string{CoreDir, TemplatesDir}

	if len(dirs) != len(expectedDirs) {
		t.Errorf("Expected %d directories, got %d", len(expectedDirs), len(dirs))
	}

	for i, expected := range expectedDirs {
		if i >= len(dirs) || dirs[i] != expected {
			t.Errorf("Expected directory %s at index %d, got %s", expected, i, dirs[i])
		}
	}

	// Verify these are the directories that get replaced during updates
	if !contains(dirs, CoreDir) {
		t.Error("Expected CoreDir in core directories")
	}
	if !contains(dirs, TemplatesDir) {
		t.Error("Expected TemplatesDir in core directories")
	}

	// Verify user directories are not in core directories
	if contains(dirs, ArchivesDir) {
		t.Error("ArchivesDir should not be in core directories")
	}
	if contains(dirs, IssuesDir) {
		t.Error("IssuesDir should not be in core directories")
	}
}

func TestGetUserPreservedDirectories(t *testing.T) {
	dirs := GetUserPreservedDirectories()

	expectedDirs := []string{ArchivesDir, DecisionsDir, IssuesDir, PlanDir, ProductDir, ResearchDir, SummaryDir, ToolsDir}

	if len(dirs) != len(expectedDirs) {
		t.Errorf("Expected %d directories, got %d", len(expectedDirs), len(dirs))
	}

	for i, expected := range expectedDirs {
		if i >= len(dirs) || dirs[i] != expected {
			t.Errorf("Expected directory %s at index %d, got %s", expected, i, dirs[i])
		}
	}

	// Verify these are user content directories
	if !contains(dirs, ArchivesDir) {
		t.Error("Expected ArchivesDir in user preserved directories")
	}
	if !contains(dirs, IssuesDir) {
		t.Error("Expected IssuesDir in user preserved directories")
	}

	// Verify core directories are not in user preserved directories
	if contains(dirs, CoreDir) {
		t.Error("CoreDir should not be in user preserved directories")
	}
	if contains(dirs, TemplatesDir) {
		t.Error("TemplatesDir should not be in user preserved directories")
	}
}

func TestGetRequiredSymlinks(t *testing.T) {
	symlinks := GetRequiredSymlinks()

	// Test expected symlink paths
	expectedPaths := []string{
		"agents/strategic",
		"commands/strategic",
		"hooks/strategic",
	}

	if len(symlinks) != len(expectedPaths) {
		t.Errorf("Expected %d symlinks, got %d", len(expectedPaths), len(symlinks))
	}

	for _, path := range expectedPaths {
		target, exists := symlinks[path]
		if !exists {
			t.Errorf("Expected symlink path %s not found", path)
		}

		// Verify target format
		expectedTarget := "../../" + StrategicClaudeBasicDir + "/core/" + strings.Split(path, "/")[0]
		if target != expectedTarget {
			t.Errorf("Expected target %s for path %s, got %s", expectedTarget, path, target)
		}
	}

	// Test specific symlinks
	if target, exists := symlinks["agents/strategic"]; exists {
		expectedTarget := "../../" + StrategicClaudeBasicDir + "/core/agents"
		if target != expectedTarget {
			t.Errorf("Expected agents target %s, got %s", expectedTarget, target)
		}
	} else {
		t.Error("agents/strategic symlink not found")
	}

	if target, exists := symlinks["commands/strategic"]; exists {
		expectedTarget := "../../" + StrategicClaudeBasicDir + "/core/commands"
		if target != expectedTarget {
			t.Errorf("Expected commands target %s, got %s", expectedTarget, target)
		}
	} else {
		t.Error("commands/strategic symlink not found")
	}

	if target, exists := symlinks["hooks/strategic"]; exists {
		expectedTarget := "../../" + StrategicClaudeBasicDir + "/core/hooks"
		if target != expectedTarget {
			t.Errorf("Expected hooks target %s, got %s", expectedTarget, target)
		}
	} else {
		t.Error("hooks/strategic symlink not found")
	}
}

func TestGetBackupDirName(t *testing.T) {
	// Get a backup directory name
	backupName := GetBackupDirName()

	// Should start with backup prefix
	if !strings.HasPrefix(backupName, BackupDirPrefix) {
		t.Errorf("Expected backup name to start with %s, got %s", BackupDirPrefix, backupName)
	}

	// Should have timestamp format: prefix + YYYYMMDD-HHMMSS
	expectedLen := len(BackupDirPrefix) + 15 // 8 digits + dash + 6 digits
	if len(backupName) != expectedLen {
		t.Errorf("Expected backup name length %d, got %d (%s)", expectedLen, len(backupName), backupName)
	}

	// Test that multiple calls produce different names (due to timestamp)
	time.Sleep(time.Second + 10*time.Millisecond) // Ensure different timestamp
	backupName2 := GetBackupDirName()

	if backupName == backupName2 {
		t.Error("Expected different backup names for different timestamps")
	}

	// Test timestamp format
	timestamp := strings.TrimPrefix(backupName, BackupDirPrefix)
	if len(timestamp) != 15 { // YYYYMMDD-HHMMSS
		t.Errorf("Expected timestamp length 15, got %d (%s)", len(timestamp), timestamp)
	}

	// Check dash in the right position
	if timestamp[8] != '-' {
		t.Errorf("Expected dash at position 8 in timestamp, got %c", timestamp[8])
	}

	// Check all other characters are digits
	for i, char := range timestamp {
		if i == 8 { // Skip the dash
			continue
		}
		if char < '0' || char > '9' {
			t.Errorf("Expected digit at position %d in timestamp, got %c", i, char)
		}
	}
}

func TestIsUserPreservedPath(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		expected bool
	}{
		{
			name:     "archives directory",
			path:     ArchivesDir,
			expected: true,
		},
		{
			name:     "issues directory",
			path:     IssuesDir,
			expected: true,
		},
		{
			name:     "file in archives",
			path:     ArchivesDir + "/file.txt",
			expected: true,
		},
		{
			name:     "nested path in research",
			path:     ResearchDir + "/subfolder/file.md",
			expected: true,
		},
		{
			name:     "core directory (not preserved)",
			path:     CoreDir,
			expected: false,
		},
		{
			name:     "templates directory (not preserved)",
			path:     TemplatesDir,
			expected: false,
		},
		{
			name:     "file in core (not preserved)",
			path:     CoreDir + "/agents/file.txt",
			expected: false,
		},
		{
			name:     "random path (not preserved)",
			path:     "random/path",
			expected: false,
		},
		{
			name:     "similar but not exact match",
			path:     "archivesx/file.txt", // Similar to archives but not exact
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsUserPreservedPath(tt.path)
			if result != tt.expected {
				t.Errorf("Expected %v for path %s, got %v", tt.expected, tt.path, result)
			}
		})
	}
}

func TestIsCoreFile(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		expected bool
	}{
		{
			name:     "core directory",
			path:     CoreDir,
			expected: true,
		},
		{
			name:     "templates directory",
			path:     TemplatesDir,
			expected: true,
		},
		{
			name:     "file in core",
			path:     CoreDir + "/agents/agent.py",
			expected: true,
		},
		{
			name:     "nested path in templates",
			path:     TemplatesDir + "/hooks/hook.sh",
			expected: true,
		},
		{
			name:     "archives directory (not core)",
			path:     ArchivesDir,
			expected: false,
		},
		{
			name:     "issues directory (not core)",
			path:     IssuesDir,
			expected: false,
		},
		{
			name:     "file in archives (not core)",
			path:     ArchivesDir + "/old-file.txt",
			expected: false,
		},
		{
			name:     "random path (not core)",
			path:     "random/path",
			expected: false,
		},
		{
			name:     "similar but not exact match",
			path:     "corex/file.txt", // Similar to core but not exact
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsCoreFile(tt.path)
			if result != tt.expected {
				t.Errorf("Expected %v for path %s, got %v", tt.expected, tt.path, result)
			}
		})
	}
}

func TestConstants_Values(t *testing.T) {
	// Test that important constants have expected values
	tests := []struct {
		name     string
		constant interface{}
		nonEmpty bool
	}{
		{"DefaultRepoURL", DefaultRepoURL, true},
		{"FixedCommit", FixedCommit, true},
		{"Branch", Branch, true},
		{"StrategicClaudeBasicDir", StrategicClaudeBasicDir, true},
		{"ClaudeDir", ClaudeDir, true},
		{"BackupDirPrefix", BackupDirPrefix, true},
		{"TempDirPrefix", TempDirPrefix, true},
		{"CoreDir", CoreDir, true},
		{"TemplatesDir", TemplatesDir, true},
		{"ArchivesDir", ArchivesDir, true},
		{"IssuesDir", IssuesDir, true},
		{"AgentsDir", AgentsDir, true},
		{"CommandsDir", CommandsDir, true},
		{"HooksDir", HooksDir, true},
		{"AppName", AppName, true},
		{"AppDescription", AppDescription, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.nonEmpty {
				str, ok := tt.constant.(string)
				if !ok {
					t.Errorf("Expected %s to be a string", tt.name)
				} else if strings.TrimSpace(str) == "" {
					t.Errorf("Expected %s to be non-empty", tt.name)
				}
			}
		})
	}

	// Test specific value formats
	if !strings.HasPrefix(DefaultRepoURL, "https://") {
		t.Errorf("Expected DefaultRepoURL to be HTTPS, got %s", DefaultRepoURL)
	}

	if !strings.HasSuffix(DefaultRepoURL, ".git") {
		t.Errorf("Expected DefaultRepoURL to end with .git, got %s", DefaultRepoURL)
	}

	if Branch != "main" {
		t.Errorf("Expected Branch to be 'main', got %s", Branch)
	}

	if StrategicClaudeBasicDir != ".strategic-claude-basic" {
		t.Errorf("Expected StrategicClaudeBasicDir to be '.strategic-claude-basic', got %s", StrategicClaudeBasicDir)
	}

	if ClaudeDir != ".claude" {
		t.Errorf("Expected ClaudeDir to be '.claude', got %s", ClaudeDir)
	}
}

func TestTimeoutConstants(t *testing.T) {
	// Test that timeout constants are reasonable
	if DefaultGitTimeout <= 0 {
		t.Errorf("Expected DefaultGitTimeout to be positive, got %v", DefaultGitTimeout)
	}

	if DefaultNetworkTimeout <= 0 {
		t.Errorf("Expected DefaultNetworkTimeout to be positive, got %v", DefaultNetworkTimeout)
	}

	// Test that timeouts are in reasonable range (not too short, not too long)
	if DefaultGitTimeout < 5*time.Second {
		t.Errorf("Expected DefaultGitTimeout to be at least 5s, got %v", DefaultGitTimeout)
	}

	if DefaultGitTimeout > 5*time.Minute {
		t.Errorf("Expected DefaultGitTimeout to be at most 5m, got %v", DefaultGitTimeout)
	}
}

func TestPermissionConstants(t *testing.T) {
	// Test file permissions are reasonable
	if DirPermissions < 0o700 || DirPermissions > 0o777 {
		t.Errorf("Expected DirPermissions to be valid octal, got %o", DirPermissions)
	}

	if FilePermissions < 0o600 || FilePermissions > 0o666 {
		t.Errorf("Expected FilePermissions to be valid octal, got %o", FilePermissions)
	}

	// Test that directory permissions allow owner read/write/execute
	if DirPermissions&0o700 != 0o700 {
		t.Errorf("Expected DirPermissions to allow owner rwx, got %o", DirPermissions)
	}

	// Test that file permissions allow owner read/write
	if FilePermissions&0o600 != 0o600 {
		t.Errorf("Expected FilePermissions to allow owner rw, got %o", FilePermissions)
	}
}

func TestExitCodeConstants(t *testing.T) {
	// Test that exit codes are in valid range (0-255)
	exitCodes := []int{
		ExitSuccess, ExitGeneralError, ExitValidationError, ExitPermissionError,
		ExitNetworkError, ExitUserCancellation, ExitInstallationError,
		ExitAlreadyInstalled, ExitNotInstalled,
	}

	for i, code := range exitCodes {
		if code < 0 || code > 255 {
			t.Errorf("Exit code %d is out of range (0-255): %d", i, code)
		}
	}

	// Test that success is 0
	if ExitSuccess != 0 {
		t.Errorf("Expected ExitSuccess to be 0, got %d", ExitSuccess)
	}

	// Test that error codes are non-zero
	errorCodes := exitCodes[1:] // All except ExitSuccess
	for i, code := range errorCodes {
		if code == 0 {
			t.Errorf("Error exit code %d should be non-zero", i+1)
		}
	}
}

func TestBackupConstants(t *testing.T) {
	// Test backup configuration is reasonable
	if MaxBackupAge <= 0 {
		t.Errorf("Expected MaxBackupAge to be positive, got %v", MaxBackupAge)
	}

	if MaxBackups <= 0 {
		t.Errorf("Expected MaxBackups to be positive, got %d", MaxBackups)
	}

	// Test reasonable ranges
	if MaxBackupAge < 24*time.Hour {
		t.Errorf("Expected MaxBackupAge to be at least 1 day, got %v", MaxBackupAge)
	}

	if MaxBackups > 100 {
		t.Errorf("Expected MaxBackups to be reasonable (<= 100), got %d", MaxBackups)
	}
}

// Helper function for slice contains check
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

// Test the helper function itself
func TestContainsHelper(t *testing.T) {
	slice := []string{"a", "b", "c"}

	if !contains(slice, "a") {
		t.Error("Expected contains to find 'a'")
	}

	if !contains(slice, "b") {
		t.Error("Expected contains to find 'b'")
	}

	if !contains(slice, "c") {
		t.Error("Expected contains to find 'c'")
	}

	if contains(slice, "d") {
		t.Error("Expected contains to not find 'd'")
	}

	if contains([]string{}, "a") {
		t.Error("Expected contains to not find item in empty slice")
	}
}
