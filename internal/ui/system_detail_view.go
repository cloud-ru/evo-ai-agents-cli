package ui

import (
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

// SystemDetailViewModel представляет интерактивную модель детального просмотра системы агентов
type SystemDetailViewModel struct {
	*SystemDetailModel
}

// NewSystemDetailViewModel создает новую интерактивную модель детального просмотра системы агентов
func NewSystemDetailViewModel(detailModel *SystemDetailModel) *SystemDetailViewModel {
	return &SystemDetailViewModel{
		SystemDetailModel: detailModel,
	}
}

// Init инициализирует модель
func (m SystemDetailViewModel) Init() tea.Cmd {
	return m.Tabs.Init()
}

// Update обновляет модель
func (m SystemDetailViewModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c", "esc":
			return m, tea.Quit
		case "b", "backspace":
			return m, tea.Quit
		}
	}

	// Обновляем табы
	var cmd tea.Cmd
	m.Tabs, cmd = m.Tabs.Update(msg)
	return m, cmd
}

// View отображает модель
func (m SystemDetailViewModel) View() string {
	return m.Render()
}

// Start запускает интерактивный просмотр
func (m SystemDetailViewModel) Start() error {
	if !isInteractive() {
		// Если не интерактивный режим, просто выводим статичную версию
		os.Stdout.WriteString(m.Render())
		return nil
	}

	p := tea.NewProgram(m, tea.WithAltScreen())
	_, err := p.Run()
	return err
}
