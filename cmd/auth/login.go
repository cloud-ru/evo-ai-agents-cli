package auth

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/huh"
	"github.com/cloud-ru/evo-ai-agents-cli/internal/auth"
	"github.com/cloud-ru/evo-ai-agents-cli/internal/errors"
	"github.com/spf13/cobra"
)

var (
	devMode bool
)

// loginCmd представляет команду входа в систему
var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Войти в систему AI Agents",
	Long: `Команда для входа в систему AI Agents.

Команда запросит у вас учетные данные и сохранит их для последующего использования.
После успешного входа вам не нужно будет каждый раз указывать переменные окружения.

Примеры использования:
  ai-agents-cli auth login
  ai-agents-cli auth login --dev`,
	Run: func(cmd *cobra.Command, args []string) {
		// Создаем обработчик ошибок
		errorHandler := errors.NewHandler()

		// Создаем менеджер учетных данных
		credentialsManager := auth.NewCredentialsManager()

		// Простая форма входа
		var loginData struct {
			IAMKeyID     string
			IAMSecretKey string
			IAMEndpoint  string
		}

		// Устанавливаем значения по умолчанию только в dev режиме
		if devMode {
			loginData.IAMEndpoint = "https://iam.api.cloud.ru"
		}

		// Простая форма без лишних вопросов
		form := huh.NewForm(
			huh.NewGroup(
				huh.NewInput().
					Title("🔑 IAM Key ID").
					Description("Введите ваш IAM Key ID").
					Value(&loginData.IAMKeyID).
					Validate(func(str string) error {
						if str == "" {
							return errors.ValidationError("MISSING_KEY_ID", "IAM Key ID обязателен")
						}
						return nil
					}),

				huh.NewInput().
					Title("🔐 IAM Secret Key").
					Description("Введите ваш IAM Secret Key").
					Value(&loginData.IAMSecretKey).
					Password(true).
					Validate(func(str string) error {
						if str == "" {
							return errors.ValidationError("MISSING_SECRET_KEY", "IAM Secret Key обязателен")
						}
						return nil
					}),

				huh.NewInput().
					Title("🌐 IAM Endpoint").
					Description("Введите IAM Endpoint URL").
					Value(&loginData.IAMEndpoint).
					Validate(func(str string) error {
						if str == "" {
							return errors.ValidationError("MISSING_ENDPOINT", "IAM Endpoint обязателен")
						}
						if !strings.HasPrefix(str, "http") {
							return errors.ValidationError("INVALID_ENDPOINT", "Endpoint должен начинаться с http:// или https://")
						}
						return nil
					}),
			),
		).WithTheme(huh.ThemeCharm()).
			WithWidth(120).
			WithHeight(40)

		if err := form.Run(); err != nil {
			appErr := errorHandler.WrapUserError(err, "FORM_ERROR", "Ошибка заполнения формы входа")
			fmt.Println(errorHandler.HandlePlain(appErr))
			os.Exit(1)
		}

		// Создаем объект учетных данных
		creds := &auth.Credentials{
			IAMKeyID:     loginData.IAMKeyID,
			IAMSecretKey: loginData.IAMSecretKey,
			IAMEndpoint:  loginData.IAMEndpoint,
			LastLogin:    time.Now().Format("2006-01-02 15:04:05"),
		}

		// Сохраняем учетные данные
		if err := credentialsManager.SaveCredentials(creds); err != nil {
			appErr := errorHandler.WrapFileSystemError(err, "CREDENTIALS_SAVE_FAILED", "Ошибка сохранения учетных данных")
			appErr = appErr.WithSuggestions(
				"Проверьте права доступа к домашней директории",
				"Убедитесь что у вас есть права на создание файлов",
				"Попробуйте запустить команду от имени администратора",
				"📚 Подробная документация: https://cloud.ru/docs/ai-agents/ug/index?source-platform=Evolution",
			)
			fmt.Println(errorHandler.HandlePlain(appErr))
			os.Exit(1)
		}

		// Устанавливаем переменные окружения
		if err := credentialsManager.SetEnvironmentVariables(); err != nil {
			appErr := errorHandler.WrapConfigurationError(err, "ENV_SET_FAILED", "Ошибка установки переменных окружения")
			fmt.Println(errorHandler.HandlePlain(appErr))
			os.Exit(1)
		}

		// Простое успешное сообщение
		fmt.Printf("✅ Успешный вход в систему!\n")
		fmt.Printf("🔑 Key ID: %s\n", maskString(loginData.IAMKeyID))
		fmt.Printf("🌐 Endpoint: %s\n", loginData.IAMEndpoint)
		fmt.Printf("💡 Теперь вы можете использовать команды без указания переменных окружения!\n")
	},
}

func init() {
	loginCmd.Flags().BoolVar(&devMode, "dev", false, "Режим разработки (устанавливает endpoint по умолчанию)")
}

// maskString маскирует строку для безопасности
func maskString(s string) string {
	if len(s) <= 8 {
		return strings.Repeat("*", len(s))
	}
	return s[:4] + strings.Repeat("*", len(s)-8) + s[len(s)-4:]
}
