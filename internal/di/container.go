package di

import (
	"fmt"

	"github.com/cloud-ru/evo-ai-agents-cli/internal/api"
	"github.com/cloud-ru/evo-ai-agents-cli/internal/auth"
	"github.com/cloud-ru/evo-ai-agents-cli/internal/config"
	"github.com/samber/do/v2"
	"github.com/samber/oops"
)

// Container представляет DI контейнер для всех сервисов
type Container struct {
	injector     do.Injector
	errorHandler *ErrorHandler
}

// NewContainer создает новый DI контейнер с зарегистрированными сервисами
func NewContainer() *Container {
	injector := do.New()

	// Регистрируем конфигурацию как singleton
	do.Provide(injector, func(i do.Injector) (*config.Config, error) {
		return config.LoadWithCredentials()
	})

	// Регистрируем IAM сервис как singleton
	do.Provide(injector, func(i do.Injector) (auth.IAMAuthServiceInterface, error) {
		cfg, err := do.Invoke[*config.Config](i)
		if err != nil {
			return nil, fmt.Errorf("failed to get config: %w", err)
		}

		if cfg.IAMKeyID == "" {
			return nil, oops.Errorf("IAM_KEY_ID environment variable is required")
		}
		if cfg.IAMSecret == "" {
			return nil, oops.Errorf("IAM_SECRET environment variable is required")
		}

		return auth.NewIAMAuthService(cfg.IAMKeyID, cfg.IAMSecret, cfg.IAMEndpoint), nil
	})

	// Регистрируем API клиент как singleton
	do.Provide(injector, func(i do.Injector) (*api.API, error) {
		cfg, err := do.Invoke[*config.Config](i)
		if err != nil {
			return nil, oops.Errorf("failed to get config: %w", err)
		}

		authService, err := do.Invoke[auth.IAMAuthServiceInterface](i)
		if err != nil {
			return nil, oops.Errorf("failed to get auth service: %w", err)
		}

		if cfg.ProjectID == "" {
			return nil, oops.Errorf("PROJECT_ID environment variable is required")
		}

		baseURL := "https://" + cfg.IntegrationApiGrpcAddr
		return api.NewAPI(baseURL, cfg.ProjectID, authService), nil
	})

	return &Container{
		injector:     injector,
		errorHandler: NewErrorHandler(),
	}
}

// GetConfig возвращает конфигурацию
func (c *Container) GetConfig() (*config.Config, error) {
	config, err := do.Invoke[*config.Config](c.injector)
	if err != nil {
		return nil, c.errorHandler.HandleConfigError(err)
	}
	return config, nil
}

// GetAuthService возвращает IAM сервис аутентификации
func (c *Container) GetAuthService() (auth.IAMAuthServiceInterface, error) {
	authService, err := do.Invoke[auth.IAMAuthServiceInterface](c.injector)
	if err != nil {
		return nil, c.errorHandler.HandleAuthError(err)
	}
	return authService, nil
}

// GetAPI возвращает API клиент
func (c *Container) GetAPI() (*api.API, error) {
	api, err := do.Invoke[*api.API](c.injector)
	if err != nil {
		return nil, c.errorHandler.HandleAPIError(err)
	}
	return api, nil
}

// Close закрывает контейнер и освобождает ресурсы
func (c *Container) Close() error {
	// В do/v2 нет метода Close, но можно добавить логику очистки если нужно
	return nil
}
