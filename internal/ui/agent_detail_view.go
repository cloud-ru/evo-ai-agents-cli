package ui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/cloud-ru/evo-ai-agents-cli/internal/api"
)

// AgentDetailViewModel представляет интерактивную модель детального просмотра
type AgentDetailViewModel struct {
	agent    *api.Agent
	detail   *AgentDetailModel
	quitting bool
}

// NewAgentDetailViewModel создает новую интерактивную модель
func NewAgentDetailViewModel(agent *api.Agent) *AgentDetailViewModel {
	return &AgentDetailViewModel{
		agent:  agent,
		detail: NewAgentDetailModel(agent),
	}
}

// Init инициализирует модель
func (m AgentDetailViewModel) Init() tea.Cmd {
	return nil
}

// Update обновляет модель
func (m AgentDetailViewModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q", "esc":
			m.quitting = true
			return m, tea.Quit
		case "right", "l", "n", "tab":
			m.detail.Tabs.NextTab()
			return m, nil
		case "left", "h", "p", "shift+tab":
			m.detail.Tabs.PrevTab()
			return m, nil
		case "1", "2", "3", "4", "5", "6", "7", "8":
			// Переключение по номерам табов
			tabIndex := int(msg.String()[0] - '1')
			if tabIndex >= 0 && tabIndex < len(m.detail.Tabs.Tabs) {
				m.detail.Tabs.SetActiveTab(tabIndex)
			}
			return m, nil
		}
	}
	return m, nil
}

// View отображает модель
func (m AgentDetailViewModel) View() string {
	if m.quitting {
		return ""
	}

	help := "\n\n" + lipgloss.NewStyle().
		Foreground(lipgloss.Color("240")).
		Render("←/→ или h/l: переключение табов • 1-8: быстрый переход • q/esc: выход")

	return m.detail.Render() + help
}

// Start запускает интерактивный просмотр
func (m *AgentDetailViewModel) Start() error {
	if !isInteractive() {
		// Если не интерактивный режим, просто выводим статичную версию
		fmt.Println(m.detail.Render())
		return nil
	}

	program := tea.NewProgram(m)

	if _, err := program.Run(); err != nil {
		return err
	}

	return nil
}

// ShowAgentDetail отображает детальную информацию об агенте
func ShowAgentDetail(agent *api.Agent) error {
	// Проверяем, что мы в интерактивном режиме
	if !isInteractive() {
		// Если не интерактивный режим, просто показываем как текст
		detail := NewAgentDetailModel(agent)
		fmt.Println(detail.Render())
		return nil
	}

	model := NewAgentDetailViewModel(agent)
	program := tea.NewProgram(model)

	if _, err := program.Run(); err != nil {
		return err
	}

	return nil
}
