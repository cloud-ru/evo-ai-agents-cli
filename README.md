# AI Agents CLI

CLI инструмент для управления AI Agents в облачной платформе Cloud.ru.

## 🚀 Быстрый старт

### Установка

1. **Скачайте бинарный файл:**
   ```bash
   # Скачайте последнюю версию с GitHub Releases
   wget https://github.com/cloudru/ai-agents-cli/releases/latest/download/ai-agents-cli-linux-amd64
   chmod +x ai-agents-cli-linux-amd64
   sudo mv ai-agents-cli-linux-amd64 /usr/local/bin/ai-agents-cli
   ```

2. **Или соберите из исходников:**
   ```bash
   git clone https://github.com/cloudru/ai-agents-cli.git
   cd ai-agents-cli
   go build -o bin/ai-agents-cli .
   ```

### Настройка

1. **Скопируйте файл конфигурации:**
   ```bash
   cp env.example .env
   ```

2. **Отредактируйте `.env` с вашими данными:**
   ```bash
   nano .env
   ```

3. **Обязательные переменные:**
   ```bash
   IAM_KEY_ID=your-iam-key-id
   IAM_SECRET=your-iam-secret
   PROJECT_ID=your-project-id
   ```

### Использование

```bash
# Валидация конфигурационных файлов
./bin/ai-agents-cli validate examples/

# Справка по командам
./bin/ai-agents-cli --help

# Управление MCP серверами
./bin/ai-agents-cli mcp-servers --help

# Управление агентами
./bin/ai-agents-cli agents --help
```

## 📋 Команды

### Валидация
- `validate` - Валидация конфигурационных файлов
- `validate [file]` - Валидация конкретного файла
- `validate [dir]` - Валидация всех файлов в директории

### MCP Серверы
- `mcp-servers list` - Список MCP серверов
- `mcp-servers get <id>` - Получить информацию о сервере
- `mcp-servers create` - Создать новый сервер
- `mcp-servers update <id>` - Обновить сервер
- `mcp-servers delete <id>` - Удалить сервер

### Агенты
- `agents list` - Список агентов
- `agents get <id>` - Получить информацию об агенте
- `agents create` - Создать нового агента
- `agents update <id>` - Обновить агента
- `agents delete <id>` - Удалить агента

### Системы агентов
- `system list` - Список систем агентов
- `system get <id>` - Получить информацию о системе
- `system create` - Создать новую систему
- `system update <id>` - Обновить систему
- `system delete <id>` - Удалить систему

## 📁 Структура проекта

```
ai-agents-cli/
├── cmd/                    # CLI команды
│   ├── validate/          # Валидация конфигураций
│   ├── mcp_server/        # Управление MCP серверами
│   ├── agent/             # Управление агентами
│   └── system/            # Управление системами
├── internal/              # Внутренние пакеты
│   ├── api/               # API клиент
│   ├── auth/              # IAM аутентификация
│   ├── config/            # Конфигурация
│   └── validator/         # Валидатор конфигураций
├── examples/              # Примеры конфигураций
├── schemas/               # JSON схемы для валидации
├── env.example            # Шаблон переменных окружения
└── README.md              # Этот файл
```

## 🔧 Конфигурация

### Переменные окружения

| Переменная | Описание | Обязательная |
|------------|----------|--------------|
| `IAM_KEY_ID` | IAM Key ID для аутентификации | ✅ |
| `IAM_SECRET` | IAM Secret для аутентификации | ✅ |
| `PROJECT_ID` | ID проекта AI Agents | ✅ |
| `IAM_ENDPOINT` | IAM API endpoint | ❌ |
| `PUBLIC_API_ENDPOINT` | AI Agents API endpoint | ❌ |

### Файлы конфигурации

CLI поддерживает следующие форматы:
- **YAML** (`.yaml`, `.yml`)
- **JSON** (`.json`)

Примеры конфигураций находятся в папке `examples/`:
- `agents.yaml` - Конфигурация агентов
- `mcp-servers.yaml` - Конфигурация MCP серверов
- `agent-systems.yaml` - Конфигурация систем агентов

## 🧪 Тестирование

```bash
# Запуск всех тестов
go test ./...

# Запуск тестов с покрытием
go test -cover ./...

# Запуск конкретных тестов
go test ./internal/validator -v
```

## 🚀 Разработка

### Требования
- Go 1.24.3+
- Git

### Сборка
```bash
# Установка зависимостей
go mod tidy

# Сборка
go build -o bin/ai-agents-cli .

# Запуск тестов
go test ./...
```

### Добавление новых команд

1. Создайте новый пакет в `cmd/`
2. Реализуйте команду с помощью Cobra
3. Добавьте тесты
4. Обновите документацию

## 📚 Документация

- [Настройка переменных окружения](ENV_SETUP.md)
- [Руководство по валидации](TESTING.md)
- [Отчет о тестировании](TEST_REPORT.md)

## 🤝 Вклад в проект

1. Форкните репозиторий
2. Создайте ветку для новой функции
3. Внесите изменения
4. Добавьте тесты
5. Создайте Pull Request

## 📄 Лицензия

Этот проект лицензирован под MIT License - см. файл [LICENSE](LICENSE) для деталей.

## 🆘 Поддержка

- 📧 Email: support@cloud.ru
- 📖 Документация: https://docs.cloud.ru
- 🐛 Баги: https://github.com/cloudru/ai-agents-cli/issues

## 🔄 История изменений

### v1.0.0
- ✅ Базовая функциональность CLI
- ✅ IAM аутентификация
- ✅ Валидация конфигураций
- ✅ Управление MCP серверами
- ✅ Управление агентами
- ✅ Управление системами агентов
- ✅ Полное покрытие тестами