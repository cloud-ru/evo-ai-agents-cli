package api

// API представляет основной API клиент со всеми сервисами
type API struct {
	Client       *Client
	MCPServers   *MCPServerService
	Agents       *AgentService
	AgentSystems *AgentSystemService
}

// NewAPI создает новый экземпляр API с всеми сервисами
func NewAPI(baseURL, apiKey, projectID string) *API {
	client := NewClient(baseURL, apiKey, projectID)

	return &API{
		Client:       client,
		MCPServers:   NewMCPServerService(client),
		Agents:       NewAgentService(client),
		AgentSystems: NewAgentSystemService(client),
	}
}
