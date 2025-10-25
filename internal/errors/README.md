# Система обработки ошибок

Этот пакет предоставляет улучшенную систему обработки ошибок для ai-agents-cli с поддержкой структурированных ошибок, контекстной информации и красивого отображения в UI.

## Основные возможности

- **Структурированные ошибки** с типом, серьезностью и кодом
- **Контекстная информация** для отладки
- **Красивое отображение** в терминале с цветами и иконками
- **Логирование** с разными уровнями
- **Предложения по исправлению** ошибок
- **Обработка паники** с восстановлением
- **Повторные попытки** для временных ошибок

## Типы ошибок

- `ErrorTypeValidation` - ошибки валидации
- `ErrorTypeConfiguration` - ошибки конфигурации
- `ErrorTypeAuthentication` - ошибки аутентификации
- `ErrorTypeAPI` - ошибки API
- `ErrorTypeNetwork` - сетевые ошибки
- `ErrorTypeFileSystem` - ошибки файловой системы
- `ErrorTypeTemplate` - ошибки шаблонов
- `ErrorTypeUser` - пользовательские ошибки
- `ErrorTypeSystem` - системные ошибки

## Уровни серьезности

- `SeverityLow` - низкая серьезность
- `SeverityMedium` - средняя серьезность
- `SeverityHigh` - высокая серьезность
- `SeverityCritical` - критическая серьезность

## Использование

### Создание ошибок

```go
// Создание новой ошибки
err := ValidationError("INVALID_EMAIL", "Некорректный формат email")

// Добавление контекста
err = err.WithContext("input", "invalid-email")

// Добавление деталей
err = err.WithDetails("Email должен содержать символ @")
```

### Оборачивание существующих ошибок

```go
// Оборачивание обычной ошибки
originalErr := fmt.Errorf("file not found")
wrappedErr := handler.WrapFileSystemError(originalErr, "FILE_NOT_FOUND", "Файл не найден")
```

### Обработка ошибок

```go
// Создание обработчика
handler := errors.NewHandler()

// Обработка ошибки
message := handler.Handle(err)
fmt.Println(message)

// Обработка с восстановлением
handler.HandleWithRecovery(err, recoveryFunc)

// Критическая ошибка
handler.HandleFatal(err)
```

### Логирование

```go
// Получение логгера
logger := handler.GetLogger()

// Логирование структурированной ошибки
logger.LogAppError(appErr, "Operation failed")

// Логирование с контекстом
logger.WithContext(map[string]interface{}{
    "user_id": "123",
    "operation": "create_project",
}).LogError(err, "Failed to create project")
```

### Отображение в UI

```go
// Форматирование ошибки для UI
message := errors.FormatError(err)

// Различные типы сообщений
success := errors.FormatSuccess("Операция выполнена успешно")
warning := errors.FormatWarning("Внимание: проверьте настройки")
info := errors.FormatInfo("Информация: используйте --help для справки")
```

### Обработка паники

```go
// Обработка паники с завершением программы
defer handler.HandlePanic()

// Обработка паники с восстановлением
defer handler.HandlePanicWithRecovery(func() {
    fmt.Println("Выполняется восстановление...")
})
```

## Примеры

### Валидация пользовательского ввода

```go
func validateEmail(email string) error {
    if !strings.Contains(email, "@") {
        return ValidationError("INVALID_EMAIL", "Некорректный формат email").
            WithDetails("Email должен содержать символ @").
            WithContext("input", email)
    }
    return nil
}
```

### Обработка API ошибок

```go
func callAPI() error {
    resp, err := http.Get("https://api.example.com/data")
    if err != nil {
        return handler.WrapNetworkError(err, "API_REQUEST_FAILED", "Ошибка запроса к API")
    }
    defer resp.Body.Close()

    if resp.StatusCode >= 400 {
        return APIError("API_ERROR", "Ошибка API").
            WithDetails(fmt.Sprintf("HTTP статус: %d", resp.StatusCode)).
            WithContext("status_code", resp.StatusCode)
    }

    return nil
}
```

### Обработка файловых операций

```go
func readConfigFile(path string) error {
    data, err := os.ReadFile(path)
    if err != nil {
        return handler.WrapFileSystemError(err, "CONFIG_READ_FAILED", "Ошибка чтения конфигурации").
            WithContext("file_path", path).
            WithDetails("Проверьте существование файла и права доступа")
    }

    // Дополнительная валидация
    if len(data) == 0 {
        return ConfigurationError("EMPTY_CONFIG", "Конфигурационный файл пуст").
            WithContext("file_path", path)
    }

    return nil
}
```

## Настройка

### Уровень логирования

```go
handler := errors.NewHandler()
handler.SetLogLevel(log.InfoLevel)
```

### Форматтер логгера

```go
handler.SetLogFormatter(log.JSONFormatter)
```

### Отчеты

```go
handler.SetLogReportTimestamp(true)
handler.SetLogReportCaller(true)
```

## Интеграция с командами

В командах используйте обработчик ошибок вместо прямого вызова `log.Fatal`:

```go
func (cmd *cobra.Command) Run(cmd *cobra.Command, args []string) {
    handler := errors.NewHandler()
    
    // Вместо log.Fatal
    if err := doSomething(); err != nil {
        appErr := handler.WrapUserError(err, "OPERATION_FAILED", "Ошибка выполнения операции")
        fmt.Println(handler.Handle(appErr))
        os.Exit(1)
    }
}
```

## Преимущества

1. **Структурированность** - все ошибки имеют тип, код и контекст
2. **Отладка** - легко найти источник ошибки по контексту
3. **Пользовательский опыт** - красивые сообщения с предложениями
4. **Логирование** - детальные логи для мониторинга
5. **Восстановление** - возможность восстановления после ошибок
6. **Повторные попытки** - автоматические повторы для временных ошибок
