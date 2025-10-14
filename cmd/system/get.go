package system

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
	"github.com/cloudru/ai-agents-cli/internal/di"
	"github.com/cloudru/ai-agents-cli/internal/ui"
	"github.com/spf13/cobra"
)

var (
	systemGetOutputFormat string
)

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:   "get <system-id>",
	Short: "Получить информацию о системе агентов",
	Long:  "Показывает подробную информацию о конкретной системе агентов",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		systemID := args[0]

		// Получаем API клиент из DI контейнера
		container := di.GetContainer()
		apiClient := container.GetAPI()

		// Получаем информацию о системе
		system, err := apiClient.AgentSystems.Get(ctx, systemID)
		if err != nil {
			log.Fatal("Failed to get system", "error", err, "system_id", systemID)
		}

		if systemGetOutputFormat == "json" {
			// Выводим в JSON формате
			encoder := json.NewEncoder(os.Stdout)
			encoder.SetIndent("", "  ")
			if err := encoder.Encode(system); err != nil {
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

		// Выводим заголовок
		fmt.Println(headerStyle.Render("🏢 Информация о системе агентов"))
		fmt.Println()

		// Основная информация
		fmt.Printf("%s: %s\n", labelStyle.Render("ID"), valueStyle.Render(system.ID))
		fmt.Printf("%s: %s\n", labelStyle.Render("Название"), valueStyle.Render(system.Name))

		if system.Description != "" {
			fmt.Printf("%s: %s\n", labelStyle.Render("Описание"), valueStyle.Render(system.Description))
		}

		// Статус
		status := ui.FormatStatus(system.Status)
		fmt.Printf("%s: %s\n", labelStyle.Render("Статус"), status)

		// Даты
		fmt.Printf("%s: %s\n", labelStyle.Render("Создана"), valueStyle.Render(system.CreatedAt.Format("02.01.2006 15:04:05")))
		fmt.Printf("%s: %s\n", labelStyle.Render("Обновлена"), valueStyle.Render(system.UpdatedAt.Format("02.01.2006 15:04:05")))

		// Агенты в системе
		if len(system.Agents) > 0 {
			fmt.Println()
			fmt.Println(labelStyle.Render("🤖 Агенты в системе:"))
			for i, agentID := range system.Agents {
				fmt.Printf("  %d. %s\n", i+1, valueStyle.Render(agentID))
			}
		} else {
			fmt.Println()
			fmt.Println(labelStyle.Render("🤖 Агенты в системе:") + " " + valueStyle.Render("Нет агентов"))
		}

		// Опции
		if len(system.Options) > 0 {
			fmt.Println()
			fmt.Println(labelStyle.Render("⚙️  Опции:"))
			for key, value := range system.Options {
				valueStr := fmt.Sprintf("%v", value)
				if len(valueStr) > 50 {
					valueStr = valueStr[:50] + "..."
				}
				fmt.Printf("  %s: %s\n", labelStyle.Render(key), valueStyle.Render(valueStr))
			}
		}
	},
}

func init() {
	RootCMD.AddCommand(getCmd)

	getCmd.Flags().StringVarP(&systemGetOutputFormat, "output", "o", "table", "Формат вывода (table, json)")
}
