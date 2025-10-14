package mcp_server

import (
	"context"
	"fmt"
	"os"

	"github.com/charmbracelet/log"
	"github.com/cloud-ru/evo-ai-agents-cli/internal/deployer"
	"github.com/cloud-ru/evo-ai-agents-cli/internal/di"
	"github.com/cloud-ru/evo-ai-agents-cli/internal/ui"
	"github.com/cloud-ru/evo-ai-agents-cli/localizations"
	"github.com/spf13/cobra"
)

var (
	deployFile   string
	dryRun       bool
	validateOnly bool
)

// deployCmd represents the deploy command
var deployCmd = &cobra.Command{
	Use:   "deploy [config-file]",
	Short: localizations.Localization.Get("deploy_short"),
	Long:  localizations.Localization.Get("deploy_long"),
	Args:  cobra.MaximumNArgs(1),
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
				"mcp-servers.yaml",
				"mcp-servers.yml",
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
		mcpDeployer := deployer.NewMCPDeployer(apiClient)

		// Валидация конфигурации
		fmt.Println(ui.FormatInfo("Validating configuration..."))
		if err := mcpDeployer.ValidateMCPServers(configFile); err != nil {
			log.Error("Configuration validation failed", "error", err)
			fmt.Println(ui.CheckAndDisplayError(err))
			return
		}
		fmt.Println(ui.FormatSuccess("Configuration is valid"))

		if validateOnly {
			fmt.Println(ui.FormatInfo("Validation completed successfully"))
			return
		}

		// Развертывание
		fmt.Println(ui.FormatInfo("Starting deployment..."))
		results, err := mcpDeployer.DeployMCPServers(ctx, configFile, dryRun)
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

	deployCmd.Flags().StringVarP(&deployFile, "file", "f", "", "Путь к файлу конфигурации")
	deployCmd.Flags().BoolVarP(&dryRun, "dry-run", "d", false, "Режим предварительного просмотра без создания ресурсов")
	deployCmd.Flags().BoolVar(&validateOnly, "validate-only", false, "Только валидация конфигурации без развертывания")
}
