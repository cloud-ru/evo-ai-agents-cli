package errors

import (
	"fmt"

	"github.com/charmbracelet/log"
)

// ExampleUsage демонстрирует использование новой системы обработки ошибок
func ExampleUsage() {
	// Создаем обработчик ошибок
	handler := NewHandler()

	// Пример 1: Создание структурированной ошибки
	validationErr := ValidationError("INVALID_EMAIL", "Некорректный формат email адреса").
		WithDetails("Email должен содержать символ @").
		WithContext("input", "invalid-email")

	fmt.Println("Пример 1 - Структурированная ошибка:")
	fmt.Println(FormatError(validationErr))
	fmt.Println()

	// Пример 2: Оборачивание существующей ошибки
	originalErr := fmt.Errorf("file not found")
	wrappedErr := handler.WrapFileSystemError(originalErr, "FILE_NOT_FOUND", "Файл не найден").
		WithContext("file_path", "/path/to/file.txt").
		WithContext("operation", "read")

	fmt.Println("Пример 2 - Оборачивание ошибки:")
	fmt.Println(FormatError(wrappedErr))
	fmt.Println()

	// Пример 3: Обработка ошибки с предложениями
	fmt.Println("Пример 3 - Предложения по исправлению:")
	suggestions := GetErrorSuggestions(validationErr)
	for i, suggestion := range suggestions {
		fmt.Printf("%d. %s\n", i+1, suggestion)
	}
	fmt.Println()

	// Пример 4: Логирование ошибки
	fmt.Println("Пример 4 - Логирование ошибки:")
	handler.GetLogger().LogAppError(validationErr, "Validation failed during user input")
	fmt.Println()

	// Пример 5: Обработка паники
	fmt.Println("Пример 5 - Обработка паники:")
	defer handler.HandlePanic()

	// Пример 6: Различные типы ошибок
	fmt.Println("Пример 6 - Различные типы ошибок:")

	configErr := ConfigurationError("MISSING_ENV_VAR", "Отсутствует переменная окружения").
		WithDetails("Переменная IAM_KEY_ID не установлена")

	authErr := AuthenticationError("INVALID_TOKEN", "Недействительный токен аутентификации").
		WithDetails("Токен истек или имеет неверный формат")

	apiErr := APIError("SERVICE_UNAVAILABLE", "Сервис временно недоступен").
		WithDetails("Попробуйте повторить запрос через несколько минут")

	fmt.Println("Конфигурация:")
	fmt.Println(FormatError(configErr))
	fmt.Println()

	fmt.Println("Аутентификация:")
	fmt.Println(FormatError(authErr))
	fmt.Println()

	fmt.Println("API:")
	fmt.Println(FormatError(apiErr))
	fmt.Println()

	// Пример 7: Проверка типа ошибки
	fmt.Println("Пример 7 - Проверка типа ошибки:")
	if handler.IsErrorType(validationErr, ErrorTypeValidation) {
		fmt.Println("✅ Это ошибка валидации")
	}
	if handler.IsErrorSeverity(validationErr, SeverityMedium) {
		fmt.Println("✅ Это ошибка средней серьезности")
	}
	if handler.IsErrorCode(validationErr, "INVALID_EMAIL") {
		fmt.Println("✅ Это ошибка с кодом INVALID_EMAIL")
	}
	fmt.Println()

	// Пример 8: Успешные сообщения
	fmt.Println("Пример 8 - Успешные сообщения:")
	fmt.Println(FormatSuccess("Проект успешно создан!"))
	fmt.Println(FormatInfo("Используйте 'make help' для просмотра доступных команд"))
	fmt.Println(FormatWarning("Не забудьте настроить переменные окружения"))
	fmt.Println()

	// Пример 9: Обработка с восстановлением
	fmt.Println("Пример 9 - Обработка с восстановлением:")
	recoveryFunc := func() {
		fmt.Println("🔄 Выполняется восстановление...")
	}

	handler.HandleWithRecovery(validationErr, recoveryFunc)
	fmt.Println()

	// Пример 10: Настройка логгера
	fmt.Println("Пример 10 - Настройка логгера:")
	handler.SetLogLevel(log.DebugLevel)
	handler.SetLogReportTimestamp(true)
	handler.SetLogReportCaller(true)
	fmt.Println("Логгер настроен для отладки")
}
