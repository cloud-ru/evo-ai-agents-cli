package agent

import (
	"context"
	"fmt"

	"github.com/charmbracelet/log"
	"github.com/cloudru/ai-agents-cli/internal/ui"
	"github.com/spf13/cobra"
)

var (
	agentLimit  int
	agentOffset int
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Список AI агентов",
	Long: `Показывает список всех AI агентов в проекте.

Команда отображает таблицу с информацией о всех агентах:
• ID агента
• Название агента
• Текущий статус
• Дата создания
• Дата последнего обновления

Поддерживает постраничную навигацию с помощью флагов --limit и --offset.

Примеры использования:
  ai-agents-cli agents list
  ai-agents-cli agents list --limit 10
  ai-agents-cli agents list --offset 20 --limit 5`,
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()

		// Проверяем размер терминала
		if err := ui.CheckTerminalSize(); err != nil {
			log.Error("Ошибка размера терминала", "error", err)
			fmt.Println("❌", err)
			return
		}

		// Показываем таблицу агентов
		if err := ui.ShowAgentsListFromAPI(ctx, agentLimit, agentOffset); err != nil {
			log.Error("Ошибка отображения таблицы агентов", "error", err)
			log.Fatal("Failed to show agents table", "error", err)
		}
	},
}

func init() {
	RootCMD.AddCommand(listCmd)

	listCmd.Flags().IntVarP(&agentLimit, "limit", "l", 20, "Количество записей для отображения")
	listCmd.Flags().IntVarP(&agentOffset, "offset", "o", 0, "Смещение для постраничной навигации")
}
