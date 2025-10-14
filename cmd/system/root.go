package system

import (
	"github.com/charmbracelet/log"
	"github.com/spf13/cobra"
)

var (
	isVerbose bool
)

// RootCMD represents the base command when called without any subcommands
var RootCMD = &cobra.Command{
	Use:   "system",
	Short: "Управление системами агентов",
	Long: `Управление системами AI агентов в облачной платформе.

Системы агентов - это группы агентов, которые работают вместе
для выполнения сложных задач. Системы могут координировать
работу нескольких агентов и управлять их взаимодействием.

Доступные операции:
• list - Просмотр списка систем
• get - Получение информации о системе
• create - Создание новой системы
• update - Обновление существующей системы
• delete - Удаление системы
• resume - Возобновление работы системы
• suspend - Приостановка системы

Примеры использования:
  ai-agents-cli system list
  ai-agents-cli system get system-id
  ai-agents-cli system create --name my-system`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Info("Команда систем агентов вызвана без подкоманды")
		// Показываем справку если нет подкоманд
		cmd.Help()
	},
	Args: cobra.ArbitraryArgs,
}

func init() {
	// Добавляем подкоманды
	RootCMD.AddCommand(deployCmd)
	// TODO: Добавить остальные команды по мере реализации
	// RootCMD.AddCommand(listCmd)
	// RootCMD.AddCommand(getCmd)
	// RootCMD.AddCommand(createCmd)
	// RootCMD.AddCommand(updateCmd)
	// RootCMD.AddCommand(deleteCmd)
	// RootCMD.AddCommand(resumeCmd)
	// RootCMD.AddCommand(suspendCmd)
}
