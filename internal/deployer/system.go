package deployer

import (
	"context"
	"fmt"
	"log"

	"github.com/cloud-ru/evo-ai-agents-cli/internal/api"
	"github.com/cloud-ru/evo-ai-agents-cli/internal/parser"
	"github.com/cloud-ru/evo-ai-agents-cli/internal/ui"
	"github.com/cloud-ru/evo-ai-agents-cli/internal/validator"
)

// SystemDeployer обрабатывает развертывание систем агентов
type SystemDeployer struct {
	api *api.API
}

// NewSystemDeployer создает новый деплойер систем агентов
func NewSystemDeployer(api *api.API) *SystemDeployer {
	return &SystemDeployer{
		api: api,
	}
}

// SystemConfig представляет конфигурацию системы агентов из YAML
type SystemConfig struct {
	Name        string                 `yaml:"name"`
	Description string                 `yaml:"description"`
	Agents      []string               `yaml:"agents"`
	Options     map[string]interface{} `yaml:"options"`
}

// ValidateSystems валидирует конфигурацию систем агентов
func (d *SystemDeployer) ValidateSystems(configFile string) error {
	// Обрабатываем includes
	processedConfig, err := parser.ProcessYAMLFile(configFile)
	if err != nil {
		return fmt.Errorf("failed to process YAML file with includes: %w", err)
	}

	// Валидируем по схеме
	schemaPath := "schemas/systems.schema.json"
	if err := validator.ValidateConfig(processedConfig, schemaPath); err != nil {
		return fmt.Errorf("configuration validation failed: %w", err)
	}

	return nil
}

// DeploySystems развертывает системы агентов на основе YAML конфигурации
func (d *SystemDeployer) DeploySystems(ctx context.Context, configFile string, dryRun bool) ([]DeployResult, error) {
	results := []DeployResult{}

	// Обрабатываем includes
	processedConfig, err := parser.ProcessYAMLFile(configFile)
	if err != nil {
		return nil, fmt.Errorf("failed to process YAML file with includes: %w", err)
	}

	// Извлекаем системы агентов из обработанной конфигурации
	systemsConfig, ok := processedConfig["agent-systems"].([]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid 'agent-systems' section in config file")
	}

	fmt.Println(ui.FormatInfo(fmt.Sprintf("Found %d agent systems to deploy.", len(systemsConfig))))

	for i, systemConfigRaw := range systemsConfig {
		systemConfigMap, ok := systemConfigRaw.(map[string]interface{})
		if !ok {
			results = append(results, DeployResult{
				Success: false,
				Message: fmt.Sprintf("Invalid agent system configuration format for system %d", i+1),
			})
			continue
		}

		name, _ := systemConfigMap["name"].(string)
		description, _ := systemConfigMap["description"].(string)
		options, _ := systemConfigMap["options"].(map[string]interface{})
		agents, _ := systemConfigMap["agents"].([]interface{})

		// Преобразуем agents в []string
		var agentNames []string
		for _, agent := range agents {
			if agentName, ok := agent.(string); ok {
				agentNames = append(agentNames, agentName)
			}
		}

		if dryRun {
			results = append(results, DeployResult{
				Success: true,
				Message: fmt.Sprintf("Would deploy agent system: %s", name),
			})
			fmt.Println(ui.FormatInfo(fmt.Sprintf("[%d/%d] Dry run for agent system: %s", i+1, len(systemsConfig), name)))
			continue
		}

		fmt.Println(ui.FormatInfo(fmt.Sprintf("[%d/%d] Deploying agent system: %s", i+1, len(systemsConfig), name)))

		// Создаем запрос для создания системы агентов
		createReq := &api.AgentSystemCreateRequest{
			Name:        name,
			Description: description,
			Options:     options,
		}

		// Создаем систему агентов
		system, err := d.api.AgentSystems.Create(ctx, createReq)
		if err != nil {
			results = append(results, DeployResult{
				Success: false,
				Message: fmt.Sprintf("Failed to create agent system %s: %v", name, err),
			})
			fmt.Println(ui.FormatError(fmt.Sprintf("[%d/%d] Failed to deploy agent system %s: %v", i+1, len(systemsConfig), name, err)))
			continue
		}

		// Если есть агенты, привязываем их
		if len(agentNames) > 0 {
			err := d.attachAgents(ctx, system.ID, agentNames)
			if err != nil {
				results = append(results, DeployResult{
					Success: false,
					Message: fmt.Sprintf("Agent system %s created but failed to attach agents: %v", name, err),
				})
				fmt.Println(ui.FormatError(fmt.Sprintf("[%d/%d] Agent system %s created but failed to attach agents: %v", i+1, len(systemsConfig), name, err)))
				continue
			}
		}

		results = append(results, DeployResult{
			Success: true,
			Message: fmt.Sprintf("Successfully deployed agent system %s (ID: %s)", name, system.ID[:8]),
		})
		fmt.Println(ui.FormatSuccess(fmt.Sprintf("[%d/%d] Successfully deployed agent system %s (ID: %s)", i+1, len(systemsConfig), name, system.ID[:8])))
	}

	return results, nil
}

// attachAgents привязывает агентов к системе агентов
func (d *SystemDeployer) attachAgents(ctx context.Context, systemID string, agentNames []string) error {
	// Получаем список всех агентов
	agents, err := d.api.Agents.List(ctx, 1000, 0) // Получаем много агентов
	if err != nil {
		return fmt.Errorf("failed to list agents: %w", err)
	}

	// Создаем карту имен агентов к их ID
	agentMap := make(map[string]string)
	for _, agent := range agents.Data {
		agentMap[agent.Name] = agent.ID
	}

	// Находим ID агентов по именам
	var agentIDs []string
	for _, agentName := range agentNames {
		if agentID, exists := agentMap[agentName]; exists {
			agentIDs = append(agentIDs, agentID)
		} else {
			return fmt.Errorf("agent '%s' not found", agentName)
		}
	}

	// Привязываем агентов к системе
	// Здесь нужно будет добавить метод в API для привязки агентов
	// Пока что просто логируем
	log.Printf("Would attach agents %v to system %s", agentIDs, systemID)

	return nil
}
