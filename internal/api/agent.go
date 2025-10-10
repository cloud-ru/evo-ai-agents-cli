package api

import (
	"context"
	"fmt"
	"time"
)

// Agent представляет агента
type Agent struct {
	ID          string                 `json:"id"`
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Status      string                 `json:"status"`
	CreatedAt   time.Time              `json:"created_at"`
	UpdatedAt   time.Time              `json:"updated_at"`
	Options     map[string]interface{} `json:"options"`
	LLMOptions  map[string]interface{} `json:"llm_options"`
	MCPs        []string               `json:"mcp_servers,omitempty"`
}

// AgentCreateRequest представляет запрос на создание агента
type AgentCreateRequest struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description,omitempty"`
	Options     map[string]interface{} `json:"options"`
	LLMOptions  map[string]interface{} `json:"llm_options"`
	MCPs        []string               `json:"mcp_servers,omitempty"`
}

// AgentUpdateRequest представляет запрос на обновление агента
type AgentUpdateRequest struct {
	Name        string                 `json:"name,omitempty"`
	Description string                 `json:"description,omitempty"`
	Options     map[string]interface{} `json:"options,omitempty"`
	LLMOptions  map[string]interface{} `json:"llm_options,omitempty"`
	MCPs        []string               `json:"mcp_servers,omitempty"`
}

// AgentListResponse представляет ответ со списком агентов
type AgentListResponse struct {
	Data  []Agent `json:"data"`
	Total int     `json:"total"`
}

// AgentHistoryResponse представляет ответ с историей агента
type AgentHistoryResponse struct {
	Data []HistoryEntry `json:"data"`
}

// MarketplaceAgent представляет агента из маркетплейса
type MarketplaceAgent struct {
	ID          string                 `json:"id"`
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Status      string                 `json:"status"`
	Type        string                 `json:"type"`
	Categories  []string               `json:"categories"`
	Tags        []string               `json:"tags"`
	Options     map[string]interface{} `json:"options"`
}

// MarketplaceAgentListResponse представляет ответ со списком агентов из маркетплейса
type MarketplaceAgentListResponse struct {
	Data       []MarketplaceAgent `json:"data"`
	Total      int                `json:"total"`
	Categories []string           `json:"categories"`
	Tags       []string           `json:"tags"`
}

// AgentService предоставляет методы для работы с агентами
type AgentService struct {
	client *Client
}

// NewAgentService создает новый сервис для работы с агентами
func NewAgentService(client *Client) *AgentService {
	return &AgentService{client: client}
}

// List возвращает список агентов
func (s *AgentService) List(ctx context.Context, limit, offset int) (*AgentListResponse, error) {
	query := map[string]string{
		"limit":  fmt.Sprintf("%d", limit),
		"offset": fmt.Sprintf("%d", offset),
	}

	var result AgentListResponse
	err := s.client.Get(ctx, fmt.Sprintf("/api/v1/%s/agents", s.client.projectID), query, &result)
	return &result, err
}

// Get возвращает информацию о конкретном агенте
func (s *AgentService) Get(ctx context.Context, agentID string) (*Agent, error) {
	var result Agent
	err := s.client.Get(ctx, fmt.Sprintf("/api/v1/%s/agents/%s", s.client.projectID, agentID), nil, &result)
	return &result, err
}

// Create создает нового агента
func (s *AgentService) Create(ctx context.Context, req *AgentCreateRequest) (*Agent, error) {
	var result Agent
	err := s.client.Post(ctx, fmt.Sprintf("/api/v1/%s/agents", s.client.projectID), req, &result)
	return &result, err
}

// Update обновляет существующего агента
func (s *AgentService) Update(ctx context.Context, agentID string, req *AgentUpdateRequest) (*Agent, error) {
	var result Agent
	err := s.client.Put(ctx, fmt.Sprintf("/api/v1/%s/agents/%s", s.client.projectID, agentID), req, &result)
	return &result, err
}

// Delete удаляет агента
func (s *AgentService) Delete(ctx context.Context, agentID string) error {
	return s.client.Delete(ctx, fmt.Sprintf("/api/v1/%s/agents/%s", s.client.projectID, agentID), nil)
}

// Resume возобновляет работу агента
func (s *AgentService) Resume(ctx context.Context, agentID string) error {
	return s.client.Post(ctx, fmt.Sprintf("/api/v1/%s/agents/resume/%s", s.client.projectID, agentID), nil, nil)
}

// Suspend приостанавливает работу агента
func (s *AgentService) Suspend(ctx context.Context, agentID string) error {
	return s.client.Post(ctx, fmt.Sprintf("/api/v1/%s/agents/suspend/%s", s.client.projectID, agentID), nil, nil)
}

// GetHistory возвращает историю операций агента
func (s *AgentService) GetHistory(ctx context.Context, agentID string) (*AgentHistoryResponse, error) {
	var result AgentHistoryResponse
	err := s.client.Get(ctx, fmt.Sprintf("/api/v1/%s/agents/%s/history", s.client.projectID, agentID), nil, &result)
	return &result, err
}

// SearchMarketplace ищет агентов в маркетплейсе
func (s *AgentService) SearchMarketplace(ctx context.Context, req *MarketplaceSearchRequest) (*MarketplaceAgentListResponse, error) {
	query := make(map[string]string)
	if req.Limit > 0 {
		query["limit"] = fmt.Sprintf("%d", req.Limit)
	}
	if req.Offset > 0 {
		query["offset"] = fmt.Sprintf("%d", req.Offset)
	}
	if req.Name != "" {
		query["name"] = req.Name
	}
	if len(req.Tags) > 0 {
		for _, tag := range req.Tags {
			query["tags"] = tag
		}
	}
	if len(req.Categories) > 0 {
		for _, category := range req.Categories {
			query["categories"] = category
		}
	}
	if len(req.Statuses) > 0 {
		for _, status := range req.Statuses {
			query["statuses"] = status
		}
	}
	if len(req.Types) > 0 {
		for _, agentType := range req.Types {
			query["types"] = agentType
		}
	}

	var result MarketplaceAgentListResponse
	err := s.client.Get(ctx, "/api/v1/marketplace/agents", query, &result)
	return &result, err
}

// GetMarketplaceAgent возвращает информацию об агенте из маркетплейса
func (s *AgentService) GetMarketplaceAgent(ctx context.Context, agentID string) (*MarketplaceAgent, error) {
	var result struct {
		PredefinedAgent MarketplaceAgent `json:"predefined_agent"`
	}
	err := s.client.Get(ctx, fmt.Sprintf("/api/v1/marketplace/agents/%s", agentID), nil, &result)
	return &result.PredefinedAgent, err
}

// MarketplaceSearchRequest представляет запрос поиска в маркетплейсе
type MarketplaceSearchRequest struct {
	Limit      int      `json:"limit,omitempty"`
	Offset     int      `json:"offset,omitempty"`
	Name       string   `json:"name,omitempty"`
	Tags       []string `json:"tags,omitempty"`
	Categories []string `json:"categories,omitempty"`
	Statuses   []string `json:"statuses,omitempty"`
	Types      []string `json:"types,omitempty"`
}
