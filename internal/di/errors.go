package di

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// ErrorHandler обрабатывает ошибки DI контейнера с красивым выводом
type ErrorHandler struct{}

// NewErrorHandler создает новый обработчик ошибок
func NewErrorHandler() *ErrorHandler {
	return &ErrorHandler{}
}

// HandleConfigError обрабатывает ошибки конфигурации
func (h *ErrorHandler) HandleConfigError(err error) error {
	if err == nil {
		return nil
	}

	// Проверяем, является ли ошибка связанной с переменными окружения
	if strings.Contains(err.Error(), "environment variable is required") {
		return h.formatEnvironmentError(err)
	}

	return h.formatGenericError("Configuration", err)
}

// HandleAuthError обрабатывает ошибки аутентификации
func (h *ErrorHandler) HandleAuthError(err error) error {
	if err == nil {
		return nil
	}

	return h.formatGenericError("Authentication", err)
}

// HandleAPIError обрабатывает ошибки API
func (h *ErrorHandler) HandleAPIError(err error) error {
	if err == nil {
		return nil
	}

	return h.formatGenericError("API Client", err)
}

// formatEnvironmentError форматирует ошибки переменных окружения
func (h *ErrorHandler) formatEnvironmentError(err error) error {
	errorStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("1")).
		Border(lipgloss.RoundedBorder()).
		Padding(1, 2)

	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("3")).
		MarginBottom(1)

	helpStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("8")).
		MarginTop(1)

	title := titleStyle.Render("🔧 Configuration Error")
	message := errorStyle.Render(err.Error())

	help := helpStyle.Render(`
Для работы CLI необходимо настроить переменные окружения:

1. Создайте файл .env в корне проекта:
   cp env.example .env

2. Заполните необходимые переменные:
   IAM_KEY_ID=your_iam_key_id
   IAM_SECRET=your_iam_secret
   PROJECT_ID=your_project_id
   CUSTOMER_ID=your_customer_id

3. Или экспортируйте переменные в текущей сессии:
   export IAM_KEY_ID=your_iam_key_id
   export IAM_SECRET=your_iam_secret
   export PROJECT_ID=your_project_id
   export CUSTOMER_ID=your_customer_id

4. Проверьте настройки:
   ai-agents-cli --help
`)

	return fmt.Errorf("%s\n\n%s\n\n%s", title, message, help)
}

// formatGenericError форматирует общие ошибки
func (h *ErrorHandler) formatGenericError(service string, err error) error {
	errorStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("1")).
		Border(lipgloss.RoundedBorder()).
		Padding(1, 2)

	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("3")).
		MarginBottom(1)

	title := titleStyle.Render(fmt.Sprintf("❌ %s Error", service))
	message := errorStyle.Render(err.Error())

	return fmt.Errorf("%s\n\n%s", title, message)
}

// HandleContainerError обрабатывает ошибки контейнера
func (h *ErrorHandler) HandleContainerError(err error) error {
	if err == nil {
		return nil
	}

	// Проверяем тип ошибки
	if strings.Contains(err.Error(), "environment variable is required") {
		return h.HandleConfigError(err)
	}

	if strings.Contains(err.Error(), "auth") || strings.Contains(err.Error(), "IAM") {
		return h.HandleAuthError(err)
	}

	if strings.Contains(err.Error(), "API") {
		return h.HandleAPIError(err)
	}

	return h.formatGenericError("Container", err)
}
