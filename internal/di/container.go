package di

import (
	"github.com/cloud-ru/evo-ai-agents-cli/internal/api"
	"github.com/cloud-ru/evo-ai-agents-cli/internal/auth"
	"github.com/cloud-ru/evo-ai-agents-cli/internal/config"
	"github.com/samber/do/v2"
)

// Container представляет DI контейнер для всех сервисов
type Container struct {
	injector do.Injector
}

// NewContainer создает новый DI контейнер с зарегистрированными сервисами
func NewContainer() *Container {
	injector := do.New()

	// Регистрируем конфигурацию как singleton
	do.Provide(injector, func(i do.Injector) (*config.Config, error) {
		return config.Load()
	})

	// Регистрируем IAM сервис как singleton
	do.Provide(injector, func(i do.Injector) (auth.IAMAuthServiceInterface, error) {
		cfg := do.MustInvoke[*config.Config](i)

		if cfg.IAMKeyID == "" {
			panic("IAM_KEY_ID environment variable is required")
		}
		if cfg.IAMSecret == "" {
			panic("IAM_SECRET environment variable is required")
		}

		return auth.NewIAMAuthService(cfg.IAMKeyID, cfg.IAMSecret, cfg.IAMEndpoint), nil
	})

	// Регистрируем API клиент как singleton
	do.Provide(injector, func(i do.Injector) (*api.API, error) {
		cfg := do.MustInvoke[*config.Config](i)
		authService := do.MustInvoke[auth.IAMAuthServiceInterface](i)

		if cfg.ProjectID == "" {
			panic("PROJECT_ID environment variable is required")
		}

		baseURL := "https://" + cfg.IntegrationApiGrpcAddr
		return api.NewAPI(baseURL, cfg.ProjectID, authService), nil
	})

	return &Container{injector: injector}
}

// GetConfig возвращает конфигурацию
func (c *Container) GetConfig() *config.Config {
	return do.MustInvoke[*config.Config](c.injector)
}

// GetAuthService возвращает IAM сервис аутентификации
func (c *Container) GetAuthService() auth.IAMAuthServiceInterface {
	return do.MustInvoke[auth.IAMAuthServiceInterface](c.injector)
}

// GetAPI возвращает API клиент
func (c *Container) GetAPI() *api.API {
	return do.MustInvoke[*api.API](c.injector)
}

// Close закрывает контейнер и освобождает ресурсы
func (c *Container) Close() error {
	// В do/v2 нет метода Close, но можно добавить логику очистки если нужно
	return nil
}
