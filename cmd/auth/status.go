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

		fmt.Println("🔍 Проверка статуса аутентификации...\n")

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
				"Проверьте права доступа к файлу: " + credentialsManager.GetCredentialsPath(),
				"📚 Подробная документация: https://cloud.ru/docs/ai-agents/ug/index?source-platform=Evolution",
			)
			fmt.Println(errorHandler.HandlePlain(appErr))
			os.Exit(1)
		}

		// Показываем информацию о сохраненных учетных данных
		fmt.Println("✅ Учетные данные найдены:")
		fmt.Printf("📧 Email: %s\n", creds.UserEmail)
		fmt.Printf("🔑 Key ID: %s\n", maskString(creds.IAMKeyID))
		fmt.Printf("🌐 Endpoint: %s\n", creds.IAMEndpoint)
		fmt.Printf("⏰ Последний вход: %s\n", creds.LastLogin)
		fmt.Printf("📁 Файл: %s\n\n", credentialsManager.GetCredentialsPath())

		// Проверяем переменные окружения
		fmt.Println("🔍 Проверка переменных окружения:")
		
		keyID := os.Getenv("IAM_KEY_ID")
		secretKey := os.Getenv("IAM_SECRET_KEY")
		endpoint := os.Getenv("IAM_ENDPOINT")

		if keyID != "" {
			fmt.Printf("✅ IAM_KEY_ID: %s\n", maskString(keyID))
		} else {
			fmt.Println("❌ IAM_KEY_ID: не установлена")
		}

		if secretKey != "" {
			fmt.Printf("✅ IAM_SECRET_KEY: %s\n", maskString(secretKey))
		} else {
			fmt.Println("❌ IAM_SECRET_KEY: не установлена")
		}

		if endpoint != "" {
			fmt.Printf("✅ IAM_ENDPOINT: %s\n", endpoint)
		} else {
			fmt.Println("❌ IAM_ENDPOINT: не установлена")
		}

		// Проверяем соответствие переменных окружения сохраненным учетным данным
		fmt.Println("\n🔍 Проверка соответствия:")
		if keyID == creds.IAMKeyID {
			fmt.Println("✅ IAM_KEY_ID соответствует сохраненным данным")
		} else {
			fmt.Println("⚠️  IAM_KEY_ID не соответствует сохраненным данным")
		}

		if secretKey == creds.IAMSecretKey {
			fmt.Println("✅ IAM_SECRET_KEY соответствует сохраненным данным")
		} else {
			fmt.Println("⚠️  IAM_SECRET_KEY не соответствует сохраненным данным")
		}

		if endpoint == creds.IAMEndpoint {
			fmt.Println("✅ IAM_ENDPOINT соответствует сохраненным данным")
		} else {
			fmt.Println("⚠️  IAM_ENDPOINT не соответствует сохраненным данным")
		}

		// Рекомендации
		fmt.Println("\n💡 Рекомендации:")
		if keyID == "" || secretKey == "" || endpoint == "" {
			fmt.Println("🔄 Для установки переменных окружения выполните: ai-agents-cli auth login")
		} else {
			fmt.Println("✅ Все настроено корректно! Можете использовать команды без дополнительной настройки.")
		}
	},
}
