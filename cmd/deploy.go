package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/charmbracelet/log"
	"github.com/cloud-ru/evo-ai-agents-cli/internal/deployer"
	"github.com/cloud-ru/evo-ai-agents-cli/internal/di"
	"github.com/cloud-ru/evo-ai-agents-cli/internal/ui"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var (
	deployFile         string
	deployDryRun       bool
	deployValidateOnly bool
)

// deployCmd represents the deploy command
var deployCmd = &cobra.Command{
	Use:   "deploy [config-file]",
	Short: "Универсальное развертывание из YAML конфигурации",
	Long: `Универсальное развертывание всех типов ресурсов из YAML конфигурации.

Команда автоматически определяет тип конфигурации и развертывает ресурсы в правильном порядке:
1. MCP серверы
2. Агенты (с привязкой к MCP серверам)
3. Системы агентов (с привязкой к агентам)

Поддерживает:
• Включения других файлов через !include
• Валидации конфигурации по JSON схемам
• Автоматическое разрешение зависимостей
• Режим предварительного просмотра (dry-run)
• Только валидации без развертывания

Примеры использования:
  ai-agents-cli deploy config.yaml
  ai-agents-cli deploy --file config.yaml --dry-run
  ai-agents-cli deploy --validate-only`,
	Args: cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()

		// Определяем файл конфигурации
		configFile := deployFile
		if len(args) > 0 {
			configFile = args[0]
		}

		if configFile == "" {
			// Ищем файл конфигурации по умолчанию
			defaultFiles := []string{
				"ai-agents.yaml",
				"ai-agents.yml",
				"deploy.yaml",
				"deploy.yml",
				"config.yaml",
				"config.yml",
			}

			for _, file := range defaultFiles {
				if _, err := os.Stat(file); err == nil {
					configFile = file
					fmt.Printf("📁 Using configuration file: %s\n", file)
					break
				}
			}

			if configFile == "" {
				fmt.Println("❌ No configuration file found. Looking for:")
				for _, file := range defaultFiles {
					fmt.Printf("   - %s\n", file)
				}
				fmt.Println("\n💡 Create one of these files or specify with: ai-agents-cli deploy <file>")
				os.Exit(1)
			}
		}

		// Получаем API клиент из DI контейнера
		container := di.GetContainer()
		apiClient := container.GetAPI()

		// Создаем деплойеры
		mcpDeployer := deployer.NewMCPDeployer(apiClient)
		agentDeployer := deployer.NewAgentDeployer(apiClient)
		systemDeployer := deployer.NewSystemDeployer(apiClient)

		// Валидация конфигурации
		fmt.Println(ui.FormatInfo("Validating configuration..."))

		// Определяем тип конфигурации и валидируем
		configType, err := detectConfigType(configFile)
		if err != nil {
			log.Error("Failed to detect configuration type", "error", err)
			fmt.Println(ui.CheckAndDisplayError(err))
			return
		}

		// Валидируем в зависимости от типа
		switch configType {
		case "mcp":
			if err := mcpDeployer.ValidateMCPServers(configFile); err != nil {
				log.Error("MCP configuration validation failed", "error", err)
				fmt.Println(ui.CheckAndDisplayError(err))
				return
			}
		case "agent":
			if err := agentDeployer.ValidateAgents(configFile); err != nil {
				log.Error("Agent configuration validation failed", "error", err)
				fmt.Println(ui.CheckAndDisplayError(err))
				return
			}
		case "system":
			if err := systemDeployer.ValidateSystems(configFile); err != nil {
				log.Error("System configuration validation failed", "error", err)
				fmt.Println(ui.CheckAndDisplayError(err))
				return
			}
		case "mixed":
			// Для смешанных конфигураций валидируем все секции
			if err := mcpDeployer.ValidateMCPServers(configFile); err != nil {
				log.Warn("MCP section validation failed", "error", err)
			}
			if err := agentDeployer.ValidateAgents(configFile); err != nil {
				log.Warn("Agent section validation failed", "error", err)
			}
			if err := systemDeployer.ValidateSystems(configFile); err != nil {
				log.Warn("System section validation failed", "error", err)
			}
		}

		fmt.Println(ui.FormatSuccess("Configuration is valid"))

		if deployValidateOnly {
			fmt.Println(ui.FormatInfo("Validation completed successfully"))
			return
		}

		// Развертывание в правильном порядке
		fmt.Println(ui.FormatInfo("Starting deployment..."))

		var allResults []deployer.DeployResult

		// 1. Развертываем MCP серверы
		if configType == "mcp" || configType == "mixed" {
			fmt.Println(ui.FormatInfo("Deploying MCP servers..."))
			mcpResults, err := mcpDeployer.DeployMCPServers(ctx, configFile, deployDryRun)
			if err != nil {
				log.Error("MCP deployment failed", "error", err)
				fmt.Println(ui.FormatError("MCP deployment failed: " + err.Error()))
			} else {
				allResults = append(allResults, mcpResults...)
			}
		}

		// 2. Развертываем агентов
		if configType == "agent" || configType == "mixed" {
			fmt.Println(ui.FormatInfo("Deploying agents..."))
			agentResults, err := agentDeployer.DeployAgents(ctx, configFile, deployDryRun)
			if err != nil {
				log.Error("Agent deployment failed", "error", err)
				fmt.Println(ui.FormatError("Agent deployment failed: " + err.Error()))
			} else {
				allResults = append(allResults, agentResults...)
			}
		}

		// 3. Развертываем системы агентов
		if configType == "system" || configType == "mixed" {
			fmt.Println(ui.FormatInfo("Deploying agent systems..."))
			systemResults, err := systemDeployer.DeploySystems(ctx, configFile, deployDryRun)
			if err != nil {
				log.Error("System deployment failed", "error", err)
				fmt.Println(ui.FormatError("System deployment failed: " + err.Error()))
			} else {
				allResults = append(allResults, systemResults...)
			}
		}

		// Показываем общие результаты
		fmt.Println(ui.FormatInfo("Deployment completed!"))
		deployer.ShowDeployResults(allResults)
	},
}

// detectConfigType определяет тип конфигурации
func detectConfigType(configFile string) (string, error) {
	// Читаем файл
	data, err := os.ReadFile(configFile)
	if err != nil {
		return "", fmt.Errorf("failed to read config file: %w", err)
	}

	// Парсим YAML
	var config map[string]interface{}
	if err := yaml.Unmarshal(data, &config); err != nil {
		return "", fmt.Errorf("failed to parse YAML: %w", err)
	}

	// Определяем тип по наличию секций
	hasMCP := config["mcp-servers"] != nil
	hasAgents := config["agents"] != nil
	hasSystems := config["agent-systems"] != nil

	sectionCount := 0
	if hasMCP {
		sectionCount++
	}
	if hasAgents {
		sectionCount++
	}
	if hasSystems {
		sectionCount++
	}

	if sectionCount == 0 {
		return "", fmt.Errorf("no recognized configuration sections found")
	}

	if sectionCount == 1 {
		if hasMCP {
			return "mcp", nil
		}
		if hasAgents {
			return "agent", nil
		}
		if hasSystems {
			return "system", nil
		}
	}

	return "mixed", nil
}

func init() {
	RootCMD.AddCommand(deployCmd)

	deployCmd.Flags().StringVarP(&deployFile, "file", "f", "", "Путь к файлу конфигурации")
	deployCmd.Flags().BoolVarP(&deployDryRun, "dry-run", "d", false, "Режим предварительного просмотра без создания ресурсов")
	deployCmd.Flags().BoolVar(&deployValidateOnly, "validate-only", false, "Только валидация конфигурации без развертывания")
}
