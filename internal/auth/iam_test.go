package auth

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestNewIAMAuthService(t *testing.T) {
	service := NewIAMAuthService("test-key-id", "test-secret", "https://iam.test.com")

	if service.keyID != "test-key-id" {
		t.Errorf("Expected keyID 'test-key-id', got '%s'", service.keyID)
	}

	if service.secret != "test-secret" {
		t.Errorf("Expected secret 'test-secret', got '%s'", service.secret)
	}

	if service.endpoint != "https://iam.test.com" {
		t.Errorf("Expected endpoint 'https://iam.test.com', got '%s'", service.endpoint)
	}

	if service.client == nil {
		t.Errorf("Expected http client to be initialized")
	}

	if service.client.Timeout != 30*time.Second {
		t.Errorf("Expected timeout 30s, got %v", service.client.Timeout)
	}
}

func TestIAMAuthService_GetToken(t *testing.T) {
	// Создаем тестовый сервер
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Проверяем URL
		expectedPath := "/api/v1/auth/token"
		if r.URL.Path != expectedPath {
			t.Errorf("Expected path '%s', got '%s'", expectedPath, r.URL.Path)
		}

		// Проверяем метод
		if r.Method != "POST" {
			t.Errorf("Expected method POST, got %s", r.Method)
		}

		// Проверяем заголовки
		if r.Header.Get("Content-Type") != "application/json" {
			t.Errorf("Expected Content-Type 'application/json', got '%s'", r.Header.Get("Content-Type"))
		}

		// Проверяем тело запроса
		var req IAMTokenRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Errorf("Failed to decode request body: %v", err)
		}

		if req.KeyID != "test-key-id" {
			t.Errorf("Expected KeyID 'test-key-id', got '%s'", req.KeyID)
		}

		if req.Secret != "test-secret" {
			t.Errorf("Expected Secret 'test-secret', got '%s'", req.Secret)
		}

		// Отправляем ответ
		response := IAMTokenResponse{
			AccessToken: "test-access-token",
			TokenType:   "Bearer",
			ExpiresIn:   3600, // 1 час
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	// Создаем сервис с тестовым сервером
	service := NewIAMAuthService("test-key-id", "test-secret", server.URL)

	// Тестируем получение токена
	ctx := context.Background()
	token, err := service.GetToken(ctx)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if token != "test-access-token" {
		t.Errorf("Expected token 'test-access-token', got '%s'", token)
	}

	// Проверяем, что токен кэшируется
	token2, err := service.GetToken(ctx)
	if err != nil {
		t.Fatalf("Unexpected error on second call: %v", err)
	}

	if token2 != "test-access-token" {
		t.Errorf("Expected cached token 'test-access-token', got '%s'", token2)
	}
}

func TestIAMAuthService_GetToken_Error(t *testing.T) {
	// Создаем тестовый сервер, который возвращает ошибку
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Invalid credentials",
		})
	}))
	defer server.Close()

	// Создаем сервис с тестовым сервером
	service := NewIAMAuthService("invalid-key", "invalid-secret", server.URL)

	// Тестируем получение токена с ошибкой
	ctx := context.Background()
	_, err := service.GetToken(ctx)
	if err == nil {
		t.Errorf("Expected error but got none")
	}

	expectedError := "IAM API error (status 401)"
	if err.Error()[:len(expectedError)] != expectedError {
		t.Errorf("Expected error to start with '%s', got '%s'", expectedError, err.Error())
	}
}

func TestIAMAuthService_IsAuthenticated(t *testing.T) {
	service := NewIAMAuthService("test-key-id", "test-secret", "https://iam.test.com")

	// Изначально не аутентифицирован
	if service.IsAuthenticated() {
		t.Errorf("Expected not authenticated initially")
	}

	// Устанавливаем токен вручную для тестирования
	service.mutex.Lock()
	service.token = "test-token"
	service.expiresAt = time.Now().Add(time.Hour)
	service.mutex.Unlock()

	// Теперь должен быть аутентифицирован
	if !service.IsAuthenticated() {
		t.Errorf("Expected authenticated after setting token")
	}

	// Очищаем токен
	service.ClearToken()

	// Снова не аутентифицирован
	if service.IsAuthenticated() {
		t.Errorf("Expected not authenticated after clearing token")
	}
}

func TestIAMAuthService_ClearToken(t *testing.T) {
	service := NewIAMAuthService("test-key-id", "test-secret", "https://iam.test.com")

	// Устанавливаем токен
	service.mutex.Lock()
	service.token = "test-token"
	service.expiresAt = time.Now().Add(time.Hour)
	service.mutex.Unlock()

	// Проверяем, что токен установлен
	if !service.IsAuthenticated() {
		t.Errorf("Expected authenticated before clearing")
	}

	// Очищаем токен
	service.ClearToken()

	// Проверяем, что токен очищен
	if service.IsAuthenticated() {
		t.Errorf("Expected not authenticated after clearing")
	}

	// Проверяем, что поля очищены
	service.mutex.RLock()
	if service.token != "" {
		t.Errorf("Expected empty token, got '%s'", service.token)
	}
	if !service.expiresAt.IsZero() {
		t.Errorf("Expected zero expiresAt, got %v", service.expiresAt)
	}
	service.mutex.RUnlock()
}
