package registry

import (
	"context"
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
	"github.com/cloud-ru/evo-ai-agents-cli/internal/api"
	"github.com/cloud-ru/evo-ai-agents-cli/internal/di"
	"github.com/spf13/cobra"
)

var (
	listLimit  int
	listOffset int
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Показать список реестров",
	Long:  "Выводит список всех реестров в проекте",
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()

		// Получаем API клиент из DI контейнера
		container := di.GetContainer()
		apiClient, err := container.GetAPI()
		if err != nil {
			log.Fatal("Failed to get API client", "error", err)
		}

		// Получаем список реестров
		response, err := apiClient.Registries.List(ctx, listLimit, listOffset)
		if err != nil {
			log.Fatal("Failed to list registries", "error", err)
		}

		// Создаем стили для вывода
		headerStyle := lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("205")).
			Border(lipgloss.RoundedBorder()).
			Padding(0, 1)

		fmt.Println(headerStyle.Render(fmt.Sprintf("📋 Реестры (%d)", len(response.Registries))))
		fmt.Println()

		if len(response.Registries) == 0 {
			fmt.Println("Реестры не найдены. Создайте реестр с помощью команды:")
			fmt.Println("  ai-agents-cli registry create --name my-registry")
			return
		}

		// Выводим список реестров в табличном формате
		fmt.Println("ID\t\tНазвание\tТип\t\tСтатус\t\tПубличный")
		fmt.Println("────────────────────────────────────────────────────────────────")

		for _, registry := range response.Registries {
			public := "Да"
			if !registry.IsPublic {
				public = "Нет"
			}

			statusIcon := "🟢"
			switch registry.Status {
			case api.RegistryStatusCreating:
				statusIcon = "🟡"
			case api.RegistryStatusError:
				statusIcon = "🔴"
			}

			fmt.Printf("%s\t%s\t%s\t%s %s\t%s\n",
				registry.ID[:8], registry.Name, string(registry.RegistryType), statusIcon, string(registry.Status), public)
		}

		if response.NextPageToken != "" {
			fmt.Println()
			fmt.Println("💡 Для загрузки следующей страницы используйте:")
			fmt.Printf("  ai-agents-cli registry list --offset %s\n", response.NextPageToken)
		}
	},
}

func init() {
	listCmd.Flags().IntVarP(&listLimit, "limit", "l", 100, "Лимит количества результатов")
	listCmd.Flags().IntVarP(&listOffset, "offset", "o", 0, "Смещение для постраничной навигации")
}
