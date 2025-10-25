package mcp_server

import (
	"context"
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
	"github.com/cloud-ru/evo-ai-agents-cli/internal/di"
	"github.com/spf13/cobra"
)

// suspendCmd represents the suspend command
var suspendCmd = &cobra.Command{
	Use:   "suspend <server-id>",
	Short: "Приостановить работу MCP сервера",
	Long:  "Приостанавливает работу активного MCP сервера",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		serverID := args[0]

		// Приостанавливаем работу MCP сервера
		// Получаем API клиент из DI контейнера
		container := di.GetContainer()
		apiClient, err := container.GetAPI()
		if err != nil {
			log.Fatal("Failed to get API client", "error", err)
		}

		if err := apiClient.MCPServers.Suspend(ctx, serverID); err != nil {
			log.Fatal("Failed to suspend MCP server", "error", err, "server_id", serverID)
		}

		// Создаем стили для вывода
		warningStyle := lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("3")).
			Border(lipgloss.RoundedBorder()).
			Padding(0, 1)

		// Выводим результат
		fmt.Println(warningStyle.Render("⏸️  MCP сервер приостановлен"))
		fmt.Printf("ID: %s\n", serverID)
	},
}

func init() {
	RootCMD.AddCommand(suspendCmd)
}
