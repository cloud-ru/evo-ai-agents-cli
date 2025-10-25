package auth

import (
	"fmt"
	"os"

	"github.com/cloud-ru/evo-ai-agents-cli/internal/auth"
	"github.com/cloud-ru/evo-ai-agents-cli/internal/errors"
	"github.com/spf13/cobra"
)

// statusCmd представляет команду проверки статуса аутентификации
var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Проверить статус аутентификации",
	Long: `Команда для проверки текущего статуса аутентификации.

Показывает информацию о сохраненных учетных данных и переменных окружения.

Примеры использования:
  ai-agents-cli auth status`,
	Run: func(cmd *cobra.Command, args []string) {
		// Создаем обработчик ошибок
		errorHandler := errors.NewHandler()

		// Создаем менеджер учетных данных
		credentialsManager := auth.NewCredentialsManager()

		fmt.Println("🔍 Проверка статуса аутентификации...")

		// Проверяем сохраненные учетные данные
		if !credentialsManager.HasCredentials() {
			fmt.Println("❌ Учетные данные не найдены")
			fmt.Println("💡 Для входа выполните: ai-agents-cli auth login")
			return
		}

		// Загружаем учетные данные
		creds, err := credentialsManager.LoadCredentials()
		if err != nil {
			appErr := errorHandler.WrapFileSystemError(err, "CREDENTIALS_LOAD_FAILED", "Ошибка загрузки учетных данных")
			appErr = appErr.WithSuggestions(
				"Попробуйте перелогиниться: ai-agents-cli auth logout && ai-agents-cli auth login",
				"Проверьте права доступа к файлу: "+credentialsManager.GetCredentialsPath(),
				"📚 Подробная документация: https://cloud.ru/docs/ai-agents/ug/index?source-platform=Evolution",
			)
			fmt.Println(errorHandler.HandlePlain(appErr))
			os.Exit(1)
		}

		// Простая проверка статуса
		fmt.Println("✅ Учетные данные найдены:")
		fmt.Printf("🔑 Key ID: %s\n", maskString(creds.IAMKeyID))
		fmt.Printf("🌐 Endpoint: %s\n", creds.IAMEndpoint)
		fmt.Printf("📋 Project ID: %s\n", creds.ProjectID)
		if creds.CustomerID != "" {
			fmt.Printf("👤 Customer ID: %s\n", creds.CustomerID)
		}
		fmt.Printf("⏰ Последний вход: %s\n", creds.LastLogin)

		// Проверяем, что учетные данные доступны для использования
		fmt.Println("\n✅ Учетные данные готовы к использованию!")
		fmt.Println("💡 CLI автоматически читает конфигурацию из файла ~/.ai-agents-cli/credentials.json")
		fmt.Println("💡 Переменные окружения больше не требуются - все работает из файла конфигурации")
	},
}
