package errors

import (
	"fmt"
	"time"
)

// ErrorType определяет тип ошибки
type ErrorType string

const (
	// ErrorTypeValidation - ошибки валидации
	ErrorTypeValidation ErrorType = "validation"
	// ErrorTypeConfiguration - ошибки конфигурации
	ErrorTypeConfiguration ErrorType = "configuration"
	// ErrorTypeAuthentication - ошибки аутентификации
	ErrorTypeAuthentication ErrorType = "authentication"
	// ErrorTypeAPI - ошибки API
	ErrorTypeAPI ErrorType = "api"
	// ErrorTypeNetwork - сетевые ошибки
	ErrorTypeNetwork ErrorType = "network"
	// ErrorTypeFileSystem - ошибки файловой системы
	ErrorTypeFileSystem ErrorType = "filesystem"
	// ErrorTypeTemplate - ошибки шаблонов
	ErrorTypeTemplate ErrorType = "template"
	// ErrorTypeUser - пользовательские ошибки
	ErrorTypeUser ErrorType = "user"
	// ErrorTypeSystem - системные ошибки
	ErrorTypeSystem ErrorType = "system"
)

// Severity определяет серьезность ошибки
type Severity string

const (
	// SeverityLow - низкая серьезность
	SeverityLow Severity = "low"
	// SeverityMedium - средняя серьезность
	SeverityMedium Severity = "medium"
	// SeverityHigh - высокая серьезность
	SeverityHigh Severity = "high"
	// SeverityCritical - критическая серьезность
	SeverityCritical Severity = "critical"
)

// AppError представляет структурированную ошибку приложения
type AppError struct {
	Type      ErrorType              `json:"type"`
	Severity  Severity               `json:"severity"`
	Code      string                 `json:"code"`
	Message   string                 `json:"message"`
	Details   string                 `json:"details,omitempty"`
	Context   map[string]interface{} `json:"context,omitempty"`
	Timestamp time.Time              `json:"timestamp"`
	Original  error                  `json:"-"`
	Suggestions []string             `json:"suggestions,omitempty"`
}

// Error реализует интерфейс error
func (e *AppError) Error() string {
	if e.Details != "" {
		return fmt.Sprintf("%s: %s (%s)", e.Message, e.Details, e.Code)
	}
	return fmt.Sprintf("%s (%s)", e.Message, e.Code)
}

// Unwrap возвращает оригинальную ошибку
func (e *AppError) Unwrap() error {
	return e.Original
}

// Is проверяет, является ли ошибка определенного типа
func (e *AppError) Is(target error) bool {
	if t, ok := target.(*AppError); ok {
		return e.Type == t.Type && e.Code == t.Code
	}
	return false
}

// WithContext добавляет контекст к ошибке
func (e *AppError) WithContext(key string, value interface{}) *AppError {
	if e.Context == nil {
		e.Context = make(map[string]interface{})
	}
	e.Context[key] = value
	return e
}

// WithDetails добавляет детали к ошибке
func (e *AppError) WithDetails(details string) *AppError {
	e.Details = details
	return e
}

// WithSuggestions добавляет подсказки к ошибке
func (e *AppError) WithSuggestions(suggestions ...string) *AppError {
	e.Suggestions = suggestions
	return e
}

// New создает новую структурированную ошибку
func New(errorType ErrorType, severity Severity, code, message string) *AppError {
	return &AppError{
		Type:      errorType,
		Severity:  severity,
		Code:      code,
		Message:   message,
		Timestamp: time.Now(),
	}
}

// Wrap оборачивает существующую ошибку
func Wrap(err error, errorType ErrorType, severity Severity, code, message string) *AppError {
	if err == nil {
		return nil
	}

	appErr := &AppError{
		Type:      errorType,
		Severity:  severity,
		Code:      code,
		Message:   message,
		Timestamp: time.Now(),
		Original:  err,
	}

	// Если это уже AppError, сохраняем контекст
	if existing, ok := err.(*AppError); ok {
		appErr.Context = existing.Context
		appErr.Details = existing.Details
	}

	return appErr
}

// ValidationError создает ошибку валидации
func ValidationError(code, message string) *AppError {
	return New(ErrorTypeValidation, SeverityMedium, code, message)
}

// ConfigurationError создает ошибку конфигурации
func ConfigurationError(code, message string) *AppError {
	return New(ErrorTypeConfiguration, SeverityHigh, code, message)
}

// AuthenticationError создает ошибку аутентификации
func AuthenticationError(code, message string) *AppError {
	return New(ErrorTypeAuthentication, SeverityHigh, code, message)
}

// APIError создает ошибку API
func APIError(code, message string) *AppError {
	return New(ErrorTypeAPI, SeverityMedium, code, message)
}

// NetworkError создает сетевую ошибку
func NetworkError(code, message string) *AppError {
	return New(ErrorTypeNetwork, SeverityMedium, code, message)
}

// FileSystemError создает ошибку файловой системы
func FileSystemError(code, message string) *AppError {
	return New(ErrorTypeFileSystem, SeverityMedium, code, message)
}

// TemplateError создает ошибку шаблона
func TemplateError(code, message string) *AppError {
	return New(ErrorTypeTemplate, SeverityMedium, code, message)
}

// UserError создает пользовательскую ошибку
func UserError(code, message string) *AppError {
	return New(ErrorTypeUser, SeverityLow, code, message)
}

// SystemError создает системную ошибку
func SystemError(code, message string) *AppError {
	return New(ErrorTypeSystem, SeverityCritical, code, message)
}
