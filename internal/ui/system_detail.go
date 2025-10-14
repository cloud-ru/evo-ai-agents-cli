package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/cloud-ru/evo-ai-agents-cli/internal/api"
)

// SystemDetailModel представляет детальный просмотр системы агентов с табами
type SystemDetailModel struct {
	System *api.AgentSystem
	Tabs   *TabModel
}

// NewSystemDetailModel создает новую модель детального просмотра системы агентов
func NewSystemDetailModel(system *api.AgentSystem) *SystemDetailModel {
	tabs := []string{
		"📋 Общая информация",
		"🤖 Агенты",
		"🧠 Оркестратор",
		"⚙️ Опции",
		"🔗 Интеграция",
		"📊 Статистика",
		"🔄 Обновления",
		"📈 Масштабирование",
	}

	content := []string{
		formatSystemGeneralInfo(system),
		formatSystemAgents(system),
		formatSystemOrchestrator(system),
		formatSystemOptions(system),
		formatSystemIntegration(system),
		formatSystemStatistics(system),
		formatSystemUpdates(system),
		formatSystemScaling(system),
	}

	return &SystemDetailModel{
		System: system,
		Tabs:   NewTabModel(tabs, content),
	}
}

// Render отображает детальную информацию о системе агентов
func (m *SystemDetailModel) Render() string {
	header := fmt.Sprintf("🏗️ Детальная информация: %s", m.System.Name)
	headerStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("39")).
		MarginBottom(1)

	return headerStyle.Render(header) + "\n" + m.Tabs.Render()
}

// formatSystemGeneralInfo форматирует общую информацию системы агентов
func formatSystemGeneralInfo(system *api.AgentSystem) string {
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
	leftInfo.WriteString(fmt.Sprintf("🆔 ID: %s\n", system.ID))
	leftInfo.WriteString(fmt.Sprintf("📝 Название: %s\n", system.Name))
	leftInfo.WriteString(fmt.Sprintf("📄 Описание: %s\n", getDescription(system.Description)))
	leftInfo.WriteString(fmt.Sprintf("📊 Статус: %s\n", FormatSystemStatus(system.Status)))

	// Добавляем информацию о Project ID
	if system.ProjectID != "" {
		leftInfo.WriteString(fmt.Sprintf("🏢 Проект: %s\n", system.ProjectID))
	}

	// Добавляем информацию о типе инстанса
	if system.InstanceType.ID != "" {
		leftInfo.WriteString(fmt.Sprintf("💻 Инстанс: %s (%s)\n", system.InstanceType.Name, system.InstanceType.SKUCode))
	}

	// Правая колонка - временные метки и пользователи
	rightInfo := strings.Builder{}
	rightInfo.WriteString(fmt.Sprintf("📅 Создан: %s\n", system.CreatedAt.Format("02.01.2006 15:04:05")))
	rightInfo.WriteString(fmt.Sprintf("🔄 Обновлен: %s\n", system.UpdatedAt.Format("02.01.2006 15:04:05")))

	if system.CreatedBy != "" {
		rightInfo.WriteString(fmt.Sprintf("👤 Создал: %s\n", system.CreatedBy))
	}
	if system.UpdatedBy != "" {
		rightInfo.WriteString(fmt.Sprintf("✏️ Изменил: %s\n", system.UpdatedBy))
	}

	// Центральная секция - URL и дополнительная информация
	centerInfo := strings.Builder{}
	if system.PublicURL != "" {
		centerInfo.WriteString(fmt.Sprintf("🌐 Публичный URL: %s\n", system.PublicURL))
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

// formatSystemAgents форматирует информацию об агентах в системе
func formatSystemAgents(system *api.AgentSystem) string {
	var info strings.Builder

	info.WriteString("🤖 Агенты в системе:\n\n")

	if len(system.Agents) > 0 {
		for i, agent := range system.Agents {
			info.WriteString(fmt.Sprintf("🤖 %d. %s\n", i+1, agent.Name))
			info.WriteString(fmt.Sprintf("   🆔 ID: %s\n", agent.ID))
			info.WriteString(fmt.Sprintf("   📊 Статус: %s\n", FormatAgentStatus(agent.Status)))
			info.WriteString("\n")
		}
	} else {
		info.WriteString("Нет агентов в системе\n")
	}

	return info.String()
}

// formatSystemOrchestrator форматирует настройки оркестратора
func formatSystemOrchestrator(system *api.AgentSystem) string {
	var info strings.Builder

	info.WriteString("🧠 Настройки оркестратора:\n\n")

	if system.OrchestratorOptions != nil {
		for key, value := range system.OrchestratorOptions {
			info.WriteString(fmt.Sprintf("⚙️ %s: %v\n", key, value))
		}
	} else {
		info.WriteString("Нет настроек оркестратора\n")
	}

	return info.String()
}

// formatSystemOptions форматирует опции системы
func formatSystemOptions(system *api.AgentSystem) string {
	var info strings.Builder

	info.WriteString("⚙️ Опции системы:\n\n")

	if system.Options != nil {
		for key, value := range system.Options {
			info.WriteString(fmt.Sprintf("🔧 %s: %v\n", key, value))
		}
	} else {
		info.WriteString("Нет дополнительных опций\n")
	}

	return info.String()
}

// formatSystemIntegration форматирует настройки интеграции
func formatSystemIntegration(system *api.AgentSystem) string {
	var info strings.Builder

	info.WriteString("🔗 Настройки интеграции:\n\n")

	if system.IntegrationOptions != nil {
		for key, value := range system.IntegrationOptions {
			info.WriteString(fmt.Sprintf("⚙️ %s: %v\n", key, value))
		}
	} else {
		info.WriteString("Нет настроек интеграции\n")
	}

	return info.String()
}

// formatSystemStatistics форматирует статистику системы
func formatSystemStatistics(system *api.AgentSystem) string {
	var info strings.Builder

	info.WriteString("📊 Статистика системы:\n\n")
	info.WriteString(fmt.Sprintf("🤖 Количество агентов: %d\n", len(system.Agents)))
	info.WriteString(fmt.Sprintf("📅 Время работы: %s\n", system.UpdatedAt.Sub(system.CreatedAt).String()))

	return info.String()
}

// formatSystemUpdates форматирует информацию об обновлениях
func formatSystemUpdates(system *api.AgentSystem) string {
	var info strings.Builder

	info.WriteString("🔄 Информация об обновлениях:\n\n")
	info.WriteString(fmt.Sprintf("📅 Последнее обновление: %s\n", system.UpdatedAt.Format("02.01.2006 15:04:05")))

	if system.UpdatedBy != "" {
		info.WriteString(fmt.Sprintf("👤 Обновил: %s\n", system.UpdatedBy))
	}

	return info.String()
}

// formatSystemScaling форматирует настройки масштабирования
func formatSystemScaling(system *api.AgentSystem) string {
	var info strings.Builder

	info.WriteString("📈 Настройки масштабирования:\n\n")

	// Здесь можно добавить логику для парсинга настроек масштабирования из Options
	info.WriteString("Настройки масштабирования будут добавлены в будущих версиях\n")

	return info.String()
}

// FormatSystemStatus форматирует статус системы агентов
func FormatSystemStatus(status string) string {
	switch status {
	case "ACTIVE":
		return "🟢 Активна"
	case "SUSPENDED":
		return "🟡 Приостановлена"
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
