package agent

import (
	"context"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
	"github.com/spf13/cobra"
)

var (
	agentLimit  int
	agentOffset int
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Список агентов",
	Long:  "Показывает список всех агентов в проекте",
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()

		// Получаем список агентов
		agents, err := apiClient.Agents.List(ctx, agentLimit, agentOffset)
		if err != nil {
			log.Fatal("Failed to list agents", "error", err)
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
		fmt.Println(headerStyle.Render(fmt.Sprintf("🤖 Агенты (всего: %d)", agents.Total)))
		fmt.Println()

		if len(agents.Data) == 0 {
			fmt.Println("🔍 Агенты не найдены")
			return
		}

		// Создаем таблицу
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "ID\tНазвание\tСтатус\tСоздан\tОбновлен")
		fmt.Fprintln(w, "---\t--------\t------\t------\t--------")

		for _, agent := range agents.Data {
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

			fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n",
				agent.ID[:8]+"...",
				agent.Name,
				status,
				agent.CreatedAt.Format("02.01.2006 15:04"),
				agent.UpdatedAt.Format("02.01.2006 15:04"),
			)
		}

		w.Flush()
	},
}

func init() {
	RootCMD.AddCommand(listCmd)

	listCmd.Flags().IntVarP(&agentLimit, "limit", "l", 20, "Количество записей для отображения")
	listCmd.Flags().IntVarP(&agentOffset, "offset", "o", 0, "Смещение для постраничной навигации")
}
