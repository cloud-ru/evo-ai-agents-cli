package auth

import (
	"fmt"
	"os"

	"github.com/cloud-ru/evo-ai-agents-cli/internal/auth"
	"github.com/cloud-ru/evo-ai-agents-cli/internal/errors"
	"github.com/spf13/cobra"
)

// configCmd представляет команду настройки аутентификации
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Настроить параметры аутентификации",
	Long: `Команда для настройки параметров аутентификации.

Позволяет изменить настройки без полного перелогина.

Примеры использования:
  ai-agents-cli auth config
  ai-agents-cli auth config --endpoint https://api.cloud.ru`,
	Run: func(cmd *cobra.Command, args []string) {
		// Создаем обработчик ошибок
		errorHandler := errors.NewHandler()

		// Создаем менеджер учетных данных
		credentialsManager := auth.NewCredentialsManager()

		// Проверяем, есть ли сохраненные учетные данные
		if !credentialsManager.HasCredentials() {
			fmt.Println("❌ Учетные данные не найдены")
			fmt.Println("💡 Для входа выполните: ai-agents-cli auth login")
			return
		}

		// Загружаем существующие учетные данные
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

		// Показываем текущие настройки
		fmt.Println("🔧 Текущие настройки аутентификации:")
		fmt.Printf("🔑 Key ID: %s\n", maskString(creds.IAMKeyID))
		fmt.Printf("🌐 Endpoint: %s\n", creds.IAMEndpoint)
		fmt.Printf("⏰ Последний вход: %s\n", creds.LastLogin)
		fmt.Printf("📁 Файл: %s\n\n", credentialsManager.GetCredentialsPath())

		// Показываем переменные окружения
		fmt.Println("🔍 Текущие переменные окружения:")
		keyID := os.Getenv("IAM_KEY_ID")
		secretKey := os.Getenv("IAM_SECRET") // API клиент использует IAM_SECRET
		endpoint := os.Getenv("IAM_ENDPOINT")
		projectID := os.Getenv("PROJECT_ID")
		customerID := os.Getenv("CUSTOMER_ID")

		if keyID != "" {
			fmt.Printf("✅ IAM_KEY_ID: %s\n", maskString(keyID))
		} else {
			fmt.Println("❌ IAM_KEY_ID: не установлена")
		}

		if secretKey != "" {
			fmt.Printf("✅ IAM_SECRET: %s\n", maskString(secretKey))
		} else {
			fmt.Println("❌ IAM_SECRET: не установлена")
		}

		if endpoint != "" {
			fmt.Printf("✅ IAM_ENDPOINT: %s\n", endpoint)
		} else {
			fmt.Println("❌ IAM_ENDPOINT: не установлена")
		}

		if projectID != "" {
			fmt.Printf("✅ PROJECT_ID: %s\n", projectID)
		} else {
			fmt.Println("❌ PROJECT_ID: не установлена")
		}

		if customerID != "" {
			fmt.Printf("✅ CUSTOMER_ID: %s\n", customerID)
		} else {
			fmt.Println("❌ CUSTOMER_ID: не установлена")
		}

		// Рекомендации
		fmt.Println("\n💡 Доступные действия:")
		fmt.Println("🔄 Перелогиниться: ai-agents-cli auth login")
		fmt.Println("🚪 Выйти из системы: ai-agents-cli auth logout")
		fmt.Println("📊 Проверить статус: ai-agents-cli auth status")
		fmt.Println("📚 Документация: https://cloud.ru/docs/ai-agents/ug/index?source-platform=Evolution")
	},
}
