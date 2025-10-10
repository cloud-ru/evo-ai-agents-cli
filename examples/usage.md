# Примеры использования AI Agents CLI

## 🚀 Быстрый старт

### 1. Настройка окружения

```bash
# Установите переменные окружения
export API_KEY="your-api-key-here"
export PROJECT_ID="your-project-id"

# Или создайте .env файл
echo "API_KEY=your-api-key-here" > .env
echo "PROJECT_ID=your-project-id" >> .env
```

### 2. Проверка подключения

```bash
# Проверьте статус системы
./ai-agents-cli ci status
```

## 📋 Управление MCP серверами

### Создание MCP сервера

```bash
# Простое создание
./ai-agents-cli mcp-servers create --name "my-database" --description "Database MCP server"

# Создание из конфигурации
./ai-agents-cli mcp-servers create --config database-config.json
```

### Просмотр MCP серверов

```bash
# Список всех серверов
./ai-agents-cli mcp-servers list

# С ограничением
./ai-agents-cli mcp-servers list --limit 10 --offset 0

# Подробная информация
./ai-agents-cli mcp-servers get server-id-123

# В JSON формате
./ai-agents-cli mcp-servers get server-id-123 --output json
```

### Управление состоянием

```bash
# Приостановка сервера
./ai-agents-cli mcp-servers suspend server-id-123

# Возобновление работы
./ai-agents-cli mcp-servers resume server-id-123

# История операций
./ai-agents-cli mcp-servers history server-id-123
```

### Развертывание

```bash
# Развертывание из файла
./ai-agents-cli mcp-servers deploy mcp-servers.yaml

# Предварительный просмотр
./ai-agents-cli mcp-servers deploy --dry-run

# Указание конкретного файла
./ai-agents-cli mcp-servers deploy --file production-config.yaml
```

## 🤖 Управление агентами

### Поиск в маркетплейсе

```bash
# Все доступные агенты
./ai-agents-cli agents marketplace

# Поиск по названию
./ai-agents-cli agents marketplace --name "assistant"

# Фильтрация по тегам
./ai-agents-cli agents marketplace --tags "ai,chat,support"

# Фильтрация по категориям
./ai-agents-cli agents marketplace --categories "customer-service,data-analysis"

# Комбинированный поиск
./ai-agents-cli agents marketplace \
  --name "support" \
  --tags "ai" \
  --categories "customer-service" \
  --types "AGENT_PREDEFINED_TYPE_FREE_TIER"
```

### Управление агентами

```bash
# Список агентов
./ai-agents-cli agents list

# Информация об агенте
./ai-agents-cli agents get agent-id-123

# Управление состоянием
./ai-agents-cli agents suspend agent-id-123
./ai-agents-cli agents resume agent-id-123
```

## 🔧 CI/CD функции

### Проверка статуса

```bash
# Общий статус системы
./ai-agents-cli ci status

# Статус конкретного ресурса
./ai-agents-cli ci status mcp-server server-id-123
./ai-agents-cli ci status agent agent-id-123
./ai-agents-cli ci status agent-system system-id-123

# Статус всех ресурсов определенного типа
./ai-agents-cli ci status mcp-servers
./ai-agents-cli ci status agents
./ai-agents-cli ci status agent-systems
```

### Мониторинг логов

```bash
# Последние логи
./ai-agents-cli ci logs

# Логи конкретного ресурса
./ai-agents-cli ci logs mcp-server server-id-123

# Мониторинг в реальном времени
./ai-agents-cli ci logs --follow

# Последние 100 записей
./ai-agents-cli ci logs --tail 100

# Логи за определенный период
./ai-agents-cli ci logs --since "2024-01-01" --until "2024-01-31"
```

## ✅ Валидация конфигурации

### Валидация файлов

```bash
# Валидация конкретного файла
./ai-agents-cli validate mcp-servers.yaml

# Валидация всех файлов в директории
./ai-agents-cli validate --dir ./configs

# Валидация с указанием файла
./ai-agents-cli validate --file production-config.json
```

## 🎨 Красивый вывод

CLI использует цветовую схему для лучшего восприятия:

- 🟢 **Зеленый** - успешные операции, активные ресурсы
- 🔴 **Красный** - ошибки, неактивные ресурсы  
- 🟡 **Желтый** - предупреждения, приостановленные ресурсы
- 🔵 **Синий** - информация, заголовки
- ⚪ **Серый** - нейтральная информация

## 📊 Примеры конфигураций

### MCP серверы (mcp-servers.yaml)

```yaml
mcp-servers:
  - name: "database_mcp"
    description: "MCP сервер для работы с PostgreSQL"
    options:
      host: "localhost"
      port: 5432
      database: "myapp"
      username: "${DB_USER}"
      password: "${DB_PASSWORD}"
      ssl_mode: "require"
      max_connections: 10
      timeout: 30
      
  - name: "api_mcp"
    description: "MCP сервер для внешних API"
    options:
      base_url: "https://api.example.com"
      api_key: "${API_KEY}"
      rate_limit: 100
      timeout: 30
      retries: 3
```

### Агенты (agents.yaml)

```yaml
agents:
  - name: "customer_support"
    description: "AI агент для поддержки клиентов"
    options:
      personality: "helpful and professional"
      response_style: "conversational"
      max_conversation_turns: 10
    llm_options:
      provider: "openai"
      model: "gpt-4"
      temperature: 0.7
      max_tokens: 1000
      api_key: "${OPENAI_API_KEY}"
    mcp_servers:
      - "database_mcp"
      - "ticket_system_mcp"
```

## 🔄 CI/CD интеграция

### GitHub Actions

```yaml
name: Deploy AI Agents
on:
  push:
    branches: [main]

jobs:
  validate:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - name: Validate configuration
        run: ./ai-agents-cli validate
        env:
          API_KEY: ${{ secrets.API_KEY }}
          PROJECT_ID: ${{ secrets.PROJECT_ID }}

  deploy:
    needs: validate
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - name: Deploy MCP servers
        run: ./ai-agents-cli mcp-servers deploy
        env:
          API_KEY: ${{ secrets.API_KEY }}
          PROJECT_ID: ${{ secrets.PROJECT_ID }}
      - name: Deploy agents
        run: ./ai-agents-cli agents deploy
        env:
          API_KEY: ${{ secrets.API_KEY }}
          PROJECT_ID: ${{ secrets.PROJECT_ID }}

  verify:
    needs: deploy
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - name: Check deployment status
        run: ./ai-agents-cli ci status
        env:
          API_KEY: ${{ secrets.API_KEY }}
          PROJECT_ID: ${{ secrets.PROJECT_ID }}
```

### GitLab CI

```yaml
stages:
  - validate
  - deploy
  - verify

validate:
  stage: validate
  script:
    - ./ai-agents-cli validate
  only:
    - main

deploy:
  stage: deploy
  script:
    - ./ai-agents-cli mcp-servers deploy
    - ./ai-agents-cli agents deploy
  only:
    - main

verify:
  stage: verify
  script:
    - ./ai-agents-cli ci status
  only:
    - main
```

## 🐛 Отладка

### Включение детального логирования

```bash
# Детальные логи
./ai-agents-cli --verbose mcp-servers list

# Логи в JSON формате
./ai-agents-cli mcp-servers get server-id --output json
```

### Проверка подключения

```bash
# Проверка статуса API
./ai-agents-cli ci status

# Проверка конкретного ресурса
./ai-agents-cli mcp-servers get server-id
```

## 📝 Лучшие практики

1. **Всегда валидируйте конфигурацию** перед развертыванием
2. **Используйте dry-run** для предварительного просмотра
3. **Мониторьте логи** после развертывания
4. **Используйте переменные окружения** для секретов
5. **Версионируйте конфигурации** в Git
6. **Настройте CI/CD** для автоматического развертывания

## 🆘 Решение проблем

### Частые ошибки

```bash
# Ошибка аутентификации
Error: API_KEY environment variable is required
# Решение: Установите переменную API_KEY

# Ошибка проекта
Error: PROJECT_ID environment variable is required  
# Решение: Установите переменную PROJECT_ID

# Ошибка валидации
Error: Configuration validation failed
# Решение: Проверьте синтаксис конфигурационного файла
```

### Получение помощи

```bash
# Общая справка
./ai-agents-cli --help

# Справка по команде
./ai-agents-cli mcp-servers --help
./ai-agents-cli agents --help
./ai-agents-cli ci --help
```

---

**AI Agents CLI** - ваш надежный помощник в управлении AI агентами! 🚀
