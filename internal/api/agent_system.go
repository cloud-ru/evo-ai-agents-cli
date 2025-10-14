package api

import (
	"context"
	"fmt"
	"time"
)

// AgentSystem представляет систему агентов
type AgentSystem struct {
	ID                  string                 `json:"id"`
	ProjectID           string                 `json:"projectId,omitempty"`
	Name                string                 `json:"name"`
	Description         string                 `json:"description"`
	Status              string                 `json:"status"`
	StatusReason        StatusReason           `json:"statusReason,omitempty"`
	InstanceType        InstanceType           `json:"instanceType,omitempty"`
	Agents              []AgentSystemAgent     `json:"agents,omitempty"`
	OrchestratorOptions map[string]interface{} `json:"orchestratorOptions,omitempty"`
	Options             map[string]interface{} `json:"options,omitempty"`
	IntegrationOptions  map[string]interface{} `json:"integrationOptions,omitempty"`
	PublicURL           string                 `json:"publicUrl,omitempty"`
	CreatedAt           time.Time              `json:"createdAt"`
	UpdatedAt           time.Time              `json:"updatedAt"`
	CreatedBy           string                 `json:"createdBy,omitempty"`
	UpdatedBy           string                 `json:"updatedBy,omitempty"`
}

// AgentSystemAgent представляет агента в системе
type AgentSystemAgent struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Status string `json:"status,omitempty"`
}

// HistoryEntry представляет запись в истории
type HistoryEntry struct {
	ID        string    `json:"id"`
	Action    string    `json:"action"`
	Status    string    `json:"status"`
	Message   string    `json:"message"`
	CreatedAt time.Time `json:"created_at"`
	Timestamp time.Time `json:"timestamp"`
	Level     string    `json:"level"`
	Source    string    `json:"source"`
}

// AgentSystemCreateRequest представляет запрос на создание системы агентов
type AgentSystemCreateRequest struct {
	Name                string                 `json:"name"`
	Description         string                 `json:"description,omitempty"`
	InstanceTypeID      string                 `json:"instance_type_id,omitempty"`
	OrchestratorOptions map[string]interface{} `json:"orchestratorOptions,omitempty"`
	Options             map[string]interface{} `json:"options,omitempty"`
	IntegrationOptions  map[string]interface{} `json:"integrationOptions,omitempty"`
	Agents              []string               `json:"agents,omitempty"`
}

// AgentSystemUpdateRequest представляет запрос на обновление системы агентов
type AgentSystemUpdateRequest struct {
	Name                string                 `json:"name,omitempty"`
	Description         string                 `json:"description,omitempty"`
	InstanceTypeID      string                 `json:"instance_type_id,omitempty"`
	OrchestratorOptions map[string]interface{} `json:"orchestratorOptions,omitempty"`
	Options             map[string]interface{} `json:"options,omitempty"`
	IntegrationOptions  map[string]interface{} `json:"integrationOptions,omitempty"`
	Agents              []string               `json:"agents,omitempty"`
}

// AgentSystemListResponse представляет ответ со списком систем агентов
type AgentSystemListResponse struct {
	Data  []AgentSystem `json:"data"`
	Total int           `json:"total"`
}

// AgentSystemService предоставляет методы для работы с системами агентов
type AgentSystemService struct {
	client *Client
}

// NewAgentSystemService создает новый сервис для работы с системами агентов
func NewAgentSystemService(client *Client) *AgentSystemService {
	return &AgentSystemService{client: client}
}

// List возвращает список систем агентов
func (s *AgentSystemService) List(ctx context.Context, limit, offset int) (*AgentSystemListResponse, error) {
	query := map[string]string{
		"limit":  fmt.Sprintf("%d", limit),
		"offset": fmt.Sprintf("%d", offset),
	}

	var result AgentSystemListResponse
	err := s.client.Get(ctx, fmt.Sprintf("/api/v1/%s/agentSystems", s.client.projectID), query, &result)
	return &result, err
}

// Get возвращает информацию о конкретной системе агентов
func (s *AgentSystemService) Get(ctx context.Context, systemID string) (*AgentSystem, error) {
	var result AgentSystem
	err := s.client.Get(ctx, fmt.Sprintf("/api/v1/%s/agentSystems/%s", s.client.projectID, systemID), nil, &result)
	return &result, err
}

// Create создает новую систему агентов
func (s *AgentSystemService) Create(ctx context.Context, req *AgentSystemCreateRequest) (*AgentSystem, error) {
	var result AgentSystem
	err := s.client.Post(ctx, fmt.Sprintf("/api/v1/%s/agentSystems", s.client.projectID), req, &result)
	return &result, err
}

// Update обновляет существующую систему агентов
func (s *AgentSystemService) Update(ctx context.Context, systemID string, req *AgentSystemUpdateRequest) (*AgentSystem, error) {
	var result AgentSystem
	err := s.client.Put(ctx, fmt.Sprintf("/api/v1/%s/agentSystems/%s", s.client.projectID, systemID), req, &result)
	return &result, err
}

// Delete удаляет систему агентов
func (s *AgentSystemService) Delete(ctx context.Context, systemID string) error {
	return s.client.Delete(ctx, fmt.Sprintf("/api/v1/%s/agentSystems/%s", s.client.projectID, systemID), nil)
}

// GetHistory возвращает историю системы агентов
func (s *AgentSystemService) GetHistory(ctx context.Context, systemID string, limit, offset int) (*AgentSystemListResponse, error) {
	query := map[string]string{
		"limit":  fmt.Sprintf("%d", limit),
		"offset": fmt.Sprintf("%d", offset),
	}

	var result AgentSystemListResponse
	err := s.client.Get(ctx, fmt.Sprintf("/api/v1/%s/agentSystems/%s/history", s.client.projectID, systemID), query, &result)
	return &result, err
}

// Resume возобновляет работу системы агентов
func (s *AgentSystemService) Resume(ctx context.Context, systemID string) error {
	return s.client.Post(ctx, fmt.Sprintf("/api/v1/%s/agentSystems/%s/resume", s.client.projectID, systemID), nil, nil)
}

// Suspend приостанавливает работу системы агентов
func (s *AgentSystemService) Suspend(ctx context.Context, systemID string) error {
	return s.client.Post(ctx, fmt.Sprintf("/api/v1/%s/agentSystems/%s/suspend", s.client.projectID, systemID), nil, nil)
}
