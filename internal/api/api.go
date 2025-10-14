package api

import (
	"github.com/cloudru/ai-agents-cli/internal/auth"
)

// API представляет основной API клиент со всеми сервисами
type API struct {
	Client       *Client
	MCPServers   *MCPServerService
	Agents       *AgentService
	AgentSystems *AgentSystemService
	Users        *UserService
}

// NewAPI создает новый экземпляр API с всеми сервисами
func NewAPI(baseURL, projectID string, authService auth.IAMAuthServiceInterface) *API {
	client := NewClient(baseURL, projectID, authService)

	return &API{
		Client:       client,
		MCPServers:   NewMCPServerService(client),
		Agents:       NewAgentService(client),
		AgentSystems: NewAgentSystemService(client),
		Users:        NewUserService(client),
	}
}
