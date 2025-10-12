# Тестирование AI Agents CLI

## Быстрый старт

### 1. Установка зависимостей
```bash
go mod tidy
```

### 2. Сборка проекта
```bash
go build -o bin/ai-agents-cli .
```

### 3. Запуск тестов
```bash
# Все тесты
go test ./internal/... ./cmd_test -v

# С покрытием кода
go test ./internal/... ./cmd_test -coverprofile=coverage.out
go tool cover -html=coverage.out -o coverage.html
```

## Тестирование функциональности

### Валидация конфигурации
```bash
# Установить переменные окружения
export API_KEY=test
export PROJECT_ID=test

# Валидация примеров
./bin/ai-agents-cli validate examples/agents.yaml
./bin/ai-agents-cli validate examples/mcp-servers.yaml
./bin/ai-agents-cli validate examples/agent-systems.yaml

# Валидация всех файлов в директории
./bin/ai-agents-cli validate examples/
```

### CLI команды
```bash
# Справка
./bin/ai-agents-cli --help

# Verbose режим
./bin/ai-agents-cli --verbose

# Валидация
./bin/ai-agents-cli validate --help
```

## Структура тестов

```
├── internal/
│   ├── config/
│   │   └── config_test.go          # Тесты конфигурации
│   ├── validator/
│   │   └── validator_test.go       # Тесты валидатора
│   └── api/
│       ├── api_test.go             # Тесты API
│       ├── client_test.go          # Тесты HTTP клиента
│       └── mcp_server_test.go      # Тесты MCP серверов
├── cmd_test/
│   └── validate_test.go            # Тесты CLI команд
└── TEST_REPORT.md                  # Подробный отчет
```

## Покрытие кода

| Компонент | Покрытие |
|-----------|----------|
| `internal/config` | 66.7% |
| `internal/validator` | 70.9% |
| `internal/api` | 49.4% |
| **Общее** | **62.3%** |

## Примеры тестов

### Тест валидации
```go
func TestConfigValidator_ValidateFile(t *testing.T) {
    validator := validator.NewConfigValidator()
    
    // Валидный файл
    result, err := validator.ValidateFile("valid.yaml")
    assert.NoError(t, err)
    assert.True(t, result.Valid)
    
    // Невалидный файл
    result, err = validator.ValidateFile("invalid.yaml")
    assert.NoError(t, err)
    assert.False(t, result.Valid)
    assert.NotEmpty(t, result.Errors)
}
```

### Тест API клиента
```go
func TestClient_Get(t *testing.T) {
    server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusOK)
        json.NewEncoder(w).Encode(map[string]string{"status": "success"})
    }))
    defer server.Close()
    
    client := NewClient(server.URL, "test-key", "test-project")
    
    var result map[string]string
    err := client.Get(context.Background(), "/test", nil, &result)
    assert.NoError(t, err)
    assert.Equal(t, "success", result["status"])
}
```

## Отладка

### Включить verbose логирование
```bash
./bin/ai-agents-cli --verbose validate examples/agents.yaml
```

### Проверить переменные окружения
```bash
echo $API_KEY
echo $PROJECT_ID
```

### Запустить тесты с детальным выводом
```bash
go test ./internal/... -v -count=1
```

## CI/CD

### GitHub Actions
```yaml
name: Test
on: [push, pull_request]
jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-go@v3
        with:
          go-version: '1.24'
      - run: go test ./internal/... ./cmd_test -v
      - run: go test ./internal/... ./cmd_test -coverprofile=coverage.out
```

### Makefile команды
```bash
make test          # Запустить тесты
make test-coverage # Тесты с покрытием
make build         # Собрать проект
make validate      # Валидация конфигурации
```
