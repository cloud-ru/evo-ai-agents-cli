package mcp_server

import (
	"context"
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
	"github.com/cloudru/ai-agents-cli/internal/di"
	"github.com/spf13/cobra"
)

// resumeCmd represents the resume command
var resumeCmd = &cobra.Command{
	Use:   "resume <server-id>",
	Short: "Возобновить работу MCP сервера",
	Long:  "Возобновляет работу приостановленного MCP сервера",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		serverID := args[0]

		// Возобновляем работу MCP сервера
		// Получаем API клиент из DI контейнера
		container := di.GetContainer()
		apiClient := container.GetAPI()

		err := apiClient.MCPServers.Resume(ctx, serverID)
		if err != nil {
			log.Fatal("Failed to resume MCP server", "error", err, "server_id", serverID)
		}

		// Создаем стили для вывода
		successStyle := lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("2")).
			Border(lipgloss.RoundedBorder()).
			Padding(0, 1)

		// Выводим результат
		fmt.Println(successStyle.Render("✅ MCP сервер возобновлен"))
		fmt.Printf("ID: %s\n", serverID)
	},
}

func init() {
	RootCMD.AddCommand(resumeCmd)
}
