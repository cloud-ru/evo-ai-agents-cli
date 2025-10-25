package mcp_server

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/charmbracelet/log"
	"github.com/cloud-ru/evo-ai-agents-cli/internal/di"
	"github.com/cloud-ru/evo-ai-agents-cli/internal/ui"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

var (
	outputFormat string
)

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:   "get <server-id>",
	Short: "Получить информацию о MCP сервере",
	Long:  "Показывает подробную информацию о конкретном MCP сервере",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		serverID := args[0]

		// Получаем API клиент из DI контейнера
		container := di.GetContainer()
		apiClient, err := container.GetAPI()
	if err != nil {
		log.Fatal("Failed to get API client", "error", err)
	}

		// Получаем информацию о MCP сервере
		server, err := apiClient.MCPServers.Get(ctx, serverID)
		if err != nil {
			log.Fatal("Failed to get MCP server", "error", err, "server_id", serverID)
		}

		if outputFormat == "json" {
			// Выводим в JSON формате
			encoder := json.NewEncoder(os.Stdout)
			encoder.SetIndent("", "  ")
			if err := encoder.Encode(server); err != nil {
				log.Fatal("Failed to encode JSON", "error", err)
			}
			return
		}

		// Показываем детальную информацию с табами
		if isTerminal() {
			// Интерактивная версия с табами
			program := ui.NewMCPDetailViewModel(ui.NewMCPDetailModel(server))
			if err := program.Start(); err != nil {
				log.Fatal("Failed to start detail view", "error", err)
			}
		} else {
			// Простая версия для не-терминала
			fmt.Printf("🔧 MCP Сервер: %s\n", server.Name)
			fmt.Printf("🆔 ID: %s\n", server.ID)
			fmt.Printf("📊 Статус: %s\n", server.Status)
		}
	},
}

// isTerminal проверяет, является ли терминал терминалом
func isTerminal() bool {
	return term.IsTerminal(int(os.Stdout.Fd()))
}

func init() {
	RootCMD.AddCommand(getCmd)

	getCmd.Flags().StringVarP(&outputFormat, "output", "o", "table", "Формат вывода (table, json)")
}
