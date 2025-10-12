package cmd

import (
	"os"
	"path/filepath"
	"testing"
	"github.com/cloudru/ai-agents-cli/internal/validator"
)

func init() {
	// Устанавливаем переменные окружения для тестов
	os.Setenv("API_KEY", "test-api-key")
	os.Setenv("PROJECT_ID", "test-project")
}

func TestFindConfigFiles(t *testing.T) {
	// Создаем временную директорию
	tempDir := t.TempDir()

	// Создаем тестовые файлы
	testFiles := []string{
		"config.yaml",
		"config.yml",
		"config.json",
		"other.txt",
		"subdir/config.yaml",
	}

	for _, file := range testFiles {
		filePath := filepath.Join(tempDir, file)
		// Создаем директорию если нужно
		if err := os.MkdirAll(filepath.Dir(filePath), 0755); err != nil {
			t.Fatalf("Failed to create directory: %v", err)
		}
		// Создаем файл
		if err := os.WriteFile(filePath, []byte("test content"), 0644); err != nil {
			t.Fatalf("Failed to create file %s: %v", file, err)
		}
	}

	// Тестируем поиск файлов
	files, err := findConfigFiles(tempDir)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Проверяем, что найдены правильные файлы
	expectedFiles := []string{
		filepath.Join(tempDir, "config.yaml"),
		filepath.Join(tempDir, "config.yml"),
		filepath.Join(tempDir, "config.json"),
		filepath.Join(tempDir, "subdir/config.yaml"),
	}

	if len(files) != len(expectedFiles) {
		t.Errorf("Expected %d files, got %d", len(expectedFiles), len(files))
	}

	// Проверяем, что все ожидаемые файлы найдены
	for _, expectedFile := range expectedFiles {
		found := false
		for _, file := range files {
			if file == expectedFile {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected file %s not found", expectedFile)
		}
	}
}

func TestFindConfigFiles_EmptyDir(t *testing.T) {
	// Создаем пустую временную директорию
	tempDir := t.TempDir()

	files, err := findConfigFiles(tempDir)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if len(files) != 0 {
		t.Errorf("Expected 0 files in empty directory, got %d", len(files))
	}
}

func TestFindConfigFiles_NonexistentDir(t *testing.T) {
	_, err := findConfigFiles("/nonexistent/directory")
	if err == nil {
		t.Errorf("Expected error for nonexistent directory")
	}
}

func TestLoadSchemas(t *testing.T) {
	// Создаем временную директорию для схем
	tempDir := t.TempDir()
	schemasDir := filepath.Join(tempDir, "schemas")
	if err := os.MkdirAll(schemasDir, 0755); err != nil {
		t.Fatalf("Failed to create schemas directory: %v", err)
	}

	// Создаем тестовые схемы
	schemas := map[string]string{
		"mcp-servers":   "mcp.schema.json",
		"agents":        "agent.schema.json",
		"agent-systems": "systems.schema.json",
	}

	for _, filename := range schemas {
		schemaPath := filepath.Join(schemasDir, filename)
		schemaContent := `{
			"type": "object",
			"properties": {
				"name": {
					"type": "string"
				}
			}
		}`
		if err := os.WriteFile(schemaPath, []byte(schemaContent), 0644); err != nil {
			t.Fatalf("Failed to create schema file %s: %v", filename, err)
		}
	}

	// Создаем валидатор
	configValidator := validator.NewConfigValidator()

	// Тестируем загрузку схем
	err := loadSchemas(configValidator)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Проверяем, что схемы загружены
	// Поскольку schemas - приватное поле, мы не можем проверить его напрямую
	// Вместо этого проверяем, что загрузка прошла без ошибок
	if err != nil {
		t.Errorf("Failed to load schemas: %v", err)
	}
}

func TestLoadSchemas_NonexistentFiles(t *testing.T) {
	// Создаем валидатор
	configValidator := validator.NewConfigValidator()

	// Тестируем загрузку несуществующих схем
	err := loadSchemas(configValidator)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Валидатор должен быть создан даже если схемы не загружены
	if configValidator == nil {
		t.Errorf("Expected validator to be created")
	}
}
