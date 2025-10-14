package validator

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/xeipuuv/gojsonschema"
)

// ConfigValidator представляет валидатор конфигурации
type ConfigValidator struct {
	schemas map[string]map[string]interface{}
}

// NewConfigValidator создает новый валидатор конфигурации
func NewConfigValidator() *ConfigValidator {
	return &ConfigValidator{
		schemas: make(map[string]map[string]interface{}),
	}
}

// LoadSchema загружает схему из файла
func (v *ConfigValidator) LoadSchema(name, filePath string) error {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read schema file %s: %w", filePath, err)
	}

	var schema map[string]interface{}
	if err := json.Unmarshal(data, &schema); err != nil {
		return fmt.Errorf("failed to parse schema %s: %w", filePath, err)
	}

	v.schemas[name] = schema
	return nil
}

// ValidationResult представляет результат валидации
type ValidationResult struct {
	Valid  bool
	Errors []string
}

// ValidateFileWithResult валидирует файл и возвращает результат
func (v *ConfigValidator) ValidateFileWithResult(filePath string) (*ValidationResult, error) {
	// Читаем файл
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file %s: %w", filePath, err)
	}

	// Парсим YAML в JSON
	var config map[string]interface{}
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse configuration file: %w", err)
	}

	// Определяем тип схемы по содержимому файла
	var schemaName string
	if _, hasAgents := config["agents"]; hasAgents {
		schemaName = "agent"
	} else if _, hasMCPServers := config["mcp-servers"]; hasMCPServers {
		schemaName = "mcp"
	} else if _, hasSystems := config["agent-systems"]; hasSystems {
		schemaName = "system"
	} else {
		return nil, fmt.Errorf("unknown configuration type")
	}

	// Получаем схему
	schema, exists := v.schemas[schemaName]
	if !exists {
		return nil, fmt.Errorf("schema %s not found", schemaName)
	}

	// Создаем JSON документ из конфигурации
	configJSON, err := json.Marshal(config)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal config to JSON: %w", err)
	}

	// Создаем загрузчики
	schemaLoader := gojsonschema.NewGoLoader(schema)
	documentLoader := gojsonschema.NewBytesLoader(configJSON)

	// Валидируем
	result, err := gojsonschema.Validate(schemaLoader, documentLoader)
	if err != nil {
		return nil, fmt.Errorf("validation error: %w", err)
	}

	// Преобразуем ошибки
	var errors []string
	if !result.Valid() {
		for _, desc := range result.Errors() {
			errors = append(errors, desc.String())
		}
	}

	return &ValidationResult{
		Valid:  result.Valid(),
		Errors: errors,
	}, nil
}

// ValidateFile валидирует файл (старый API для совместимости)
func (v *ConfigValidator) ValidateFile(filePath string) (*ValidationResult, error) {
	return v.ValidateFileWithResult(filePath)
}

// PrintErrors выводит ошибки валидации
func (v *ConfigValidator) PrintErrors(result *ValidationResult) {
	if result.Valid {
		fmt.Println("✅ Файл валиден")
		return
	}

	fmt.Println("❌ Ошибки валидации:")
	for i, err := range result.Errors {
		fmt.Printf("  %d. %s\n", i+1, err)
	}
}

// ValidateConfig валидирует конфигурацию по JSON схеме
func ValidateConfig(config map[string]interface{}, schemaPath string) error {
	// Читаем схему
	schemaData, err := ioutil.ReadFile(schemaPath)
	if err != nil {
		return fmt.Errorf("failed to read schema file %s: %w", schemaPath, err)
	}

	// Парсим схему
	var schema map[string]interface{}
	if err := json.Unmarshal(schemaData, &schema); err != nil {
		return fmt.Errorf("failed to parse schema: %w", err)
	}

	// Создаем JSON документ из конфигурации
	configJSON, err := json.Marshal(config)
	if err != nil {
		return fmt.Errorf("failed to marshal config to JSON: %w", err)
	}

	// Создаем загрузчик схемы
	schemaLoader := gojsonschema.NewGoLoader(schema)
	documentLoader := gojsonschema.NewBytesLoader(configJSON)

	// Валидируем
	result, err := gojsonschema.Validate(schemaLoader, documentLoader)
	if err != nil {
		return fmt.Errorf("validation error: %w", err)
	}

	if !result.Valid() {
		var errors string
		for _, desc := range result.Errors() {
			errors += fmt.Sprintf("- %s\n", desc)
		}
		return fmt.Errorf("validation failed:\n%s", errors)
	}

	return nil
}

// ValidateFile валидирует файл конфигурации
func ValidateFile(filePath string, schemaPath string) error {
	// Проверяем существование файла
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return fmt.Errorf("configuration file %s does not exist", filePath)
	}

	// Читаем файл
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read file %s: %w", filePath, err)
	}

	// Парсим YAML в JSON
	var config map[string]interface{}
	if err := json.Unmarshal(data, &config); err != nil {
		return fmt.Errorf("failed to parse configuration file: %w", err)
	}

	// Валидируем
	return ValidateConfig(config, schemaPath)
}
