package system

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
	systemDeployFile   string
	systemDryRun       bool
	systemValidateOnly bool
)

// deployCmd represents the deploy command
var deployCmd = &cobra.Command{
	Use:   "deploy [config-file]",
	Short: "Развертывание систем агентов из YAML конфигурации",
	Long: `Развертывание систем AI агентов из YAML конфигурации.

Команда позволяет развертывать системы агентов из YAML файла с поддержкой:
• Включения других файлов через !include
• Валидации конфигурации по JSON схеме
• Привязки агентов к системам
• Режима предварительного просмотра (dry-run)
• Только валидации без развертывания

Примеры использования:
  ai-agents-cli system deploy systems.yaml
  ai-agents-cli system deploy --file systems.yaml --dry-run
  ai-agents-cli system deploy --validate-only`,
	Args: cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()

		// Определяем файл конфигурации
		configFile := systemDeployFile
		if len(args) > 0 {
			configFile = args[0]
		}

		if configFile == "" {
			// Ищем файл конфигурации по умолчанию
			defaultFiles := []string{
				"systems.yaml",
				"systems.yml",
				"agent-systems.yaml",
				"agent-systems.yml",
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
		apiClient, err := container.GetAPI()
	if err != nil {
		log.Fatal("Failed to get API client", "error", err)
	}

		// Создаем деплойер
		systemDeployer := deployer.NewSystemDeployer(apiClient)

		// Валидация конфигурации
		fmt.Println(ui.FormatInfo("Validating configuration..."))
		if err := systemDeployer.ValidateSystems(configFile); err != nil {
			log.Error("Configuration validation failed", "error", err)
			fmt.Println(ui.CheckAndDisplayError(err))
			return
		}
		fmt.Println(ui.FormatSuccess("Configuration is valid"))

		if systemValidateOnly {
			fmt.Println(ui.FormatInfo("Validation completed successfully"))
			return
		}

		// Развертывание
		fmt.Println(ui.FormatInfo("Starting deployment..."))
		results, err := systemDeployer.DeploySystems(ctx, configFile, systemDryRun)
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

	deployCmd.Flags().StringVarP(&systemDeployFile, "file", "f", "", "Путь к файлу конфигурации")
	deployCmd.Flags().BoolVarP(&systemDryRun, "dry-run", "d", false, "Режим предварительного просмотра без создания ресурсов")
	deployCmd.Flags().BoolVar(&systemValidateOnly, "validate-only", false, "Только валидация конфигурации без развертывания")
}
