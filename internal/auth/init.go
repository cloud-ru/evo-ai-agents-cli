package auth

import (
	"os"

	"github.com/cloud-ru/evo-ai-agents-cli/internal/errors"
)

// InitCredentials инициализирует учетные данные из сохраненного файла
func InitCredentials() error {
	credentialsManager := NewCredentialsManager()

	// Если есть сохраненные учетные данные, загружаем их
	if credentialsManager.HasCredentials() {
		creds, err := credentialsManager.LoadCredentials()
		if err != nil {
			// Если не удалось загрузить, не критично
			return nil
		}

		// Устанавливаем переменные окружения
		os.Setenv("IAM_KEY_ID", creds.IAMKeyID)
		os.Setenv("IAM_SECRET", creds.IAMSecretKey) // API клиент ожидает IAM_SECRET
		os.Setenv("IAM_ENDPOINT", creds.IAMEndpoint)

		// Устанавливаем дополнительные переменные если они есть
		if creds.ProjectID != "" {
			os.Setenv("PROJECT_ID", creds.ProjectID)
		}
		if creds.CustomerID != "" {
			os.Setenv("CUSTOMER_ID", creds.CustomerID)
		}
	}

	return nil
}

// CheckCredentials проверяет наличие учетных данных
func CheckCredentials() error {
	keyID := os.Getenv("IAM_KEY_ID")
	secretKey := os.Getenv("IAM_SECRET") // API клиент ожидает IAM_SECRET
	endpoint := os.Getenv("IAM_ENDPOINT")

	if keyID == "" || secretKey == "" || endpoint == "" {
		return errors.New(errors.ErrorTypeAuthentication, errors.SeverityHigh, "MISSING_CREDENTIALS", "Учетные данные не найдены")
	}

	return nil
}

// GetCredentialsPath возвращает путь к файлу с учетными данными
func GetCredentialsPath() string {
	credentialsManager := NewCredentialsManager()
	return credentialsManager.GetCredentialsPath()
}
