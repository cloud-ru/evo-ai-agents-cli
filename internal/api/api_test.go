package api

import (
	"testing"
	"github.com/cloudru/ai-agents-cli/internal/auth"
)

func TestNewAPI(t *testing.T) {
	// Создаем мок IAM сервис
	authService := auth.NewIAMAuthService("test-key-id", "test-secret", "https://iam.test.com")
	
	api := NewAPI("https://api.test.com", "test-project", authService)

	if api.Client == nil {
		t.Errorf("Expected Client to be initialized")
	}

	if api.MCPServers == nil {
		t.Errorf("Expected MCPServers to be initialized")
	}

	if api.Agents == nil {
		t.Errorf("Expected Agents to be initialized")
	}

	if api.AgentSystems == nil {
		t.Errorf("Expected AgentSystems to be initialized")
	}

	// Проверяем, что все сервисы используют один и тот же клиент
	if api.MCPServers.client != api.Client {
		t.Errorf("Expected MCPServers to use the same client")
	}

	if api.Agents.client != api.Client {
		t.Errorf("Expected Agents to use the same client")
	}

	if api.AgentSystems.client != api.Client {
		t.Errorf("Expected AgentSystems to use the same client")
	}
}

func TestAPI_ClientConfiguration(t *testing.T) {
	baseURL := "https://api.test.com"
	projectID := "test-project"
	authService := auth.NewIAMAuthService("test-key-id", "test-secret", "https://iam.test.com")

	api := NewAPI(baseURL, projectID, authService)

	// Проверяем конфигурацию клиента
	if api.Client.baseURL != baseURL {
		t.Errorf("Expected baseURL '%s', got '%s'", baseURL, api.Client.baseURL)
	}

	if api.Client.projectID != projectID {
		t.Errorf("Expected projectID '%s', got '%s'", projectID, api.Client.projectID)
	}

	if api.Client.auth == nil {
		t.Errorf("Expected auth service to be initialized")
	}
}
