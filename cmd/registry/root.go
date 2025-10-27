package registry

import (
	"github.com/charmbracelet/log"
	"github.com/spf13/cobra"
)

var (
	isVerbose bool
)

// RootCMD represents the base command when called without any subcommands
var RootCMD = &cobra.Command{
	Use:   "registry",
	Short: "Управление Artifact Registry для образов контейнеров",
	Long: `Команды для работы с реестрами образов контейнеров Artifact Registry.

Реестры используются для хранения образов контейнеров агентов и MCP серверов.
Эта группа команд позволяет создавать и управлять реестрами в Cloud.ru.

Доступные операции:
• create - Создать новый реестр
• list - Показать список реестров
• get - Получить информацию о реестре
• delete - Удалить реестр

Примеры использования:
  ai-agents-cli registry create --name my-registry
  ai-agents-cli registry list
  ai-agents-cli registry get my-registry`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Debug("Команда registry вызвана без подкоманды")
		// Показываем справку если нет подкоманд
		cmd.Help()
	},
	Args: cobra.ArbitraryArgs,
}

func init() {
	log.Debug("Инициализация команды registry")

	// Добавляем подкоманды
	RootCMD.AddCommand(createCmd)
	RootCMD.AddCommand(listCmd)
	RootCMD.AddCommand(getCmd)
	RootCMD.AddCommand(deleteCmd)
}
