package ui

import (
	"fmt"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// SpinnerModel represents a spinner model
type SpinnerModel struct {
	spinner spinner.Model
	message string
}

// NewSpinnerModel creates a new spinner model
func NewSpinnerModel(message string) *SpinnerModel {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

	return &SpinnerModel{
		spinner: s,
		message: message,
	}
}

// Init initializes the spinner
func (m *SpinnerModel) Init() tea.Cmd {
	return m.spinner.Tick
}

// Update updates the spinner
func (m *SpinnerModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c", "esc":
			return m, tea.Quit
		}
	}

	var cmd tea.Cmd
	m.spinner, cmd = m.spinner.Update(msg)
	return m, cmd
}

// View renders the spinner
func (m *SpinnerModel) View() string {
	style := lipgloss.NewStyle().
		Foreground(lipgloss.Color("252")).
		Margin(1, 0)

	spinnerView := m.spinner.View()
	message := m.message

	return style.Render(fmt.Sprintf("%s %s", spinnerView, message))
}

// RunSpinner runs the spinner for a specified duration
func RunSpinner(message string, duration time.Duration) error {
	model := NewSpinnerModel(message)
	program := tea.NewProgram(model)

	// Run for specified duration
	go func() {
		time.Sleep(duration)
		program.Quit()
	}()

	_, err := program.Run()
	return err
}

// ShowLoadingSpinner shows a loading spinner with a message
func ShowLoadingSpinner(message string) string {
	model := NewSpinnerModel(message)
	return model.View()
}

// LoadingMessage represents a loading message
type LoadingMessage struct {
	Message string
}

// ShowLoadingMessage shows a simple loading message without animation
func ShowLoadingMessage(message string) string {
	style := lipgloss.NewStyle().
		Foreground(lipgloss.Color("205")).
		Bold(true).
		Margin(1, 0)

	return style.Render(fmt.Sprintf("⏳ %s", message))
}
