package system

import (
	"context"
	"fmt"

	"github.com/charmbracelet/log"
	"github.com/cloud-ru/evo-ai-agents-cli/internal/di"
	"github.com/cloud-ru/evo-ai-agents-cli/internal/ui"
	"github.com/spf13/cobra"
)

var (
	systemOutputFormat string
	systemLimit        int
	systemOffset       int
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Просмотр списка систем агентов",
	Long:  "Показывает список всех агентных систем с возможностью пагинации",
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()

		// Получаем API клиент из DI контейнера
		container := di.GetContainer()
		apiClient, err := container.GetAPI()
		if err != nil {
			log.Fatal("Failed to get API client", "error", err)
		}

		if systemOutputFormat == "json" {
			// Выводим в JSON формате
			systems, err := apiClient.AgentSystems.List(ctx, systemLimit, systemOffset)
			if err != nil {
				log.Fatal("Failed to get systems list", "error", err)
			}

			// Выводим JSON
			fmt.Printf(`{"data":%d,"total":%d}`, len(systems.Data), systems.Total)
			return
		}

		// Показываем интерактивную таблицу

		if err = ui.ShowAgentSystemsListFromAPI(ctx, systemLimit, systemOffset); err != nil {
			log.Fatal("Failed to show systems table", "error", err)
		}
	},
}

func init() {
	RootCMD.AddCommand(listCmd)

	listCmd.Flags().StringVarP(&systemOutputFormat, "output", "o", "table", "Формат вывода (table, json)")
	listCmd.Flags().IntVarP(&systemLimit, "limit", "l", 20, "Количество систем на странице")
	listCmd.Flags().IntVarP(&systemOffset, "offset", "", 0, "Смещение для пагинации")
}
