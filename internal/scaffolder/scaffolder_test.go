package scaffolder

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestNewScaffolder(t *testing.T) {
	scaffolder := NewScaffolder()
	if scaffolder == nil {
		t.Fatal("NewScaffolder() returned nil")
	}
}

func TestValidateTemplate(t *testing.T) {
	scaffolder := NewScaffolder()

	// Test valid templates
	validTemplates := []string{"mcp", "agent"}
	for _, template := range validTemplates {
		if err := scaffolder.ValidateTemplate(template); err != nil {
			t.Errorf("ValidateTemplate(%s) failed: %v", template, err)
		}
	}

	// Test invalid template
	if err := scaffolder.ValidateTemplate("invalid"); err == nil {
		t.Error("ValidateTemplate('invalid') should have failed")
	}
}

func TestGetAvailableTemplates(t *testing.T) {
	scaffolder := NewScaffolder()
	templates, err := scaffolder.GetAvailableTemplates()
	if err != nil {
		t.Fatalf("GetAvailableTemplates() failed: %v", err)
	}

	if len(templates) == 0 {
		t.Error("GetAvailableTemplates() returned empty list")
	}

	// Check that we have expected templates
	expectedTemplates := []string{"mcp", "agent"}
	for _, expected := range expectedTemplates {
		found := false
		for _, template := range templates {
			if template == expected {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected template '%s' not found in available templates", expected)
		}
	}
}

func TestCreateProject(t *testing.T) {
	scaffolder := NewScaffolder()

	// Create a temporary directory for testing
	tempDir := t.TempDir()
	projectPath := filepath.Join(tempDir, "test-project")

	// Test creating MCP project
	err := scaffolder.CreateProject("mcp", "test-mcp", projectPath, "both")
	if err != nil {
		t.Fatalf("CreateProject('mcp') failed: %v", err)
	}

	// Check that project directory was created
	if _, err := os.Stat(projectPath); os.IsNotExist(err) {
		t.Error("Project directory was not created")
	}

	// Check for some expected files
	expectedFiles := []string{
		"README.md",
		"Makefile",
		"Dockerfile",
		"docker-compose.yml",
		"requirements.txt",
		"pyproject.toml",
		"src/main.py",
	}

	for _, file := range expectedFiles {
		filePath := filepath.Join(projectPath, file)
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			t.Errorf("Expected file %s was not created", file)
		}
	}

	// Test creating Agent project
	agentPath := filepath.Join(tempDir, "test-agent")
	err = scaffolder.CreateProject("agent", "test-agent", agentPath, "github")
	if err != nil {
		t.Fatalf("CreateProject('agent') failed: %v", err)
	}

	// Check that agent project directory was created
	if _, err := os.Stat(agentPath); os.IsNotExist(err) {
		t.Error("Agent project directory was not created")
	}

	// Check for agent-specific files
	agentFiles := []string{
		"src/agent.py",
	}

	for _, file := range agentFiles {
		filePath := filepath.Join(agentPath, file)
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			t.Errorf("Expected agent file %s was not created", file)
		}
	}
}

func TestCreateProjectValidation(t *testing.T) {
	scaffolder := NewScaffolder()
	tempDir := t.TempDir()

	// Test invalid project type
	err := scaffolder.CreateProject("invalid", "test", tempDir, "both")
	if err == nil {
		t.Error("CreateProject with invalid type should have failed")
	}

	// Test empty project name
	err = scaffolder.CreateProject("mcp", "", tempDir, "both")
	if err == nil {
		t.Error("CreateProject with empty name should have failed")
	}

	// Test existing directory
	existingDir := filepath.Join(tempDir, "existing")
	os.MkdirAll(existingDir, 0755)
	err = scaffolder.CreateProject("mcp", "test", existingDir, "both")
	if err == nil {
		t.Error("CreateProject in existing directory should have failed")
	}
}

func TestProcessTemplate(t *testing.T) {
	scaffolder := NewScaffolder()
	data := &ProjectData{
		ProjectName: "test-project",
		ProjectType: "mcp",
		Author:      "Test Author",
		Year:        "2024",
		CICDType:    "both",
	}

	// Test template processing
	templateContent := "Project: {{.ProjectName}}, Type: {{.ProjectType}}, Author: {{.Author}}, Year: {{.Year}}"
	result, err := scaffolder.processTemplate(templateContent, data)
	if err != nil {
		t.Fatalf("processTemplate failed: %v", err)
	}

	expected := "Project: test-project, Type: mcp, Author: Test Author, Year: 2024"
	if result != expected {
		t.Errorf("processTemplate result mismatch. Expected: %s, Got: %s", expected, result)
	}
}

func TestGetGitAuthor(t *testing.T) {
	// Test git author detection
	author := getGitAuthor()

	// If git is available and configured, we should get a non-empty result
	// If git is not available, we should get an empty string
	if author != "" {
		// Should contain name and email in format "Name <email>"
		if !strings.Contains(author, "<") {
			t.Errorf("Git author should contain email in format 'Name <email>', got: %s", author)
		}
	}
}

func TestGetAuthorPriority(t *testing.T) {
	// Test author priority: config > git config > environment > default

	// Test with custom config (highest priority)
	config := &ScaffolderConfig{
		Author: "Custom Author",
	}
	scaffolder := NewScaffolderWithConfig(config)
	author := scaffolder.getAuthor()
	if author != "Custom Author" {
		t.Errorf("Expected 'Custom Author', got: %s", author)
	}

	// Test with empty config (should fall back to git config or default)
	config = &ScaffolderConfig{
		Author: "",
	}
	scaffolder = NewScaffolderWithConfig(config)
	author = scaffolder.getAuthor()
	// Should be either git author or default
	if author == "" {
		t.Error("Author should not be empty")
	}
}
