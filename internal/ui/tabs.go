package ui

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// TabModel представляет модель с табами
type TabModel struct {
	Tabs       []string
	TabContent []string
	ActiveTab  int
}

// NewTabModel создает новую модель с табами
func NewTabModel(tabs []string, content []string) *TabModel {
	return &TabModel{
		Tabs:       tabs,
		TabContent: content,
		ActiveTab:  0,
	}
}

// Init инициализирует модель
func (m *TabModel) Init() tea.Cmd {
	return nil
}

// Update обновляет модель
func (m *TabModel) Update(msg tea.Msg) (*TabModel, tea.Cmd) {
	return m, nil
}

// tabBorderWithBottom создает границу для табов
func tabBorderWithBottom(left, middle, right string) lipgloss.Border {
	border := lipgloss.RoundedBorder()
	border.BottomLeft = left
	border.Bottom = middle
	border.BottomRight = right
	return border
}

var (
	inactiveTabBorder = tabBorderWithBottom("┴", "─", "┴")
	activeTabBorder   = tabBorderWithBottom("┘", " ", "└")
	docStyle          = lipgloss.NewStyle().Padding(1, 2, 1, 2)
	highlightColor    = lipgloss.AdaptiveColor{Light: "#874BFD", Dark: "#7D56F4"}
	inactiveTabStyle  = lipgloss.NewStyle().
				Border(inactiveTabBorder, true).
				BorderForeground(highlightColor).
				Padding(0, 1)
	activeTabStyle = inactiveTabStyle.Border(activeTabBorder, true)
	windowStyle    = lipgloss.NewStyle().
			BorderForeground(highlightColor).
			Padding(2, 0).
			Align(lipgloss.Center).
			Border(lipgloss.NormalBorder()).
			UnsetBorderTop()
)

// Render отображает табы
func (m *TabModel) Render() string {
	doc := strings.Builder{}
	var renderedTabs []string

	for i, t := range m.Tabs {
		var style lipgloss.Style
		isFirst, isLast, isActive := i == 0, i == len(m.Tabs)-1, i == m.ActiveTab

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
		renderedTabs = append(renderedTabs, style.Render(t))
	}

	row := lipgloss.JoinHorizontal(lipgloss.Top, renderedTabs...)
	doc.WriteString(row)
	doc.WriteString("\n")
	doc.WriteString(windowStyle.Width((lipgloss.Width(row) - windowStyle.GetHorizontalFrameSize())).Render(m.TabContent[m.ActiveTab]))

	return docStyle.Render(doc.String())
}

// SetActiveTab устанавливает активный таб
func (m *TabModel) SetActiveTab(index int) {
	if index >= 0 && index < len(m.Tabs) {
		m.ActiveTab = index
	}
}

// NextTab переключает на следующий таб
func (m *TabModel) NextTab() {
	m.ActiveTab = (m.ActiveTab + 1) % len(m.Tabs)
}

// PrevTab переключает на предыдущий таб
func (m *TabModel) PrevTab() {
	m.ActiveTab = (m.ActiveTab - 1 + len(m.Tabs)) % len(m.Tabs)
}

// FormatAgentDetails форматирует детали агента в табы
func FormatAgentDetails(agent interface{}) *TabModel {
	// Здесь будет логика форматирования деталей агента
	// Пока создаем пример
	tabs := []string{"Общая информация", "LLM настройки", "Scaling", "Аутентификация", "Логирование"}

	content := []string{
		"Общая информация об агенте...",
		"Настройки языковой модели...",
		"Настройки масштабирования...",
		"Настройки аутентификации...",
		"Настройки логирования...",
	}

	return NewTabModel(tabs, content)
}
