package auth

import (
	"github.com/spf13/cobra"
)

// RootCMD представляет корневую команду для аутентификации
var RootCMD = &cobra.Command{
	Use:   "auth",
	Short: "Управление аутентификацией",
	Long: `Команды для управления аутентификацией в AI Agents CLI.

Доступные команды:
  login    - Войти в систему
  logout   - Выйти из системы
  status   - Проверить статус аутентификации
  config   - Настроить параметры аутентификации`,
}

func init() {
	RootCMD.AddCommand(loginCmd)
	RootCMD.AddCommand(logoutCmd)
	RootCMD.AddCommand(statusCmd)
	RootCMD.AddCommand(configCmd)
}
