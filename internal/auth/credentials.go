package auth

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/cloud-ru/evo-ai-agents-cli/internal/errors"
)

// Credentials представляет учетные данные пользователя
type Credentials struct {
	IAMKeyID     string `json:"iam_key_id"`
	IAMSecretKey string `json:"iam_secret_key"`
	IAMEndpoint  string `json:"iam_endpoint"`
	ProjectID    string `json:"project_id,omitempty"`
	CustomerID   string `json:"customer_id,omitempty"`
	UserEmail    string `json:"user_email,omitempty"`
	LastLogin    string `json:"last_login,omitempty"`
}

// CredentialsManager управляет сохранением и загрузкой учетных данных
type CredentialsManager struct {
	credentialsPath string
}

// NewCredentialsManager создает новый менеджер учетных данных
func NewCredentialsManager() *CredentialsManager {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		// Fallback на текущую директорию
		homeDir = "."
	}

	credentialsPath := filepath.Join(homeDir, ".ai-agents-cli", "credentials.json")

	return &CredentialsManager{
		credentialsPath: credentialsPath,
	}
}

// SaveCredentials сохраняет учетные данные в файл
func (cm *CredentialsManager) SaveCredentials(creds *Credentials) error {
	// Создаем директорию если не существует
	dir := filepath.Dir(cm.credentialsPath)
	if err := os.MkdirAll(dir, 0700); err != nil {
		return errors.Wrap(err, errors.ErrorTypeFileSystem, errors.SeverityMedium, "DIRECTORY_CREATION_FAILED", "Ошибка создания директории для учетных данных")
	}

	// Сохраняем учетные данные
	data, err := json.MarshalIndent(creds, "", "  ")
	if err != nil {
		return errors.Wrap(err, errors.ErrorTypeFileSystem, errors.SeverityMedium, "CREDENTIALS_ENCODE_FAILED", "Ошибка кодирования учетных данных")
	}

	// Записываем с правами только для владельца
	if err := os.WriteFile(cm.credentialsPath, data, 0600); err != nil {
		return errors.Wrap(err, errors.ErrorTypeFileSystem, errors.SeverityMedium, "CREDENTIALS_SAVE_FAILED", "Ошибка сохранения учетных данных")
	}

	return nil
}

// LoadCredentials загружает учетные данные из файла
func (cm *CredentialsManager) LoadCredentials() (*Credentials, error) {
	// Проверяем существование файла
	if _, err := os.Stat(cm.credentialsPath); os.IsNotExist(err) {
		return nil, errors.New(errors.ErrorTypeAuthentication, errors.SeverityMedium, "CREDENTIALS_NOT_FOUND", "Учетные данные не найдены")
	}

	// Читаем файл
	data, err := os.ReadFile(cm.credentialsPath)
	if err != nil {
		return nil, errors.Wrap(err, errors.ErrorTypeFileSystem, errors.SeverityMedium, "CREDENTIALS_READ_FAILED", "Ошибка чтения учетных данных")
	}

	// Декодируем JSON
	var creds Credentials
	if err := json.Unmarshal(data, &creds); err != nil {
		return nil, errors.Wrap(err, errors.ErrorTypeFileSystem, errors.SeverityMedium, "CREDENTIALS_DECODE_FAILED", "Ошибка декодирования учетных данных")
	}

	return &creds, nil
}

// DeleteCredentials удаляет сохраненные учетные данные
func (cm *CredentialsManager) DeleteCredentials() error {
	if _, err := os.Stat(cm.credentialsPath); os.IsNotExist(err) {
		return errors.New(errors.ErrorTypeAuthentication, errors.SeverityLow, "CREDENTIALS_NOT_FOUND", "Учетные данные не найдены")
	}

	if err := os.Remove(cm.credentialsPath); err != nil {
		return errors.Wrap(err, errors.ErrorTypeFileSystem, errors.SeverityMedium, "CREDENTIALS_DELETE_FAILED", "Ошибка удаления учетных данных")
	}

	return nil
}

// HasCredentials проверяет наличие сохраненных учетных данных
func (cm *CredentialsManager) HasCredentials() bool {
	_, err := os.Stat(cm.credentialsPath)
	return !os.IsNotExist(err)
}

// SetEnvironmentVariables устанавливает переменные окружения из сохраненных учетных данных
func (cm *CredentialsManager) SetEnvironmentVariables() error {
	creds, err := cm.LoadCredentials()
	if err != nil {
		return err
	}

	// Устанавливаем переменные окружения
	os.Setenv("IAM_KEY_ID", creds.IAMKeyID)
	os.Setenv("IAM_SECRET", creds.IAMSecretKey) // API клиент ожидает IAM_SECRET, а не IAM_SECRET_KEY
	os.Setenv("IAM_ENDPOINT", creds.IAMEndpoint)

	// Устанавливаем дополнительные переменные если они есть
	if creds.ProjectID != "" {
		os.Setenv("PROJECT_ID", creds.ProjectID)
	}
	if creds.CustomerID != "" {
		os.Setenv("CUSTOMER_ID", creds.CustomerID)
	}

	return nil
}

// GetCredentialsPath возвращает путь к файлу с учетными данными
func (cm *CredentialsManager) GetCredentialsPath() string {
	return cm.credentialsPath
}
