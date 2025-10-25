package auth

import (
	"os"

	"github.com/cloud-ru/evo-ai-agents-cli/internal/errors"
)

// InitCredentials инициализирует учетные данные из сохраненного файла
func InitCredentials() error {
	credentialsManager := NewCredentialsManager()
	
	// Проверяем, есть ли сохраненные учетные данные
	if !credentialsManager.HasCredentials() {
		// Если нет сохраненных учетных данных, проверяем переменные окружения
		keyID := os.Getenv("IAM_KEY_ID")
		secretKey := os.Getenv("IAM_SECRET_KEY")
		endpoint := os.Getenv("IAM_ENDPOINT")
		
		// Если переменные окружения установлены, сохраняем их
		if keyID != "" && secretKey != "" && endpoint != "" {
			creds := &Credentials{
				IAMKeyID:     keyID,
				IAMSecretKey: secretKey,
				IAMEndpoint:  endpoint,
				UserEmail:    os.Getenv("USER_EMAIL"),
			}
			
			if err := credentialsManager.SaveCredentials(creds); err != nil {
				// Не критично, если не удалось сохранить
				return nil
			}
		}
		
		return nil
	}
	
	// Загружаем сохраненные учетные данные
	creds, err := credentialsManager.LoadCredentials()
	if err != nil {
		// Если не удалось загрузить, не критично
		return nil
	}
	
	// Устанавливаем переменные окружения
	os.Setenv("IAM_KEY_ID", creds.IAMKeyID)
	os.Setenv("IAM_SECRET_KEY", creds.IAMSecretKey)
	os.Setenv("IAM_ENDPOINT", creds.IAMEndpoint)
	
	return nil
}

// CheckCredentials проверяет наличие учетных данных
func CheckCredentials() error {
	keyID := os.Getenv("IAM_KEY_ID")
	secretKey := os.Getenv("IAM_SECRET_KEY")
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
