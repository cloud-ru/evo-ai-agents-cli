package ui

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
	"github.com/cloudru/ai-agents-cli/internal/api"
	"github.com/cloudru/ai-agents-cli/internal/di"
)

// TableInterface определяет интерфейс для работы с таблицами
type TableInterface interface {
	View() string
	Update(msg tea.Msg) (tea.Model, tea.Cmd)
	Init() tea.Cmd
	GetSelectedRow() table.Row
}

// TableProgram представляет программу для отображения таблиц
type TableProgram struct {
	table         TableInterface
	showDetails   bool
	selectedAgent *api.Agent
}

// NewTableProgram создает новую программу таблицы
func NewTableProgram(table TableInterface) *TableProgram {
	return &TableProgram{
		table:         table,
		showDetails:   false,
		selectedAgent: nil,
	}
}

// Run запускает программу таблицы
func (p *TableProgram) Run() error {
	// Проверяем, что мы в интерактивном режиме
	if !isInteractive() {
		// Если не интерактивный режим, просто показываем таблицу как текст
		fmt.Println(p.table.View())
		return nil
	}

	program := tea.NewProgram(p)
	if _, err := program.Run(); err != nil {
		return err
	}
	return nil
}

// isInteractive проверяет, что мы в интерактивном режиме
func isInteractive() bool {
	// Проверяем, что stdout подключен к терминалу
	return isTerminal()
}

// isTerminal проверяет, что файл является терминалом
func isTerminal() bool {
	// Проверяем, что stdout подключен к терминалу
	fileInfo, err := os.Stdout.Stat()
	if err != nil {
		return false
	}

	// Проверяем, что это устройство символьного ввода-вывода (терминал)
	return (fileInfo.Mode() & os.ModeCharDevice) != 0
}

// Init инициализирует программу
func (p *TableProgram) Init() tea.Cmd {
	return nil
}

// Update обновляет программу
func (p *TableProgram) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c", "esc":
			if p.showDetails {
				p.showDetails = false
				p.selectedAgent = nil
				return p, nil
			}
			return p, tea.Quit
		case "enter":
			if !p.showDetails {
				// Получаем выбранную строку и показываем детали
				selectedRow := p.table.GetSelectedRow()
				if len(selectedRow) > 0 {
					agentID := selectedRow[0] // ID из первого элемента

					// Получаем полную информацию об агенте через API
					ctx := context.Background()
					container := di.GetContainer()
					apiClient := container.GetAPI()

					agent, err := apiClient.Agents.Get(ctx, agentID)
					if err != nil {
						log.Error("Ошибка получения информации об агенте", "error", err, "agent_id", agentID)
						return p, nil
					}

					p.selectedAgent = agent
					p.showDetails = true
					log.Debug("Выбран агент", "agent_id", agentID)
				}
				return p, nil
			}
		case "backspace", "b":
			if p.showDetails {
				p.showDetails = false
				p.selectedAgent = nil
				return p, nil
			}
		}
	}

	var cmd tea.Cmd
	updatedTable, cmd := p.table.Update(msg)
	if tableInterface, ok := updatedTable.(TableInterface); ok {
		p.table = tableInterface
	}
	return p, cmd
}

// View отображает программу
func (p *TableProgram) View() string {
	if p.showDetails {
		return p.renderDetails()
	}
	return p.table.View()
}

// renderDetails отображает детали выбранного агента
func (p *TableProgram) renderDetails() string {
	if p.selectedAgent == nil {
		return "Ошибка: агент не выбран"
	}

	// Используем сохраненного агента
	ctx := context.Background()
	container := di.GetContainer()

	// Используем общую функцию для отображения деталей
	result := RenderAgentDetails(p.selectedAgent, ctx, container)

	// Добавляем инструкцию для возврата к таблице
	result += fmt.Sprintf("\n\n%s", lipgloss.NewStyle().Foreground(lipgloss.Color("240")).Render("Нажмите 'b' или 'Backspace' для возврата к таблице"))

	return result
}

// ShowAgentsTable показывает таблицу агентов
func ShowAgentsTable(agents []api.Agent, title string) error {
	table := CreateAgentsTable(agents, title)
	program := NewTableProgram(table)
	return program.Run()
}

// ShowMCPServersTable показывает таблицу MCP серверов
func ShowMCPServersTable(servers []api.MCPServer, title string) error {
	table := CreateMCPServersTable(servers, title)
	program := NewTableProgram(table)
	return program.Run()
}

// ShowAgentsListFromAPI показывает список агентов из API
func ShowAgentsListFromAPI(ctx context.Context, limit, offset int) error {
	container := di.GetContainer()
	apiClient := container.GetAPI()

	// Создаем функцию загрузки данных
	dataLoader := func(ctx context.Context, limit, offset int) ([]table.Row, int, error) {
		log.Debug("Запрос списка агентов", "limit", limit, "offset", offset)

		agents, err := apiClient.Agents.List(ctx, limit, offset)
		if err != nil {
			log.Error("Ошибка получения списка агентов", "error", err)
			return nil, 0, fmt.Errorf("failed to list agents: %w", err)
		}

		log.Debug("Список агентов получен", "total", agents.Total, "count", len(agents.Data))

		// Преобразуем агентов в строки таблицы
		var rows []table.Row
		for _, agent := range agents.Data {
			// Получаем тип агента с переводом
			agentType := FormatAgentType(agent.AgentType)

			// Показываем только ID пользователей в списке (без API запросов)
			createdByInfo := agent.CreatedBy
			updatedByInfo := agent.UpdatedBy

			rows = append(rows, table.Row{
				agent.ID,
				agent.Name,
				FormatStatus(agent.Status),
				agentType,
				createdByInfo,
				updatedByInfo,
				agent.CreatedAt.Time.Format("02.01.2006 15:04"),
				agent.UpdatedAt.Time.Format("02.01.2006 15:04"),
			})
		}

		return rows, agents.Total, nil
	}

	// Создаем колонки таблицы
	columns := []table.Column{
		{Title: "ID", Width: 200},
		{Title: "Название", Width: 100},
		{Title: "Статус", Width: 25},
		{Title: "Тип", Width: 50},
		{Title: "Создал", Width: 200},
		{Title: "Изменил", Width: 200},
		{Title: "Создан", Width: 16},
		{Title: "Обновлен", Width: 16},
	}

	// Создаем модель таблицы с серверной пагинацией
	tableModel := NewServerPaginatedTableModel(ctx, "🤖 Агенты", columns, limit, dataLoader)

	// Создаем программу таблицы
	program := NewTableProgram(tableModel)
	return program.Run()
}

// ShowMCPServersListFromAPI показывает список MCP серверов из API
func ShowMCPServersListFromAPI(ctx context.Context, limit, offset int) error {
	container := di.GetContainer()
	apiClient := container.GetAPI()

	log.Debug("Запрос списка MCP серверов", "limit", limit, "offset", offset)

	servers, err := apiClient.MCPServers.List(ctx, limit, offset)
	if err != nil {
		log.Error("Ошибка получения списка MCP серверов", "error", err)
		return fmt.Errorf("failed to list MCP servers: %w", err)
	}

	log.Debug("Список MCP серверов получен", "total", servers.Total, "count", len(servers.Data))

	title := fmt.Sprintf("📋 MCP Серверы (всего: %d)", servers.Total)
	return ShowMCPServersTable(servers.Data, title)
}

// ShowAgentSystemsListFromAPI показывает список систем агентов из API
func ShowAgentSystemsListFromAPI(ctx context.Context, limit, offset int) error {
	container := di.GetContainer()
	apiClient := container.GetAPI()

	// Создаем функцию загрузки данных
	dataLoader := func(ctx context.Context, limit, offset int) ([]table.Row, int, error) {
		log.Debug("Запрос списка систем агентов", "limit", limit, "offset", offset)

		systems, err := apiClient.AgentSystems.List(ctx, limit, offset)
		if err != nil {
			log.Error("Ошибка получения списка систем агентов", "error", err)
			return nil, 0, fmt.Errorf("failed to list agent systems: %w", err)
		}

		log.Debug("Список систем агентов получен", "total", systems.Total, "count", len(systems.Data))

		// Преобразуем системы в строки таблицы
		var rows []table.Row
		for _, system := range systems.Data {
			rows = append(rows, table.Row{
				system.ID,
				system.Name,
				FormatStatus(system.Status),
				fmt.Sprintf("%d", len(system.Agents)),
				system.CreatedAt.Format("02.01.2006 15:04"),
				system.UpdatedAt.Format("02.01.2006 15:04"),
			})
		}

		return rows, systems.Total, nil
	}

	// Создаем колонки таблицы
	columns := []table.Column{
		{Title: "ID", Width: 40},
		{Title: "Название", Width: 50},
		{Title: "Статус", Width: 25},
		{Title: "Агентов", Width: 10},
		{Title: "Создана", Width: 16},
		{Title: "Обновлена", Width: 16},
	}

	// Создаем модель таблицы с серверной пагинацией
	tableModel := NewServerPaginatedTableModel(ctx, "🏢 Системы агентов", columns, limit, dataLoader)

	// Создаем программу таблицы
	program := NewTableProgram(tableModel)
	return program.Run()
}

// CheckTerminalSize проверяет размер терминала
func CheckTerminalSize() error {
	// Проверяем, что терминал достаточно большой для таблицы
	width, height := 80, 24 // Минимальные размеры

	// В реальном приложении можно использовать termenv для получения реального размера
	if width < 80 || height < 24 {
		return fmt.Errorf("терминал слишком мал. Минимальный размер: 80x24")
	}

	return nil
}

// getCreatedByInfo получает информацию о создателе агента для таблицы
func getCreatedByInfo(ctx context.Context, container *di.Container, userID string) string {
	if userID == "" {
		return "Не указан"
	}

	config := container.GetConfig()
	if config.CustomerID == "" {
		return fmt.Sprintf("ID: %s", userID)
	}

	apiClient := container.GetAPI()
	user, err := apiClient.Users.Get(ctx, config.CustomerID, userID)
	if err != nil {
		return fmt.Sprintf("ID: %s", userID)
	}

	return FormatUserName(user.ID, user.FirstName, user.LastName, user.Email)
}

// getUpdatedByInfo получает информацию об изменяющем агента для таблицы
func getUpdatedByInfo(ctx context.Context, container *di.Container, userID string) string {
	if userID == "" {
		return "Не указан"
	}

	config := container.GetConfig()
	if config.CustomerID == "" {
		return fmt.Sprintf("ID: %s", userID)
	}

	apiClient := container.GetAPI()
	user, err := apiClient.Users.Get(ctx, config.CustomerID, userID)
	if err != nil {
		return fmt.Sprintf("ID: %s", userID)
	}

	return FormatUserName(user.ID, user.FirstName, user.LastName, user.Email)
}

// RenderAgentDetails отображает полную информацию об агенте
func RenderAgentDetails(agent *api.Agent, ctx context.Context, container *di.Container) string {
	// Создаем стили для вывода
	headerStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("205")).
		Border(lipgloss.RoundedBorder()).
		Padding(0, 1).
		Margin(0, 0, 1, 0)

	labelStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("39"))

	valueStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("252"))

	sectionStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("214")).
		Margin(1, 0)

	tabStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("214")).
		Margin(0, 0, 0, 2)

	// Получаем информацию о пользователях
	createdByInfo := getCreatedByInfoForUI(ctx, container, agent.CreatedBy)
	updatedByInfo := getUpdatedByInfoForUI(ctx, container, agent.UpdatedBy)

	// Формируем результат
	result := headerStyle.Render("🤖 Информация об агенте")
	result += "\n\n"

	// ===== ОБЩАЯ ИНФОРМАЦИЯ =====
	result += sectionStyle.Render("📋 ОБЩАЯ ИНФОРМАЦИЯ")
	result += "\n"

	// Основная информация
	result += fmt.Sprintf("%s: %s\n", labelStyle.Render("ID"), valueStyle.Render(agent.ID))
	result += fmt.Sprintf("%s: %s\n", labelStyle.Render("Название"), valueStyle.Render(agent.Name))

	if agent.ProjectID != "" {
		result += fmt.Sprintf("%s: %s\n", labelStyle.Render("Project ID"), valueStyle.Render(agent.ProjectID))
	}

	if agent.Description != "" {
		result += fmt.Sprintf("%s: %s\n", labelStyle.Render("Описание"), valueStyle.Render(agent.Description))
	}

	// Тип агента с переводом
	agentType := FormatAgentType(agent.AgentType)
	result += fmt.Sprintf("%s: %s\n", labelStyle.Render("Тип"), valueStyle.Render(agentType))

	// Статус с полной информацией
	status := FormatStatus(agent.Status)
	result += fmt.Sprintf("%s: %s\n", labelStyle.Render("Статус"), status)

	// Причина статуса
	if agent.StatusReason.ReasonType != "" {
		result += fmt.Sprintf("%s: %s\n", labelStyle.Render("Причина статуса"), valueStyle.Render(agent.StatusReason.ReasonType))
		if agent.StatusReason.Message != "" {
			result += fmt.Sprintf("%s: %s\n", labelStyle.Render("Сообщение"), valueStyle.Render(agent.StatusReason.Message))
		}
		if agent.StatusReason.Key != "" {
			result += fmt.Sprintf("%s: %s\n", labelStyle.Render("Ключ"), valueStyle.Render(agent.StatusReason.Key))
		}
	}

	// Даты и авторы
	result += fmt.Sprintf("%s: %s\n", labelStyle.Render("Создан"), valueStyle.Render(agent.CreatedAt.Time.Format("02.01.2006 15:04:05")))
	result += fmt.Sprintf("%s: %s\n", labelStyle.Render("Создал"), valueStyle.Render(createdByInfo))
	result += fmt.Sprintf("%s: %s\n", labelStyle.Render("Обновлен"), valueStyle.Render(agent.UpdatedAt.Time.Format("02.01.2006 15:04:05")))
	result += fmt.Sprintf("%s: %s\n", labelStyle.Render("Изменил"), valueStyle.Render(updatedByInfo))

	// URLs
	if agent.PublicURL != "" {
		result += fmt.Sprintf("%s: %s\n", labelStyle.Render("Публичный URL"), valueStyle.Render(agent.PublicURL))
	}
	if agent.ArizePhoenixPublicURL != "" {
		result += fmt.Sprintf("%s: %s\n", labelStyle.Render("Arize Phoenix URL"), valueStyle.Render(agent.ArizePhoenixPublicURL))
	}

	// ===== MCP СЕРВЕРА =====
	result += fmt.Sprintf("\n%s\n", sectionStyle.Render("🔌 MCP СЕРВЕРА"))

	// MCP серверы (новые)
	if len(agent.MCPServers) > 0 {
		result += fmt.Sprintf("\n%s\n", tabStyle.Render("📡 Подключенные серверы:"))
		for i, mcp := range agent.MCPServers {
			result += fmt.Sprintf("  %d. %s (%s) - %s\n", i+1, mcp.Name, mcp.ID, mcp.Status)
			if mcp.Source != nil && len(mcp.Source) > 0 {
				for key, value := range mcp.Source {
					result += fmt.Sprintf("     %s: %v\n", labelStyle.Render("Источник"), valueStyle.Render(fmt.Sprintf("%s: %v", key, value)))
				}
			}
			if len(mcp.Tools) > 0 {
				result += fmt.Sprintf("     %s (%d):\n", labelStyle.Render("Инструменты"), len(mcp.Tools))
				for j, tool := range mcp.Tools {
					result += fmt.Sprintf("       %d. %s\n", j+1, tool.Name)
					if tool.Description != "" {
						// Обрезаем описание если слишком длинное
						desc := tool.Description
						if len(desc) > 100 {
							desc = desc[:100] + "..."
						}
						result += fmt.Sprintf("          %s\n", valueStyle.Render(desc))
					}
					result += fmt.Sprintf("          %s: %d\n", labelStyle.Render("Аргументы"), len(tool.Args))
				}
			}
		}
	} else {
		result += fmt.Sprintf("\n%s\n", tabStyle.Render("❌ MCP серверы не подключены"))
	}

	// MCP серверы (старые)
	if len(agent.MCPs) > 0 {
		result += fmt.Sprintf("\n%s\n", tabStyle.Render("📡 Старые MCP серверы:"))
		for i, mcp := range agent.MCPs {
			result += fmt.Sprintf("  %d. %s\n", i+1, mcp)
		}
	}

	// ===== ДОПОЛНИТЕЛЬНЫЕ ОПЦИИ =====
	result += fmt.Sprintf("\n%s\n", sectionStyle.Render("⚙️ ДОПОЛНИТЕЛЬНЫЕ ОПЦИИ"))

	// Статистика
	result += fmt.Sprintf("\n%s\n", tabStyle.Render("📊 Статистика:"))
	mcpCount := len(agent.MCPServers)
	if len(agent.MCPs) > 0 {
		mcpCount = len(agent.MCPs)
	}
	result += fmt.Sprintf("  %s: %d\n", labelStyle.Render("MCP серверов"), mcpCount)
	result += fmt.Sprintf("  %s: %d\n", labelStyle.Render("Опций"), len(agent.Options))
	result += fmt.Sprintf("  %s: %d\n", labelStyle.Render("Интеграций"), len(agent.IntegrationOptions))
	if len(agent.UsedInAgentSystems) > 0 {
		result += fmt.Sprintf("  %s: %d\n", labelStyle.Render("Используется в системах"), len(agent.UsedInAgentSystems))
	}

	// Тип инстанса
	if agent.InstanceType.ID != "" {
		result += fmt.Sprintf("\n%s\n", tabStyle.Render("💻 Тип инстанса:"))
		result += fmt.Sprintf("  %s: %s\n", labelStyle.Render("ID"), valueStyle.Render(agent.InstanceType.ID))
		result += fmt.Sprintf("  %s: %s\n", labelStyle.Render("Название"), valueStyle.Render(agent.InstanceType.Name))
		if agent.InstanceType.SKUCode != "" {
			result += fmt.Sprintf("  %s: %s\n", labelStyle.Render("SKU код"), valueStyle.Render(agent.InstanceType.SKUCode))
		}
		if agent.InstanceType.ResourceCode != "" {
			result += fmt.Sprintf("  %s: %s\n", labelStyle.Render("Код ресурса"), valueStyle.Render(agent.InstanceType.ResourceCode))
		}
		result += fmt.Sprintf("  %s: %s\n", labelStyle.Render("Активен"), valueStyle.Render(fmt.Sprintf("%t", agent.InstanceType.IsActive)))
		if agent.InstanceType.MCPU > 0 {
			result += fmt.Sprintf("  %s: %s\n", labelStyle.Render("CPU"), valueStyle.Render(fmt.Sprintf("%d мCPU", agent.InstanceType.MCPU)))
		}
		if agent.InstanceType.MibRAM > 0 {
			result += fmt.Sprintf("  %s: %s\n", labelStyle.Render("RAM"), valueStyle.Render(fmt.Sprintf("%d МБ", agent.InstanceType.MibRAM)))
		}
		if agent.InstanceType.CreatedAt != "" {
			result += fmt.Sprintf("  %s: %s\n", labelStyle.Render("Создан"), valueStyle.Render(agent.InstanceType.CreatedAt))
		}
		if agent.InstanceType.UpdatedAt != "" {
			result += fmt.Sprintf("  %s: %s\n", labelStyle.Render("Обновлен"), valueStyle.Render(agent.InstanceType.UpdatedAt))
		}
		if agent.InstanceType.CreatedBy != "" {
			result += fmt.Sprintf("  %s: %s\n", labelStyle.Render("Создал"), valueStyle.Render(agent.InstanceType.CreatedBy))
		}
		if agent.InstanceType.UpdatedBy != "" {
			result += fmt.Sprintf("  %s: %s\n", labelStyle.Render("Изменил"), valueStyle.Render(agent.InstanceType.UpdatedBy))
		}
	}

	// System Prompt из Options
	if systemPrompt, ok := agent.Options["systemPrompt"]; ok {
		result += fmt.Sprintf("\n%s\n", tabStyle.Render("📝 System Prompt:"))
		result += fmt.Sprintf("  %s\n", valueStyle.Render(fmt.Sprintf("%v", systemPrompt)))
	}

	// LLM настройки из Options
	if llm, ok := agent.Options["llm"]; ok {
		if llmMap, ok := llm.(map[string]interface{}); ok {
			result += fmt.Sprintf("\n%s\n", tabStyle.Render("🧠 LLM настройки:"))
			if foundationModels, ok := llmMap["foundationModels"]; ok {
				if fmMap, ok := foundationModels.(map[string]interface{}); ok {
					if modelName, ok := fmMap["modelName"]; ok {
						result += fmt.Sprintf("  %s: %s\n", labelStyle.Render("Модель"), valueStyle.Render(fmt.Sprintf("%v", modelName)))
					}
					if gcInstanceId, ok := fmMap["gcInstanceId"]; ok {
						result += fmt.Sprintf("  %s: %s\n", labelStyle.Render("GC Instance ID"), valueStyle.Render(fmt.Sprintf("%v", gcInstanceId)))
					}
				}
			}
		}
	}

	// Scaling настройки из Options
	if scaling, ok := agent.Options["scaling"]; ok {
		if scalingMap, ok := scaling.(map[string]interface{}); ok {
			result += fmt.Sprintf("\n%s\n", tabStyle.Render("📈 Scaling настройки:"))
			if minScale, ok := scalingMap["minScale"]; ok {
				result += fmt.Sprintf("  %s: %s\n", labelStyle.Render("Min Scale"), valueStyle.Render(fmt.Sprintf("%v", minScale)))
			}
			if maxScale, ok := scalingMap["maxScale"]; ok {
				result += fmt.Sprintf("  %s: %s\n", labelStyle.Render("Max Scale"), valueStyle.Render(fmt.Sprintf("%v", maxScale)))
			}
			if isKeepAlive, ok := scalingMap["isKeepAlive"]; ok {
				result += fmt.Sprintf("  %s: %s\n", labelStyle.Render("Keep Alive"), valueStyle.Render(fmt.Sprintf("%v", isKeepAlive)))
			}
			if keepAliveDuration, ok := scalingMap["keepAliveDuration"]; ok {
				if durationMap, ok := keepAliveDuration.(map[string]interface{}); ok {
					hours := durationMap["hours"]
					minutes := durationMap["minutes"]
					seconds := durationMap["seconds"]
					result += fmt.Sprintf("  %s: %s\n", labelStyle.Render("Keep Alive Duration"), valueStyle.Render(fmt.Sprintf("%vh %vm %vs", hours, minutes, seconds)))
				}
			}
		}
	}

	// Аутентификация из IntegrationOptions
	if authOptions, ok := agent.IntegrationOptions["authOptions"]; ok {
		if authMap, ok := authOptions.(map[string]interface{}); ok {
			result += fmt.Sprintf("\n%s\n", tabStyle.Render("🔐 Аутентификация:"))
			if isEnabled, ok := authMap["isEnabled"]; ok {
				result += fmt.Sprintf("  %s: %s\n", labelStyle.Render("Включена"), valueStyle.Render(fmt.Sprintf("%v", isEnabled)))
			}
			if authType, ok := authMap["type"]; ok {
				result += fmt.Sprintf("  %s: %s\n", labelStyle.Render("Тип"), valueStyle.Render(fmt.Sprintf("%v", authType)))
			}
			if serviceAccountId, ok := authMap["serviceAccountId"]; ok {
				result += fmt.Sprintf("  %s: %s\n", labelStyle.Render("Service Account ID"), valueStyle.Render(fmt.Sprintf("%v", serviceAccountId)))
			}
		}
	}

	// Логирование из IntegrationOptions
	if logging, ok := agent.IntegrationOptions["logging"]; ok {
		if loggingMap, ok := logging.(map[string]interface{}); ok {
			result += fmt.Sprintf("\n%s\n", tabStyle.Render("📊 Логирование:"))
			if isEnabled, ok := loggingMap["isEnabledLogging"]; ok {
				result += fmt.Sprintf("  %s: %s\n", labelStyle.Render("Включено"), valueStyle.Render(fmt.Sprintf("%v", isEnabled)))
			}
			if logGroupId, ok := loggingMap["logGroupId"]; ok {
				result += fmt.Sprintf("  %s: %s\n", labelStyle.Render("Log Group ID"), valueStyle.Render(fmt.Sprintf("%v", logGroupId)))
			}
		}
	}

	// Автообновление из IntegrationOptions
	if autoUpdate, ok := agent.IntegrationOptions["autoUpdateOptions"]; ok {
		if autoUpdateMap, ok := autoUpdate.(map[string]interface{}); ok {
			result += fmt.Sprintf("\n%s\n", tabStyle.Render("🔄 Автообновление:"))
			if isEnabled, ok := autoUpdateMap["isEnabled"]; ok {
				result += fmt.Sprintf("  %s: %s\n", labelStyle.Render("Включено"), valueStyle.Render(fmt.Sprintf("%v", isEnabled)))
			}
		}
	}

	// Источник образа
	if agent.ImageSource != nil {
		result += fmt.Sprintf("\n%s\n", tabStyle.Render("🐳 Источник образа:"))
		for key, value := range agent.ImageSource {
			result += fmt.Sprintf("  %s: %s\n", labelStyle.Render(key), valueStyle.Render(fmt.Sprintf("%v", value)))
		}
	}

	// Интеграционные настройки
	if len(agent.IntegrationOptions) > 0 {
		result += fmt.Sprintf("\n%s\n", tabStyle.Render("🔗 Интеграционные настройки:"))
		for key, value := range agent.IntegrationOptions {
			result += fmt.Sprintf("  %s: %v\n", labelStyle.Render(key), valueStyle.Render(fmt.Sprintf("%v", value)))
		}
	}

	// Опции
	if len(agent.Options) > 0 {
		result += fmt.Sprintf("\n%s\n", tabStyle.Render("⚙️ Опции:"))
		for key, value := range agent.Options {
			result += fmt.Sprintf("  %s: %v\n", labelStyle.Render(key), valueStyle.Render(fmt.Sprintf("%v", value)))
		}
	}

	return result
}

// getCreatedByInfoForUI получает информацию о создателе агента для UI
func getCreatedByInfoForUI(ctx context.Context, container *di.Container, userID string) string {
	if userID == "" {
		return "Не указан"
	}

	config := container.GetConfig()
	if config.CustomerID == "" {
		return fmt.Sprintf("ID: %s", userID)
	}

	apiClient := container.GetAPI()
	user, err := apiClient.Users.Get(ctx, config.CustomerID, userID)
	if err != nil {
		return fmt.Sprintf("ID: %s", userID)
	}

	return FormatUserName(user.ID, user.FirstName, user.LastName, user.Email)
}

// getUpdatedByInfoForUI получает информацию об изменяющем агента для UI
func getUpdatedByInfoForUI(ctx context.Context, container *di.Container, userID string) string {
	if userID == "" {
		return "Не указан"
	}

	config := container.GetConfig()
	if config.CustomerID == "" {
		return fmt.Sprintf("ID: %s", userID)
	}

	apiClient := container.GetAPI()
	user, err := apiClient.Users.Get(ctx, config.CustomerID, userID)
	if err != nil {
		return fmt.Sprintf("ID: %s", userID)
	}

	return FormatUserName(user.ID, user.FirstName, user.LastName, user.Email)
}

// FormatAgentType форматирует тип агента с переводом
func FormatAgentType(agentType string) string {
	if agentType == "" {
		return "Неизвестно"
	}

	// Маппинг типов агентов на русские названия
	typeTranslations := map[string]string{
		"AGENT_TYPE_FROM_HUB":        "Из маркетплейса",
		"AGENT_TYPE_CUSTOM":          "Пользовательский",
		"AGENT_TYPE_PREDEFINED":      "Предустановленный",
		"AGENT_TYPE_TEMPLATE":        "Шаблон",
		"AGENT_TYPE_WORKFLOW":        "Рабочий процесс",
		"AGENT_TYPE_INTEGRATION":     "Интеграция",
		"AGENT_TYPE_ANALYTICS":       "Аналитика",
		"AGENT_TYPE_AUTOMATION":      "Автоматизация",
		"AGENT_TYPE_CHAT":            "Чат-бот",
		"AGENT_TYPE_ASSISTANT":       "Ассистент",
		"AGENT_TYPE_WORKER":          "Работник",
		"AGENT_TYPE_ANALYZER":        "Анализатор",
		"AGENT_TYPE_GENERATOR":       "Генератор",
		"AGENT_TYPE_CLASSIFIER":      "Классификатор",
		"AGENT_TYPE_TRANSLATOR":      "Переводчик",
		"AGENT_TYPE_SUMMARIZER":      "Суммаризатор",
		"AGENT_TYPE_EXTRACTOR":       "Извлекатель",
		"AGENT_TYPE_VALIDATOR":       "Валидатор",
		"AGENT_TYPE_OPTIMIZER":       "Оптимизатор",
		"AGENT_TYPE_MONITOR":         "Монитор",
		"AGENT_TYPE_SCHEDULER":       "Планировщик",
		"AGENT_TYPE_ROUTER":          "Маршрутизатор",
		"AGENT_TYPE_AGGREGATOR":      "Агрегатор",
		"AGENT_TYPE_FILTER":          "Фильтр",
		"AGENT_TYPE_TRANSFORMER":     "Трансформер",
		"AGENT_TYPE_ENRICHER":        "Обогатитель",
		"AGENT_TYPE_NOTIFIER":        "Уведомитель",
		"AGENT_TYPE_ARCHIVER":        "Архиватор",
		"AGENT_TYPE_BACKUP":          "Резервное копирование",
		"AGENT_TYPE_SYNC":            "Синхронизатор",
		"AGENT_TYPE_MIGRATOR":        "Мигратор",
		"AGENT_TYPE_CLEANER":         "Очиститель",
		"AGENT_TYPE_SECURITY":        "Безопасность",
		"AGENT_TYPE_AUDIT":           "Аудит",
		"AGENT_TYPE_COMPLIANCE":      "Соответствие",
		"AGENT_TYPE_REPORTING":       "Отчетность",
		"AGENT_TYPE_PREDICTION":      "Прогнозирование",
		"AGENT_TYPE_RECOMMENDATION":  "Рекомендации",
		"AGENT_TYPE_PERSONALIZATION": "Персонализация",
		"AGENT_TYPE_API":             "API",
		"AGENT_TYPE_WEBHOOK":         "Webhook",
		"AGENT_TYPE_QUEUE":           "Очередь",
		"AGENT_TYPE_CACHE":           "Кэш",
		"AGENT_TYPE_DATABASE":        "База данных",
		"AGENT_TYPE_FILE":            "Файл",
		"AGENT_TYPE_EMAIL":           "Email",
		"AGENT_TYPE_SMS":             "SMS",
		"AGENT_TYPE_PUSH":            "Push-уведомления",
		"AGENT_TYPE_VOICE":           "Голос",
		"AGENT_TYPE_VIDEO":           "Видео",
		"AGENT_TYPE_IMAGE":           "Изображение",
		"AGENT_TYPE_DOCUMENT":        "Документ",
		"AGENT_TYPE_SPREADSHEET":     "Таблица",
		"AGENT_TYPE_PRESENTATION":    "Презентация",
		"AGENT_TYPE_PDF":             "PDF",
		"AGENT_TYPE_XML":             "XML",
		"AGENT_TYPE_JSON":            "JSON",
		"AGENT_TYPE_CSV":             "CSV",
		"AGENT_TYPE_EXCEL":           "Excel",
		"AGENT_TYPE_WORD":            "Word",
		"AGENT_TYPE_POWERPOINT":      "PowerPoint",
		"AGENT_TYPE_UNKNOWN":         "Неизвестно",
		"AGENT_TYPE_OTHER":           "Другое",
		// Обратная совместимость
		"CHAT":            "Чат-бот",
		"ASSISTANT":       "Ассистент",
		"WORKER":          "Работник",
		"ANALYZER":        "Анализатор",
		"GENERATOR":       "Генератор",
		"CLASSIFIER":      "Классификатор",
		"TRANSLATOR":      "Переводчик",
		"SUMMARIZER":      "Суммаризатор",
		"EXTRACTOR":       "Извлекатель",
		"VALIDATOR":       "Валидатор",
		"OPTIMIZER":       "Оптимизатор",
		"MONITOR":         "Монитор",
		"SCHEDULER":       "Планировщик",
		"ROUTER":          "Маршрутизатор",
		"AGGREGATOR":      "Агрегатор",
		"FILTER":          "Фильтр",
		"TRANSFORMER":     "Трансформер",
		"ENRICHER":        "Обогатитель",
		"NOTIFIER":        "Уведомитель",
		"ARCHIVER":        "Архиватор",
		"BACKUP":          "Резервное копирование",
		"SYNC":            "Синхронизатор",
		"MIGRATOR":        "Мигратор",
		"CLEANER":         "Очиститель",
		"SECURITY":        "Безопасность",
		"AUDIT":           "Аудит",
		"COMPLIANCE":      "Соответствие",
		"REPORTING":       "Отчетность",
		"ANALYTICS":       "Аналитика",
		"PREDICTION":      "Прогнозирование",
		"RECOMMENDATION":  "Рекомендации",
		"PERSONALIZATION": "Персонализация",
		"AUTOMATION":      "Автоматизация",
		"INTEGRATION":     "Интеграция",
		"API":             "API",
		"WEBHOOK":         "Webhook",
		"QUEUE":           "Очередь",
		"CACHE":           "Кэш",
		"DATABASE":        "База данных",
		"FILE":            "Файл",
		"EMAIL":           "Email",
		"SMS":             "SMS",
		"PUSH":            "Push-уведомления",
		"VOICE":           "Голос",
		"VIDEO":           "Видео",
		"IMAGE":           "Изображение",
		"DOCUMENT":        "Документ",
		"SPREADSHEET":     "Таблица",
		"PRESENTATION":    "Презентация",
		"PDF":             "PDF",
		"XML":             "XML",
		"JSON":            "JSON",
		"CSV":             "CSV",
		"EXCEL":           "Excel",
		"WORD":            "Word",
		"POWERPOINT":      "PowerPoint",
		"UNKNOWN":         "Неизвестно",
		"OTHER":           "Другое",
	}

	if translation, exists := typeTranslations[strings.ToUpper(agentType)]; exists {
		return translation
	}

	// Если тип не найден в переводах, возвращаем оригинальное значение
	return agentType
}

// FormatUserName форматирует имя пользователя для отображения
func FormatUserName(userID, firstName, lastName, email string) string {
	if firstName != "" && lastName != "" {
		return fmt.Sprintf("%s %s (%s)", firstName, lastName, userID)
	}
	if firstName != "" {
		return fmt.Sprintf("%s (%s)", firstName, userID)
	}
	if lastName != "" {
		return fmt.Sprintf("%s (%s)", lastName, userID)
	}
	if email != "" {
		return fmt.Sprintf("%s (%s)", email, userID)
	}
	// Если нет дополнительной информации, показываем только ID
	return fmt.Sprintf("ID: %s", userID)
}

// TruncateString обрезает строку до указанной длины
func TruncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}

// FormatStatus форматирует статус агента
func FormatStatus(status string) string {
	switch status {
	case "AGENT_STATUS_ACTIVE":
		return "🟢 Запущен"
	case "AGENT_STATUS_INACTIVE":
		return "🔴 Остановлен"
	case "AGENT_STATUS_PENDING":
		return "⏳ Ожидает"
	case "AGENT_STATUS_ERROR":
		return "❌ Ошибка"
	case "AGENT_STATUS_DELETING":
		return "🔴 Удаляется"
	case "AGENT_STATUS_DELETED":
		return "⚫ Удален"
	case "AGENT_STATUS_COOLED":
		return "💤 Ожидает запроса"
	default:
		return status
	}
}

// CreateAgentsTable создает таблицу агентов
func CreateAgentsTable(agents []api.Agent, title string) *TableModel {
	columns := []table.Column{
		{Title: "ID", Width: 200},
		{Title: "Название", Width: 100},
		{Title: "Статус", Width: 25},
		{Title: "Тип", Width: 50},
		{Title: "Создал", Width: 200},
		{Title: "Изменил", Width: 200},
		{Title: "Создан", Width: 16},
		{Title: "Обновлен", Width: 16},
	}

	var rows []table.Row
	for _, agent := range agents {
		agentType := FormatAgentType(agent.AgentType)
		rows = append(rows, table.Row{
			agent.ID,
			agent.Name,
			FormatStatus(agent.Status),
			agentType,
			"ID: " + agent.CreatedBy,
			"ID: " + agent.UpdatedBy,
			agent.CreatedAt.Time.Format("02.01.2006 15:04"),
			agent.UpdatedAt.Time.Format("02.01.2006 15:04"),
		})
	}

	return NewTableModel(title, columns, rows)
}

// CreateMCPServersTable создает таблицу MCP серверов
func CreateMCPServersTable(servers []api.MCPServer, title string) *TableModel {
	columns := []table.Column{
		{Title: "ID", Width: 40},
		{Title: "Название", Width: 50},
		{Title: "Статус", Width: 25},
		{Title: "Описание", Width: 60},
		{Title: "Создан", Width: 16},
		{Title: "Обновлен", Width: 16},
	}

	var rows []table.Row
	for _, server := range servers {
		rows = append(rows, table.Row{
			server.ID,
			server.Name,
			FormatStatus(server.Status),
			server.Description,
			server.CreatedAt.Format("02.01.2006 15:04"),
			server.UpdatedAt.Format("02.01.2006 15:04"),
		})
	}

	return NewTableModel(title, columns, rows)
}

// TableModel представляет модель таблицы
type TableModel struct {
	table table.Model
	title string
}

// NewTableModel создает новую модель таблицы
func NewTableModel(title string, columns []table.Column, rows []table.Row) *TableModel {
	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(20),
	)

	return &TableModel{
		table: t,
		title: title,
	}
}

// View отображает таблицу
func (m *TableModel) View() string {
	return fmt.Sprintf("%s\n\n%s", m.title, m.table.View())
}

// Update обновляет модель таблицы
func (m *TableModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

// Init инициализирует модель таблицы
func (m *TableModel) Init() tea.Cmd {
	return nil
}

// GetSelectedRow возвращает выбранную строку
func (m *TableModel) GetSelectedRow() table.Row {
	return m.table.SelectedRow()
}

// ServerPaginatedTableModel представляет модель таблицы с серверной пагинацией
type ServerPaginatedTableModel struct {
	table      table.Model
	title      string
	ctx        context.Context
	limit      int
	offset     int
	total      int
	page       int
	pages      int
	dataLoader func(ctx context.Context, limit, offset int) ([]table.Row, int, error)
	loading    bool
}

// NewServerPaginatedTableModel создает новую модель таблицы с серверной пагинацией
func NewServerPaginatedTableModel(ctx context.Context, title string, columns []table.Column, limit int, dataLoader func(ctx context.Context, limit, offset int) ([]table.Row, int, error)) *ServerPaginatedTableModel {
	t := table.New(
		table.WithColumns(columns),
		table.WithFocused(true),
		table.WithHeight(20),
	)

	model := &ServerPaginatedTableModel{
		table:      t,
		title:      title,
		ctx:        ctx,
		limit:      limit,
		offset:     0,
		total:      0,
		page:       1,
		pages:      0,
		dataLoader: dataLoader,
		loading:    false,
	}

	// Загружаем первую страницу
	model.loadPage(1)

	return model
}

// loadPage загружает указанную страницу
func (m *ServerPaginatedTableModel) loadPage(page int) {
	m.loading = true
	m.page = page
	m.offset = (page - 1) * m.limit

	rows, total, err := m.dataLoader(m.ctx, m.limit, m.offset)
	if err != nil {
		log.Error("Ошибка загрузки данных", "error", err)
		return
	}

	m.total = total
	m.pages = (total + m.limit - 1) / m.limit
	m.table.SetRows(rows)
	m.loading = false
}

// View отображает таблицу
func (m *ServerPaginatedTableModel) View() string {
	if m.loading {
		return fmt.Sprintf("%s\n\nЗагрузка...", m.title)
	}

	title := fmt.Sprintf("%s (страница %d из %d, всего: %d)", m.title, m.page, m.pages, m.total)
	return fmt.Sprintf("%s\n\n%s", title, m.table.View())
}

// Update обновляет модель таблицы
func (m *ServerPaginatedTableModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "h", "left":
			if m.page > 1 {
				m.loadPage(m.page - 1)
			}
		case "l", "right":
			if m.page < m.pages {
				m.loadPage(m.page + 1)
			}
		}
	}

	var cmd tea.Cmd
	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

// Init инициализирует модель таблицы
func (m *ServerPaginatedTableModel) Init() tea.Cmd {
	return nil
}

// GetSelectedRow возвращает выбранную строку
func (m *ServerPaginatedTableModel) GetSelectedRow() table.Row {
	return m.table.SelectedRow()
}
