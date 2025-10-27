package api

import (
	"context"
	"fmt"
	"os"
)

// RegistryCreateRequest представляет запрос на создание реестра
type RegistryCreateRequest struct {
	Name         string       `json:"name"`         // Название реестра (обязательно)
	IsPublic     bool         `json:"isPublic"`     // Публичность реестра
	RegistryType RegistryType `json:"registryType"` // Тип реестра (обязательно)
}

// Registry представляет реестр образов контейнеров
type Registry struct {
	ID                       string         `json:"id"`                       // UUID реестра
	Name                     string         `json:"name"`                     // Название реестра
	CreatedAt                string         `json:"createdAt"`                // Время создания (ISO 8601)
	UpdatedAt                string         `json:"updatedAt"`                // Время обновления
	RegistryType             RegistryType   `json:"registryType"`             // Тип реестра (DOCKER/DEBIAN/RPM)
	RetentionPolicyIsEnabled bool           `json:"retentionPolicyIsEnabled"` // Включена ли политика удаления
	RetentionPolicy          string         `json:"retentionPolicy"`          // Настройки политики удаления
	Status                   RegistryStatus `json:"status"`                   // Статус реестра
	IsPublic                 bool           `json:"isPublic"`                 // Публичность
	QuarantineMode           QuarantineMode `json:"quarantineMode"`           // Уровень карантина
}

// RegistryListResponse представляет ответ со списком реестров
type RegistryListResponse struct {
	Registries    []Registry `json:"registries"`    // Список реестров
	NextPageToken string     `json:"nextPageToken"` // Токен для следующей страницы
}

// RegistryType представляет тип реестра
type RegistryType string

const (
	RegistryTypeDocker RegistryType = "DOCKER"
	RegistryTypeDebian RegistryType = "DEBIAN"
	RegistryTypeRPM    RegistryType = "RPM"
)

// RegistryStatus представляет статус реестра
type RegistryStatus string

const (
	RegistryStatusCreating RegistryStatus = "CREATING"
	RegistryStatusActive   RegistryStatus = "ACTIVE"
	RegistryStatusError    RegistryStatus = "ERROR"
)

// QuarantineMode представляет уровень карантина
type QuarantineMode string

const (
	QuarantineModeDisabled QuarantineMode = "DISABLED"
	QuarantineModeLow      QuarantineMode = "LOW"
	QuarantineModeMedium   QuarantineMode = "MEDIUM"
	QuarantineModeHigh     QuarantineMode = "HIGH"
	QuarantineModeCritical QuarantineMode = "CRITICAL"
)

// RegistryService предоставляет методы для работы с реестрами
type RegistryService struct {
	client      *Client
	registryURL string
}

// NewRegistryService создает новый сервис для работы с реестрами
func NewRegistryService(client *Client) *RegistryService {
	// Получаем URL Artifact Registry из переменных окружения
	registryURL := os.Getenv("ARTIFACT_REGISTRY_URL")
	if registryURL == "" {
		registryURL = "https://ar.api.cloud.ru"
	}

	return &RegistryService{
		client:      client,
		registryURL: registryURL,
	}
}

// Create создает новый реестр
// Возвращает операцию создания, так как это длительная операция
func (s *RegistryService) Create(ctx context.Context, req *RegistryCreateRequest) (*Operation, error) {
	path := fmt.Sprintf("/v1/projects/%s/registries", s.client.projectID)

	// Используем специальный URL для Artifact Registry
	var response Operation
	originalURL := s.client.baseURL
	s.client.baseURL = s.registryURL

	err := s.client.Post(ctx, path, req, &response)
	s.client.baseURL = originalURL

	if err != nil {
		return nil, fmt.Errorf("failed to create registry: %w", err)
	}

	return &response, nil
}

// List возвращает список реестров
func (s *RegistryService) List(ctx context.Context, limit, offset int) (*RegistryListResponse, error) {
	path := fmt.Sprintf("/v1/projects/%s/registries", s.client.projectID)

	query := make(map[string]string)
	if limit > 0 {
		query["pageSize"] = fmt.Sprintf("%d", limit)
	}
	if offset > 0 {
		query["pageToken"] = fmt.Sprintf("%d", offset)
	}

	// Используем специальный URL для Artifact Registry
	var response RegistryListResponse
	originalURL := s.client.baseURL
	s.client.baseURL = s.registryURL

	err := s.client.Get(ctx, path, query, &response)
	s.client.baseURL = originalURL

	if err != nil {
		return nil, fmt.Errorf("failed to list registries: %w", err)
	}

	return &response, nil
}

// Get возвращает информацию о реестре по ID
func (s *RegistryService) Get(ctx context.Context, registryID string) (*Registry, error) {
	path := fmt.Sprintf("/v1/projects/%s/registries/%s", s.client.projectID, registryID)

	// Используем специальный URL для Artifact Registry
	var response Registry
	originalURL := s.client.baseURL
	s.client.baseURL = s.registryURL

	err := s.client.Get(ctx, path, nil, &response)
	s.client.baseURL = originalURL

	if err != nil {
		return nil, fmt.Errorf("failed to get registry: %w", err)
	}

	return &response, nil
}

// Delete удаляет реестр по ID
// Возвращает операцию удаления, так как это длительная операция
func (s *RegistryService) Delete(ctx context.Context, registryID string) (*Operation, error) {
	path := fmt.Sprintf("/v1/projects/%s/registries/%s", s.client.projectID, registryID)

	// Используем специальный URL для Artifact Registry
	var response Operation
	originalURL := s.client.baseURL
	s.client.baseURL = s.registryURL

	err := s.client.Delete(ctx, path, &response)
	s.client.baseURL = originalURL

	if err != nil {
		return nil, fmt.Errorf("failed to delete registry: %w", err)
	}

	return &response, nil
}

// PatchQuarantineMode изменяет уровень карантина реестра
func (s *RegistryService) PatchQuarantineMode(ctx context.Context, registryID string, mode QuarantineMode) (*Registry, error) {
	path := fmt.Sprintf("/v1/projects/%s/registry/%s/quarantine", s.client.projectID, registryID)

	req := map[string]interface{}{
		"quarantineMode": mode,
	}

	// Используем специальный URL для Artifact Registry
	var response Registry
	originalURL := s.client.baseURL
	s.client.baseURL = s.registryURL

	err := s.client.Put(ctx, path, req, &response)
	s.client.baseURL = originalURL

	if err != nil {
		return nil, fmt.Errorf("failed to patch registry quarantine mode: %w", err)
	}

	return &response, nil
}

// Operation представляет длительную операцию
type Operation struct {
	ID           string      `json:"id"`           // UUID операции
	ResourceName string      `json:"resourceName"` // Название ресурса
	ResourceID   string      `json:"resourceId"`   // ID ресурса
	CreatedAt    string      `json:"createdAt"`    // Время создания
	UpdatedAt    string      `json:"updatedAt"`    // Время обновления
	Done         bool        `json:"done"`         // Завершена ли операция
	Description  string      `json:"description"`  // Описание операции
	Error        interface{} `json:"error"`        // Ошибка если есть
}
