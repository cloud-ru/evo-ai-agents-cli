package agent

import (
	"context"
	"fmt"
	"os"

	"github.com/charmbracelet/log"
	"github.com/cloud-ru/evo-ai-agents-cli/internal/deployer"
	"github.com/cloud-ru/evo-ai-agents-cli/internal/di"
	"github.com/cloud-ru/evo-ai-agents-cli/internal/ui"
	"github.com/spf13/cobra"
)

var (
	agentDeployFile   string
	agentDryRun       bool
	agentValidateOnly bool
)

// deployCmd represents the deploy command
var deployCmd = &cobra.Command{
	Use:   "deploy [config-file]",
	Short: "Развертывание агентов из YAML конфигурации",
	Long: `Развертывание AI агентов из YAML конфигурации.

Команда позволяет развертывать агентов из YAML файла с поддержкой:
• Включения других файлов через !include
• Валидации конфигурации по JSON схеме
• Привязки MCP серверов к агентам
• Режима предварительного просмотра (dry-run)
• Только валидации без развертывания

Примеры использования:
  ai-agents-cli agents deploy agents.yaml
  ai-agents-cli agents deploy --file agents.yaml --dry-run
  ai-agents-cli agents deploy --validate-only`,
	Args: cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()

		// Определяем файл конфигурации
		configFile := agentDeployFile
		if len(args) > 0 {
			configFile = args[0]
		}

		if configFile == "" {
			// Ищем файл конфигурации по умолчанию
			defaultFiles := []string{
				"agents.yaml",
				"agents.yml",
			}

			for _, file := range defaultFiles {
				if _, err := os.Stat(file); err == nil {
					configFile = file
					break
				}
			}

			if configFile == "" {
				log.Fatal("No configuration file specified and no default file found")
			}
		}

		// Получаем API клиент из DI контейнера
		container := di.GetContainer()
		apiClient := container.GetAPI()

		// Создаем деплойер
		agentDeployer := deployer.NewAgentDeployer(apiClient)

		// Валидация конфигурации
		fmt.Println(ui.FormatInfo("Validating configuration..."))
		if err := agentDeployer.ValidateAgents(configFile); err != nil {
			log.Error("Configuration validation failed", "error", err)
			fmt.Println(ui.CheckAndDisplayError(err))
			return
		}
		fmt.Println(ui.FormatSuccess("Configuration is valid"))

		if agentValidateOnly {
			fmt.Println(ui.FormatInfo("Validation completed successfully"))
			return
		}

		// Развертывание
		fmt.Println(ui.FormatInfo("Starting deployment..."))
		results, err := agentDeployer.DeployAgents(ctx, configFile, agentDryRun)
		if err != nil {
			log.Error("Deployment failed", "error", err)
			fmt.Println(ui.CheckAndDisplayError(err))
			return
		}

		// Показываем результаты
		deployer.ShowDeployResults(results)
	},
}

func init() {
	RootCMD.AddCommand(deployCmd)

	deployCmd.Flags().StringVarP(&agentDeployFile, "file", "f", "", "Путь к файлу конфигурации")
	deployCmd.Flags().BoolVarP(&agentDryRun, "dry-run", "d", false, "Режим предварительного просмотра без создания ресурсов")
	deployCmd.Flags().BoolVar(&agentValidateOnly, "validate-only", false, "Только валидация конфигурации без развертывания")
}
