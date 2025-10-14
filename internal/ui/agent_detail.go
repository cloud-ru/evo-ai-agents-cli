package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/cloud-ru/evo-ai-agents-cli/internal/api"
)

// AgentDetailModel представляет детальный просмотр агента с табами
type AgentDetailModel struct {
	Agent *api.Agent
	Tabs  *TabModel
}

// NewAgentDetailModel создает новую модель детального просмотра агента
func NewAgentDetailModel(agent *api.Agent) *AgentDetailModel {
	tabs := []string{
		"📋 Общая информация",
		"🧠 LLM настройки",
		"📈 Scaling",
		"🔐 Аутентификация",
		"📊 Логирование",
		"🔄 Автообновление",
		"🐳 Образ",
		"⚙️ Опции",
	}

	content := []string{
		formatGeneralInfo(agent),
		formatLLMSettings(agent),
		formatScalingSettings(agent),
		formatAuthSettings(agent),
		formatLoggingSettings(agent),
		formatAutoUpdateSettings(agent),
		formatImageSource(agent),
		formatOptions(agent),
	}

	return &AgentDetailModel{
		Agent: agent,
		Tabs:  NewTabModel(tabs, content),
	}
}

// Render отображает детальную информацию об агенте
func (m *AgentDetailModel) Render() string {
	header := fmt.Sprintf("🤖 Детальная информация: %s", m.Agent.Name)
	headerStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("39")).
		MarginBottom(1)

	return headerStyle.Render(header) + "\n" + m.Tabs.Render()
}

// formatGeneralInfo форматирует общую информацию с гармоничным layout
func formatGeneralInfo(agent *api.Agent) string {
	// Стили для разных секций
	leftStyle := lipgloss.NewStyle().
		Width(50).
		Align(lipgloss.Left)

	rightStyle := lipgloss.NewStyle().
		Width(50).
		Align(lipgloss.Right)

	centerStyle := lipgloss.NewStyle().
		Width(100).
		Align(lipgloss.Center).
		MarginBottom(1)

	// Левая колонка - основная информация
	leftInfo := strings.Builder{}
	leftInfo.WriteString(fmt.Sprintf("🆔 ID: %s\n", agent.ID))
	leftInfo.WriteString(fmt.Sprintf("📝 Название: %s\n", agent.Name))
	leftInfo.WriteString(fmt.Sprintf("📄 Описание: %s\n", getDescription(agent.Description)))
	leftInfo.WriteString(fmt.Sprintf("📊 Статус: %s\n", FormatStatus(agent.Status)))
	leftInfo.WriteString(fmt.Sprintf("🏷️ Тип: %s\n", FormatAgentType(agent.AgentType)))

	// Добавляем информацию о Project ID
	if agent.ProjectID != "" {
		leftInfo.WriteString(fmt.Sprintf("🏢 Проект: %s\n", agent.ProjectID))
	}

	// Добавляем информацию о типе инстанса
	if agent.InstanceType.ID != "" {
		leftInfo.WriteString(fmt.Sprintf("💻 Инстанс: %s (%s)\n", agent.InstanceType.Name, agent.InstanceType.SKUCode))
	}

	// Добавляем информацию о MCP серверах
	if len(agent.MCPServers) > 0 {
		leftInfo.WriteString(fmt.Sprintf("🔌 MCP серверы: %d\n", len(agent.MCPServers)))
		for i, mcp := range agent.MCPServers {
			if i < 3 { // Показываем только первые 3
				leftInfo.WriteString(fmt.Sprintf("  • %s (%s)\n", mcp.Name, mcp.Status))
			}
		}
		if len(agent.MCPServers) > 3 {
			leftInfo.WriteString(fmt.Sprintf("  ... и еще %d\n", len(agent.MCPServers)-3))
		}
	}

	// Правая колонка - временные метки и пользователи
	rightInfo := strings.Builder{}
	rightInfo.WriteString(fmt.Sprintf("📅 Создан: %s\n", agent.CreatedAt.Time.Format("02.01.2006 15:04:05")))
	rightInfo.WriteString(fmt.Sprintf("🔄 Обновлен: %s\n", agent.UpdatedAt.Time.Format("02.01.2006 15:04:05")))

	if agent.CreatedBy != "" {
		rightInfo.WriteString(fmt.Sprintf("👤 Создал: %s\n", agent.CreatedBy))
	}
	if agent.UpdatedBy != "" {
		rightInfo.WriteString(fmt.Sprintf("✏️ Изменил: %s\n", agent.UpdatedBy))
	}

	// Центральная секция - URL и дополнительная информация
	centerInfo := strings.Builder{}
	if agent.PublicURL != "" {
		centerInfo.WriteString(fmt.Sprintf("🌐 Публичный URL: %s\n", agent.PublicURL))
	}
	if agent.ArizePhoenixPublicURL != "" {
		centerInfo.WriteString(fmt.Sprintf("📊 Arize Phoenix URL: %s\n", agent.ArizePhoenixPublicURL))
	}

	// Объединяем все секции
	var result strings.Builder
	result.WriteString(centerStyle.Render(centerInfo.String()))

	// Создаем две колонки
	leftContent := leftStyle.Render(leftInfo.String())
	rightContent := rightStyle.Render(rightInfo.String())

	result.WriteString(lipgloss.JoinHorizontal(lipgloss.Top, leftContent, rightContent))

	return result.String()
}

// formatLLMSettings форматирует настройки LLM
func formatLLMSettings(agent *api.Agent) string {
	var info strings.Builder

	info.WriteString("🧠 Настройки языковой модели:\n\n")

	// LLM настройки из Options
	if llmOpts, ok := agent.Options["llm"].(map[string]interface{}); ok {
		if foundationModels, ok := llmOpts["foundationModels"].(map[string]interface{}); ok {
			if modelName, ok := foundationModels["modelName"].(string); ok {
				info.WriteString(fmt.Sprintf("📝 Модель: %s\n", modelName))
			}
			if gcInstanceId, ok := foundationModels["gcInstanceId"].(string); ok && gcInstanceId != "" {
				info.WriteString(fmt.Sprintf("🆔 GC Instance ID: %s\n", gcInstanceId))
			}
		}
	}

	// System prompt
	if systemPrompt, ok := agent.Options["systemPrompt"].(string); ok && systemPrompt != "" {
		info.WriteString(fmt.Sprintf("💬 System Prompt: %s\n", systemPrompt))
	}

	return info.String()
}

// formatScalingSettings форматирует настройки масштабирования
func formatScalingSettings(agent *api.Agent) string {
	var info strings.Builder

	info.WriteString("📈 Настройки масштабирования:\n\n")

	if scaling, ok := agent.Options["scaling"].(map[string]interface{}); ok {
		if minScale, ok := scaling["minScale"].(float64); ok {
			info.WriteString(fmt.Sprintf("📉 Min Scale: %.0f\n", minScale))
		}
		if maxScale, ok := scaling["maxScale"].(float64); ok {
			info.WriteString(fmt.Sprintf("📈 Max Scale: %.0f\n", maxScale))
		}
		if isKeepAlive, ok := scaling["isKeepAlive"].(bool); ok {
			info.WriteString(fmt.Sprintf("💤 Keep Alive: %t\n", isKeepAlive))
		}
		if isScaleUpAllSystem, ok := scaling["isScaleUpAllSystem"].(bool); ok {
			info.WriteString(fmt.Sprintf("🔄 Scale Up All System: %t\n", isScaleUpAllSystem))
		}
	}

	return info.String()
}

// formatAuthSettings форматирует настройки аутентификации
func formatAuthSettings(agent *api.Agent) string {
	var info strings.Builder

	info.WriteString("🔐 Настройки аутентификации:\n\n")

	if authOptions, ok := agent.Options["authOptions"].(map[string]interface{}); ok {
		if isEnabled, ok := authOptions["isEnabled"].(bool); ok {
			info.WriteString(fmt.Sprintf("✅ Включена: %t\n", isEnabled))
		}
		if authType, ok := authOptions["type"].(string); ok {
			info.WriteString(fmt.Sprintf("🔑 Тип: %s\n", authType))
		}
		if serviceAccountId, ok := authOptions["serviceAccountId"].(string); ok {
			info.WriteString(fmt.Sprintf("👤 Service Account ID: %s\n", serviceAccountId))
		}
	}

	return info.String()
}

// formatLoggingSettings форматирует настройки логирования
func formatLoggingSettings(agent *api.Agent) string {
	var info strings.Builder

	info.WriteString("📊 Настройки логирования:\n\n")

	if logging, ok := agent.Options["logging"].(map[string]interface{}); ok {
		if isEnabled, ok := logging["isEnabledLogging"].(bool); ok {
			info.WriteString(fmt.Sprintf("✅ Включено: %t\n", isEnabled))
		}
		if logGroupId, ok := logging["logGroupId"].(string); ok {
			info.WriteString(fmt.Sprintf("📁 Log Group ID: %s\n", logGroupId))
		}
	}

	return info.String()
}

// formatAutoUpdateSettings форматирует настройки автообновления
func formatAutoUpdateSettings(agent *api.Agent) string {
	var info strings.Builder

	info.WriteString("🔄 Настройки автообновления:\n\n")

	if autoUpdate, ok := agent.Options["autoUpdateOptions"].(map[string]interface{}); ok {
		if isEnabled, ok := autoUpdate["isEnabled"].(bool); ok {
			info.WriteString(fmt.Sprintf("✅ Включено: %t\n", isEnabled))
		}
	}

	return info.String()
}

// formatImageSource форматирует информацию об образе
func formatImageSource(agent *api.Agent) string {
	var info strings.Builder

	info.WriteString("🐳 Источник образа:\n\n")

	if marketplaceAgentId, ok := agent.Options["marketplaceAgentId"].(string); ok {
		info.WriteString(fmt.Sprintf("🏪 Marketplace Agent ID: %s\n", marketplaceAgentId))
	}

	return info.String()
}

// formatOptions форматирует все опции
func formatOptions(agent *api.Agent) string {
	var info strings.Builder

	info.WriteString("⚙️ Все опции:\n\n")

	for key, value := range agent.Options {
		info.WriteString(fmt.Sprintf("🔧 %s: %v\n", key, value))
	}

	return info.String()
}

// getDescription возвращает описание или прочерк
func getDescription(desc string) string {
	if desc == "" {
		return "—"
	}
	return desc
}

// FormatAgentStatus форматирует статус агента
func FormatAgentStatus(status string) string {
	switch status {
	case "ACTIVE":
		return "🟢 Активен"
	case "SUSPENDED":
		return "🟡 Приостановлен"
	case "ERROR":
		return "🔴 Ошибка"
	case "PENDING":
		return "⏳ Ожидает"
	case "CREATING":
		return "🔨 Создается"
	case "DELETING":
		return "🗑️ Удаляется"
	default:
		return "⚪ " + status
	}
}

// FormatAgentType форматирует тип агента
func FormatAgentType(agentType string) string {
	switch agentType {
	case "marketplace":
		return "Из маркетплейса"
	case "custom":
		return "Пользовательский"
	case "":
		return "Не указан"
	default:
		return agentType
	}
}
