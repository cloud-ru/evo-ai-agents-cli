package system

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/charmbracelet/log"
	"github.com/cloud-ru/evo-ai-agents-cli/internal/di"
	"github.com/cloud-ru/evo-ai-agents-cli/internal/ui"
	"github.com/spf13/cobra"
	"golang.org/x/term"
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
		apiClient, err := container.GetAPI()
	if err != nil {
		log.Fatal("Failed to get API client", "error", err)
	}

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

		// Показываем детальную информацию с табами
		if isTerminal() {
			// Интерактивная версия с табами
			detailModel := ui.NewSystemDetailModel(system)
			program := ui.NewSystemDetailViewModel(detailModel)
			if err := program.Start(); err != nil {
				log.Fatal("Failed to start detail view", "error", err)
			}
		} else {
			// Простая версия для не-терминала
			fmt.Printf("🏗️ Система агентов: %s\n", system.Name)
			fmt.Printf("🆔 ID: %s\n", system.ID)
			fmt.Printf("📊 Статус: %s\n", system.Status)
		}
	},
}

// isTerminal проверяет, является ли терминал терминалом
func isTerminal() bool {
	return term.IsTerminal(int(os.Stdout.Fd()))
}

func init() {
	RootCMD.AddCommand(getCmd)

	getCmd.Flags().StringVarP(&systemGetOutputFormat, "output", "o", "table", "Формат вывода (table, json)")
}
