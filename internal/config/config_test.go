package config

import (
	"os"
	"testing"
)

func TestLoad(t *testing.T) {
	// Сохраняем оригинальные переменные окружения
	originalIAMKeyID := os.Getenv("IAM_KEY_ID")
	originalIAMSecret := os.Getenv("IAM_SECRET")
	originalIAMEndpoint := os.Getenv("IAM_ENDPOINT")
	originalProjectID := os.Getenv("PROJECT_ID")
	originalEndpoint := os.Getenv("PUBLIC_API_ENDPOINT")

	// Очищаем переменные окружения
	os.Unsetenv("IAM_KEY_ID")
	os.Unsetenv("IAM_SECRET")
	os.Unsetenv("IAM_ENDPOINT")
	os.Unsetenv("PROJECT_ID")
	os.Unsetenv("PUBLIC_API_ENDPOINT")

	// Восстанавливаем переменные после теста
	defer func() {
		if originalIAMKeyID != "" {
			os.Setenv("IAM_KEY_ID", originalIAMKeyID)
		}
		if originalIAMSecret != "" {
			os.Setenv("IAM_SECRET", originalIAMSecret)
		}
		if originalIAMEndpoint != "" {
			os.Setenv("IAM_ENDPOINT", originalIAMEndpoint)
		}
		if originalProjectID != "" {
			os.Setenv("PROJECT_ID", originalProjectID)
		}
		if originalEndpoint != "" {
			os.Setenv("PUBLIC_API_ENDPOINT", originalEndpoint)
		}
	}()

	tests := []struct {
		name           string
		envVars        map[string]string
		expectError    bool
		expectedConfig *Config
	}{
		{
			name: "valid config with all required fields",
			envVars: map[string]string{
				"IAM_KEY_ID":          "test-key-id",
				"IAM_SECRET":          "test-secret",
				"IAM_ENDPOINT":        "https://iam.test.com",
				"PROJECT_ID":          "test-project",
				"PUBLIC_API_ENDPOINT": "https://api.test.com",
			},
			expectError: false,
			expectedConfig: &Config{
				IAMKeyID:               "test-key-id",
				IAMSecret:              "test-secret",
				IAMEndpoint:            "https://iam.test.com",
				ProjectID:              "test-project",
				IntegrationApiGrpcAddr: "https://api.test.com",
			},
		},
		{
			name: "config with default values",
			envVars: map[string]string{
				"IAM_KEY_ID": "test-key-id",
				"IAM_SECRET": "test-secret",
				"PROJECT_ID": "test-project",
			},
			expectError: false,
			expectedConfig: &Config{
				IAMKeyID:               "test-key-id",
				IAMSecret:              "test-secret",
				IAMEndpoint:            "https://iam.api.cloud.ru",
				ProjectID:              "test-project",
				IntegrationApiGrpcAddr: "ai-agents.api.cloud.ru",
			},
		},
		{
			name:        "missing IAM_KEY_ID (uses default empty string)",
			envVars:     map[string]string{},
			expectError: false,
			expectedConfig: &Config{
				IAMKeyID:               "",
				IAMSecret:              "",
				IAMEndpoint:            "https://iam.api.cloud.ru",
				ProjectID:              "",
				IntegrationApiGrpcAddr: "ai-agents.api.cloud.ru",
			},
		},
		{
			name: "missing IAM_SECRET (uses default empty string)",
			envVars: map[string]string{
				"IAM_KEY_ID": "test-key-id",
			},
			expectError: false,
			expectedConfig: &Config{
				IAMKeyID:               "test-key-id",
				IAMSecret:              "",
				IAMEndpoint:            "https://iam.api.cloud.ru",
				ProjectID:              "",
				IntegrationApiGrpcAddr: "ai-agents.api.cloud.ru",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Устанавливаем переменные окружения для теста
			for key, value := range tt.envVars {
				os.Setenv(key, value)
			}

			// Очищаем переменные после теста
			defer func() {
				for key := range tt.envVars {
					os.Unsetenv(key)
				}
			}()

			config, err := Load()

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if config == nil {
				t.Errorf("Expected config but got nil")
				return
			}

			if tt.expectedConfig != nil {
				if config.IAMKeyID != tt.expectedConfig.IAMKeyID {
					t.Errorf("Expected IAMKeyID %s, got %s", tt.expectedConfig.IAMKeyID, config.IAMKeyID)
				}
				if config.IAMSecret != tt.expectedConfig.IAMSecret {
					t.Errorf("Expected IAMSecret %s, got %s", tt.expectedConfig.IAMSecret, config.IAMSecret)
				}
				if config.IAMEndpoint != tt.expectedConfig.IAMEndpoint {
					t.Errorf("Expected IAMEndpoint %s, got %s", tt.expectedConfig.IAMEndpoint, config.IAMEndpoint)
				}
				if config.ProjectID != tt.expectedConfig.ProjectID {
					t.Errorf("Expected ProjectID %s, got %s", tt.expectedConfig.ProjectID, config.ProjectID)
				}
				if config.IntegrationApiGrpcAddr != tt.expectedConfig.IntegrationApiGrpcAddr {
					t.Errorf("Expected IntegrationApiGrpcAddr %s, got %s", tt.expectedConfig.IntegrationApiGrpcAddr, config.IntegrationApiGrpcAddr)
				}
			}
		})
	}
}

func TestServiceNameConstants(t *testing.T) {
	tests := []struct {
		name     string
		service  ServiceName
		expected string
	}{
		{"SenderApp", SenderApp, "ai-assistant-billing-sender"},
		{"HandlerApp", HandlerApp, "ai-assistant-billing"},
		{"SchedulerApp", SchedulerApp, "ai-assistant-billing-scheduler"},
		{"ConsumerApp", ConsumerApp, "ai-assistant-billing-consumer"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if string(tt.service) != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, string(tt.service))
			}
		})
	}
}
