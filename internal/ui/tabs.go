package ui

import (
	"context"
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/cloudru/ai-agents-cli/internal/api"
	"github.com/cloudru/ai-agents-cli/internal/di"
)

// RenderTabs создает простые читабельные табы
func RenderTabs(tabs []string, activeTab int, content string) string {
	// Простые стили для читабельности
	activeTabStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("39")).
		Background(lipgloss.Color("236")).
		Padding(0, 2)

	inactiveTabStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("240")).
		Padding(0, 2)

	contentStyle := lipgloss.NewStyle().
		Padding(1, 0).
		Margin(1, 0)

	// Рендерим табы в простом формате
	var renderedTabs []string
	for i, tab := range tabs {
		if i == activeTab {
			renderedTabs = append(renderedTabs, activeTabStyle.Render(tab))
		} else {
			renderedTabs = append(renderedTabs, inactiveTabStyle.Render(tab))
		}
	}

	// Объединяем табы
	row := strings.Join(renderedTabs, " | ")

	// Содержимое
	content = contentStyle.Render(content)

	// Результат
	result := strings.Builder{}
	result.WriteString("┌─ Табы: ")
	result.WriteString(row)
	result.WriteString("\n")
	result.WriteString("└─\n\n")
	result.WriteString(content)

	return result.String()
}

// RenderAgentDetailsWithTabs отображает детальную информацию об агенте с табами
func RenderAgentDetailsWithTabs(agent *api.Agent, ctx context.Context, container *di.Container) string {
	// Стили для красивого отображения
	headerStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("39")).
		Align(lipgloss.Center).
		Padding(1, 2).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("39"))

	labelStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("39"))

	valueStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("252"))

	tabStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("39")).
		Margin(0, 0, 1, 0)

	sectionStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("214")).
		Margin(1, 0)

	// Отображаем заголовок
	result := ""
	result += headerStyle.Render("🤖 Информация об агенте")
	result += "\n\n"

	// Показываем все секции сразу (как табы, но статично)

	// Секция 1: Общая информация
	result += sectionStyle.Render("📋 ОБЩАЯ ИНФОРМАЦИЯ")
	result += "\n"
	generalInfo := renderGeneralInfo(agent, ctx, container, labelStyle, valueStyle, tabStyle)
	result += generalInfo
	result += "\n\n"

	// Секция 2: MCP Серверы
	result += sectionStyle.Render("🔌 MCP СЕРВЕРА")
	result += "\n"
	mcpInfo := renderMCPInfo(agent, labelStyle, valueStyle, tabStyle)
	result += mcpInfo
	result += "\n\n"

	// Секция 3: Дополнительные опции
	result += sectionStyle.Render("⚙️ ДОПОЛНИТЕЛЬНЫЕ ОПЦИИ")
	result += "\n"
	optionsInfo := renderOptionsInfo(agent, labelStyle, valueStyle, tabStyle)
	result += optionsInfo

	return result
}

// AgentDetailsTabs представляет интерактивную модель с табами
type AgentDetailsTabs struct {
	agent      *api.Agent
	ctx        context.Context
	container  *di.Container
	tabs       []string
	tabContent []string
	activeTab  int
	width      int
	height     int
}

// NewAgentDetailsTabs создает новую интерактивную модель с табами
func NewAgentDetailsTabs(agent *api.Agent, ctx context.Context, container *di.Container) *AgentDetailsTabs {
	tabs := []string{
		"📋 Общая информация",
		"🔌 MCP Серверы",
		"⚙️ Дополнительные опции",
	}

	model := &AgentDetailsTabs{
		agent:     agent,
		ctx:       ctx,
		container: container,
		tabs:      tabs,
		activeTab: 0,
		width:     80,
		height:    24,
	}

	// Генерируем содержимое табов
	model.generateTabContent()

	return model
}

// Init инициализирует модель
func (m *AgentDetailsTabs) Init() tea.Cmd {
	return nil
}

// Update обновляет модель
func (m *AgentDetailsTabs) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q", "esc":
			return m, tea.Quit
		case "right", "l", "n", "tab":
			m.activeTab = min(m.activeTab+1, len(m.tabs)-1)
			return m, nil
		case "left", "h", "p", "shift+tab":
			m.activeTab = max(m.activeTab-1, 0)
			return m, nil
		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil
	}
	return m, nil
}

// View отображает модель
func (m *AgentDetailsTabs) View() string {
	// Стили для табов (как в примере Bubble Tea)
	tabBorder := func(left, middle, right string) lipgloss.Border {
		border := lipgloss.RoundedBorder()
		border.BottomLeft = left
		border.Bottom = middle
		border.BottomRight = right
		return border
	}

	inactiveTabBorder := tabBorder("┴", "─", "┴")
	activeTabBorder := tabBorder("┘", " ", "└")
	highlightColor := lipgloss.AdaptiveColor{Light: "#874BFD", Dark: "#7D56F4"}

	inactiveTabStyle := lipgloss.NewStyle().
		Border(inactiveTabBorder, true).
		BorderForeground(highlightColor).
		Padding(0, 1)

	activeTabStyle := inactiveTabStyle.Border(activeTabBorder, true)

	windowStyle := lipgloss.NewStyle().
		BorderForeground(highlightColor).
		Padding(2, 0).
		Align(lipgloss.Center).
		Border(lipgloss.NormalBorder()).
		UnsetBorderTop()

	docStyle := lipgloss.NewStyle().Padding(1, 2, 1, 2)

	// Рендерим табы
	var renderedTabs []string
	for i, tab := range m.tabs {
		var style lipgloss.Style
		isFirst, isLast, isActive := i == 0, i == len(m.tabs)-1, i == m.activeTab

		if isActive {
			style = activeTabStyle
		} else {
			style = inactiveTabStyle
		}

		border, _, _, _, _ := style.GetBorder()
		if isFirst && isActive {
			border.BottomLeft = "│"
		} else if isFirst && !isActive {
			border.BottomLeft = "├"
		} else if isLast && isActive {
			border.BottomRight = "│"
		} else if isLast && !isActive {
			border.BottomRight = "┤"
		}
		style = style.Border(border)

		renderedTabs = append(renderedTabs, style.Render(tab))
	}

	row := lipgloss.JoinHorizontal(lipgloss.Top, renderedTabs...)

	// Содержимое активного таба
	content := m.tabContent[m.activeTab]
	windowStyle = windowStyle.Width(lipgloss.Width(row) - windowStyle.GetHorizontalFrameSize())
	content = windowStyle.Render(content)

	// Объединяем табы и содержимое
	result := strings.Builder{}
	result.WriteString(row)
	result.WriteString("\n")
	result.WriteString(content)

	return docStyle.Render(result.String())
}

// generateTabContent генерирует содержимое для всех табов
func (m *AgentDetailsTabs) generateTabContent() {
	// Стили
	labelStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("39"))

	valueStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("252"))

	tabStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("39")).
		Margin(0, 0, 1, 0)

	// Таб 1: Общая информация
	generalInfo := renderGeneralInfo(m.agent, m.ctx, m.container, labelStyle, valueStyle, tabStyle)
	m.tabContent = append(m.tabContent, generalInfo)

	// Таб 2: MCP Серверы
	mcpInfo := renderMCPInfo(m.agent, labelStyle, valueStyle, tabStyle)
	m.tabContent = append(m.tabContent, mcpInfo)

	// Таб 3: Дополнительные опции
	optionsInfo := renderOptionsInfo(m.agent, labelStyle, valueStyle, tabStyle)
	m.tabContent = append(m.tabContent, optionsInfo)
}

// min и max функции для совместимости
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// renderGeneralInfo рендерит общую информацию
func renderGeneralInfo(agent *api.Agent, ctx context.Context, container *di.Container, labelStyle, valueStyle, tabStyle lipgloss.Style) string {
	// Получаем информацию о пользователях
	createdByInfo := getCreatedByInfoForUI(ctx, container, agent.CreatedBy)
	updatedByInfo := getUpdatedByInfoForUI(ctx, container, agent.UpdatedBy)

	result := ""
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

	// Статус
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

	return result
}

// renderMCPInfo рендерит информацию о MCP серверах
func renderMCPInfo(agent *api.Agent, labelStyle, valueStyle, tabStyle lipgloss.Style) string {
	result := ""

	// MCP серверы (новые)
	if len(agent.MCPServers) > 0 {
		result += fmt.Sprintf("%s\n", tabStyle.Render("📡 Подключенные серверы:"))
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
		result += fmt.Sprintf("%s\n", tabStyle.Render("❌ MCP серверы не подключены"))
	}

	// MCP серверы (старые)
	if len(agent.MCPs) > 0 {
		result += fmt.Sprintf("\n%s\n", tabStyle.Render("📡 Старые MCP серверы:"))
		for i, mcp := range agent.MCPs {
			result += fmt.Sprintf("  %d. %s\n", i+1, mcp)
		}
	}

	return result
}

// renderOptionsInfo рендерит дополнительную информацию
func renderOptionsInfo(agent *api.Agent, labelStyle, valueStyle, tabStyle lipgloss.Style) string {
	result := ""

	// Статистика
	result += fmt.Sprintf("%s\n", tabStyle.Render("📊 Статистика:"))
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

	return result
}
