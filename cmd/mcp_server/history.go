package mcp_server

import (
	"context"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
	"github.com/spf13/cobra"
)

// historyCmd represents the history command
var historyCmd = &cobra.Command{
	Use:   "history <server-id>",
	Short: "История операций MCP сервера",
	Long:  "Показывает историю операций для указанного MCP сервера",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		serverID := args[0]

		// Получаем историю MCP сервера
		history, err := apiClient.MCPServers.GetHistory(ctx, serverID)
		if err != nil {
			log.Fatal("Failed to get MCP server history", "error", err, "server_id", serverID)
		}

		// Создаем стили для вывода
		headerStyle := lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("205")).
			Border(lipgloss.RoundedBorder()).
			Padding(0, 1)

		statusStyle := lipgloss.NewStyle().
			Bold(true)

		// Выводим заголовок
		fmt.Println(headerStyle.Render(fmt.Sprintf("📜 История MCP сервера %s", serverID[:8]+"...")))
		fmt.Println()

		if len(history.Data) == 0 {
			fmt.Println("🔍 История операций не найдена")
			return
		}

		// Создаем таблицу
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "Время\tДействие\tСтатус\tСообщение")
		fmt.Fprintln(w, "-----\t--------\t------\t--------")

		for _, entry := range history.Data {
			status := entry.Status
			switch status {
			case "SUCCESS":
				status = statusStyle.Copy().Foreground(lipgloss.Color("2")).Render("✅ Успех")
			case "ERROR":
				status = statusStyle.Copy().Foreground(lipgloss.Color("1")).Render("❌ Ошибка")
			case "PENDING":
				status = statusStyle.Copy().Foreground(lipgloss.Color("3")).Render("⏳ В процессе")
			default:
				status = statusStyle.Copy().Foreground(lipgloss.Color("8")).Render("⚪ " + status)
			}

			fmt.Fprintf(w, "%s\t%s\t%s\t%s\n",
				entry.CreatedAt.Format("02.01.2006 15:04:05"),
				entry.Action,
				status,
				entry.Message,
			)
		}

		w.Flush()
	},
}

func init() {
	RootCMD.AddCommand(historyCmd)
}
