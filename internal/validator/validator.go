package validator

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/log"
	"gopkg.in/yaml.v3"
)

// ValidationError представляет ошибку валидации
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
	Value   string `json:"value,omitempty"`
}

// ValidationResult представляет результат валидации
type ValidationResult struct {
	Valid  bool              `json:"valid"`
	Errors []ValidationError `json:"errors,omitempty"`
}

// ConfigValidator валидирует конфигурационные файлы
type ConfigValidator struct {
	schemas map[string]interface{}
}

// NewConfigValidator создает новый валидатор конфигурации
func NewConfigValidator() *ConfigValidator {
	return &ConfigValidator{
		schemas: make(map[string]interface{}),
	}
}

// LoadSchema загружает JSON схему
func (v *ConfigValidator) LoadSchema(name, schemaPath string) error {
	data, err := ioutil.ReadFile(schemaPath)
	if err != nil {
		return fmt.Errorf("failed to read schema file %s: %w", schemaPath, err)
	}

	var schema interface{}
	if err := json.Unmarshal(data, &schema); err != nil {
		return fmt.Errorf("failed to parse schema file %s: %w", schemaPath, err)
	}

	v.schemas[name] = schema
	return nil
}

// ValidateFile валидирует конфигурационный файл
func (v *ConfigValidator) ValidateFile(filePath string) (*ValidationResult, error) {
	// Читаем файл
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return &ValidationResult{
			Valid: false,
			Errors: []ValidationError{{
				Field:   "file",
				Message: fmt.Sprintf("Failed to read file: %v", err),
			}},
		}, nil
	}

	// Определяем тип файла
	ext := strings.ToLower(filepath.Ext(filePath))

	var config interface{}
	switch ext {
	case ".json":
		if err := json.Unmarshal(data, &config); err != nil {
			return &ValidationResult{
				Valid: false,
				Errors: []ValidationError{{
					Field:   "json",
					Message: fmt.Sprintf("Invalid JSON: %v", err),
				}},
			}, nil
		}
	case ".yaml", ".yml":
		// Используем полноценный YAML парсер
		if err := yaml.Unmarshal(data, &config); err != nil {
			return &ValidationResult{
				Valid: false,
				Errors: []ValidationError{{
					Field:   "yaml",
					Message: fmt.Sprintf("Invalid YAML: %v", err),
				}},
			}, nil
		}
	default:
		return &ValidationResult{
			Valid: false,
			Errors: []ValidationError{{
				Field:   "format",
				Message: fmt.Sprintf("Unsupported file format: %s", ext),
			}},
		}, nil
	}

	// Валидируем конфигурацию
	return v.validateConfig(config, filePath)
}

// validateYAMLContent выполняет базовую валидацию YAML контента
func (v *ConfigValidator) validateYAMLContent(data []byte) (*ValidationResult, error) {
	// Простая проверка на наличие обязательных полей
	content := string(data)

	var errors []ValidationError

	// Проверяем наличие основных секций
	if !strings.Contains(content, "mcp-servers") &&
		!strings.Contains(content, "agents") &&
		!strings.Contains(content, "agent-systems") {
		errors = append(errors, ValidationError{
			Field:   "content",
			Message: "Configuration file must contain at least one of: mcp-servers, agents, agent-systems",
		})
	}

	// Проверяем базовую структуру YAML
	if strings.Count(content, ":") == 0 {
		errors = append(errors, ValidationError{
			Field:   "yaml",
			Message: "Invalid YAML structure: no key-value pairs found",
		})
	}

	return &ValidationResult{
		Valid:  len(errors) == 0,
		Errors: errors,
	}, nil
}

// validateConfig валидирует конфигурацию
func (v *ConfigValidator) validateConfig(config interface{}, filePath string) (*ValidationResult, error) {
	var errors []ValidationError

	// Преобразуем в map для удобства
	configMap, ok := config.(map[string]interface{})
	if !ok {
		return &ValidationResult{
			Valid: false,
			Errors: []ValidationError{{
				Field:   "config",
				Message: "Configuration must be a JSON object",
			}},
		}, nil
	}

	// Валидируем MCP серверы
	if mcpServers, exists := configMap["mcp-servers"]; exists {
		if errs := v.validateMCPServers(mcpServers); len(errs) > 0 {
			errors = append(errors, errs...)
		}
	}

	// Валидируем агентов
	if agents, exists := configMap["agents"]; exists {
		if errs := v.validateAgents(agents); len(errs) > 0 {
			errors = append(errors, errs...)
		}
	}

	// Валидируем агентные системы
	if systems, exists := configMap["agent-systems"]; exists {
		if errs := v.validateAgentSystems(systems); len(errs) > 0 {
			errors = append(errors, errs...)
		}
	}

	return &ValidationResult{
		Valid:  len(errors) == 0,
		Errors: errors,
	}, nil
}

// validateMCPServers валидирует секцию MCP серверов
func (v *ConfigValidator) validateMCPServers(mcpServers interface{}) []ValidationError {
	var errors []ValidationError

	servers, ok := mcpServers.([]interface{})
	if !ok {
		return []ValidationError{{
			Field:   "mcp-servers",
			Message: "mcp-servers must be an array",
		}}
	}

	for i, server := range servers {
		serverMap, ok := server.(map[string]interface{})
		if !ok {
			errors = append(errors, ValidationError{
				Field:   fmt.Sprintf("mcp-servers[%d]", i),
				Message: "MCP server must be an object",
			})
			continue
		}

		// Проверяем обязательные поля
		if name, exists := serverMap["name"]; !exists || name == "" {
			errors = append(errors, ValidationError{
				Field:   fmt.Sprintf("mcp-servers[%d].name", i),
				Message: "Name is required",
			})
		} else if nameStr, ok := name.(string); ok {
			// Проверяем формат имени
			if !v.isValidName(nameStr) {
				errors = append(errors, ValidationError{
					Field:   fmt.Sprintf("mcp-servers[%d].name", i),
					Message: "Name must contain only lowercase letters, numbers, and hyphens",
					Value:   nameStr,
				})
			}
		}

		// Проверяем опции
		if options, exists := serverMap["options"]; exists {
			if _, ok := options.(map[string]interface{}); !ok {
				errors = append(errors, ValidationError{
					Field:   fmt.Sprintf("mcp-servers[%d].options", i),
					Message: "Options must be an object",
				})
			}
		}
	}

	return errors
}

// validateAgents валидирует секцию агентов
func (v *ConfigValidator) validateAgents(agents interface{}) []ValidationError {
	var errors []ValidationError

	agentsList, ok := agents.([]interface{})
	if !ok {
		return []ValidationError{{
			Field:   "agents",
			Message: "agents must be an array",
		}}
	}

	for i, agent := range agentsList {
		agentMap, ok := agent.(map[string]interface{})
		if !ok {
			errors = append(errors, ValidationError{
				Field:   fmt.Sprintf("agents[%d]", i),
				Message: "Agent must be an object",
			})
			continue
		}

		// Проверяем обязательные поля
		if name, exists := agentMap["name"]; !exists || name == "" {
			errors = append(errors, ValidationError{
				Field:   fmt.Sprintf("agents[%d].name", i),
				Message: "Name is required",
			})
		} else if nameStr, ok := name.(string); ok {
			if !v.isValidName(nameStr) {
				errors = append(errors, ValidationError{
					Field:   fmt.Sprintf("agents[%d].name", i),
					Message: "Name must contain only lowercase letters, numbers, and hyphens",
					Value:   nameStr,
				})
			}
		}

		// Проверяем LLM опции
		if llmOptions, exists := agentMap["llm_options"]; exists {
			if llmMap, ok := llmOptions.(map[string]interface{}); ok {
				if provider, exists := llmMap["provider"]; !exists || provider == "" {
					errors = append(errors, ValidationError{
						Field:   fmt.Sprintf("agents[%d].llm_options.provider", i),
						Message: "LLM provider is required",
					})
				}
			} else {
				errors = append(errors, ValidationError{
					Field:   fmt.Sprintf("agents[%d].llm_options", i),
					Message: "LLM options must be an object",
				})
			}
		} else {
			errors = append(errors, ValidationError{
				Field:   fmt.Sprintf("agents[%d].llm_options", i),
				Message: "LLM options are required",
			})
		}
	}

	return errors
}

// validateAgentSystems валидирует секцию агентных систем
func (v *ConfigValidator) validateAgentSystems(systems interface{}) []ValidationError {
	var errors []ValidationError

	systemsList, ok := systems.([]interface{})
	if !ok {
		return []ValidationError{{
			Field:   "agent-systems",
			Message: "agent-systems must be an array",
		}}
	}

	for i, system := range systemsList {
		systemMap, ok := system.(map[string]interface{})
		if !ok {
			errors = append(errors, ValidationError{
				Field:   fmt.Sprintf("agent-systems[%d]", i),
				Message: "Agent system must be an object",
			})
			continue
		}

		// Проверяем обязательные поля
		if name, exists := systemMap["name"]; !exists || name == "" {
			errors = append(errors, ValidationError{
				Field:   fmt.Sprintf("agent-systems[%d].name", i),
				Message: "Name is required",
			})
		} else if nameStr, ok := name.(string); ok {
			if !v.isValidName(nameStr) {
				errors = append(errors, ValidationError{
					Field:   fmt.Sprintf("agent-systems[%d].name", i),
					Message: "Name must contain only lowercase letters, numbers, and hyphens",
					Value:   nameStr,
				})
			}
		}

		// Проверяем агентов
		if agents, exists := systemMap["agents"]; exists {
			if agentsList, ok := agents.([]interface{}); ok {
				if len(agentsList) == 0 {
					errors = append(errors, ValidationError{
						Field:   fmt.Sprintf("agent-systems[%d].agents", i),
						Message: "Agent system must have at least one agent",
					})
				}
			} else {
				errors = append(errors, ValidationError{
					Field:   fmt.Sprintf("agent-systems[%d].agents", i),
					Message: "Agents must be an array",
				})
			}
		} else {
			errors = append(errors, ValidationError{
				Field:   fmt.Sprintf("agent-systems[%d].agents", i),
				Message: "Agents are required",
			})
		}
	}

	return errors
}

// isValidName проверяет валидность имени
func (v *ConfigValidator) isValidName(name string) bool {
	if len(name) < 3 || len(name) > 50 {
		return false
	}

	// Проверяем, что имя содержит только строчные буквы, цифры и дефисы
	for _, char := range name {
		if !((char >= 'a' && char <= 'z') ||
			(char >= '0' && char <= '9') ||
			char == '-') {
			return false
		}
	}

	// Имя не должно начинаться или заканчиваться дефисом
	if name[0] == '-' || name[len(name)-1] == '-' {
		return false
	}

	return true
}

// PrintErrors выводит ошибки валидации в красивом формате
func (v *ConfigValidator) PrintErrors(result *ValidationResult) {
	if result.Valid {
		log.Info("✅ Configuration is valid")
		return
	}

	log.Error("❌ Configuration validation failed")

	for _, err := range result.Errors {
		log.Errorf("  %s: %s", err.Field, err.Message)
		if err.Value != "" {
			log.Errorf("    Value: %s", err.Value)
		}
	}
}
