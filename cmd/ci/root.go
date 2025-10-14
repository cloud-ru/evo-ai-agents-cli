package ci

import (
	"github.com/charmbracelet/log"
	"github.com/cloudru/ai-agents-cli/internal/di"
	"github.com/spf13/cobra"
)

var (
	isVerbose bool
)

// RootCMD represents the base command when called without any subcommands
var RootCMD = &cobra.Command{
	Use:   "ci",
	Short: "CI/CD функции",
	Long: `Команды для интеграции с CI/CD процессами.

Эта группа команд предназначена для автоматизации развертывания
и управления AI агентами в рамках CI/CD пайплайнов.

Доступные операции:
• deploy - Развертывание конфигураций
• validate - Валидация конфигураций в CI
• status - Проверка статуса развертывания
• rollback - Откат к предыдущей версии

Примеры использования:
  ai-agents-cli ci deploy --config examples/
  ai-agents-cli ci validate --config examples/
  ai-agents-cli ci status --deployment-id 123`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Debug("Команда CI/CD вызвана без подкоманды")
		// Показываем справку если нет подкоманд
		cmd.Help()
	},
	Args: cobra.ArbitraryArgs,
}

func init() {
	log.Debug("Инициализация команды CI/CD")

	// Инициализируем DI контейнер
	container := di.GetContainer()

	// Получаем API клиент из контейнера (для инициализации)
	_ = container.GetAPI()

	log.Debug("Команда CI/CD инициализирована успешно")

	// Добавляем подкоманды
	RootCMD.AddCommand(statusCmd)
	RootCMD.AddCommand(logsCmd)
}
