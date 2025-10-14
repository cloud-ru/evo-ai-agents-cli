package ci

import (
	"context"
	"fmt"
	"os"
	"text/tabwriter"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
	"github.com/cloudru/ai-agents-cli/internal/di"
	"github.com/spf13/cobra"
)

var (
	statusTimeout int
	statusFormat  string
)

// statusCmd represents the status command
var statusCmd = &cobra.Command{
	Use:   "status [resource-type] [resource-id]",
	Short: "Проверка статуса ресурсов",
	Long:  "Проверяет статус MCP серверов, агентов или агентных систем",
	Args:  cobra.RangeArgs(0, 2),
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()

		// Если указан конкретный ресурс
		if len(args) == 2 {
			resourceType := args[0]
			resourceID := args[1]

			switch resourceType {
			case "mcp-server", "mcp":
				checkMCPServerStatus(ctx, resourceID)
			case "agent":
				checkAgentStatus(ctx, resourceID)
			case "agent-system", "system":
				checkAgentSystemStatus(ctx, resourceID)
			default:
				log.Fatal("Unknown resource type. Use: mcp-server, agent, or agent-system")
			}
			return
		}

		// Если указан только тип ресурса, показываем все
		if len(args) == 1 {
			resourceType := args[0]
			switch resourceType {
			case "mcp-servers", "mcp":
				checkAllMCPServersStatus(ctx)
			case "agents":
				checkAllAgentsStatus(ctx)
			case "agent-systems", "systems":
				checkAllAgentSystemsStatus(ctx)
			default:
				log.Fatal("Unknown resource type. Use: mcp-servers, agents, or agent-systems")
			}
			return
		}

		// Если не указаны аргументы, показываем общий статус
		checkOverallStatus(ctx)
	},
}

func checkMCPServerStatus(ctx context.Context, serverID string) {
	container := di.GetContainer()
	apiClient := container.GetAPI()

	server, err := apiClient.MCPServers.Get(ctx, serverID)
	if err != nil {
		log.Fatal("Failed to get MCP server", "error", err, "server_id", serverID)
	}

	printResourceStatus("MCP Server", serverID, server.Status, server.UpdatedAt.Time)
}

func checkAgentStatus(ctx context.Context, agentID string) {
	container := di.GetContainer()
	apiClient := container.GetAPI()

	agent, err := apiClient.Agents.Get(ctx, agentID)
	if err != nil {
		log.Fatal("Failed to get agent", "error", err, "agent_id", agentID)
	}

	printResourceStatus("Agent", agentID, agent.Status, agent.UpdatedAt.Time)
}

func checkAgentSystemStatus(ctx context.Context, systemID string) {
	container := di.GetContainer()
	apiClient := container.GetAPI()

	system, err := apiClient.AgentSystems.Get(ctx, systemID)
	if err != nil {
		log.Fatal("Failed to get agent system", "error", err, "system_id", systemID)
	}

	printResourceStatus("Agent System", systemID, system.Status, system.UpdatedAt)
}

func checkAllMCPServersStatus(ctx context.Context) {
	container := di.GetContainer()
	apiClient := container.GetAPI()

	servers, err := apiClient.MCPServers.List(ctx, 100, 0)
	if err != nil {
		log.Fatal("Failed to list MCP servers", "error", err)
	}

	printResourcesStatus("MCP Servers", servers.Data)
}

func checkAllAgentsStatus(ctx context.Context) {
	container := di.GetContainer()
	apiClient := container.GetAPI()

	agents, err := apiClient.Agents.List(ctx, 100, 0)
	if err != nil {
		log.Fatal("Failed to list agents", "error", err)
	}

	printResourcesStatus("Agents", agents.Data)
}

func checkAllAgentSystemsStatus(ctx context.Context) {
	container := di.GetContainer()
	apiClient := container.GetAPI()

	systems, err := apiClient.AgentSystems.List(ctx, 100, 0)
	if err != nil {
		log.Fatal("Failed to list agent systems", "error", err)
	}

	printResourcesStatus("Agent Systems", systems.Data)
}

func checkOverallStatus(ctx context.Context) {
	container := di.GetContainer()
	apiClient := container.GetAPI()

	// Получаем статус всех ресурсов
	servers, err := apiClient.MCPServers.List(ctx, 100, 0)
	if err != nil {
		log.Fatal("Failed to list MCP servers", "error", err)
	}

	agents, err := apiClient.Agents.List(ctx, 100, 0)
	if err != nil {
		log.Fatal("Failed to list agents", "error", err)
	}

	systems, err := apiClient.AgentSystems.List(ctx, 100, 0)
	if err != nil {
		log.Fatal("Failed to list agent systems", "error", err)
	}

	// Создаем стили для вывода
	headerStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("205")).
		Border(lipgloss.RoundedBorder()).
		Padding(0, 1)

	// Выводим общий статус
	fmt.Println(headerStyle.Render("📊 Общий статус системы"))
	fmt.Println()

	// Статистика MCP серверов
	activeServers := 0
	errorServers := 0
	for _, server := range servers.Data {
		switch server.Status {
		case "ACTIVE":
			activeServers++
		case "ERROR":
			errorServers++
		}
	}

	// Статистика агентов
	activeAgents := 0
	errorAgents := 0
	for _, agent := range agents.Data {
		switch agent.Status {
		case "ACTIVE":
			activeAgents++
		case "ERROR":
			errorAgents++
		}
	}

	// Статистика систем
	activeSystems := 0
	errorSystems := 0
	for _, system := range systems.Data {
		switch system.Status {
		case "ACTIVE":
			activeSystems++
		case "ERROR":
			errorSystems++
		}
	}

	// Выводим таблицу
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "Тип ресурса\tВсего\tАктивных\tОшибок\tСтатус")
	fmt.Fprintln(w, "-----------\t-----\t--------\t-------\t------")

	// MCP серверы
	status := "🟢 OK"
	if errorServers > 0 {
		status = "🔴 ERROR"
	} else if activeServers == 0 {
		status = "⚪ NO DATA"
	}
	fmt.Fprintf(w, "MCP Servers\t%d\t%d\t%d\t%s\n",
		len(servers.Data), activeServers, errorServers, status)

	// Агенты
	status = "🟢 OK"
	if errorAgents > 0 {
		status = "🔴 ERROR"
	} else if activeAgents == 0 {
		status = "⚪ NO DATA"
	}
	fmt.Fprintf(w, "Agents\t%d\t%d\t%d\t%s\n",
		len(agents.Data), activeAgents, errorAgents, status)

	// Системы
	status = "🟢 OK"
	if errorSystems > 0 {
		status = "🔴 ERROR"
	} else if activeSystems == 0 {
		status = "⚪ NO DATA"
	}
	fmt.Fprintf(w, "Agent Systems\t%d\t%d\t%d\t%s\n",
		len(systems.Data), activeSystems, errorSystems, status)

	w.Flush()

	// Общий статус
	fmt.Println()
	totalErrors := errorServers + errorAgents + errorSystems
	if totalErrors == 0 {
		fmt.Println(lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("2")).Render("✅ Все системы работают нормально"))
	} else {
		fmt.Println(lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("1")).Render(fmt.Sprintf("❌ Обнаружено %d ошибок", totalErrors)))
	}
}

func printResourceStatus(resourceType, resourceID, status string, updatedAt time.Time) {
	headerStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("205")).
		Border(lipgloss.RoundedBorder()).
		Padding(0, 1)

	statusStyle := lipgloss.NewStyle().Bold(true)

	// Определяем цвет статуса
	var statusColor lipgloss.Color
	var statusIcon string
	switch status {
	case "ACTIVE":
		statusColor = lipgloss.Color("2")
		statusIcon = "🟢"
	case "SUSPENDED":
		statusColor = lipgloss.Color("3")
		statusIcon = "🟡"
	case "ERROR":
		statusColor = lipgloss.Color("1")
		statusIcon = "🔴"
	default:
		statusColor = lipgloss.Color("8")
		statusIcon = "⚪"
	}

	fmt.Println(headerStyle.Render(fmt.Sprintf("%s %s", statusIcon, resourceType)))
	fmt.Printf("ID: %s\n", resourceID)
	fmt.Printf("Статус: %s\n", statusStyle.Copy().Foreground(statusColor).Render(status))
	fmt.Printf("Обновлен: %s\n", updatedAt.Format("02.01.2006 15:04:05"))
}

func printResourcesStatus(resourceType string, resources interface{}) {
	// Эта функция будет реализована для каждого типа ресурса
	fmt.Printf("Status check for %s (implementation needed)\n", resourceType)
}

func init() {
	RootCMD.AddCommand(statusCmd)

	statusCmd.Flags().IntVarP(&statusTimeout, "timeout", "t", 30, "Таймаут для проверки статуса (секунды)")
	statusCmd.Flags().StringVarP(&statusFormat, "format", "f", "table", "Формат вывода (table, json)")
}
