package mcp_server

import (
	"github.com/charmbracelet/log"
	"github.com/cloudru/ai-agents-cli/internal/api"
	"github.com/cloudru/ai-agents-cli/internal/config"
	"github.com/cloudru/ai-agents-cli/localizations"
	"github.com/spf13/cobra"
)

var (
	isVerbose bool
	apiClient *api.API
)

// RootCMD represents the base command when called without any subcommands
var RootCMD = &cobra.Command{
	Use:   "mcp-servers",
	Short: localizations.Localization.Get("root_short"),
	Long:  localizations.Localization.Get("root_long"),
	Run: func(cmd *cobra.Command, args []string) {
		// Показываем справку если нет подкоманд
		cmd.Help()
	},
	Args: cobra.ArbitraryArgs,
}

func init() {
	// Инициализируем API клиент
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Failed to load config", "error", err)
	}

	if cfg.APIKey == "" {
		log.Fatal("API_KEY environment variable is required")
	}

	if cfg.ProjectID == "" {
		log.Fatal("PROJECT_ID environment variable is required")
	}

	apiClient = api.NewAPI("https://"+cfg.IntegrationApiGrpcAddr, cfg.APIKey, cfg.ProjectID)
}
