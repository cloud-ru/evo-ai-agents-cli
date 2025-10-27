package api

import (
	"context"
)

// MockIAMService - мок для IAM сервиса для тестов
type MockIAMService struct {
	token string
}

func (m *MockIAMService) GetToken(ctx context.Context) (string, error) {
	return m.token, nil
}

func (m *MockIAMService) IsAuthenticated() bool {
	return m.token != ""
}

func (m *MockIAMService) ClearToken() {
	m.token = ""
}

