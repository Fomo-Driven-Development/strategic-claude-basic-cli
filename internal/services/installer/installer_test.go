package installer

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/Fomo-Driven-Development/strategic-claude-basic-cli/internal/config"
	"github.com/Fomo-Driven-Development/strategic-claude-basic-cli/internal/models"
	"github.com/Fomo-Driven-Development/strategic-claude-basic-cli/internal/templates"
)

func TestNew(t *testing.T) {
	service := New()

	if service == nil {
		t.Fatal("New() returned nil")
	}

	if service.gitService == nil {
		t.Error("Git service not initialized")
	}

	if service.filesystemService == nil {
		t.Error("Filesystem service not initialized")
	}

	if service.statusService == nil {
		t.Error("Status service not initialized")
	}

	if service.symlinkService == nil {
		t.Error("Symlink service not initialized")
	}
}

func TestAnalyzeInstallation(t *testing.T) {
	service := New()

	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "installer-test-")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	tests := []struct {
		name          string
		installConfig models.InstallConfig
		setupFunc     func(string) error
		expectedType  models.InstallationType
		expectError   bool
	}{
		{
			name: "new installation in empty directory",
			installConfig: models.InstallConfig{
				TargetDir:   tempDir,
				TemplateID:  "main",
				Force:       false,
				ForceCore:   false,
				SkipConfirm: true,
				NoBackup:    true,
			},
			setupFunc:    nil, // No setup needed for empty directory
			expectedType: models.InstallationTypeNew,
			expectError:  false,
		},
		{
			name: "force installation",
			installConfig: models.InstallConfig{
				TargetDir:   tempDir,
				TemplateID:  "main",
				Force:       true,
				ForceCore:   false,
				SkipConfirm: true,
				NoBackup:    true,
			},
			setupFunc: func(dir string) error {
				// Create existing strategic-claude-basic directory
				strategicDir := filepath.Join(dir, config.StrategicClaudeBasicDir)
				return os.MkdirAll(strategicDir, 0755)
			},
			expectedType: models.InstallationTypeOverwrite,
			expectError:  false,
		},
		{
			name: "force-core installation",
			installConfig: models.InstallConfig{
				TargetDir:   tempDir,
				TemplateID:  "main",
				Force:       false,
				ForceCore:   true,
				SkipConfirm: true,
				NoBackup:    true,
			},
			setupFunc: func(dir string) error {
				// Create existing strategic-claude-basic directory with some content
				strategicDir := filepath.Join(dir, config.StrategicClaudeBasicDir)
				coreDir := filepath.Join(strategicDir, config.CoreDir)
				return os.MkdirAll(coreDir, 0755)
			},
			expectedType: models.InstallationTypeUpdate,
			expectError:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clean up the temp directory
			os.RemoveAll(tempDir)
			if err := os.MkdirAll(tempDir, 0755); err != nil {
				t.Fatalf("Failed to create temp dir: %v", err)
			}

			// Run setup if provided
			if tt.setupFunc != nil {
				if err := tt.setupFunc(tempDir); err != nil {
					t.Fatalf("Setup failed: %v", err)
				}
			}

			plan, err := service.AnalyzeInstallation(tt.installConfig)

			if tt.expectError && err == nil {
				t.Error("Expected error but got none")
				return
			}

			if !tt.expectError && err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if !tt.expectError {
				if plan == nil {
					t.Error("Expected plan but got nil")
					return
				}

				if plan.InstallationType != tt.expectedType {
					t.Errorf("Expected installation type %s, got %s", tt.expectedType, plan.InstallationType)
				}

				if plan.TargetDir != tempDir {
					t.Errorf("Expected target dir %s, got %s", tempDir, plan.TargetDir)
				}
			}
		})
	}
}

func TestDetermineInstallationType(t *testing.T) {
	service := New()

	tests := []struct {
		name          string
		status        *models.StatusInfo
		installConfig models.InstallConfig
		expectedType  models.InstallationType
	}{
		{
			name:          "force flag set",
			status:        &models.StatusInfo{IsInstalled: true},
			installConfig: models.InstallConfig{Force: true},
			expectedType:  models.InstallationTypeOverwrite,
		},
		{
			name:          "force-core flag set",
			status:        &models.StatusInfo{IsInstalled: true},
			installConfig: models.InstallConfig{ForceCore: true},
			expectedType:  models.InstallationTypeUpdate,
		},
		{
			name:          "not installed",
			status:        &models.StatusInfo{IsInstalled: false},
			installConfig: models.InstallConfig{},
			expectedType:  models.InstallationTypeNew,
		},
		{
			name:          "installed with no flags",
			status:        &models.StatusInfo{IsInstalled: true},
			installConfig: models.InstallConfig{},
			expectedType:  models.InstallationTypeOverwrite,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := service.determineInstallationType(tt.status, tt.installConfig)

			if result != tt.expectedType {
				t.Errorf("Expected %s, got %s", tt.expectedType, result)
			}
		})
	}
}

func TestNeedsBackup(t *testing.T) {
	service := New()

	tests := []struct {
		name          string
		plan          *models.InstallationPlan
		installConfig models.InstallConfig
		expected      bool
	}{
		{
			name:          "no backup flag set",
			plan:          &models.InstallationPlan{WillReplace: []string{"some-file"}},
			installConfig: models.InstallConfig{NoBackup: true},
			expected:      false,
		},
		{
			name:          "files will be replaced",
			plan:          &models.InstallationPlan{WillReplace: []string{"some-file"}},
			installConfig: models.InstallConfig{NoBackup: false},
			expected:      true,
		},
		{
			name:          "no files to replace",
			plan:          &models.InstallationPlan{WillReplace: []string{}},
			installConfig: models.InstallConfig{NoBackup: false},
			expected:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := service.needsBackup(tt.plan, tt.installConfig)

			if result != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestAnalyzeFileOperations(t *testing.T) {
	service := New()
	tempDir, err := os.MkdirTemp("", "installer-test-")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	tests := []struct {
		name             string
		installType      models.InstallationType
		setupFunc        func(string) error
		expectedCreate   []string
		expectedReplace  []string
		expectedPreserve []string
	}{
		{
			name:             "new installation",
			installType:      models.InstallationTypeNew,
			setupFunc:        nil,
			expectedCreate:   []string{config.StrategicClaudeBasicDir},
			expectedReplace:  []string{},
			expectedPreserve: []string{},
		},
		{
			name:        "update installation",
			installType: models.InstallationTypeUpdate,
			setupFunc: func(dir string) error {
				strategicDir := filepath.Join(dir, config.StrategicClaudeBasicDir)
				coreDir := filepath.Join(strategicDir, config.CoreDir)
				return os.MkdirAll(coreDir, 0755)
			},
			expectedCreate: []string{
				filepath.Join(config.StrategicClaudeBasicDir, config.TemplatesDir),
			},
			expectedReplace: []string{
				filepath.Join(config.StrategicClaudeBasicDir, config.CoreDir),
			},
			expectedPreserve: []string{
				filepath.Join(config.StrategicClaudeBasicDir, config.ArchivesDir),
				filepath.Join(config.StrategicClaudeBasicDir, config.DecisionsDir),
				filepath.Join(config.StrategicClaudeBasicDir, config.IssuesDir),
				filepath.Join(config.StrategicClaudeBasicDir, config.PlanDir),
				filepath.Join(config.StrategicClaudeBasicDir, config.ProductDir),
				filepath.Join(config.StrategicClaudeBasicDir, config.ResearchDir),
				filepath.Join(config.StrategicClaudeBasicDir, config.SummaryDir),
				filepath.Join(config.StrategicClaudeBasicDir, config.ToolsDir),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clean up and setup
			os.RemoveAll(tempDir)
			if err := os.MkdirAll(tempDir, 0755); err != nil {
				t.Fatalf("Failed to create temp dir: %v", err)
			}

			if tt.setupFunc != nil {
				if err := tt.setupFunc(tempDir); err != nil {
					t.Fatalf("Setup failed: %v", err)
				}
			}

			// Use a default template for testing
			defaultTemplate := templates.Template{
				ID:      "main",
				Name:    "Test Template",
				RepoURL: "https://test.com/repo.git",
				Branch:  "main",
				Commit:  "test123",
			}
			plan := models.NewInstallationPlan(tempDir, tt.installType, defaultTemplate)
			status := &models.StatusInfo{TargetDir: tempDir}

			service.analyzeFileOperations(plan, status)

			// Check created files
			if len(plan.WillCreate) != len(tt.expectedCreate) {
				t.Errorf("Expected %d files to create, got %d", len(tt.expectedCreate), len(plan.WillCreate))
			}

			// Check replaced files
			if len(plan.WillReplace) != len(tt.expectedReplace) {
				t.Errorf("Expected %d files to replace, got %d", len(tt.expectedReplace), len(plan.WillReplace))
			}

			// Check preserved files
			if len(plan.WillPreserve) != len(tt.expectedPreserve) {
				t.Errorf("Expected %d files to preserve, got %d", len(tt.expectedPreserve), len(plan.WillPreserve))
			}
		})
	}
}
