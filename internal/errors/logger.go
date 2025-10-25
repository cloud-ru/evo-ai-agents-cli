package errors

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/charmbracelet/log"
)

// Logger представляет улучшенный логгер с обработкой ошибок
type Logger struct {
	*log.Logger
	context map[string]interface{}
}

// NewLogger создает новый логгер
func NewLogger() *Logger {
	// Настраиваем логгер
	logger := log.New(os.Stderr)
	logger.SetLevel(log.InfoLevel)
	logger.SetFormatter(log.JSONFormatter)
	logger.SetReportTimestamp(true)
	logger.SetReportCaller(true)

	return &Logger{
		Logger:  logger,
		context: make(map[string]interface{}),
	}
}

// WithContext добавляет контекст к логгеру
func (l *Logger) WithContext(ctx map[string]interface{}) *Logger {
	newLogger := &Logger{
		Logger:  l.Logger,
		context: make(map[string]interface{}),
	}
	
	// Копируем существующий контекст
	for k, v := range l.context {
		newLogger.context[k] = v
	}
	
	// Добавляем новый контекст
	for k, v := range ctx {
		newLogger.context[k] = v
	}
	
	return newLogger
}

// WithField добавляет поле к контексту
func (l *Logger) WithField(key string, value interface{}) *Logger {
	return l.WithContext(map[string]interface{}{key: value})
}

// WithError добавляет ошибку к контексту
func (l *Logger) WithError(err error) *Logger {
	fields := map[string]interface{}{"error": err.Error()}
	
	// Если это AppError, добавляем дополнительную информацию
	if appErr, ok := err.(*AppError); ok {
		fields["error_type"] = string(appErr.Type)
		fields["error_severity"] = string(appErr.Severity)
		fields["error_code"] = appErr.Code
		if appErr.Details != "" {
			fields["error_details"] = appErr.Details
		}
		if len(appErr.Context) > 0 {
			fields["error_context"] = appErr.Context
		}
	}
	
	return l.WithContext(fields)
}

// LogError логирует ошибку с соответствующим уровнем
func (l *Logger) LogError(err error, message string, fields ...interface{}) {
	if err == nil {
		return
	}

	// Определяем уровень логирования на основе типа ошибки
	level := l.getLogLevel(err)
	
	// Подготавливаем поля для логирования
	logFields := append([]interface{}{"error", err}, fields...)
	
	// Добавляем контекст
	for k, v := range l.context {
		logFields = append(logFields, k, v)
	}

	// Логируем с соответствующим уровнем
	switch level {
	case log.DebugLevel:
		l.Debug(message, logFields...)
	case log.InfoLevel:
		l.Info(message, logFields...)
	case log.WarnLevel:
		l.Warn(message, logFields...)
	case log.ErrorLevel:
		l.Error(message, logFields...)
	case log.FatalLevel:
		l.Fatal(message, logFields...)
	}
}

// LogAppError логирует структурированную ошибку
func (l *Logger) LogAppError(err *AppError, message string, fields ...interface{}) {
	if err == nil {
		return
	}

	// Подготавливаем поля
	logFields := []interface{}{
		"error_type", string(err.Type),
		"error_severity", string(err.Severity),
		"error_code", err.Code,
		"error_message", err.Message,
	}
	
	if err.Details != "" {
		logFields = append(logFields, "error_details", err.Details)
	}
	
	if len(err.Context) > 0 {
		logFields = append(logFields, "error_context", err.Context)
	}
	
	if err.Original != nil {
		logFields = append(logFields, "original_error", err.Original.Error())
	}
	
	// Добавляем дополнительные поля
	logFields = append(logFields, fields...)
	
	// Добавляем контекст логгера
	for k, v := range l.context {
		logFields = append(logFields, k, v)
	}

	// Логируем с соответствующим уровнем
	level := l.getSeverityLevel(err.Severity)
	switch level {
	case log.DebugLevel:
		l.Debug(message, logFields...)
	case log.InfoLevel:
		l.Info(message, logFields...)
	case log.WarnLevel:
		l.Warn(message, logFields...)
	case log.ErrorLevel:
		l.Error(message, logFields...)
	case log.FatalLevel:
		l.Fatal(message, logFields...)
	}
}

// getLogLevel определяет уровень логирования на основе ошибки
func (l *Logger) getLogLevel(err error) log.Level {
	if appErr, ok := err.(*AppError); ok {
		return l.getSeverityLevel(appErr.Severity)
	}
	
	// Для обычных ошибок используем Error уровень
	return log.ErrorLevel
}

// getSeverityLevel преобразует серьезность в уровень логирования
func (l *Logger) getSeverityLevel(severity Severity) log.Level {
	switch severity {
	case SeverityLow:
		return log.InfoLevel
	case SeverityMedium:
		return log.WarnLevel
	case SeverityHigh:
		return log.ErrorLevel
	case SeverityCritical:
		return log.FatalLevel
	default:
		return log.ErrorLevel
	}
}

// LogOperation логирует начало и конец операции
func (l *Logger) LogOperation(operation string, fn func() error) error {
	start := time.Now()
	l.Info("Starting operation", "operation", operation)
	
	err := fn()
	
	duration := time.Since(start)
	
	if err != nil {
		l.LogError(err, "Operation failed", 
			"operation", operation, 
			"duration", duration.String())
		return err
	}
	
	l.Info("Operation completed", 
		"operation", operation, 
		"duration", duration.String())
	
	return nil
}

// LogOperationWithContext логирует операцию с контекстом
func (l *Logger) LogOperationWithContext(ctx context.Context, operation string, fn func() error) error {
	start := time.Now()
	l.Info("Starting operation", "operation", operation)
	
	err := fn()
	
	duration := time.Since(start)
	
	if err != nil {
		l.LogError(err, "Operation failed", 
			"operation", operation, 
			"duration", duration.String())
		return err
	}
	
	l.Info("Operation completed", 
		"operation", operation, 
		"duration", duration.String())
	
	return nil
}

// RecoverAndLog восстанавливается от паники и логирует ошибку
func (l *Logger) RecoverAndLog() {
	if r := recover(); r != nil {
		err := fmt.Errorf("panic recovered: %v", r)
		l.LogError(err, "Panic recovered")
	}
}

// RecoverAndLogWithHandler восстанавливается от паники с пользовательским обработчиком
func (l *Logger) RecoverAndLogWithHandler(handler func(error)) {
	if r := recover(); r != nil {
		err := fmt.Errorf("panic recovered: %v", r)
		l.LogError(err, "Panic recovered")
		if handler != nil {
			handler(err)
		}
	}
}

// SetLevel устанавливает уровень логирования
func (l *Logger) SetLevel(level log.Level) {
	l.Logger.SetLevel(level)
}

// SetFormatter устанавливает форматтер логгера
func (l *Logger) SetFormatter(formatter log.Formatter) {
	l.Logger.SetFormatter(formatter)
}

// SetReportTimestamp включает/выключает отчет о времени
func (l *Logger) SetReportTimestamp(report bool) {
	l.Logger.SetReportTimestamp(report)
}

// SetReportCaller включает/выключает отчет о вызывающем
func (l *Logger) SetReportCaller(report bool) {
	l.Logger.SetReportCaller(report)
}
