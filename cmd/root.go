package cmd

import (
	"os"

	"github.com/charmbracelet/log"
	authCmd "github.com/cloud-ru/evo-ai-agents-cli/cmd/auth"
	"github.com/cloud-ru/evo-ai-agents-cli/cmd/create"
	"github.com/cloud-ru/evo-ai-agents-cli/internal/auth"
	"github.com/spf13/cobra"
)

var (
	isVerbose bool
)

// RootCMD represents the base command when called without any subcommands
var RootCMD = &cobra.Command{
	Use:   "ai-agents-cli",
	Short: "CLI инструмент для управления AI Agents в облачной платформе Cloud.ru",
	Long: `AI Agents CLI - это мощный инструмент командной строки для управления 
и развертывания AI агентов в облачной платформе Cloud.ru.

Основные возможности:
• Валидация конфигурационных файлов
• Управление MCP серверами
• Создание и настройка агентов
• Управление системами агентов
• Создание проектов из шаблонов
• Интеграция с CI/CD процессами

Для начала работы используйте команду 'validate' для проверки конфигурации
или '--help' для просмотра всех доступных команд.`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		// Получаем значение флага verbose
		verbose, _ := cmd.Flags().GetBool("verbose")

		// Настройка логирования
		logger := log.New(os.Stderr)
		logger.SetReportTimestamp(true)
		logger.SetReportCaller(true)

		// Установка уровня логирования
		if verbose {
			logger.SetLevel(log.DebugLevel)
			logger.Info("Включен подробный режим логирования")
		} else {
			logger.SetLevel(log.InfoLevel)
		}

		log.SetDefault(logger)
		log.Debug("AI Agents CLI запущен", "version", "1.0.0", "verbose", verbose)
		
		// Инициализируем учетные данные
		if err := auth.InitCredentials(); err != nil {
			log.Debug("Не удалось инициализировать учетные данные", "error", err)
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		// Показываем справку если нет аргументов
		if len(args) == 0 {
			log.Debug("Показ справки по командам")
			cmd.Help()
		}
	},
	Args: cobra.ArbitraryArgs,
}

func init() {
	RootCMD.PersistentFlags().
		BoolVarP(&isVerbose, "verbose", "v", false, "Детализация процесса")

	// Add commands
	RootCMD.AddCommand(authCmd.RootCMD)
	RootCMD.AddCommand(create.RootCMD)
}
