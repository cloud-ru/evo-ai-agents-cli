package agent

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
	"github.com/spf13/cobra"
)

var (
	agentOutputFormat string
)

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:   "get <agent-id>",
	Short: "Получить информацию об агенте",
	Long:  "Показывает подробную информацию о конкретном агенте",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		agentID := args[0]

		// Получаем информацию об агенте
		agent, err := apiClient.Agents.Get(ctx, agentID)
		if err != nil {
			log.Fatal("Failed to get agent", "error", err, "agent_id", agentID)
		}

		if agentOutputFormat == "json" {
			// Выводим в JSON формате
			encoder := json.NewEncoder(os.Stdout)
			encoder.SetIndent("", "  ")
			if err := encoder.Encode(agent); err != nil {
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
		fmt.Println(headerStyle.Render("🤖 Информация об агенте"))
		fmt.Println()

		// Основная информация
		fmt.Printf("%s: %s\n", labelStyle.Render("ID"), valueStyle.Render(agent.ID))
		fmt.Printf("%s: %s\n", labelStyle.Render("Название"), valueStyle.Render(agent.Name))

		if agent.Description != "" {
			fmt.Printf("%s: %s\n", labelStyle.Render("Описание"), valueStyle.Render(agent.Description))
		}

		// Статус
		status := agent.Status
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
		fmt.Printf("%s: %s\n", labelStyle.Render("Создан"), valueStyle.Render(agent.CreatedAt.Format("02.01.2006 15:04:05")))
		fmt.Printf("%s: %s\n", labelStyle.Render("Обновлен"), valueStyle.Render(agent.UpdatedAt.Format("02.01.2006 15:04:05")))

		// Опции
		if len(agent.Options) > 0 {
			fmt.Println()
			fmt.Println(labelStyle.Render("⚙️  Опции:"))
			for key, value := range agent.Options {
				valueStr := fmt.Sprintf("%v", value)
				fmt.Printf("  %s: %s\n", labelStyle.Render(key), valueStyle.Render(valueStr))
			}
		}

		// LLM опции
		if len(agent.LLMOptions) > 0 {
			fmt.Println()
			fmt.Println(labelStyle.Render("🧠 LLM настройки:"))
			for key, value := range agent.LLMOptions {
				valueStr := fmt.Sprintf("%v", value)
				fmt.Printf("  %s: %s\n", labelStyle.Render(key), valueStyle.Render(valueStr))
			}
		}

		// MCP серверы
		if len(agent.MCPs) > 0 {
			fmt.Println()
			fmt.Println(labelStyle.Render("🔌 MCP серверы:"))
			for _, mcp := range agent.MCPs {
				fmt.Printf("  • %s\n", valueStyle.Render(mcp))
			}
		}
	},
}

func init() {
	RootCMD.AddCommand(getCmd)

	getCmd.Flags().StringVarP(&agentOutputFormat, "output", "o", "table", "Формат вывода (table, json)")
}
