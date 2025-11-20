package status

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/Fomo-Driven-Development/strategic-claude-basic-cli/internal/config"
	"github.com/Fomo-Driven-Development/strategic-claude-basic-cli/internal/models"
)

// createTestDirectory creates a temporary directory structure for testing
func createTestDirectory(t *testing.T, structure map[string]interface{}) string {
	t.Helper()

	tempDir, err := os.MkdirTemp("", "status-test-")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}

	t.Cleanup(func() {
		os.RemoveAll(tempDir)
	})

	createStructure(t, tempDir, structure)
	return tempDir
}

// createStructure recursively creates files and directories
func createStructure(t *testing.T, basePath string, structure map[string]interface{}) {
	t.Helper()

	for name, content := range structure {
		path := filepath.Join(basePath, name)

		switch v := content.(type) {
		case map[string]interface{}:
			// Directory
			if err := os.MkdirAll(path, 0755); err != nil {
				t.Fatalf("Failed to create directory %s: %v", path, err)
			}
			createStructure(t, path, v)
		case string:
			// File with content
			if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
				t.Fatalf("Failed to create parent directory for %s: %v", path, err)
			}
			if err := os.WriteFile(path, []byte(v), 0644); err != nil {
				t.Fatalf("Failed to create file %s: %v", path, err)
			}
		case nil:
			// Empty directory
			if err := os.MkdirAll(path, 0755); err != nil {
				t.Fatalf("Failed to create directory %s: %v", path, err)
			}
		}
	}
}

// createSymlink creates a symlink for testing
func createSymlink(t *testing.T, oldname, newname string) {
	t.Helper()

	if err := os.MkdirAll(filepath.Dir(newname), 0755); err != nil {
		t.Fatalf("Failed to create parent directory for symlink %s: %v", newname, err)
	}

	if err := os.Symlink(oldname, newname); err != nil {
		t.Fatalf("Failed to create symlink %s -> %s: %v", newname, oldname, err)
	}
}

func TestService_CheckInstallation_NoInstallation(t *testing.T) {
	// Create empty directory
	tempDir := createTestDirectory(t, map[string]interface{}{})

	service := NewService()
	status, err := service.CheckInstallation(tempDir)

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if status.IsInstalled {
		t.Error("Expected IsInstalled to be false")
	}

	if status.StrategicClaudeDir {
		t.Error("Expected StrategicClaudeDir to be false")
	}

	if status.ClaudeDir {
		t.Error("Expected ClaudeDir to be false")
	}

	if len(status.Symlinks) != 3 { // Should still check for 3 expected symlinks
		t.Errorf("Expected 3 symlinks to be checked, got %d", len(status.Symlinks))
	}

	// All symlinks should be invalid/not found
	for _, symlink := range status.Symlinks {
		if symlink.Valid {
			t.Errorf("Expected symlink %s to be invalid", symlink.Name)
		}
		if symlink.Exists {
			t.Errorf("Expected symlink %s to not exist", symlink.Name)
		}
	}
}

func TestService_CheckInstallation_PartialInstallation(t *testing.T) {
	// Create only .strategic-claude-basic directory without .claude
	structure := map[string]interface{}{
		config.StrategicClaudeBasicDir: map[string]interface{}{
			config.CoreDir: map[string]interface{}{
				config.AgentsDir:   nil,
				config.CommandsDir: nil,
				config.HooksDir:    nil,
			},
			config.TemplatesDir: nil,
			config.ConfigDir:    nil,
		},
	}

	tempDir := createTestDirectory(t, structure)

	service := NewService()
	status, err := service.CheckInstallation(tempDir)

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if status.IsInstalled {
		t.Error("Expected IsInstalled to be false for partial installation")
	}

	if !status.StrategicClaudeDir {
		t.Error("Expected StrategicClaudeDir to be true")
	}

	if status.ClaudeDir {
		t.Error("Expected ClaudeDir to be false")
	}

	if !status.HasIssues() {
		t.Error("Expected to have issues for partial installation")
	}

	// Check for partial installation issue
	found := false
	for _, issue := range status.Issues {
		if issue == "Partial installation detected: .strategic-claude-basic exists but .claude directory is missing" {
			found = true
			break
		}
	}
	if !found {
		t.Error("Expected partial installation issue to be detected")
	}
}

func TestService_CheckInstallation_CompleteInstallation(t *testing.T) {
	// Create complete installation structure
	structure := map[string]interface{}{
		config.StrategicClaudeBasicDir: map[string]interface{}{
			config.CoreDir: map[string]interface{}{
				config.AgentsDir:   nil,
				config.CommandsDir: nil,
				config.HooksDir:    nil,
			},
			config.TemplatesDir: nil,
			config.ConfigDir:    nil,
		},
		config.ClaudeDir: map[string]interface{}{
			config.AgentsDir:   nil,
			config.CommandsDir: nil,
			config.HooksDir:    nil,
		},
		config.CodexDir: map[string]interface{}{
			config.PromptsDir: nil,
			config.HooksDir:   nil,
		},
	}

	tempDir := createTestDirectory(t, structure)

	// Create valid symlinks
	requiredSymlinks := config.GetRequiredSymlinks()
	for symlinkPath, target := range requiredSymlinks {
		symlinkFullPath := filepath.Join(tempDir, config.ClaudeDir, symlinkPath)
		// The target is relative to the symlink location, so it should be used as-is
		createSymlink(t, target, symlinkFullPath)
	}

	service := NewService()
	status, err := service.CheckInstallation(tempDir)

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if !status.IsInstalled {
		t.Error("Expected IsInstalled to be true")
	}

	if !status.StrategicClaudeDir {
		t.Error("Expected StrategicClaudeDir to be true")
	}

	if !status.ClaudeDir {
		t.Error("Expected ClaudeDir to be true")
	}

	if status.ValidSymlinks() != 3 {
		t.Errorf("Expected 3 valid symlinks, got %d", status.ValidSymlinks())
	}

	// Should have minimal or no issues for complete installation
	if len(status.Issues) > 0 {
		t.Errorf("Expected no issues for complete installation, got: %v", status.Issues)
	}
}

func TestService_CheckInstallation_BrokenSymlinks(t *testing.T) {
	// Create installation with broken symlinks
	structure := map[string]interface{}{
		config.StrategicClaudeBasicDir: map[string]interface{}{
			config.CoreDir: map[string]interface{}{
				config.AgentsDir:   nil,
				config.CommandsDir: nil,
				config.HooksDir:    nil,
			},
			config.TemplatesDir: nil,
			config.ConfigDir:    nil,
		},
		config.ClaudeDir: map[string]interface{}{
			config.AgentsDir:   nil,
			config.CommandsDir: nil,
			config.HooksDir:    nil,
		},
	}

	tempDir := createTestDirectory(t, structure)

	// Create broken symlinks (pointing to wrong targets)
	symlinkPaths := []string{
		filepath.Join(config.AgentsDir, "strategic"),
		filepath.Join(config.CommandsDir, "strategic"),
		filepath.Join(config.HooksDir, "strategic"),
	}

	for _, symlinkPath := range symlinkPaths {
		symlinkFullPath := filepath.Join(tempDir, config.ClaudeDir, symlinkPath)
		createSymlink(t, "/nonexistent/path", symlinkFullPath)
	}

	service := NewService()
	status, err := service.CheckInstallation(tempDir)

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if !status.StrategicClaudeDir {
		t.Error("Expected StrategicClaudeDir to be true")
	}

	if !status.ClaudeDir {
		t.Error("Expected ClaudeDir to be true")
	}

	if status.ValidSymlinks() != 0 {
		t.Errorf("Expected 0 valid symlinks, got %d", status.ValidSymlinks())
	}

	// All symlinks should exist but be invalid
	for _, symlink := range status.Symlinks {
		if !symlink.Exists {
			t.Errorf("Expected symlink %s to exist", symlink.Name)
		}
		if symlink.Valid {
			t.Errorf("Expected symlink %s to be invalid", symlink.Name)
		}
		if symlink.Error == "" {
			t.Errorf("Expected symlink %s to have an error message", symlink.Name)
		}
	}

	if !status.HasIssues() {
		t.Error("Expected to have issues with broken symlinks")
	}
}

func TestService_CheckInstallation_MissingFrameworkDirectories(t *testing.T) {
	// Create installation missing some framework directories
	structure := map[string]interface{}{
		config.StrategicClaudeBasicDir: map[string]interface{}{
			config.CoreDir: map[string]interface{}{
				config.AgentsDir: nil,
				// Missing commands and hooks directories
			},
			// Missing guides, templates, config directories
		},
		config.ClaudeDir: map[string]interface{}{
			config.AgentsDir:   nil,
			config.CommandsDir: nil,
			config.HooksDir:    nil,
		},
	}

	tempDir := createTestDirectory(t, structure)

	service := NewService()
	status, err := service.CheckInstallation(tempDir)

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if !status.HasIssues() {
		t.Error("Expected to have issues with missing directories")
	}

	// Check for specific missing directory issues
	expectedIssues := []string{
		"Missing framework directory: templates",
		"Missing core subdirectory: core/commands",
		"Missing core subdirectory: core/hooks",
	}

	for _, expectedIssue := range expectedIssues {
		found := false
		for _, issue := range status.Issues {
			if issue == expectedIssue {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected issue '%s' not found in: %v", expectedIssue, status.Issues)
		}
	}
}

func TestService_CheckInstallation_InvalidTargetDirectory(t *testing.T) {
	service := NewService()

	// Test with non-existent directory
	_, err := service.CheckInstallation("/nonexistent/directory")
	if err == nil {
		t.Error("Expected error for non-existent directory")
	}

	// Test with file instead of directory
	tempFile, err := os.CreateTemp("", "not-a-directory")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())
	tempFile.Close()

	_, err = service.CheckInstallation(tempFile.Name())
	if err == nil {
		t.Error("Expected error when target is a file, not directory")
	}
}

func TestService_GetStatusSummary(t *testing.T) {
	service := NewService()

	tests := []struct {
		name     string
		setupFn  func() *models.StatusInfo
		expected string
	}{
		{
			name: "Not installed",
			setupFn: func() *models.StatusInfo {
				status := models.NewStatusInfo("/test")
				status.IsInstalled = false
				return status
			},
			expected: "Strategic Claude Basic is not installed",
		},
		{
			name: "Installed with issues",
			setupFn: func() *models.StatusInfo {
				status := models.NewStatusInfo("/test")
				status.IsInstalled = true
				status.AddIssue("Test issue")
				return status
			},
			expected: "Strategic Claude Basic is installed but has 1 issue(s)",
		},
		{
			name: "Installed correctly",
			setupFn: func() *models.StatusInfo {
				status := models.NewStatusInfo("/test")
				status.IsInstalled = true
				return status
			},
			expected: "Strategic Claude Basic is installed and configured correctly",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			status := tt.setupFn()
			summary := service.GetStatusSummary(status)
			if summary != tt.expected {
				t.Errorf("Expected summary '%s', got '%s'", tt.expected, summary)
			}
		})
	}
}

func TestService_detectStrategicClaudeBasic(t *testing.T) {
	service := NewService()

	tests := []struct {
		name         string
		structure    map[string]interface{}
		expectDir    bool
		expectIssues []string
	}{
		{
			name:         "No directory",
			structure:    map[string]interface{}{},
			expectDir:    false,
			expectIssues: []string{".strategic-claude-basic directory does not exist"},
		},
		{
			name: "Complete structure",
			structure: map[string]interface{}{
				config.StrategicClaudeBasicDir: map[string]interface{}{
					config.CoreDir: map[string]interface{}{
						config.AgentsDir:   nil,
						config.CommandsDir: nil,
						config.HooksDir:    nil,
					},
					config.TemplatesDir: nil,
					config.ConfigDir:    nil,
				},
			},
			expectDir:    true,
			expectIssues: []string{},
		},
		{
			name: "Missing core subdirectories",
			structure: map[string]interface{}{
				config.StrategicClaudeBasicDir: map[string]interface{}{
					config.CoreDir: map[string]interface{}{
						config.AgentsDir: nil,
						// Missing commands and hooks
					},
					config.TemplatesDir: nil,
					config.ConfigDir:    nil,
				},
			},
			expectDir: true,
			expectIssues: []string{
				"Missing core subdirectory: core/commands",
				"Missing core subdirectory: core/hooks",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tempDir := createTestDirectory(t, tt.structure)

			status := models.NewStatusInfo(tempDir)
			status.StrategicClaudeDirPath = filepath.Join(tempDir, config.StrategicClaudeBasicDir)

			err := service.detectStrategicClaudeBasic(status)
			if err != nil {
				t.Fatalf("Expected no error, got: %v", err)
			}

			if status.StrategicClaudeDir != tt.expectDir {
				t.Errorf("Expected StrategicClaudeDir to be %v, got %v", tt.expectDir, status.StrategicClaudeDir)
			}

			if len(status.Issues) != len(tt.expectIssues) {
				t.Errorf("Expected %d issues, got %d: %v", len(tt.expectIssues), len(status.Issues), status.Issues)
			}

			for _, expectedIssue := range tt.expectIssues {
				found := false
				for _, issue := range status.Issues {
					if issue == expectedIssue {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("Expected issue '%s' not found in: %v", expectedIssue, status.Issues)
				}
			}
		})
	}
}
