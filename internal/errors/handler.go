package errors

import (
	"fmt"
	"os"
	"runtime"
	"strings"

	"github.com/charmbracelet/log"
)

// Handler представляет обработчик ошибок
type Handler struct {
	logger *Logger
}

// NewHandler создает новый обработчик ошибок
func NewHandler() *Handler {
	return &Handler{
		logger: NewLogger(),
	}
}

// Handle обрабатывает ошибку и возвращает пользовательское сообщение
func (h *Handler) Handle(err error) string {
	if err == nil {
		return ""
	}

	// Логируем ошибку
	h.logger.LogError(err, "Error occurred")

	// Если это AppError, используем структурированное отображение
	if appErr, ok := err.(*AppError); ok {
		return h.handleAppError(appErr)
	}

	// Для обычных ошибок используем простое отображение
	return h.handleGenericError(err)
}

// HandleSimple обрабатывает ошибку и возвращает простое сообщение без рамок
func (h *Handler) HandleSimple(err error) string {
	if err == nil {
		return ""
	}

	// Логируем ошибку только один раз
	if appErr, ok := err.(*AppError); ok {
		h.logger.LogAppError(appErr, "Structured error occurred")
	} else {
		h.logger.LogError(err, "Error occurred")
	}

	// Используем простое отображение
	return FormatSimpleError(err)
}

// HandlePlain обрабатывает ошибку и возвращает простой текст без стилей
func (h *Handler) HandlePlain(err error) string {
	if err == nil {
		return ""
	}

	// Логируем ошибку только один раз
	if appErr, ok := err.(*AppError); ok {
		h.logger.LogAppError(appErr, "Structured error occurred")
	} else {
		h.logger.LogError(err, "Error occurred")
	}

	// Используем простое текстовое отображение
	return FormatPlainError(err)
}

// handleAppError обрабатывает структурированную ошибку
func (h *Handler) handleAppError(err *AppError) string {
	// Логируем структурированную ошибку
	h.logger.LogAppError(err, "Structured error occurred")

	// Возвращаем отформатированное сообщение для UI
	return FormatError(err)
}

// handleGenericError обрабатывает обычную ошибку
func (h *Handler) handleGenericError(err error) string {
	// Создаем структурированную ошибку из обычной
	appErr := Wrap(err, ErrorTypeSystem, SeverityMedium, "GENERIC_ERROR", "Произошла ошибка")
	
	// Логируем
	h.logger.LogAppError(appErr, "Generic error occurred")

	// Возвращаем отформатированное сообщение
	return FormatError(appErr)
}

// HandleWithRecovery обрабатывает ошибку с восстановлением
func (h *Handler) HandleWithRecovery(err error, recovery func()) string {
	if err == nil {
		return ""
	}

	// Обрабатываем ошибку
	message := h.Handle(err)

	// Выполняем восстановление
	if recovery != nil {
		recovery()
	}

	return message
}

// HandleFatal обрабатывает критическую ошибку и завершает программу
func (h *Handler) HandleFatal(err error) {
	if err == nil {
		return
	}

	// Логируем критическую ошибку
	h.logger.LogError(err, "Fatal error occurred")

	// Отображаем ошибку
	message := h.Handle(err)
	fmt.Fprintf(os.Stderr, "%s\n", message)

	// Завершаем программу
	os.Exit(1)
}

// HandleWithExitCode обрабатывает ошибку и завершает программу с кодом выхода
func (h *Handler) HandleWithExitCode(err error, exitCode int) {
	if err == nil {
		return
	}

	// Логируем ошибку
	h.logger.LogError(err, "Error occurred with exit code", "exit_code", exitCode)

	// Отображаем ошибку
	message := h.Handle(err)
	fmt.Fprintf(os.Stderr, "%s\n", message)

	// Завершаем программу с кодом выхода
	os.Exit(exitCode)
}

// WrapError оборачивает ошибку в структурированную
func (h *Handler) WrapError(err error, errorType ErrorType, severity Severity, code, message string) *AppError {
	return Wrap(err, errorType, severity, code, message)
}

// WrapValidationError оборачивает ошибку валидации
func (h *Handler) WrapValidationError(err error, code, message string) *AppError {
	return h.WrapError(err, ErrorTypeValidation, SeverityMedium, code, message)
}

// WrapConfigurationError оборачивает ошибку конфигурации
func (h *Handler) WrapConfigurationError(err error, code, message string) *AppError {
	return h.WrapError(err, ErrorTypeConfiguration, SeverityHigh, code, message)
}

// WrapAuthenticationError оборачивает ошибку аутентификации
func (h *Handler) WrapAuthenticationError(err error, code, message string) *AppError {
	return h.WrapError(err, ErrorTypeAuthentication, SeverityHigh, code, message)
}

// WrapAPIError оборачивает ошибку API
func (h *Handler) WrapAPIError(err error, code, message string) *AppError {
	return h.WrapError(err, ErrorTypeAPI, SeverityMedium, code, message)
}

// WrapFileSystemError оборачивает ошибку файловой системы
func (h *Handler) WrapFileSystemError(err error, code, message string) *AppError {
	return h.WrapError(err, ErrorTypeFileSystem, SeverityMedium, code, message)
}

// WrapTemplateError оборачивает ошибку шаблона
func (h *Handler) WrapTemplateError(err error, code, message string) *AppError {
	return h.WrapError(err, ErrorTypeTemplate, SeverityMedium, code, message)
}

// WrapUserError оборачивает пользовательскую ошибку
func (h *Handler) WrapUserError(err error, code, message string) *AppError {
	return h.WrapError(err, ErrorTypeUser, SeverityLow, code, message)
}

// WrapSystemError оборачивает системную ошибку
func (h *Handler) WrapSystemError(err error, code, message string) *AppError {
	return h.WrapError(err, ErrorTypeSystem, SeverityCritical, code, message)
}

// GetLogger возвращает логгер
func (h *Handler) GetLogger() *Logger {
	return h.logger
}

// SetLogLevel устанавливает уровень логирования
func (h *Handler) SetLogLevel(level log.Level) {
	h.logger.SetLevel(level)
}

// SetLogFormatter устанавливает форматтер логирования
func (h *Handler) SetLogFormatter(formatter log.Formatter) {
	h.logger.SetFormatter(formatter)
}

// SetLogReportTimestamp включает/выключает отчет о времени
func (h *Handler) SetLogReportTimestamp(report bool) {
	h.logger.SetReportTimestamp(report)
}

// SetLogReportCaller включает/выключает отчет о вызывающем
func (h *Handler) SetLogReportCaller(report bool) {
	h.logger.SetReportCaller(report)
}

// HandlePanic обрабатывает панику
func (h *Handler) HandlePanic() {
	if r := recover(); r != nil {
		// Получаем стек вызовов
		stack := make([]byte, 4096)
		length := runtime.Stack(stack, false)
		
		// Создаем структурированную ошибку
		appErr := SystemError("PANIC_RECOVERED", "Программа восстановлена после паники").
			WithDetails(fmt.Sprintf("Panic value: %v", r)).
			WithContext("stack_trace", string(stack[:length]))
		
		// Логируем
		h.logger.LogAppError(appErr, "Panic recovered")
		
		// Отображаем ошибку
		message := FormatError(appErr)
		fmt.Fprintf(os.Stderr, "%s\n", message)
		
		// Завершаем программу
		os.Exit(1)
	}
}

// HandlePanicWithRecovery обрабатывает панику с восстановлением
func (h *Handler) HandlePanicWithRecovery(recovery func()) {
	if r := recover(); r != nil {
		// Логируем
		h.logger.LogError(fmt.Errorf("panic recovered: %v", r), "Panic recovered with recovery function")
		
		// Выполняем восстановление
		if recovery != nil {
			recovery()
		}
	}
}

// IsErrorType проверяет, является ли ошибка определенного типа
func (h *Handler) IsErrorType(err error, errorType ErrorType) bool {
	if appErr, ok := err.(*AppError); ok {
		return appErr.Type == errorType
	}
	return false
}

// IsErrorSeverity проверяет, является ли ошибка определенной серьезности
func (h *Handler) IsErrorSeverity(err error, severity Severity) bool {
	if appErr, ok := err.(*AppError); ok {
		return appErr.Severity == severity
	}
	return false
}

// IsErrorCode проверяет, является ли ошибка определенного кода
func (h *Handler) IsErrorCode(err error, code string) bool {
	if appErr, ok := err.(*AppError); ok {
		return appErr.Code == code
	}
	return false
}

// GetErrorType возвращает тип ошибки
func (h *Handler) GetErrorType(err error) ErrorType {
	if appErr, ok := err.(*AppError); ok {
		return appErr.Type
	}
	return ErrorTypeSystem
}

// GetErrorSeverity возвращает серьезность ошибки
func (h *Handler) GetErrorSeverity(err error) Severity {
	if appErr, ok := err.(*AppError); ok {
		return appErr.Severity
	}
	return SeverityMedium
}

// GetErrorCode возвращает код ошибки
func (h *Handler) GetErrorCode(err error) string {
	if appErr, ok := err.(*AppError); ok {
		return appErr.Code
	}
	return "UNKNOWN"
}

// GetErrorSuggestions возвращает предложения по исправлению ошибки
func (h *Handler) GetErrorSuggestions(err error) []string {
	return GetErrorSuggestions(err)
}

// ShouldRetry определяет, следует ли повторить операцию
func (h *Handler) ShouldRetry(err error) bool {
	if appErr, ok := err.(*AppError); ok {
		// Повторяем для сетевых ошибок и временных ошибок API
		return appErr.Type == ErrorTypeNetwork || 
			   appErr.Type == ErrorTypeAPI && 
			   strings.Contains(strings.ToLower(appErr.Message), "timeout")
	}
	return false
}

// GetRetryDelay возвращает задержку перед повторной попыткой
func (h *Handler) GetRetryDelay(err error, attempt int) int {
	// Экспоненциальная задержка: 1s, 2s, 4s, 8s, 16s
	delay := 1 << uint(attempt-1)
	if delay > 60 {
		delay = 60 // Максимум 60 секунд
	}
	return delay
}
