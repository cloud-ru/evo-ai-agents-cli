package deployer

import (
	"context"
	"fmt"

	"github.com/charmbracelet/log"
	"github.com/cloud-ru/evo-ai-agents-cli/internal/api"
	"github.com/cloud-ru/evo-ai-agents-cli/internal/parser"
	"github.com/cloud-ru/evo-ai-agents-cli/internal/ui"
)

// MCPServerConfig представляет конфигурацию MCP сервера из YAML
type MCPServerConfig struct {
	Name        string                 `yaml:"name"`
	Description string                 `yaml:"description"`
	Options     map[string]interface{} `yaml:"options"`
}

// MCPDeployer обрабатывает развертывание MCP серверов
type MCPDeployer struct {
	api *api.API
}

// NewMCPDeployer создает новый MCP деплойер
func NewMCPDeployer(api *api.API) *MCPDeployer {
	return &MCPDeployer{
		api: api,
	}
}

// DeployResult представляет результат развертывания
type DeployResult struct {
	Success bool
	Message string
	Error   error
}

// DeployMCPServers развертывает MCP серверы из YAML файла
func (d *MCPDeployer) DeployMCPServers(ctx context.Context, filePath string, dryRun bool) ([]DeployResult, error) {
	log.Info("Starting MCP servers deployment", "file", filePath, "dry_run", dryRun)

	// Парсим YAML файл с includes
	config, err := parser.ProcessYAMLFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to process YAML file: %w", err)
	}

	// Извлекаем MCP серверы
	mcpServers, ok := config["mcp-servers"].([]interface{})
	if !ok {
		return nil, fmt.Errorf("mcp-servers section not found or invalid in YAML file")
	}

	if len(mcpServers) == 0 {
		return []DeployResult{{
			Success: true,
			Message: "No MCP servers to deploy",
		}}, nil
	}

	var results []DeployResult

	// Обрабатываем каждый MCP сервер
	for i, serverData := range mcpServers {
		serverMap, ok := serverData.(map[string]interface{})
		if !ok {
			results = append(results, DeployResult{
				Success: false,
				Message: fmt.Sprintf("Invalid MCP server configuration at index %d", i),
				Error:   fmt.Errorf("expected map, got %T", serverData),
			})
			continue
		}

		// Конвертируем в структуру
		serverConfig := MCPServerConfig{
			Name:        getString(serverMap, "name"),
			Description: getString(serverMap, "description"),
			Options:     getMap(serverMap, "options"),
		}

		if serverConfig.Name == "" {
			results = append(results, DeployResult{
				Success: false,
				Message: fmt.Sprintf("MCP server at index %d missing name", i),
				Error:   fmt.Errorf("name is required"),
			})
			continue
		}

		// Развертываем сервер
		result := d.deployMCPServer(ctx, serverConfig, dryRun)
		results = append(results, result)
	}

	return results, nil
}

// deployMCPServer развертывает один MCP сервер
func (d *MCPDeployer) deployMCPServer(ctx context.Context, config MCPServerConfig, dryRun bool) DeployResult {
	log.Info("Deploying MCP server", "name", config.Name, "dry_run", dryRun)

	if dryRun {
		return DeployResult{
			Success: true,
			Message: fmt.Sprintf("Would deploy MCP server: %s", config.Name),
		}
	}

	// Создаем запрос для API
	createReq := &api.MCPServerCreateRequest{
		Name:        config.Name,
		Description: config.Description,
		Options:     config.Options,
	}

	// Вызываем API
	server, err := d.api.MCPServers.Create(ctx, createReq)
	if err != nil {
		log.Error("Failed to create MCP server", "name", config.Name, "error", err)
		return DeployResult{
			Success: false,
			Message: fmt.Sprintf("Failed to create MCP server: %s", config.Name),
			Error:   err,
		}
	}

	log.Info("MCP server created successfully", "name", config.Name, "id", server.ID)
	return DeployResult{
		Success: true,
		Message: fmt.Sprintf("Successfully created MCP server: %s (ID: %s)", config.Name, server.ID),
	}
}

// ValidateMCPServers валидирует конфигурацию MCP серверов
func (d *MCPDeployer) ValidateMCPServers(filePath string) error {
	log.Info("Validating MCP servers configuration", "file", filePath)

	// Проверяем includes
	if err := parser.ValidateIncludes(filePath); err != nil {
		return fmt.Errorf("include validation failed: %w", err)
	}

	// Парсим YAML файл
	config, err := parser.ProcessYAMLFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to process YAML file: %w", err)
	}

	// Проверяем наличие секции mcp-servers
	mcpServers, ok := config["mcp-servers"].([]interface{})
	if !ok {
		return fmt.Errorf("mcp-servers section not found or invalid in YAML file")
	}

	if len(mcpServers) == 0 {
		return fmt.Errorf("no MCP servers found in configuration")
	}

	// Валидируем каждый сервер
	for i, serverData := range mcpServers {
		serverMap, ok := serverData.(map[string]interface{})
		if !ok {
			return fmt.Errorf("invalid MCP server configuration at index %d: expected map, got %T", i, serverData)
		}

		name := getString(serverMap, "name")
		if name == "" {
			return fmt.Errorf("MCP server at index %d missing required field: name", i)
		}

		description := getString(serverMap, "description")
		if len(description) > 500 {
			return fmt.Errorf("MCP server '%s' description too long (max 500 characters)", name)
		}
	}

	log.Info("MCP servers configuration is valid", "count", len(mcpServers))
	return nil
}

// ShowDeployResults отображает результаты развертывания
func ShowDeployResults(results []DeployResult) {
	successCount := 0
	errorCount := 0

	for _, result := range results {
		if result.Success {
			successCount++
			fmt.Println(ui.FormatSuccess(result.Message))
		} else {
			errorCount++
			fmt.Println(ui.FormatError(result.Message))
			if result.Error != nil {
				fmt.Printf("  Error: %v\n", result.Error)
			}
		}
	}

	// Итоговая статистика
	fmt.Printf("\n📊 Deployment Summary:\n")
	fmt.Printf("  ✅ Successful: %d\n", successCount)
	fmt.Printf("  ❌ Failed: %d\n", errorCount)
	fmt.Printf("  📋 Total: %d\n", len(results))
}

// getString извлекает строку из map с проверкой типа
func getString(m map[string]interface{}, key string) string {
	if val, ok := m[key]; ok {
		if str, ok := val.(string); ok {
			return str
		}
	}
	return ""
}

// getMap извлекает map из map с проверкой типа
func getMap(m map[string]interface{}, key string) map[string]interface{} {
	if val, ok := m[key]; ok {
		if mapVal, ok := val.(map[string]interface{}); ok {
			return mapVal
		}
	}
	return make(map[string]interface{})
}
