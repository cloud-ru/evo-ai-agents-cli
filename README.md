# AI Agents CLI

Мощный инструмент командной строки для управления AI агентами, MCP серверами и агентными системами в Cloud.ru.

## 🚀 Возможности

- **Управление MCP серверами** - создание, обновление, удаление, развертывание
- **Управление AI агентами** - полный жизненный цикл агентов
- **Управление агентными системами** - создание и управление системами агентов
- **CI/CD интеграция** - мониторинг, логи, проверка статуса
- **Красивый интерфейс** - цветной вывод с использованием charmbracelet
- **Многоязычность** - поддержка русского и английского языков

## 📦 Установка

```bash
# Клонируйте репозиторий
git clone <repository-url>
cd cloud-ru-ai-agents-cli

# Соберите проект
go build -o ai-agents-cli

# Или используйте make
make build
```

## ⚙️ Конфигурация

Создайте файл `.env` или установите переменные окружения:

```bash
export API_KEY="your-api-key"
export PROJECT_ID="your-project-id"
export PUBLIC_API_ENDPOINT="ai-agents.api.cloud.ru"
```

## 🎯 Использование

### Основные команды

```bash
# Показать справку
./ai-agents-cli --help

# MCP серверы
./ai-agents-cli mcp-servers --help

# Агенты
./ai-agents-cli agents --help

# CI/CD функции
./ai-agents-cli ci --help
```

### MCP Серверы

```bash
# Список MCP серверов
./ai-agents-cli mcp-servers list

# Информация о сервере
./ai-agents-cli mcp-servers get <server-id>

# Создание сервера
./ai-agents-cli mcp-servers create --name "my-server" --description "Описание"

# Создание из конфигурации
./ai-agents-cli mcp-servers create --config config.json

# Обновление сервера
./ai-agents-cli mcp-servers update <server-id> --name "new-name"

# Удаление сервера
./ai-agents-cli mcp-servers delete <server-id>

# Развертывание из файла
./ai-agents-cli mcp-servers deploy mcp-servers.yaml

# Предварительный просмотр развертывания
./ai-agents-cli mcp-servers deploy --dry-run

# Управление состоянием
./ai-agents-cli mcp-servers resume <server-id>
./ai-agents-cli mcp-servers suspend <server-id>

# История операций
./ai-agents-cli mcp-servers history <server-id>
```

### Агенты

```bash
# Список агентов
./ai-agents-cli agents list

# Информация об агенте
./ai-agents-cli agents get <agent-id>

# Поиск в маркетплейсе
./ai-agents-cli agents marketplace

# Поиск с фильтрами
./ai-agents-cli agents marketplace --name "assistant" --tags "ai,chat"

# Управление состоянием
./ai-agents-cli agents resume <agent-id>
./ai-agents-cli agents suspend <agent-id>
```

### CI/CD Функции

```bash
# Проверка статуса системы
./ai-agents-cli ci status

# Статус конкретного ресурса
./ai-agents-cli ci status mcp-server <server-id>
./ai-agents-cli ci status agent <agent-id>
./ai-agents-cli ci status agent-system <system-id>

# Просмотр логов
./ai-agents-cli ci logs

# Логи конкретного ресурса
./ai-agents-cli ci logs mcp-server <server-id>

# Мониторинг логов в реальном времени
./ai-agents-cli ci logs --follow

# Последние 100 записей
./ai-agents-cli ci logs --tail 100
```

## 📁 Конфигурационные файлы

### MCP Серверы (mcp-servers.yaml)

```yaml
mcp-servers:
  - name: "my_simple_mcp"
    description: "Простой MCP сервер для демонстрации"
    options:
      host: "localhost"
      port: 8080
      timeout: 30
      retries: 3
      
  - name: "database_mcp"
    description: "MCP сервер для работы с базой данных"
    options:
      connection_string: "postgresql://user:pass@localhost/db"
      max_connections: 10
      query_timeout: 60
```

### Агенты (agents.yaml)

```yaml
agents:
  - name: "my_assistant"
    description: "AI ассистент для помощи пользователям"
    options:
      model: "gpt-4"
      temperature: 0.7
      max_tokens: 1000
    llm_options:
      provider: "openai"
      api_key: "${OPENAI_API_KEY}"
    mcp_servers:
      - "my_simple_mcp"
      - "database_mcp"
```

## 🔧 Переменные окружения

| Переменная | Описание | По умолчанию |
|------------|----------|--------------|
| `API_KEY` | API ключ для аутентификации | - |
| `PROJECT_ID` | ID проекта | - |
| `PUBLIC_API_ENDPOINT` | Endpoint API | `ai-agents.api.cloud.ru` |
| `BULK_OPERATIONS_CONCURRENCY` | Количество параллельных операций | `20` |

## 🎨 Цветовая схема

CLI использует красивую цветовую схему:
- 🟢 Зеленый - успешные операции, активные ресурсы
- 🔴 Красный - ошибки, неактивные ресурсы
- 🟡 Желтый - предупреждения, приостановленные ресурсы
- 🔵 Синий - информация, заголовки
- ⚪ Серый - нейтральная информация

## 🚀 CI/CD Интеграция

### GitHub Actions

```yaml
name: Deploy AI Agents
on:
  push:
    branches: [main]

jobs:
  deploy:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - name: Deploy MCP Servers
        run: |
          ./ai-agents-cli mcp-servers deploy --dry-run
          ./ai-agents-cli mcp-servers deploy
        env:
          API_KEY: ${{ secrets.API_KEY }}
          PROJECT_ID: ${{ secrets.PROJECT_ID }}
```

### GitLab CI

```yaml
deploy:
  stage: deploy
  script:
    - ./ai-agents-cli ci status
    - ./ai-agents-cli mcp-servers deploy
  only:
    - main
```

## 📝 Логирование

CLI поддерживает различные уровни логирования:

```bash
# Детальное логирование
./ai-agents-cli --verbose mcp-servers list

# Логи в JSON формате
./ai-agents-cli mcp-servers get <id> --output json
```

## 🤝 Разработка

### Структура проекта

```
cmd/                    # Команды CLI
├── agent/             # Команды для агентов
├── ci/                # CI/CD функции
├── mcp_server/        # Команды для MCP серверов
└── ...

internal/              # Внутренние пакеты
├── api/               # API клиент
├── config/            # Конфигурация
└── ...

localizations/         # Локализация
├── i18n/
│   ├── en/           # Английский
│   └── ru/           # Русский
└── ...

examples/              # Примеры конфигураций
schemas/               # JSON схемы
```

### Добавление новых команд

1. Создайте файл в соответствующей директории `cmd/`
2. Добавьте импорт в `cmd/imports.go`
3. Обновите локализацию в `localizations/i18n/`

### Тестирование

```bash
# Запуск тестов
go test ./...

# Тестирование с покрытием
go test -cover ./...

# Линтинг
golangci-lint run
```

## 📄 Лицензия

MIT License

## 🆘 Поддержка

Если у вас возникли вопросы или проблемы:

1. Проверьте [Issues](https://github.com/your-repo/issues)
2. Создайте новый Issue с подробным описанием
3. Обратитесь к команде разработки

## 🔄 Обновления

Следите за обновлениями в [Releases](https://github.com/your-repo/releases)

---

**AI Agents CLI** - ваш надежный помощник в управлении AI агентами! 🚀