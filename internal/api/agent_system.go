package api

import (
	"context"
	"fmt"
	"time"
)

// AgentSystem представляет агентную систему
type AgentSystem struct {
	ID          string                 `json:"id"`
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Status      string                 `json:"status"`
	CreatedAt   time.Time              `json:"created_at"`
	UpdatedAt   time.Time              `json:"updated_at"`
	Options     map[string]interface{} `json:"options"`
	Agents      []string               `json:"agents,omitempty"`
}

// AgentSystemCreateRequest представляет запрос на создание агентной системы
type AgentSystemCreateRequest struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description,omitempty"`
	Options     map[string]interface{} `json:"options"`
	Agents      []string               `json:"agents,omitempty"`
}

// AgentSystemUpdateRequest представляет запрос на обновление агентной системы
type AgentSystemUpdateRequest struct {
	Name        string                 `json:"name,omitempty"`
	Description string                 `json:"description,omitempty"`
	Options     map[string]interface{} `json:"options,omitempty"`
	Agents      []string               `json:"agents,omitempty"`
}

// AgentSystemListResponse представляет ответ со списком агентных систем
type AgentSystemListResponse struct {
	Data  []AgentSystem `json:"data"`
	Total int           `json:"total"`
}

// AgentSystemHistoryResponse представляет ответ с историей агентной системы
type AgentSystemHistoryResponse struct {
	Data []HistoryEntry `json:"data"`
}

// AgentSystemDiagram представляет диаграмму агентной системы
type AgentSystemDiagram struct {
	Nodes []DiagramNode `json:"nodes"`
	Edges []DiagramEdge `json:"edges"`
}

// DiagramNode представляет узел диаграммы
type DiagramNode struct {
	ID    string `json:"id"`
	Type  string `json:"type"`
	Label string `json:"label"`
	X     int    `json:"x"`
	Y     int    `json:"y"`
}

// DiagramEdge представляет связь в диаграмме
type DiagramEdge struct {
	From string `json:"from"`
	To   string `json:"to"`
	Type string `json:"type"`
}

// AgentSystemService предоставляет методы для работы с агентными системами
type AgentSystemService struct {
	client *Client
}

// NewAgentSystemService создает новый сервис для работы с агентными системами
func NewAgentSystemService(client *Client) *AgentSystemService {
	return &AgentSystemService{client: client}
}

// List возвращает список агентных систем
func (s *AgentSystemService) List(ctx context.Context, limit, offset int) (*AgentSystemListResponse, error) {
	query := map[string]string{
		"limit":  fmt.Sprintf("%d", limit),
		"offset": fmt.Sprintf("%d", offset),
	}

	var result AgentSystemListResponse
	err := s.client.Get(ctx, fmt.Sprintf("/api/v1/%s/agentSystems", s.client.projectID), query, &result)
	return &result, err
}

// Get возвращает информацию о конкретной агентной системе
func (s *AgentSystemService) Get(ctx context.Context, systemID string) (*AgentSystem, error) {
	var result AgentSystem
	err := s.client.Get(ctx, fmt.Sprintf("/api/v1/%s/agentSystems/%s", s.client.projectID, systemID), nil, &result)
	return &result, err
}

// Create создает новую агентную систему
func (s *AgentSystemService) Create(ctx context.Context, req *AgentSystemCreateRequest) (*AgentSystem, error) {
	var result AgentSystem
	err := s.client.Post(ctx, fmt.Sprintf("/api/v1/%s/agentSystems", s.client.projectID), req, &result)
	return &result, err
}

// Update обновляет существующую агентную систему
func (s *AgentSystemService) Update(ctx context.Context, systemID string, req *AgentSystemUpdateRequest) (*AgentSystem, error) {
	var result AgentSystem
	err := s.client.Put(ctx, fmt.Sprintf("/api/v1/%s/agentSystems/%s", s.client.projectID, systemID), req, &result)
	return &result, err
}

// Delete удаляет агентную систему
func (s *AgentSystemService) Delete(ctx context.Context, systemID string) error {
	return s.client.Delete(ctx, fmt.Sprintf("/api/v1/%s/agentSystems/%s", s.client.projectID, systemID), nil)
}

// Resume возобновляет работу агентной системы
func (s *AgentSystemService) Resume(ctx context.Context, systemID string) error {
	return s.client.Post(ctx, fmt.Sprintf("/api/v1/%s/agentSystems/resume/%s", s.client.projectID, systemID), nil, nil)
}

// Suspend приостанавливает работу агентной системы
func (s *AgentSystemService) Suspend(ctx context.Context, systemID string) error {
	return s.client.Post(ctx, fmt.Sprintf("/api/v1/%s/agentSystems/suspend/%s", s.client.projectID, systemID), nil, nil)
}

// GetHistory возвращает историю операций агентной системы
func (s *AgentSystemService) GetHistory(ctx context.Context, systemID string) (*AgentSystemHistoryResponse, error) {
	var result AgentSystemHistoryResponse
	err := s.client.Get(ctx, fmt.Sprintf("/api/v1/%s/agentSystems/%s/history", s.client.projectID, systemID), nil, &result)
	return &result, err
}

// GetDiagram возвращает диаграмму агентной системы
func (s *AgentSystemService) GetDiagram(ctx context.Context, systemID string) (*AgentSystemDiagram, error) {
	var result AgentSystemDiagram
	err := s.client.Get(ctx, fmt.Sprintf("/api/v1/%s/agentSystems/%s/diagram", s.client.projectID, systemID), nil, &result)
	return &result, err
}

// AddAgent добавляет агента в систему
func (s *AgentSystemService) AddAgent(ctx context.Context, systemID, agentID string) error {
	req := map[string]string{"agent_id": agentID}
	return s.client.Post(ctx, fmt.Sprintf("/api/v1/%s/agentSystems/%s/agents", s.client.projectID, systemID), req, nil)
}

// RemoveAgent удаляет агента из системы
func (s *AgentSystemService) RemoveAgent(ctx context.Context, systemID, agentID string) error {
	return s.client.Delete(ctx, fmt.Sprintf("/api/v1/%s/agentSystems/%s/agents/%s", s.client.projectID, systemID, agentID), nil)
}

// Clear удаляет всех агентов из системы
func (s *AgentSystemService) Clear(ctx context.Context, systemID string) error {
	return s.client.Post(ctx, fmt.Sprintf("/api/v1/%s/agentSystems/%s/clear", s.client.projectID, systemID), nil, nil)
}
