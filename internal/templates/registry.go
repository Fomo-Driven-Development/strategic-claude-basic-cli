package templates

import (
	"fmt"
	"sort"
)

const (
	// Default template ID
	DefaultTemplateID = "main"

	// Repository base URL
	DefaultRepoURL = "https://github.com/Fomo-Driven-Development/strategic-claude-base.git"
)

// Registry holds all available templates
var Registry = map[string]Template{
	"main": {
		ID:          "main",
		Name:        "Strategic Claude Basic",
		Description: "Main template for general development projects with comprehensive Claude Code integration",
		RepoURL:     DefaultRepoURL,
		Branch:      "main",
		Commit:      "0c3747dd81c69bad66c828175e358fa840e88227", // Latest commit - Merge refactor with codex removal and streamlined docs
		Language:    "",                                         // Language-agnostic
		Tags:        []string{"general", "default"},
	},
	"ccr": {
		ID:          "ccr",
		Name:        "CCR Template",
		Description: "Specialized template for CCR (Claude Code Router) workflows and development patterns",
		RepoURL:     DefaultRepoURL,
		Branch:      "ccr-template",
		Commit:      "2c9fa88312f7ae68747dd69bbc0075ab47b0225f", // Latest commit - Merge branch 'main' with codex-review features
		Language:    "",
		Tags:        []string{"ccr", "workflow", "specialized"},
	},
	"web-explorer": {
		ID:          "web-explorer",
		Name:        "Claude Web Explorer Template",
		Description: "A template for browser automation projects using Chromium with MCP (Model Context Protocol) integration for Claude Code",
		RepoURL:     DefaultRepoURL,
		Branch:      "web-explorer",
		Commit:      "1a91789daf511b8663e879c9e7e1f36755dfa2d6",
		Language:    "",
		Tags:        []string{"web", "explorer"},
	},
}

// GetTemplate retrieves a template by ID
func GetTemplate(id string) (Template, error) {
	template, exists := Registry[id]
	if !exists {
		return Template{}, fmt.Errorf("template '%s' not found", id)
	}

	if err := template.IsValid(); err != nil {
		return Template{}, fmt.Errorf("template '%s' is invalid: %w", id, err)
	}

	return template, nil
}

// GetDefaultTemplate returns the default template
func GetDefaultTemplate() (Template, error) {
	return GetTemplate(DefaultTemplateID)
}

// ListTemplates returns all available templates, sorted by ID
func ListTemplates() []Template {
	templates := make([]Template, 0, len(Registry))
	for _, template := range Registry {
		templates = append(templates, template)
	}

	// Sort by ID for consistent ordering
	sort.Slice(templates, func(i, j int) bool {
		return templates[i].ID < templates[j].ID
	})

	return templates
}

// ListActiveTemplates returns all non-deprecated templates
func ListActiveTemplates() []Template {
	templates := ListTemplates()
	active := make([]Template, 0, len(templates))

	for _, template := range templates {
		if !template.Deprecated {
			active = append(active, template)
		}
	}

	return active
}

// FilterTemplatesByLanguage returns templates for a specific language
func FilterTemplatesByLanguage(language string) []Template {
	templates := ListActiveTemplates()
	filtered := make([]Template, 0)

	for _, template := range templates {
		if template.Language == "" || template.Language == language {
			filtered = append(filtered, template)
		}
	}

	return filtered
}

// FilterTemplatesByTag returns templates that have a specific tag
func FilterTemplatesByTag(tag string) []Template {
	templates := ListActiveTemplates()
	filtered := make([]Template, 0)

	for _, template := range templates {
		if template.HasTag(tag) {
			filtered = append(filtered, template)
		}
	}

	return filtered
}

// ValidateTemplateID checks if a template ID exists and is valid
func ValidateTemplateID(id string) error {
	_, err := GetTemplate(id)
	return err
}

// GetTemplateIDs returns a list of all template IDs
func GetTemplateIDs() []string {
	ids := make([]string, 0, len(Registry))
	for id := range Registry {
		ids = append(ids, id)
	}

	sort.Strings(ids)
	return ids
}
