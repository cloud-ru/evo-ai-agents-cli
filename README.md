# AI Agents CLI

[![Build Status](https://github.com/cloud-ru/evo-ai-agents-cli/workflows/CI/badge.svg)](https://github.com/cloud-ru/evo-ai-agents-cli/actions)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Version](https://img.shields.io/badge/Go-1.24.3-blue.svg)](https://golang.org/)
[![Release](https://img.shields.io/github/v/release/cloud-ru/evo-ai-agents-cli)](https://github.com/cloud-ru/evo-ai-agents-cli/releases)

CLI инструмент для управления AI Agents в облачной платформе Cloud.ru.

## 🚀 Быстрый старт

### Установка

#### Windows (winget)
```bash
winget install CloudRu.AIAgentsCLI
```

#### macOS/Linux (Homebrew)
```bash
brew install cloud-ru/evo-ai-agents-cli/ai-agents-cli
```

#### Ручная установка
```bash
# Скачайте последнюю версию с GitHub Releases
wget https://github.com/cloud-ru/evo-ai-agents-cli/releases/latest/download/ai-agents-cli-linux-amd64.tar.gz
tar -xzf ai-agents-cli-linux-amd64.tar.gz
sudo mv ai-agents-cli /usr/local/bin/
```

#### Сборка из исходников
```bash
git clone https://github.com/cloud-ru/evo-ai-agents-cli.git
cd evo-ai-agents-cli
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
   # IAM аутентификация (получите в панели Cloud.ru)
   IAM_KEY_ID=your-iam-key-id
   IAM_SECRET=your-iam-secret
   
   # ID проекта AI Agents
   PROJECT_ID=your-project-id
   
   # Опциональные настройки
   IAM_ENDPOINT=https://iam.api.cloud.ru
   PUBLIC_API_ENDPOINT=ai-agents.api.cloud.ru
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

# Создание проектов из шаблонов
./bin/ai-agents-cli create mcp my-mcp-server
./bin/ai-agents-cli create agent my-ai-agent --path ./projects/
./bin/ai-agents-cli create mcp my-server --author "John Doe" --python-version "3.11"
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

### Создание проектов
- `create mcp [project-name]` - Создать проект MCP сервера из шаблона
- `create agent [project-name]` - Создать проект AI агента из шаблона
- `--path [path]` - Указать путь для создания проекта
- `--author [name]` - Указать автора проекта
- `--python-version [version]` - Указать версию Python

## 📁 Структура проекта

```
ai-agents-cli/
├── cmd/                    # CLI команды
│   ├── validate/          # Валидация конфигураций
│   ├── create/            # Создание проектов из шаблонов
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

| Переменная | Описание | Обязательная | По умолчанию |
|------------|----------|--------------|--------------|
| `IAM_KEY_ID` | IAM Key ID для аутентификации | ✅ | - |
| `IAM_SECRET` | IAM Secret для аутентификации | ✅ | - |
| `PROJECT_ID` | ID проекта AI Agents | ✅ | - |
| `IAM_ENDPOINT` | IAM API endpoint | ❌ | `https://iam.api.cloud.ru` |
| `PUBLIC_API_ENDPOINT` | AI Agents API endpoint | ❌ | `ai-agents.api.cloud.ru` |
| `SERVICE_APP_ENVIRONMENT` | Окружение приложения | ❌ | `dev` |
| `SERVICE_LOG_LEVEL` | Уровень логирования | ❌ | `debug` |
| `BULK_OPERATIONS_CONCURRENCY` | Параллельность bulk операций | ❌ | `20` |
| `SCAFFOLDER_AUTHOR` | Автор по умолчанию для новых проектов | ❌ | Автоматически из `git config user.name` и `user.email` |
| `SCAFFOLDER_PYTHON_VERSION` | Версия Python по умолчанию | ❌ | `3.9` |
| `SCAFFOLDER_DEFAULT_CICD` | CI/CD система по умолчанию | ❌ | `both` |

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

- [Руководство по валидации](TESTING.md)
- [Установка и настройка](.github/INSTALL.md)
- [Примеры использования](examples/usage.md)
- [API документация](service.swagger.json)

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
- 📖 Документация: [Cloud.ru AI Agents Documentation](https://cloud.ru/docs/ai-agents/ug/index?source-platform=Evolution)
- 🐛 Баги: [GitHub Issues](https://github.com/cloud-ru/evo-ai-agents-cli/issues)
- 💬 Обсуждения: [GitHub Discussions](https://github.com/cloud-ru/evo-ai-agents-cli/discussions)

## 🔄 История изменений

### v1.0.0
- ✅ Базовая функциональность CLI
- ✅ IAM аутентификация
- ✅ Валидация конфигураций
- ✅ Управление MCP серверами
- ✅ Управление агентами
- ✅ Управление системами агентов
- ✅ Полное покрытие тестами