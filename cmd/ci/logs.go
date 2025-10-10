package ci

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
	"github.com/cloudru/ai-agents-cli/internal/api"
	"github.com/spf13/cobra"
)

var (
	logsFollow     bool
	logsTail       int
	logsSince      string
	logsUntil      string
	logsLevel      string
	logsResource   string
	logsResourceID string
)

// logsCmd represents the logs command
var logsCmd = &cobra.Command{
	Use:   "logs [resource-type] [resource-id]",
	Short: "Просмотр логов ресурсов",
	Long:  "Показывает логи MCP серверов, агентов или агентных систем",
	Args:  cobra.RangeArgs(0, 2),
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()

		// Если указан конкретный ресурс
		if len(args) == 2 {
			logsResource = args[0]
			logsResourceID = args[1]
		}

		// Настраиваем контекст с отменой
		ctx, cancel := context.WithCancel(ctx)
		defer cancel()

		// Обработка сигналов для graceful shutdown
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		go func() {
			<-sigChan
			cancel()
		}()

		// Запускаем мониторинг логов
		if logsFollow {
			monitorLogs(ctx)
		} else {
			showLogs(ctx)
		}
	},
}

func showLogs(ctx context.Context) {
	// Создаем стили для вывода
	headerStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("205")).
		Border(lipgloss.RoundedBorder()).
		Padding(0, 1)

	// Выводим заголовок
	if logsResource != "" && logsResourceID != "" {
		fmt.Println(headerStyle.Render(fmt.Sprintf("📋 Логи %s %s", logsResource, logsResourceID)))
	} else {
		fmt.Println(headerStyle.Render("📋 Логи системы"))
	}
	fmt.Println()

	// Получаем историю ресурса
	if logsResource != "" && logsResourceID != "" {
		switch logsResource {
		case "mcp-server", "mcp":
			showMCPServerLogs(ctx, logsResourceID)
		case "agent":
			showAgentLogs(ctx, logsResourceID)
		case "agent-system", "system":
			showAgentSystemLogs(ctx, logsResourceID)
		default:
			log.Fatal("Unknown resource type. Use: mcp-server, agent, or agent-system")
		}
	} else {
		// Показываем общие логи системы
		showSystemLogs(ctx)
	}
}

func showMCPServerLogs(ctx context.Context, serverID string) {
	history, err := apiClient.MCPServers.GetHistory(ctx, serverID)
	if err != nil {
		log.Fatal("Failed to get MCP server history", "error", err, "server_id", serverID)
	}

	printLogs(history.Data)
}

func showAgentLogs(ctx context.Context, agentID string) {
	history, err := apiClient.Agents.GetHistory(ctx, agentID)
	if err != nil {
		log.Fatal("Failed to get agent history", "error", err, "agent_id", agentID)
	}

	printLogs(history.Data)
}

func showAgentSystemLogs(ctx context.Context, systemID string) {
	history, err := apiClient.AgentSystems.GetHistory(ctx, systemID)
	if err != nil {
		log.Fatal("Failed to get agent system history", "error", err, "system_id", systemID)
	}

	printLogs(history.Data)
}

func showSystemLogs(ctx context.Context) {
	// Получаем логи всех ресурсов
	servers, err := apiClient.MCPServers.List(ctx, 50, 0)
	if err != nil {
		log.Fatal("Failed to list MCP servers", "error", err)
	}

	agents, err := apiClient.Agents.List(ctx, 50, 0)
	if err != nil {
		log.Fatal("Failed to list agents", "error", err)
	}

	systems, err := apiClient.AgentSystems.List(ctx, 50, 0)
	if err != nil {
		log.Fatal("Failed to list agent systems", "error", err)
	}

	// Собираем все логи
	var allLogs []api.HistoryEntry

	// Логи MCP серверов
	for _, server := range servers.Data {
		history, err := apiClient.MCPServers.GetHistory(ctx, server.ID)
		if err == nil {
			allLogs = append(allLogs, history.Data...)
		}
	}

	// Логи агентов
	for _, agent := range agents.Data {
		history, err := apiClient.Agents.GetHistory(ctx, agent.ID)
		if err == nil {
			allLogs = append(allLogs, history.Data...)
		}
	}

	// Логи систем
	for _, system := range systems.Data {
		history, err := apiClient.AgentSystems.GetHistory(ctx, system.ID)
		if err == nil {
			allLogs = append(allLogs, history.Data...)
		}
	}

	// Сортируем по времени (новые сначала)
	sortLogsByTime(allLogs)

	// Ограничиваем количество
	if logsTail > 0 && len(allLogs) > logsTail {
		allLogs = allLogs[:logsTail]
	}

	printLogs(allLogs)
}

func printLogs(logs []api.HistoryEntry) {
	if len(logs) == 0 {
		fmt.Println("🔍 Логи не найдены")
		return
	}

	// Создаем стили для вывода
	timeStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("8")).
		Bold(false)

	actionStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("99"))

	statusStyle := lipgloss.NewStyle().
		Bold(true)

	// Выводим логи
	for _, entry := range logs {
		// Время
		timeStr := timeStyle.Render(entry.CreatedAt.Format("15:04:05"))

		// Действие
		actionStr := actionStyle.Render(entry.Action)

		// Статус
		var statusStr string
		switch entry.Status {
		case "SUCCESS":
			statusStr = statusStyle.Copy().Foreground(lipgloss.Color("2")).Render("✅")
		case "ERROR":
			statusStr = statusStyle.Copy().Foreground(lipgloss.Color("1")).Render("❌")
		case "PENDING":
			statusStr = statusStyle.Copy().Foreground(lipgloss.Color("3")).Render("⏳")
		default:
			statusStr = statusStyle.Copy().Foreground(lipgloss.Color("8")).Render("⚪")
		}

		// Сообщение
		messageStr := entry.Message
		if len(messageStr) > 100 {
			messageStr = messageStr[:100] + "..."
		}

		fmt.Printf("%s %s %s %s\n", timeStr, actionStr, statusStr, messageStr)
	}
}

func monitorLogs(ctx context.Context) {
	// Создаем стили для вывода
	headerStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("205")).
		Border(lipgloss.RoundedBorder()).
		Padding(0, 1)

	fmt.Println(headerStyle.Render("📡 Мониторинг логов (Ctrl+C для выхода)"))
	fmt.Println()

	// Показываем последние логи
	showLogs(ctx)

	// Начинаем мониторинг
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			fmt.Println("\n👋 Мониторинг остановлен")
			return
		case <-ticker.C:
			// В реальной реализации здесь бы был запрос новых логов
			// Пока просто выводим точку
			fmt.Print(".")
		}
	}
}

func sortLogsByTime(logs []api.HistoryEntry) {
	// Простая сортировка пузырьком по времени (новые сначала)
	for i := 0; i < len(logs)-1; i++ {
		for j := 0; j < len(logs)-i-1; j++ {
			if logs[j].CreatedAt.Before(logs[j+1].CreatedAt) {
				logs[j], logs[j+1] = logs[j+1], logs[j]
			}
		}
	}
}

func init() {
	RootCMD.AddCommand(logsCmd)

	logsCmd.Flags().BoolVarP(&logsFollow, "follow", "f", false, "Следить за новыми логами в реальном времени")
	logsCmd.Flags().IntVarP(&logsTail, "tail", "n", 50, "Количество последних записей для показа")
	logsCmd.Flags().StringVarP(&logsSince, "since", "s", "", "Показывать логи с указанного времени")
	logsCmd.Flags().StringVarP(&logsUntil, "until", "u", "", "Показывать логи до указанного времени")
	logsCmd.Flags().StringVarP(&logsLevel, "level", "l", "", "Фильтр по уровню логов")
}
