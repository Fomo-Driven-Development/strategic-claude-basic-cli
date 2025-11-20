package templates

import (
	"testing"
)

func TestGetTemplate(t *testing.T) {
	tests := []struct {
		name    string
		id      string
		wantErr bool
	}{
		{
			name:    "get main template",
			id:      "main",
			wantErr: false,
		},
		{
			name:    "get ccr template",
			id:      "ccr",
			wantErr: false,
		},
		{
			name:    "get web-explorer template",
			id:      "web-explorer",
			wantErr: false,
		},
		{
			name:    "get non-existent template",
			id:      "nonexistent",
			wantErr: true,
		},
		{
			name:    "empty id",
			id:      "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetTemplate(tt.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetTemplate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if got.ID != tt.id {
					t.Errorf("GetTemplate() got ID = %v, want %v", got.ID, tt.id)
				}
				if err := got.IsValid(); err != nil {
					t.Errorf("GetTemplate() returned invalid template: %v", err)
				}
			}
		})
	}
}

func TestGetDefaultTemplate(t *testing.T) {
	template, err := GetDefaultTemplate()
	if err != nil {
		t.Fatalf("GetDefaultTemplate() error = %v", err)
	}

	if template.ID != DefaultTemplateID {
		t.Errorf("GetDefaultTemplate() got ID = %v, want %v", template.ID, DefaultTemplateID)
	}

	if err := template.IsValid(); err != nil {
		t.Errorf("GetDefaultTemplate() returned invalid template: %v", err)
	}
}

func TestListTemplates(t *testing.T) {
	templates := ListTemplates()

	if len(templates) == 0 {
		t.Error("ListTemplates() returned empty list")
	}

	// Check that templates are sorted by ID
	for i := 1; i < len(templates); i++ {
		if templates[i-1].ID >= templates[i].ID {
			t.Errorf("ListTemplates() not sorted by ID: %s >= %s", templates[i-1].ID, templates[i].ID)
		}
	}

	// Check that all templates are valid
	for _, template := range templates {
		if err := template.IsValid(); err != nil {
			t.Errorf("ListTemplates() contains invalid template %s: %v", template.ID, err)
		}
	}
}

func TestListActiveTemplates(t *testing.T) {
	templates := ListActiveTemplates()

	if len(templates) == 0 {
		t.Error("ListActiveTemplates() returned empty list")
	}

	// Check that no deprecated templates are returned
	for _, template := range templates {
		if template.Deprecated {
			t.Errorf("ListActiveTemplates() contains deprecated template: %s", template.ID)
		}
	}

	// Check that all templates are valid
	for _, template := range templates {
		if err := template.IsValid(); err != nil {
			t.Errorf("ListActiveTemplates() contains invalid template %s: %v", template.ID, err)
		}
	}
}

func TestValidateTemplateID(t *testing.T) {
	tests := []struct {
		name    string
		id      string
		wantErr bool
	}{
		{
			name:    "valid main template",
			id:      "main",
			wantErr: false,
		},
		{
			name:    "valid ccr template",
			id:      "ccr",
			wantErr: false,
		},
		{
			name:    "valid web-explorer template",
			id:      "web-explorer",
			wantErr: false,
		},
		{
			name:    "invalid template",
			id:      "invalid",
			wantErr: true,
		},
		{
			name:    "empty id",
			id:      "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateTemplateID(tt.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateTemplateID() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGetTemplateIDs(t *testing.T) {
	ids := GetTemplateIDs()

	if len(ids) == 0 {
		t.Error("GetTemplateIDs() returned empty list")
	}

	// Check that IDs are sorted
	for i := 1; i < len(ids); i++ {
		if ids[i-1] >= ids[i] {
			t.Errorf("GetTemplateIDs() not sorted: %s >= %s", ids[i-1], ids[i])
		}
	}

	// Check that all IDs are valid
	for _, id := range ids {
		if err := ValidateTemplateID(id); err != nil {
			t.Errorf("GetTemplateIDs() contains invalid template ID %s: %v", id, err)
		}
	}

	// Check that registry contains expected templates
	expectedTemplates := []string{"main", "ccr", "web-explorer"}
	for _, expected := range expectedTemplates {
		found := false
		for _, id := range ids {
			if id == expected {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("GetTemplateIDs() missing expected template: %s", expected)
		}
	}
}

func TestFilterTemplatesByLanguage(t *testing.T) {
	tests := []struct {
		name          string
		language      string
		expectAtLeast int // Minimum number of templates expected
	}{
		{
			name:          "empty language (all templates)",
			language:      "",
			expectAtLeast: 3, // main, ccr, and web-explorer are all language-agnostic
		},
		{
			name:          "go language",
			language:      "go",
			expectAtLeast: 0, // May have go-specific templates in future
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			templates := FilterTemplatesByLanguage(tt.language)

			if len(templates) < tt.expectAtLeast {
				t.Errorf("FilterTemplatesByLanguage() got %d templates, want at least %d", len(templates), tt.expectAtLeast)
			}

			// Check that all returned templates match the language filter
			for _, template := range templates {
				if tt.language != "" && template.Language != "" && template.Language != tt.language {
					t.Errorf("FilterTemplatesByLanguage() returned template with wrong language: got %s, want %s", template.Language, tt.language)
				}
			}
		})
	}
}

func TestFilterTemplatesByTag(t *testing.T) {
	tests := []struct {
		name string
		tag  string
	}{
		{
			name: "general tag",
			tag:  "general",
		},
		{
			name: "ccr tag",
			tag:  "ccr",
		},
		{
			name: "nonexistent tag",
			tag:  "nonexistent",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			templates := FilterTemplatesByTag(tt.tag)

			// Check that all returned templates have the specified tag
			for _, template := range templates {
				if !template.HasTag(tt.tag) {
					t.Errorf("FilterTemplatesByTag() returned template without tag %s: %s", tt.tag, template.ID)
				}
			}
		})
	}
}
