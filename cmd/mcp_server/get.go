package mcp_server

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
	"github.com/cloudru/ai-agents-cli/internal/di"
	"github.com/spf13/cobra"
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
		apiClient := container.GetAPI()

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

		// Создаем стили для вывода
		headerStyle := lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("205")).
			Border(lipgloss.RoundedBorder()).
			Padding(0, 1)

		labelStyle := lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("99"))

		valueStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("252"))

		statusStyle := lipgloss.NewStyle().
			Bold(true)

		// Выводим заголовок
		fmt.Println(headerStyle.Render("🔧 Информация о MCP сервере"))
		fmt.Println()

		// Основная информация
		fmt.Printf("%s: %s\n", labelStyle.Render("ID"), valueStyle.Render(server.ID))
		fmt.Printf("%s: %s\n", labelStyle.Render("Название"), valueStyle.Render(server.Name))

		if server.Description != "" {
			fmt.Printf("%s: %s\n", labelStyle.Render("Описание"), valueStyle.Render(server.Description))
		}

		// Статус
		status := server.Status
		switch status {
		case "ACTIVE":
			status = statusStyle.Copy().Foreground(lipgloss.Color("2")).Render("🟢 Активен")
		case "SUSPENDED":
			status = statusStyle.Copy().Foreground(lipgloss.Color("3")).Render("🟡 Приостановлен")
		case "ERROR":
			status = statusStyle.Copy().Foreground(lipgloss.Color("1")).Render("🔴 Ошибка")
		default:
			status = statusStyle.Copy().Foreground(lipgloss.Color("8")).Render("⚪ " + status)
		}
		fmt.Printf("%s: %s\n", labelStyle.Render("Статус"), status)

		// Даты
		fmt.Printf("%s: %s\n", labelStyle.Render("Создан"), valueStyle.Render(server.CreatedAt.Time.Format("02.01.2006 15:04:05")))
		fmt.Printf("%s: %s\n", labelStyle.Render("Обновлен"), valueStyle.Render(server.UpdatedAt.Time.Format("02.01.2006 15:04:05")))

		// Статистика
		fmt.Println()
		fmt.Println(labelStyle.Render("📊 Статистика:"))
		fmt.Printf("  %s: %s\n", labelStyle.Render("Опций"), valueStyle.Render(fmt.Sprintf("%d", len(server.Options))))

		// Опции
		if len(server.Options) > 0 {
			fmt.Println()
			fmt.Println(labelStyle.Render("⚙️  Опции:"))
			for key, value := range server.Options {
				valueStr := fmt.Sprintf("%v", value)
				if len(valueStr) > 60 {
					valueStr = valueStr[:60] + "..."
				}
				fmt.Printf("  %s: %s\n", labelStyle.Render(key), valueStyle.Render(valueStr))
			}
		} else {
			fmt.Println()
			fmt.Println(labelStyle.Render("⚙️  Опции:") + " " + valueStyle.Render("Нет настроек"))
		}

		// Инструменты
		if len(server.Tools) > 0 {
			fmt.Println()
			fmt.Println(labelStyle.Render("🛠️  Инструменты:"))
			for _, tool := range server.Tools {
				fmt.Printf("  • %s: %s\n",
					labelStyle.Render(tool.Name),
					valueStyle.Render(tool.Description))
			}
		}
	},
}

func init() {
	RootCMD.AddCommand(getCmd)

	getCmd.Flags().StringVarP(&outputFormat, "output", "o", "table", "Формат вывода (table, json)")
}
