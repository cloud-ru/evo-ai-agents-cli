package deployer

import (
	"context"
	"fmt"

	"github.com/cloud-ru/evo-ai-agents-cli/internal/api"
	"github.com/cloud-ru/evo-ai-agents-cli/internal/parser"
	"github.com/cloud-ru/evo-ai-agents-cli/internal/ui"
	"github.com/cloud-ru/evo-ai-agents-cli/internal/validator"
)

// AgentDeployer обрабатывает развертывание агентов
type AgentDeployer struct {
	api *api.API
}

// NewAgentDeployer создает новый деплойер агентов
func NewAgentDeployer(api *api.API) *AgentDeployer {
	return &AgentDeployer{
		api: api,
	}
}

// AgentConfig представляет конфигурацию агента из YAML
type AgentConfig struct {
	Name        string                 `yaml:"name"`
	Description string                 `yaml:"description"`
	Options     map[string]interface{} `yaml:"options"`
	LLMOptions  map[string]interface{} `yaml:"llm_options"`
	MCPServers  []string               `yaml:"mcp_servers"`
}

// ValidateAgents валидирует конфигурацию агентов
func (d *AgentDeployer) ValidateAgents(configFile string) error {
	// Обрабатываем includes
	processedConfig, err := parser.ProcessYAMLFile(configFile)
	if err != nil {
		return fmt.Errorf("failed to process YAML file with includes: %w", err)
	}

	// Валидируем по схеме
	schemaPath := "schemas/agent.schema.json"
	if err := validator.ValidateConfig(processedConfig, schemaPath); err != nil {
		return fmt.Errorf("configuration validation failed: %w", err)
	}

	return nil
}

// DeployAgents развертывает агентов на основе YAML конфигурации
func (d *AgentDeployer) DeployAgents(ctx context.Context, configFile string, dryRun bool) ([]DeployResult, error) {
	results := []DeployResult{}

	// Обрабатываем includes
	processedConfig, err := parser.ProcessYAMLFile(configFile)
	if err != nil {
		return nil, fmt.Errorf("failed to process YAML file with includes: %w", err)
	}

	// Извлекаем агентов из обработанной конфигурации
	agentsConfig, ok := processedConfig["agents"].([]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid 'agents' section in config file")
	}

	fmt.Println(ui.FormatInfo(fmt.Sprintf("Found %d agents to deploy.", len(agentsConfig))))

	for i, agentConfigRaw := range agentsConfig {
		agentConfigMap, ok := agentConfigRaw.(map[string]interface{})
		if !ok {
			results = append(results, DeployResult{
				Success: false,
				Message: fmt.Sprintf("Invalid agent configuration format for agent %d", i+1),
			})
			continue
		}

		name, _ := agentConfigMap["name"].(string)
		description, _ := agentConfigMap["description"].(string)
		options, _ := agentConfigMap["options"].(map[string]interface{})
		llmOptions, _ := agentConfigMap["llm_options"].(map[string]interface{})
		mcpServers, _ := agentConfigMap["mcp_servers"].([]interface{})

		// Преобразуем mcp_servers в []string
		var mcpServerNames []string
		for _, server := range mcpServers {
			if serverName, ok := server.(string); ok {
				mcpServerNames = append(mcpServerNames, serverName)
			}
		}

		if dryRun {
			results = append(results, DeployResult{
				Success: true,
				Message: fmt.Sprintf("Would deploy agent: %s", name),
			})
			fmt.Println(ui.FormatInfo(fmt.Sprintf("[%d/%d] Dry run for agent: %s", i+1, len(agentsConfig), name)))
			continue
		}

		fmt.Println(ui.FormatInfo(fmt.Sprintf("[%d/%d] Deploying agent: %s", i+1, len(agentsConfig), name)))

		// Создаем запрос для создания агента
		createReq := &api.AgentCreateRequest{
			Name:           name,
			Description:    description,
			Options:        options,
			InstanceTypeID: "58a24a3d-b126-47a5-a39c-30a8aeaa4721", // Используем ID из существующего MCP сервера
		}

		// Добавляем LLM опции в Options если они есть
		if llmOptions != nil {
			if createReq.Options == nil {
				createReq.Options = make(map[string]interface{})
			}
			createReq.Options["llm"] = llmOptions
		}

		// TODO: MCP серверы пока не поддерживаются в API создания агентов
		if len(mcpServerNames) > 0 {
			fmt.Println(ui.FormatWarning(fmt.Sprintf("MCP servers %v specified but not supported yet", mcpServerNames)))
		}

		// Создаем агента
		agent, err := d.api.Agents.Create(ctx, createReq)
		if err != nil {
			results = append(results, DeployResult{
				Success: false,
				Message: fmt.Sprintf("Failed to create agent %s: %v", name, err),
			})
			fmt.Println(ui.FormatError(fmt.Sprintf("[%d/%d] Failed to deploy agent %s: %v", i+1, len(agentsConfig), name, err)))
			continue
		}

		results = append(results, DeployResult{
			Success: true,
			Message: fmt.Sprintf("Successfully deployed agent %s (ID: %s)", name, agent.ID[:8]),
		})
		fmt.Println(ui.FormatSuccess(fmt.Sprintf("[%d/%d] Successfully deployed agent %s (ID: %s)", i+1, len(agentsConfig), name, agent.ID[:8])))
	}

	return results, nil
}

// findMCPServerID находит ID MCP сервера по имени
func (d *AgentDeployer) findMCPServerID(ctx context.Context, serverName string) (string, error) {
	// Получаем список MCP серверов с разумным лимитом
	mcpServers, err := d.api.MCPServers.List(ctx, 100, 0)
	if err != nil {
		return "", fmt.Errorf("failed to list MCP servers: %w", err)
	}

	// Ищем сервер по имени
	for _, server := range mcpServers.Data {
		if server.Name == serverName {
			return server.ID, nil
		}
	}

	return "", fmt.Errorf("MCP server '%s' not found", serverName)
}
