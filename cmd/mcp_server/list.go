package mcp_server

import (
	"context"
	"fmt"

	"github.com/charmbracelet/log"
	"github.com/cloud-ru/evo-ai-agents-cli/internal/ui"
	"github.com/spf13/cobra"
)

var (
	limit  int
	offset int
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Список MCP серверов",
	Long: `Показывает список всех MCP серверов в проекте.

Команда отображает таблицу с информацией о всех MCP серверах:
• ID сервера
• Название сервера
• Текущий статус
• Дата создания
• Дата последнего обновления

Поддерживает постраничную навигацию с помощью флагов --limit и --offset.

Примеры использования:
  ai-agents-cli mcp-servers list
  ai-agents-cli mcp-servers list --limit 10
  ai-agents-cli mcp-servers list --offset 20 --limit 5`,
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()

		// Проверяем размер терминала
		if err := ui.CheckTerminalSize(); err != nil {
			log.Error("Ошибка размера терминала", "error", err)
			fmt.Println("❌", err)
			return
		}

		// Показываем таблицу MCP серверов
		if err := ui.ShowMCPServersListFromAPI(ctx, limit, offset); err != nil {
			log.Error("Ошибка отображения таблицы MCP серверов", "error", err)
			fmt.Println(ui.CheckAndDisplayError(err))
			return
		}
	},
}

func init() {
	RootCMD.AddCommand(listCmd)

	listCmd.Flags().IntVarP(&limit, "limit", "l", 20, "Количество записей для отображения")
	listCmd.Flags().IntVarP(&offset, "offset", "o", 0, "Смещение для постраничной навигации")
}
