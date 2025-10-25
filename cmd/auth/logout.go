package auth

import (
	"fmt"
	"os"

	"github.com/cloud-ru/evo-ai-agents-cli/internal/auth"
	"github.com/cloud-ru/evo-ai-agents-cli/internal/errors"
	"github.com/spf13/cobra"
)

// logoutCmd представляет команду выхода из системы
var logoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "Выйти из системы AI Agents",
	Long: `Команда для выхода из системы AI Agents.

Удаляет сохраненные учетные данные и очищает переменные окружения.
После выхода вам потребуется снова выполнить 'ai-agents-cli auth login' для входа.

Примеры использования:
  ai-agents-cli auth logout`,
	Run: func(cmd *cobra.Command, args []string) {
		// Создаем обработчик ошибок
		errorHandler := errors.NewHandler()

		// Создаем менеджер учетных данных
		credentialsManager := auth.NewCredentialsManager()

		// Проверяем, есть ли сохраненные учетные данные
		if !credentialsManager.HasCredentials() {
			fmt.Println("ℹ️  Учетные данные не найдены. Вы уже вышли из системы.")
			return
		}

		// Загружаем учетные данные для отображения информации
		creds, err := credentialsManager.LoadCredentials()
		if err != nil {
			appErr := errorHandler.WrapFileSystemError(err, "CREDENTIALS_LOAD_FAILED", "Ошибка загрузки учетных данных")
			appErr = appErr.WithSuggestions(
				"Попробуйте удалить файл вручную: " + credentialsManager.GetCredentialsPath(),
				"Проверьте права доступа к файлу",
				"📚 Подробная документация: https://cloud.ru/docs/ai-agents/ug/index?source-platform=Evolution",
			)
			fmt.Println(errorHandler.HandlePlain(appErr))
			os.Exit(1)
		}

		// Показываем информацию о том, что будет удалено
		fmt.Printf("🔐 Удаляем учетные данные для:\n")
		fmt.Printf("📧 Email: %s\n", creds.UserEmail)
		fmt.Printf("🔑 Key ID: %s\n", maskString(creds.IAMKeyID))
		fmt.Printf("🌐 Endpoint: %s\n", creds.IAMEndpoint)
		fmt.Printf("⏰ Последний вход: %s\n\n", creds.LastLogin)

		// Удаляем учетные данные
		if err := credentialsManager.DeleteCredentials(); err != nil {
			appErr := errorHandler.WrapFileSystemError(err, "CREDENTIALS_DELETE_FAILED", "Ошибка удаления учетных данных")
			appErr = appErr.WithSuggestions(
				"Попробуйте удалить файл вручную: " + credentialsManager.GetCredentialsPath(),
				"Проверьте права доступа к файлу",
				"Убедитесь что файл не заблокирован другим процессом",
				"📚 Подробная документация: https://cloud.ru/docs/ai-agents/ug/index?source-platform=Evolution",
			)
			fmt.Println(errorHandler.HandlePlain(appErr))
			os.Exit(1)
		}

		// Очищаем переменные окружения
		os.Unsetenv("IAM_KEY_ID")
		os.Unsetenv("IAM_SECRET_KEY")
		os.Unsetenv("IAM_ENDPOINT")

		// Показываем успешное сообщение
		fmt.Println("✅ Успешный выход из системы!")
		fmt.Println("💡 Для повторного входа выполните: ai-agents-cli auth login")
	},
}
