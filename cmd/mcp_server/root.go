package mcp_server

import (
	"github.com/charmbracelet/log"
	"github.com/cloudru/ai-agents-cli/internal/api"
	"github.com/cloudru/ai-agents-cli/internal/auth"
	"github.com/cloudru/ai-agents-cli/internal/config"
	"github.com/spf13/cobra"
)

var (
	isVerbose bool
	apiClient *api.API
)

// RootCMD represents the base command when called without any subcommands
var RootCMD = &cobra.Command{
	Use:   "mcp-servers",
	Short: "Управление MCP серверами",
	Long: `Управление MCP (Model Context Protocol) серверами.

MCP серверы предоставляют контекст и инструменты для AI агентов.
Эта команда позволяет создавать, настраивать и управлять MCP серверами
в вашем проекте.

Доступные операции:
• list - Просмотр списка серверов
• get - Получение информации о сервере
• create - Создание нового сервера
• update - Обновление существующего сервера
• delete - Удаление сервера
• resume - Возобновление работы сервера
• suspend - Приостановка сервера

Примеры использования:
  ai-agents-cli mcp-servers list
  ai-agents-cli mcp-servers get server-id
  ai-agents-cli mcp-servers create --name my-server`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Info("Команда MCP серверов вызвана без подкоманды")
		// Показываем справку если нет подкоманд
		cmd.Help()
	},
	Args: cobra.ArbitraryArgs,
}

func init() {
	log.Info("Инициализация MCP серверов команды")
	
	// Инициализируем API клиент
	cfg, err := config.Load()
	if err != nil {
		log.Error("Ошибка загрузки конфигурации", "error", err)
		log.Fatal("Failed to load config", "error", err)
	}

	log.Debug("Конфигурация загружена", "project_id", cfg.ProjectID, "api_endpoint", cfg.IntegrationApiGrpcAddr)

	if cfg.IAMKeyID == "" {
		log.Error("IAM_KEY_ID не установлен")
		log.Fatal("IAM_KEY_ID environment variable is required")
	}

	if cfg.IAMSecret == "" {
		log.Error("IAM_SECRET не установлен")
		log.Fatal("IAM_SECRET environment variable is required")
	}

	if cfg.ProjectID == "" {
		log.Error("PROJECT_ID не установлен")
		log.Fatal("PROJECT_ID environment variable is required")
	}

	// Создаем IAM сервис аутентификации
	log.Info("Создание IAM сервиса аутентификации", "endpoint", cfg.IAMEndpoint)
	authService := auth.NewIAMAuthService(cfg.IAMKeyID, cfg.IAMSecret, cfg.IAMEndpoint)
	
	// Создаем API клиент с IAM аутентификацией
	log.Info("Создание API клиента", "base_url", "https://"+cfg.IntegrationApiGrpcAddr, "project_id", cfg.ProjectID)
	apiClient = api.NewAPI("https://"+cfg.IntegrationApiGrpcAddr, cfg.ProjectID, authService)
	
	log.Info("MCP серверы команда инициализирована успешно")
}
