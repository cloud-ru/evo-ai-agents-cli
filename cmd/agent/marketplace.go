package agent

import (
	"context"
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
	"github.com/cloud-ru/evo-ai-agents-cli/internal/api"
	"github.com/cloud-ru/evo-ai-agents-cli/internal/di"
	"github.com/spf13/cobra"
)

var (
	marketplaceLimit      int
	marketplaceOffset     int
	marketplaceName       string
	marketplaceTags       []string
	marketplaceCategories []string
	marketplaceStatuses   []string
	marketplaceTypes      []string
)

// marketplaceCmd represents the marketplace command
var marketplaceCmd = &cobra.Command{
	Use:   "marketplace",
	Short: "Поиск агентов в маркетплейсе",
	Long:  "Показывает список агентов, доступных в маркетплейсе",
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()

		// Получаем API клиент из DI контейнера
		container := di.GetContainer()
		apiClient := container.GetAPI()

		// Формируем запрос поиска
		searchReq := &api.MarketplaceSearchRequest{
			Limit:      marketplaceLimit,
			Offset:     marketplaceOffset,
			Name:       marketplaceName,
			Tags:       marketplaceTags,
			Categories: marketplaceCategories,
			Statuses:   marketplaceStatuses,
			Types:      marketplaceTypes,
		}

		// Ищем агентов в маркетплейсе
		result, err := apiClient.Agents.SearchMarketplace(ctx, searchReq)
		if err != nil {
			log.Fatal("Failed to search marketplace", "error", err)
		}

		// Создаем стили для вывода
		headerStyle := lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("205")).
			Border(lipgloss.RoundedBorder()).
			Padding(0, 1)

		statusStyle := lipgloss.NewStyle().
			Bold(true)

		typeStyle := lipgloss.NewStyle().
			Bold(true)

		// Выводим заголовок
		fmt.Println(headerStyle.Render(fmt.Sprintf("🏪 Маркетплейс агентов (всего: %d)", result.Total)))
		fmt.Println()

		if len(result.Data) == 0 {
			fmt.Println("🔍 Агенты не найдены")
			return
		}

		// Создаем таблицу
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "ID\tНазвание\tТип\tСтатус\tКатегории\tТеги")
		fmt.Fprintln(w, "---\t--------\t---\t------\t----------\t----")

		for _, agent := range result.Data {
			// Статус
			status := agent.Status
			switch status {
			case "AGENT_PREDEFINED_STATUS_AVAILABLE":
				status = statusStyle.Copy().Foreground(lipgloss.Color("2")).Render("🟢 Доступен")
			case "AGENT_PREDEFINED_STATUS_PREVIEW":
				status = statusStyle.Copy().Foreground(lipgloss.Color("3")).Render("👁️ Превью")
			default:
				status = statusStyle.Copy().Foreground(lipgloss.Color("8")).Render("⚪ " + status)
			}

			// Тип
			agentType := agent.Type
			switch agentType {
			case "AGENT_PREDEFINED_TYPE_FREE_TIER":
				agentType = typeStyle.Copy().Foreground(lipgloss.Color("2")).Render("🆓 Бесплатный")
			case "AGENT_PREDEFINED_TYPE_PAYABLE":
				agentType = typeStyle.Copy().Foreground(lipgloss.Color("3")).Render("💰 Платный")
			case "AGENT_PREDEFINED_TYPE_INTERNAL":
				agentType = typeStyle.Copy().Foreground(lipgloss.Color("1")).Render("🏢 Внутренний")
			default:
				agentType = typeStyle.Copy().Foreground(lipgloss.Color("8")).Render("⚪ " + agentType)
			}

			// Категории и теги
			categories := strings.Join(agent.Categories, ", ")
			if len(categories) > 30 {
				categories = categories[:30] + "..."
			}

			tags := strings.Join(agent.Tags, ", ")
			if len(tags) > 30 {
				tags = tags[:30] + "..."
			}

			fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t%s\n",
				agent.ID[:8]+"...",
				agent.Name,
				agentType,
				status,
				categories,
				tags,
			)
		}

		w.Flush()

		// Показываем доступные категории и теги
		if len(result.Categories) > 0 || len(result.Tags) > 0 {
			fmt.Println()
			fmt.Println(headerStyle.Render("📋 Доступные фильтры"))

			if len(result.Categories) > 0 {
				fmt.Printf("Категории: %s\n", strings.Join(result.Categories, ", "))
			}
			if len(result.Tags) > 0 {
				fmt.Printf("Теги: %s\n", strings.Join(result.Tags, ", "))
			}
		}
	},
}

func init() {
	RootCMD.AddCommand(marketplaceCmd)

	marketplaceCmd.Flags().IntVarP(&marketplaceLimit, "limit", "l", 20, "Количество записей для отображения")
	marketplaceCmd.Flags().IntVarP(&marketplaceOffset, "offset", "o", 0, "Смещение для постраничной навигации")
	marketplaceCmd.Flags().StringVarP(&marketplaceName, "name", "n", "", "Фильтр по названию агента")
	marketplaceCmd.Flags().StringSliceVarP(&marketplaceTags, "tags", "t", []string{}, "Фильтр по тегам")
	marketplaceCmd.Flags().StringSliceVarP(&marketplaceCategories, "categories", "c", []string{}, "Фильтр по категориям")
	marketplaceCmd.Flags().StringSliceVarP(&marketplaceStatuses, "statuses", "s", []string{}, "Фильтр по статусам")
	marketplaceCmd.Flags().StringSliceVarP(&marketplaceTypes, "types", "y", []string{}, "Фильтр по типам")
}
