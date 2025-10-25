package ui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// CICDType represents the selected CI/CD system type
type CICDType string

const (
	CICDTypeGitLab CICDType = "gitlab"
	CICDTypeGitHub CICDType = "github"
	CICDTypeBoth   CICDType = "both"
	CICDTypeNone   CICDType = "none"
)

// CICDOption represents a CI/CD option in the selector
type CICDOption struct {
	Type        CICDType
	Name        string
	Description string
	Selected    bool
}

// CICDSelectorModel represents the TUI model for CI/CD selection
type CICDSelectorModel struct {
	options  []CICDOption
	cursor   int
	selected CICDType
	quitting bool
}

// NewCICDSelectorModel creates a new CI/CD selector model
func NewCICDSelectorModel() *CICDSelectorModel {
	return &CICDSelectorModel{
		options: []CICDOption{
			{
				Type:        CICDTypeGitLab,
				Name:        "GitLab CI",
				Description: "Использовать только GitLab CI (.gitlab-ci.yml)",
				Selected:    false,
			},
			{
				Type:        CICDTypeGitHub,
				Name:        "GitHub Actions",
				Description: "Использовать только GitHub Actions (.github/workflows)",
				Selected:    false,
			},
			{
				Type:        CICDTypeBoth,
				Name:        "Оба варианта",
				Description: "Создать конфигурации для GitLab CI и GitHub Actions",
				Selected:    false,
			},
			{
				Type:        CICDTypeNone,
				Name:        "Без CI/CD",
				Description: "Не создавать файлы CI/CD",
				Selected:    false,
			},
		},
		cursor:   0,
		selected: "",
		quitting: false,
	}
}

// Init initializes the model
func (m CICDSelectorModel) Init() tea.Cmd {
	return nil
}

// Update handles messages
func (m CICDSelectorModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q", "esc":
			m.quitting = true
			return m, tea.Quit

		case "enter", " ":
			// Select the current option
			m.selected = m.options[m.cursor].Type
			m.quitting = true
			return m, tea.Quit

		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}

		case "down", "j":
			if m.cursor < len(m.options)-1 {
				m.cursor++
			}
		}
	}

	return m, nil
}

// View renders the model
func (m CICDSelectorModel) View() string {
	if m.quitting {
		return ""
	}

	var s strings.Builder

	// Title
	title := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("205")).
		Margin(1, 0).
		Render("🔧 Выберите систему CI/CD")

	s.WriteString(title)
	s.WriteString("\n\n")

	// Instructions
	instructions := lipgloss.NewStyle().
		Foreground(lipgloss.Color("240")).
		Render("↑/↓: навигация • Enter: выбор • q/esc: выход")

	s.WriteString(instructions)
	s.WriteString("\n\n")

	// Options
	for i, option := range m.options {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}

		// Style for selected option
		optionStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("252"))

		if m.cursor == i {
			optionStyle = optionStyle.
				Bold(true).
				Foreground(lipgloss.Color("205"))
		}

		// Render option
		optionText := fmt.Sprintf("%s %s", cursor, option.Name)
		description := lipgloss.NewStyle().
			Foreground(lipgloss.Color("240")).
			Italic(true).
			Render(option.Description)

		s.WriteString(optionStyle.Render(optionText))
		s.WriteString("\n")
		s.WriteString("  " + description)
		s.WriteString("\n\n")
	}

	return s.String()
}

// GetSelected returns the selected CI/CD type
func (m CICDSelectorModel) GetSelected() CICDType {
	return m.selected
}

// RunCICDSelector runs the CI/CD selector TUI
func RunCICDSelector() (CICDType, error) {
	fmt.Printf("DEBUG: isInteractive() = %v\n", isInteractive())

	if !isInteractive() {
		// If not interactive, return default
		fmt.Println("DEBUG: Not interactive, returning default CICDTypeBoth")
		return CICDTypeBoth, nil
	}

	fmt.Println("DEBUG: Starting TUI selector...")
	model := NewCICDSelectorModel()
	program := tea.NewProgram(model, tea.WithAltScreen())

	finalModel, err := program.Run()
	if err != nil {
		fmt.Printf("DEBUG: TUI error: %v\n", err)
		return CICDTypeNone, err
	}

	if selectorModel, ok := finalModel.(CICDSelectorModel); ok {
		selected := selectorModel.GetSelected()
		fmt.Printf("DEBUG: User selected: %v\n", selected)
		return selected, nil
	}

	fmt.Println("DEBUG: Failed to get selection from TUI")
	return CICDTypeNone, fmt.Errorf("failed to get selection")
}

// isInteractive is defined in table.go
