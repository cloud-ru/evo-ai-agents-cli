package api

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
	"github.com/cloudru/ai-agents-cli/internal/auth"
)


func TestNewClient(t *testing.T) {
	authService := auth.NewIAMAuthService("test-key-id", "test-secret", "https://iam.test.com")
	client := NewClient("https://api.test.com", "test-project", authService)

	if client.baseURL != "https://api.test.com" {
		t.Errorf("Expected baseURL 'https://api.test.com', got '%s'", client.baseURL)
	}

	if client.projectID != "test-project" {
		t.Errorf("Expected projectID 'test-project', got '%s'", client.projectID)
	}

	if client.auth == nil {
		t.Errorf("Expected auth service to be initialized")
	}

	if client.httpClient.Timeout != 30*time.Second {
		t.Errorf("Expected timeout 30s, got %v", client.httpClient.Timeout)
	}
}

func TestClient_doRequest(t *testing.T) {
	// Создаем тестовый сервер
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Проверяем заголовки
		if r.Header.Get("Authorization") != "Bearer test-token" {
			t.Errorf("Expected Authorization header 'Bearer test-token', got '%s'", r.Header.Get("Authorization"))
		}

		if r.Header.Get("Content-Type") != "application/json" {
			t.Errorf("Expected Content-Type 'application/json', got '%s'", r.Header.Get("Content-Type"))
		}

		// Проверяем URL
		expectedPath := "/test/path"
		if r.URL.Path != expectedPath {
			t.Errorf("Expected path '%s', got '%s'", expectedPath, r.URL.Path)
		}

		// Проверяем query параметры
		if r.URL.Query().Get("param1") != "value1" {
			t.Errorf("Expected query param 'param1=value1', got '%s'", r.URL.Query().Get("param1"))
		}

		// Отправляем ответ
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"status": "success"})
	}))
	defer server.Close()

	// Создаем клиент с мок IAM сервисом
	mockAuth := &MockIAMService{token: "test-token"}
	client := NewClient(server.URL, "test-project", mockAuth)

	// Тестируем GET запрос
	resp, err := client.doRequest(context.Background(), RequestOptions{
		Method: "GET",
		Path:   "/test/path",
		Query:  map[string]string{"param1": "value1"},
	})

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	resp.Body.Close()
}

func TestClient_doRequest_WithBody(t *testing.T) {
	// Создаем тестовый сервер
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Проверяем метод
		if r.Method != "POST" {
			t.Errorf("Expected method POST, got %s", r.Method)
		}

		// Проверяем тело запроса
		var body map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Errorf("Failed to decode request body: %v", err)
		}

		if body["test"] != "data" {
			t.Errorf("Expected body to contain 'test: data', got %v", body)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"status": "success"})
	}))
	defer server.Close()

	mockAuth := &MockIAMService{token: "test-token"}
	client := NewClient(server.URL, "test-project", mockAuth)

	// Тестируем POST запрос с телом
	resp, err := client.doRequest(context.Background(), RequestOptions{
		Method: "POST",
		Path:   "/test/path",
		Body:   map[string]string{"test": "data"},
	})

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	resp.Body.Close()
}

func TestClient_parseResponse_Success(t *testing.T) {
	mockAuth := &MockIAMService{token: "test-token"}
	client := NewClient("https://api.test.com", "test-project", mockAuth)

	// Создаем тестовый ответ
	responseData := map[string]string{"status": "success", "message": "test"}
	jsonData, _ := json.Marshal(responseData)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(jsonData)
	}))
	defer server.Close()

	resp, err := http.Get(server.URL)
	if err != nil {
		t.Fatalf("Failed to create test response: %v", err)
	}

	var result map[string]string
	err = client.parseResponse(resp, &result)

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if result["status"] != "success" {
		t.Errorf("Expected status 'success', got '%s'", result["status"])
	}

	if result["message"] != "test" {
		t.Errorf("Expected message 'test', got '%s'", result["message"])
	}
}

func TestClient_parseResponse_Error(t *testing.T) {
	mockAuth := &MockIAMService{token: "test-token"}
	client := NewClient("https://api.test.com", "test-project", mockAuth)

	// Создаем тестовый ответ с ошибкой
	errorData := map[string]interface{}{
		"error": map[string]interface{}{
			"code":    400,
			"message": "Bad Request",
		},
	}
	jsonData, _ := json.Marshal(errorData)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write(jsonData)
	}))
	defer server.Close()

	resp, err := http.Get(server.URL)
	if err != nil {
		t.Fatalf("Failed to create test response: %v", err)
	}

	var result map[string]string
	err = client.parseResponse(resp, &result)

	if err == nil {
		t.Fatalf("Expected error but got none")
	}

	expectedError := "API error (status 400): Bad Request"
	if err.Error() != expectedError {
		t.Errorf("Expected error '%s', got '%s'", expectedError, err.Error())
	}
}

func TestClient_Get(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Errorf("Expected method GET, got %s", r.Method)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"status": "success"})
	}))
	defer server.Close()

	mockAuth := &MockIAMService{token: "test-token"}
	client := NewClient(server.URL, "test-project", mockAuth)

	var result map[string]string
	err := client.Get(context.Background(), "/test", map[string]string{"param": "value"}, &result)

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if result["status"] != "success" {
		t.Errorf("Expected status 'success', got '%s'", result["status"])
	}
}

func TestClient_Post(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("Expected method POST, got %s", r.Method)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"status": "created"})
	}))
	defer server.Close()

	mockAuth := &MockIAMService{token: "test-token"}
	client := NewClient(server.URL, "test-project", mockAuth)

	requestData := map[string]string{"name": "test"}
	var result map[string]string
	err := client.Post(context.Background(), "/test", requestData, &result)

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if result["status"] != "created" {
		t.Errorf("Expected status 'created', got '%s'", result["status"])
	}
}

func TestClient_Put(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "PUT" {
			t.Errorf("Expected method PUT, got %s", r.Method)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"status": "updated"})
	}))
	defer server.Close()

	mockAuth := &MockIAMService{token: "test-token"}
	client := NewClient(server.URL, "test-project", mockAuth)

	requestData := map[string]string{"name": "test"}
	var result map[string]string
	err := client.Put(context.Background(), "/test", requestData, &result)

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if result["status"] != "updated" {
		t.Errorf("Expected status 'updated', got '%s'", result["status"])
	}
}

func TestClient_Delete(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "DELETE" {
			t.Errorf("Expected method DELETE, got %s", r.Method)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"status": "deleted"})
	}))
	defer server.Close()

	mockAuth := &MockIAMService{token: "test-token"}
	client := NewClient(server.URL, "test-project", mockAuth)

	var result map[string]string
	err := client.Delete(context.Background(), "/test", &result)

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if result["status"] != "deleted" {
		t.Errorf("Expected status 'deleted', got '%s'", result["status"])
	}
}

func TestClient_Timeout(t *testing.T) {
	// Создаем сервер с задержкой
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(2 * time.Second) // Задержка больше таймаута
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	mockAuth := &MockIAMService{token: "test-token"}
	client := NewClient(server.URL, "test-project", mockAuth)
	client.httpClient.Timeout = 100 * time.Millisecond // Короткий таймаут

	_, err := client.doRequest(context.Background(), RequestOptions{
		Method: "GET",
		Path:   "/test",
	})

	if err == nil {
		t.Fatalf("Expected timeout error but got none")
	}

	// Проверяем, что это ошибка таймаута
	if !contains(err.Error(), "timeout") && !contains(err.Error(), "context deadline exceeded") {
		t.Errorf("Expected timeout error, got: %v", err)
	}
}

// Вспомогательная функция для проверки содержания строки
func contains(s, substr string) bool {
	return len(s) >= len(substr) && s[:len(substr)] == substr || 
		   len(s) > len(substr) && s[len(s)-len(substr):] == substr ||
		   len(s) >= len(substr) && containsSubstring(s, substr)
}

func containsSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
