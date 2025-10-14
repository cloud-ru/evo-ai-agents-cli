package ui

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/cloud-ru/evo-ai-agents-cli/internal/api"
)

// ShowAuthenticationError отображает ошибку аутентификации с ссылкой на документацию
func ShowAuthenticationError(err *api.AuthenticationError) string {
	// Стили для ошибки
	errorStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("1")).
		Bold(true).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("1")).
		Padding(1, 2).
		Margin(1, 0)

	titleStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("1")).
		Bold(true).
		Margin(0, 0, 1, 0)

	messageStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("252")).
		Margin(0, 0, 1, 0)

	linkStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("39")).
		Underline(true).
		Bold(true)

	detailsStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("214")).
		Italic(true).
		Margin(0, 0, 1, 0)

	// Формируем сообщение об ошибке
	result := errorStyle.Render(
		titleStyle.Render("🔐 Ошибка аутентификации") + "\n\n" +
			messageStyle.Render(fmt.Sprintf("Статус: %d", err.StatusCode)) + "\n" +
			messageStyle.Render(err.Message) + "\n\n" +
			detailsStyle.Render(err.Details) + "\n\n" +
			messageStyle.Render("Для решения проблемы ознакомьтесь с документацией:") + "\n" +
			linkStyle.Render("https://cloud.ru/docs/administration/ug/topics/api-ref__authentication") + "\n\n" +
			messageStyle.Render("Возможные причины:") + "\n" +
			messageStyle.Render("• Неверные учетные данные") + "\n" +
			messageStyle.Render("• Истек срок действия токена") + "\n" +
			messageStyle.Render("• Недостаточно прав доступа") + "\n" +
			messageStyle.Render("• Неправильно настроен PROJECT_ID"),
	)

	return result
}

// ShowGenericError отображает общую ошибку
func ShowGenericError(err error) string {
	errorStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("1")).
		Bold(true).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("1")).
		Padding(1, 2).
		Margin(1, 0)

	titleStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("1")).
		Bold(true).
		Margin(0, 0, 1, 0)

	messageStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("252")).
		Margin(0, 0, 1, 0)

	result := errorStyle.Render(
		titleStyle.Render("❌ Ошибка") + "\n\n" +
			messageStyle.Render(err.Error()),
	)

	return result
}

// CheckAndDisplayError проверяет тип ошибки и отображает соответствующее сообщение
func CheckAndDisplayError(err error) string {
	if authErr, ok := err.(*api.AuthenticationError); ok {
		return ShowAuthenticationError(authErr)
	}
	return ShowGenericError(err)
}

// FormatSuccess форматирует сообщение об успехе
func FormatSuccess(message string) string {
	style := lipgloss.NewStyle().
		Foreground(lipgloss.Color("2")).
		Bold(true)

	return style.Render("✅ " + message)
}

// FormatError форматирует сообщение об ошибке
func FormatError(message string) string {
	style := lipgloss.NewStyle().
		Foreground(lipgloss.Color("1")).
		Bold(true)

	return style.Render("❌ " + message)
}

// FormatInfo форматирует информационное сообщение
func FormatInfo(message string) string {
	style := lipgloss.NewStyle().
		Foreground(lipgloss.Color("39")).
		Bold(true)

	return style.Render("ℹ️ " + message)
}

// FormatWarning форматирует предупреждение
func FormatWarning(message string) string {
	style := lipgloss.NewStyle().
		Foreground(lipgloss.Color("214")).
		Bold(true)

	return style.Render("⚠️ " + message)
}
