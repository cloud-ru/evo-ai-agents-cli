package errors

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// UIStyles содержит стили для отображения ошибок в UI
var UIStyles = struct {
	ErrorBox     lipgloss.Style
	ErrorTitle   lipgloss.Style
	ErrorMessage lipgloss.Style
	ErrorDetails lipgloss.Style
	ErrorCode    lipgloss.Style
	ErrorContext lipgloss.Style
	SuccessBox   lipgloss.Style
	WarningBox   lipgloss.Style
	InfoBox      lipgloss.Style
}{
	ErrorBox: lipgloss.NewStyle().
		Foreground(lipgloss.Color("1")).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("1")).
		Padding(1, 2).
		Margin(1, 0),

	ErrorTitle: lipgloss.NewStyle().
		Foreground(lipgloss.Color("1")).
		Bold(true).
		Margin(0, 0, 1, 0),

	ErrorMessage: lipgloss.NewStyle().
		Foreground(lipgloss.Color("252")).
		Margin(0, 0, 1, 0),

	ErrorDetails: lipgloss.NewStyle().
		Foreground(lipgloss.Color("245")).
		Italic(true).
		Margin(0, 0, 1, 0),

	ErrorCode: lipgloss.NewStyle().
		Foreground(lipgloss.Color("8")).
		Background(lipgloss.Color("236")).
		Padding(0, 1).
		Margin(0, 0, 1, 0),

	ErrorContext: lipgloss.NewStyle().
		Foreground(lipgloss.Color("244")).
		Margin(0, 0, 0, 2),

	SuccessBox: lipgloss.NewStyle().
		Foreground(lipgloss.Color("2")).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("2")).
		Padding(1, 2).
		Margin(1, 0),

	WarningBox: lipgloss.NewStyle().
		Foreground(lipgloss.Color("214")).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("214")).
		Padding(1, 2).
		Margin(1, 0),

	InfoBox: lipgloss.NewStyle().
		Foreground(lipgloss.Color("39")).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("39")).
		Padding(1, 2).
		Margin(1, 0),
}

// ErrorIcon возвращает иконку для типа ошибки
func ErrorIcon(errorType ErrorType) string {
	switch errorType {
	case ErrorTypeValidation:
		return "⚠️"
	case ErrorTypeConfiguration:
		return "⚙️"
	case ErrorTypeAuthentication:
		return "🔐"
	case ErrorTypeAPI:
		return "🌐"
	case ErrorTypeNetwork:
		return "📡"
	case ErrorTypeFileSystem:
		return "📁"
	case ErrorTypeTemplate:
		return "📄"
	case ErrorTypeUser:
		return "👤"
	case ErrorTypeSystem:
		return "💥"
	default:
		return "❌"
	}
}

// SeverityIcon возвращает иконку для серьезности ошибки
func SeverityIcon(severity Severity) string {
	switch severity {
	case SeverityLow:
		return "ℹ️"
	case SeverityMedium:
		return "⚠️"
	case SeverityHigh:
		return "🚨"
	case SeverityCritical:
		return "💥"
	default:
		return "❌"
	}
}

// FormatError форматирует ошибку для отображения в UI
func FormatError(err error) string {
	if err == nil {
		return ""
	}

	// Если это AppError, используем структурированное отображение
	if appErr, ok := err.(*AppError); ok {
		return formatAppError(appErr)
	}

	// Для обычных ошибок используем простое отображение
	return formatGenericError(err)
}

// FormatSimpleError форматирует ошибку в простом виде без рамок
func FormatSimpleError(err error) string {
	if err == nil {
		return ""
	}

	// Если это AppError, используем простое отображение
	if appErr, ok := err.(*AppError); ok {
		return formatSimpleAppError(appErr)
	}

	// Для обычных ошибок используем простое отображение
	return formatSimpleGenericError(err)
}

// formatAppError форматирует структурированную ошибку
func formatAppError(err *AppError) string {
	var parts []string

	// Заголовок с иконками
	title := fmt.Sprintf("%s %s %s",
		ErrorIcon(err.Type),
		SeverityIcon(err.Severity),
		err.Message)
	parts = append(parts, UIStyles.ErrorTitle.Render(title))

	// Код ошибки
	if err.Code != "" {
		parts = append(parts, UIStyles.ErrorCode.Render(err.Code))
	}

	// Детали
	if err.Details != "" {
		parts = append(parts, UIStyles.ErrorDetails.Render(err.Details))
	}

	// Контекст
	if len(err.Context) > 0 {
		contextParts := []string{"Контекст:"}
		for key, value := range err.Context {
			contextParts = append(contextParts, fmt.Sprintf("  %s: %v", key, value))
		}
		parts = append(parts, UIStyles.ErrorContext.Render(strings.Join(contextParts, "\n")))
	}

	// Оригинальная ошибка
	if err.Original != nil {
		parts = append(parts, UIStyles.ErrorDetails.Render(fmt.Sprintf("Причина: %v", err.Original)))
	}

	content := strings.Join(parts, "\n")
	return UIStyles.ErrorBox.Render(content)
}

// formatGenericError форматирует обычную ошибку
func formatGenericError(err error) string {
	title := UIStyles.ErrorTitle.Render("❌ Ошибка")
	message := UIStyles.ErrorMessage.Render(err.Error())

	content := fmt.Sprintf("%s\n\n%s", title, message)
	return UIStyles.ErrorBox.Render(content)
}

// formatSimpleAppError форматирует структурированную ошибку в простом виде
func formatSimpleAppError(err *AppError) string {
	var parts []string

	// Заголовок с иконками
	title := fmt.Sprintf("%s %s %s",
		ErrorIcon(err.Type),
		SeverityIcon(err.Severity),
		err.Message)
	parts = append(parts, UIStyles.ErrorTitle.Render(title))

	// Код ошибки
	if err.Code != "" {
		parts = append(parts, UIStyles.ErrorCode.Render(err.Code))
	}

	// Детали
	if err.Details != "" {
		parts = append(parts, UIStyles.ErrorDetails.Render(err.Details))
	}

	// Оригинальная ошибка
	if err.Original != nil {
		parts = append(parts, UIStyles.ErrorDetails.Render(fmt.Sprintf("Причина: %v", err.Original)))
	}

	// Подсказки
	if len(err.Suggestions) > 0 {
		parts = append(parts, UIStyles.ErrorTitle.Render("💡 Что можно сделать:"))
		for i, suggestion := range err.Suggestions {
			// Красиво форматируем ссылки на документацию
			formattedSuggestion := formatSuggestion(suggestion)
			parts = append(parts, UIStyles.ErrorDetails.Render(fmt.Sprintf("  %d. %s", i+1, formattedSuggestion)))
		}
	}

	return strings.Join(parts, "\n")
}

// formatSimpleGenericError форматирует обычную ошибку в простом виде
func formatSimpleGenericError(err error) string {
	title := UIStyles.ErrorTitle.Render("❌ Ошибка")
	message := UIStyles.ErrorMessage.Render(err.Error())

	return fmt.Sprintf("%s\n\n%s", title, message)
}

// FormatPlainError форматирует ошибку в простом текстовом виде без стилей
func FormatPlainError(err error) string {
	if err == nil {
		return ""
	}

	// Если это AppError, используем структурированное отображение
	if appErr, ok := err.(*AppError); ok {
		return formatPlainAppError(appErr)
	}

	// Для обычных ошибок используем простое отображение
	return fmt.Sprintf("❌ Ошибка: %s", err.Error())
}

// formatPlainAppError форматирует структурированную ошибку в простом текстовом виде
func formatPlainAppError(err *AppError) string {
	var parts []string

	// Заголовок с иконками
	title := fmt.Sprintf("%s %s %s",
		ErrorIcon(err.Type),
		SeverityIcon(err.Severity),
		err.Message)
	parts = append(parts, title)

	// Код ошибки
	if err.Code != "" {
		parts = append(parts, fmt.Sprintf("Код: %s", err.Code))
	}

	// Детали
	if err.Details != "" {
		parts = append(parts, fmt.Sprintf("Детали: %s", err.Details))
	}

	// Оригинальная ошибка
	if err.Original != nil {
		parts = append(parts, fmt.Sprintf("Причина: %v", err.Original))
	}

	// Подсказки
	if len(err.Suggestions) > 0 {
		parts = append(parts, "\n💡 Что можно сделать:")
		for i, suggestion := range err.Suggestions {
			// Красиво форматируем ссылки на документацию
			formattedSuggestion := formatSuggestion(suggestion)
			parts = append(parts, fmt.Sprintf("  %d. %s", i+1, formattedSuggestion))
		}
	}

	return strings.Join(parts, "\n")
}

// FormatSuccess форматирует сообщение об успехе
func FormatSuccess(message string) string {
	title := UIStyles.ErrorTitle.Render("✅ Успех")
	msg := UIStyles.ErrorMessage.Render(message)

	content := fmt.Sprintf("%s\n\n%s", title, msg)
	return UIStyles.SuccessBox.Render(content)
}

// FormatWarning форматирует предупреждение
func FormatWarning(message string) string {
	title := UIStyles.ErrorTitle.Render("⚠️ Предупреждение")
	msg := UIStyles.ErrorMessage.Render(message)

	content := fmt.Sprintf("%s\n\n%s", title, msg)
	return UIStyles.WarningBox.Render(content)
}

// FormatInfo форматирует информационное сообщение
func FormatInfo(message string) string {
	title := UIStyles.ErrorTitle.Render("ℹ️ Информация")
	msg := UIStyles.ErrorMessage.Render(message)

	content := fmt.Sprintf("%s\n\n%s", title, msg)
	return UIStyles.InfoBox.Render(content)
}

// GetErrorSuggestions возвращает предложения по исправлению ошибки
func GetErrorSuggestions(err error) []string {
	if appErr, ok := err.(*AppError); ok {
		return getSuggestionsForType(appErr.Type, appErr.Code)
	}
	return []string{"Проверьте правильность введенных данных"}
}

// getSuggestionsForType возвращает предложения для конкретного типа ошибки
func getSuggestionsForType(errorType ErrorType, code string) []string {
	switch errorType {
	case ErrorTypeValidation:
		return []string{
			"Проверьте правильность введенных данных",
			"Убедитесь, что все обязательные поля заполнены",
			"Проверьте формат введенных значений",
		}
	case ErrorTypeConfiguration:
		return []string{
			"Проверьте переменные окружения",
			"Убедитесь, что конфигурационный файл существует",
			"Проверьте права доступа к конфигурационным файлам",
		}
	case ErrorTypeAuthentication:
		return []string{
			"Проверьте правильность учетных данных",
			"Убедитесь, что токен не истек",
			"Проверьте настройки аутентификации",
		}
	case ErrorTypeAPI:
		return []string{
			"Проверьте подключение к интернету",
			"Убедитесь, что API сервер доступен",
			"Попробуйте повторить запрос позже",
		}
	case ErrorTypeNetwork:
		return []string{
			"Проверьте подключение к интернету",
			"Убедитесь, что сетевые настройки корректны",
			"Попробуйте использовать VPN или другой сетевой интерфейс",
		}
	case ErrorTypeFileSystem:
		return []string{
			"Проверьте права доступа к файлам",
			"Убедитесь, что диск не заполнен",
			"Проверьте, что файл не заблокирован другим процессом",
		}
	case ErrorTypeTemplate:
		return []string{
			"Проверьте синтаксис шаблона",
			"Убедитесь, что все переменные определены",
			"Проверьте права доступа к файлам шаблонов",
		}
	case ErrorTypeUser:
		return []string{
			"Проверьте правильность введенных данных",
			"Обратитесь к документации",
			"Попробуйте другой подход",
		}
	case ErrorTypeSystem:
		return []string{
			"Перезапустите приложение",
			"Проверьте системные ресурсы",
			"Обратитесь к системному администратору",
		}
	default:
		return []string{"Попробуйте повторить операцию"}
	}
}

// formatSuggestion красиво форматирует подсказки, особенно ссылки на документацию
func formatSuggestion(suggestion string) string {
	// Если это ссылка на документацию, делаем её более заметной
	if strings.Contains(suggestion, "cloud.ru/docs/ai-agents") {
		// Заменяем обычную ссылку на более красивую
		return strings.Replace(suggestion,
			"📚 Подробная документация: https://cloud.ru/docs/ai-agents/ug/index?source-platform=Evolution",
			"📚 Подробная документация: https://cloud.ru/docs/ai-agents/ug/index?source-platform=Evolution",
			-1)
	}

	// Для остальных подсказок возвращаем как есть
	return suggestion
}
