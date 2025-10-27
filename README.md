# 🤖 AI Agents CLI

<div align="center">

[![Build Status](https://github.com/cloud-ru/evo-ai-agents-cli/workflows/CI/badge.svg)](https://github.com/cloud-ru/evo-ai-agents-cli/actions)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Version](https://img.shields.io/badge/Go-1.24.3+-blue.svg)](https://golang.org/)
[![Release](https://img.shields.io/github/v/release/cloud-ru/evo-ai-agents-cli)](https://github.com/cloud-ru/evo-ai-agents-cli/releases)
[![GitHub stars](https://img.shields.io/github/stars/cloud-ru/evo-ai-agents-cli.svg?style=flat-square&label=Stars)](https://github.com/cloud-ru/evo-ai-agents-cli/stargazers)
[![GitHub forks](https://img.shields.io/github/forks/cloud-ru/evo-ai-agents-cli.svg?style=flat-square&label=Forks)](https://github.com/cloud-ru/evo-ai-agents-cli/network/members)
[![Contributors](https://img.shields.io/github/contributors/cloud-ru/evo-ai-agents-cli.svg?style=flat-square)](https://github.com/cloud-ru/evo-ai-agents-cli/graphs/contributors)
[![Issues](https://img.shields.io/github/issues/cloud-ru/evo-ai-agents-cli.svg?style=flat-square)](https://github.com/cloud-ru/evo-ai-agents-cli/issues)
[![Downloads](https://img.shields.io/github/downloads/cloud-ru/evo-ai-agents-cli/total.svg?style=flat-square)](https://github.com/cloud-ru/evo-ai-agents-cli/releases)
[![Platforms](https://img.shields.io/badge/platform-Linux%20%7C%20macOS%20%7C%20Windows-lightgrey.svg)](https://github.com/cloud-ru/evo-ai-agents-cli/releases)

**Мощный CLI инструмент для управления AI агентами, MCP серверами и их системами в облачной платформе Cloud.ru**

[Установка](#-установка) • [Быстрый старт](#-быстрый-старт) • [Документация](#-документация) • [Примеры](#-примеры-использования)

</div>

---

## ✨ Ключевые возможности

- 🎯 **Управление агентами**: Создание, обновление, удаление и управление жизненным циклом AI агентов
- 🔌 **Управление MCP серверами**: Работа с Model Context Protocol серверами для интеграции
- 🏗️ **Системы агентов**: Оркестрация множества агентов в комплексные системы
- 📦 **Artifact Registry**: Управление реестрами образов контейнеров
- 🔐 **IAM аутентификация**: Безопасная интеграция с облачной платформой Cloud.ru
- ✅ **Валидация конфигураций**: Проверка корректности YAML/JSON файлов по схемам
- 🚀 **CI/CD интеграция**: Поддержка автоматизации развертывания в пайплайнах
- 🐳 **Docker интеграция**: Автоматическая сборка и загрузка образов контейнеров
- 📝 **Шаблоны проектов**: Готовые шаблоны для создания агентов и MCP серверов
- 🎨 **Красивый UI**: Цветной интерфейс с табами и emoji для лучшего UX
- 🌍 **Мультиязычность**: Поддержка русского и английского языков

---

## 🚀 Установка 

### Windows (winget)  (comming soon)

```bash
winget install CloudRu.AIAgentsCLI
```

### Windows (Scoop) (comming soon)

```powershell
scoop bucket add cloud-ru https://github.com/cloud-ru/scoop-bucket
scoop install ai-agents-cli
```

### macOS/Linux (Homebrew) (comming soon)

```bash
brew tap cloud-ru/evo-ai-agents-cli
brew install ai-agents-cli
```

### Ручная установка

```bash
# Скачайте последнюю версию
wget https://github.com/cloud-ru/evo-ai-agents-cli/releases/latest/download/ai-agents-cli-linux-amd64.tar.gz

# Распакуйте
tar -xzf ai-agents-cli-linux-amd64.tar.gz

# Установите
sudo mv ai-agents-cli /usr/local/bin/
chmod +x /usr/local/bin/ai-agents-cli
```

### Сборка из исходников

```bash
git clone https://github.com/cloud-ru/evo-ai-agents-cli.git
cd evo-ai-agents-cli
go build -o bin/ai-agents-cli .
```

---

## 🎯 Быстрый старт

### 1️⃣ Настройка переменных окружения

```bash
# Скопируйте пример конфигурации
cp env.example .env

# Отредактируйте с вашими данными
nano .env
```

Обязательные переменные:
```bash
IAM_KEY_ID=your-iam-key-id        # IAM Key ID
IAM_SECRET=your-iam-secret        # IAM Secret  
PROJECT_ID=your-project-id         # ID проекта
```

### 2️⃣ Проверка подключения

```bash
ai-agents-cli auth login           # Войти в систему
ai-agents-cli auth status          # Проверить статус авторизации
```

### 3️⃣ Валидация конфигурации

```bash
ai-agents-cli validate examples/agents.yaml        # Валидация файла
ai-agents-cli validate examples/                   # Валидация директории
```

### 4️⃣ Создание проектов из шаблонов

```bash
# Создать MCP сервер
ai-agents-cli create mcp my-mcp-server

# Создать AI агента (ADK, CrewAI, LangGraph)
ai-agents-cli create agent my-ai-agent --framework adk
ai-agents-cli create agent my-ai-agent --framework crewai
ai-agents-cli create agent my-ai-agent --framework langgraph

# С дополнительными параметрами
ai-agents-cli create agent my-ai-agent \
  --author "John Doe" \
  --python-version "3.11" \
  --framework adk
```

### 5️⃣ Развертывание

```bash
# Развертывание с автоматической сборкой и загрузкой Docker образов
ai-agents-cli agents deploy --build-image agents.yaml

# Предварительный просмотр без развертывания
ai-agents-cli agents deploy --dry-run agents.yaml

# Развертывание MCP серверов
ai-agents-cli mcp-servers deploy mcp-servers.yaml
```

---

## 📋 Доступные команды

### 🔐 Аутентификация (`auth`)

| Команда | Описание |
|---------|----------|
| `auth login` | Войти в систему с IAM учетными данными |
| `auth logout` | Выйти из системы |
| `auth status` | Проверить статус авторизации |
| `auth config` | Управление конфигурацией аутентификации |

### 🤖 Управление агентами (`agents`)

| Команда | Описание |
|---------|----------|
| `agents list` | Список всех агентов |
| `agents get <id>` | Информация об агенте |
| `agents create` | Создать нового агента |
| `agents deploy [file]` | Развертывание агентов из YAML |
| `agents marketplace` | Поиск агентов в маркетплейсе |

**Флаги для deploy:**
- `--build-image`, `-b` - Автоматическая сборка и загрузка Docker образов
- `--dry-run` - Предварительный просмотр без создания
- `--file`, `-f` - Путь к конфигурационному файлу

### 🔌 Управление MCP серверами (`mcp-servers`)

| Команда | Описание |
|---------|----------|
| `mcp-servers list` | Список всех MCP серверов |
| `mcp-servers get <id>` | Информация о сервере |
| `mcp-servers create` | Создать новый MCP сервер |
| `mcp-servers update <id>` | Обновить сервер |
| `mcp-servers delete <id>` | Удалить сервер |
| `mcp-servers deploy [file]` | Развертывание из конфигурации |
| `mcp-servers suspend <id>` | Приостановить сервер |
| `mcp-servers resume <id>` | Возобновить работу сервера |
| `mcp-servers history <id>` | История изменений |

### 🏗️ Управление системами агентов (`system`)

| Команда | Описание |
|---------|----------|
| `system list` | Список систем агентов |
| `system get <id>` | Информация о системе |
| `system create` | Создать новую систему |
| `system update <id>` | Обновить систему |
| `system delete <id>` | Удалить систему |
| `system deploy [file]` | Развертывание систем из конфигурации |
| `system suspend <id>` | Приостановить систему |
| `system resume <id>` | Возобновить работу системы |

### 📦 Управление Artifact Registry (`registry`)

| Команда | Описание |
|---------|----------|
| `registry create` | Создать новый реестр образов |
| `registry list` | Список всех реестров |
| `registry get <id>` | Информация о реестре |
| `registry delete <id>` | Удалить реестр |

### 🎨 Создание проектов (`create`)

| Команда | Описание |
|---------|----------|
| `create mcp [name]` | Создать проект MCP сервера |
| `create agent [name]` | Создать проект AI агента |

**Доступные фреймворки для агентов:**
- `adk` - Agent Development Kit (по умолчанию)
- `crewai` - CrewAI для командной работы
- `langgraph` - LangGraph с графом состояний

**Флаги:**
- `--path` - Путь для создания проекта
- `--author` - Имя автора
- `--python-version` - Версия Python
- `--framework` - Фреймворк для агента

### 🔧 CI/CD функции (`ci`)

| Команда | Описание |
|---------|----------|
| `ci status` | Проверить статус ресурсов |
| `ci logs` | Просмотр логов |

### ✅ Валидация (`validate`)

| Команда | Описание |
|---------|----------|
| `validate [file\|dir]` | Валидация конфигурационных файлов |
| `--file`, `-f` | Валидация конкретного файла |
| `--dir`, `-d` | Валидация директории |

---

## 💡 Примеры использования

### 📝 Полный workflow создания и развертывания агента

```bash
# 1. Создание проекта из шаблона
ai-agents-cli create agent my-customer-support-agent --framework adk

# 2. Переход в директорию проекта
cd my-customer-support-agent

# 3. Валидация конфигурации
ai-agents-cli validate agents.yaml

# 4. Развертывание с автоматической сборкой и загрузкой образа
ai-agents-cli agents deploy --build-image agents.yaml

# 5. Проверка статуса
ai-agents-cli agents list
```

### 🔌 Создание и развертывание MCP сервера

```bash
# 1. Создание MCP сервера
ai-agents-cli create mcp my-database-mcp

# 2. Валидация
ai-agents-cli validate mcp-servers.yaml

# 3. Развертывание
ai-agents-cli mcp-servers deploy mcp-servers.yaml

# 4. Проверка статуса
ai-agents-cli ci status mcp-servers
```

### 📦 Создание реестра и управление образами

```bash
# 1. Создание реестра
ai-agents-cli registry create --name my-images --description "My container registry"

# 2. Просмотр всех реестров
ai-agents-cli registry list

# 3. Информация о реестре
ai-agents-cli registry get my-images
```

### 🚀 CI/CD интеграция

#### GitHub Actions

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
      - name: Install AI Agents CLI
        run: brew install cloud-ru/evo-ai-agents-cli/ai-agents-cli
      
      - name: Validate configuration
        run: ai-agents-cli validate
        env:
          IAM_KEY_ID: ${{ secrets.IAM_KEY_ID }}
          IAM_SECRET: ${{ secrets.IAM_SECRET }}
          PROJECT_ID: ${{ secrets.PROJECT_ID }}
      
      - name: Deploy agents
        run: ai-agents-cli agents deploy --build-image
        env:
          IAM_KEY_ID: ${{ secrets.IAM_KEY_ID }}
          IAM_SECRET: ${{ secrets.IAM_SECRET }}
          PROJECT_ID: ${{ secrets.PROJECT_ID }}
          ARTIFACT_REGISTRY_URL: ${{ secrets.REGISTRY_URL }}
```

---

## 📁 Структура проекта

```
ai-agents-cli/
├── cmd/                    # CLI команды
│   ├── agent/             # Управление агентами
│   ├── mcp_server/        # Управление MCP серверами
│   ├── system/            # Управление системами агентов
│   ├── auth/              # Аутентификация
│   ├── registry/          # Управление реестрами
│   ├── create/            # Создание проектов из шаблонов
│   └── ci/                # CI/CD функции
├── internal/              # Внутренние пакеты
│   ├── api/               # API клиент для всех сервисов
│   ├── auth/               # IAM аутентификация
│   ├── deployer/           # Логика развертывания
│   ├── docker/             # Docker интеграция
│   ├── parser/             # Парсинг YAML с !include
│   ├── ui/                 # UI компоненты (табы, таблицы)
│   ├── validator/          # Валидатор конфигураций
│   └── scaffolder/         # Генератор проектов из шаблонов
│       └── templates/      # Шаблоны проектов
│           ├── agent-frameworks/  # ADK, CrewAI, LangGraph
│           └── mcp/               # MCP серверы
├── examples/              # Примеры конфигураций
│   ├── agents.yaml        # Пример конфигурации агентов
│   ├── mcp-servers.yaml   # Пример конфигурации MCP
│   └── agent-systems.yaml # Пример конфигурации систем
├── schemas/               # JSON схемы для валидации
│   └── schema.json        # Объединенная схема валидации
├── localizations/         # Локализация (ru/en)
├── scripts/               # Утилиты (установка, автообновление)
├── .goreleaser.yml        # Конфигурация GoReleaser
└── README.md              # Этот файл
```

---

## ⚙️ Конфигурация

### Переменные окружения

| Переменная | Описание | Обязательная | По умолчанию |
|------------|----------|--------------|--------------|
| `IAM_KEY_ID` | IAM Key ID для аутентификации | ✅ | - |
| `IAM_SECRET` | IAM Secret для аутентификации | ✅ | - |
| `PROJECT_ID` | ID проекта AI Agents | ✅ | - |
| `IAM_ENDPOINT` | IAM API endpoint | ❌ | `https://iam.api.cloud.ru` |
| `PUBLIC_API_ENDPOINT` | AI Agents API endpoint | ❌ | `ai-agents.api.cloud.ru` |
| `ARTIFACT_REGISTRY_URL` | URL Artifact Registry | ❌ | `cr.cloud.ru` |
| `SERVICE_LOG_LEVEL` | Уровень логирования | ❌ | `debug` |
| `SCAFFOLDER_PYTHON_VERSION` | Версия Python | ❌ | `3.9` |

### Поддерживаемые форматы конфигураций

CLI поддерживает YAML и JSON файлы с валидацией по JSON Schema:
- ✅ YAML (`.yaml`, `.yml`)
- ✅ JSON (`.json`)
- ✅ Директива `!include` для включения других файлов
- ✅ Валидация UUID форматов
- ✅ Валидация required полей
- ✅ Валидация массивов и их размеров

---

## 🎨 Особенности UI

### Цветовая схема
- 🟢 **Зеленый** - успешные операции, активные ресурсы
- 🔴 **Красный** - ошибки, неактивные ресурсы  
- 🟡 **Желтый** - предупреждения, приостановленные ресурсы
- 🔵 **Синий** - информация, заголовки
- ⚪ **Серый** - нейтральная информация

### Интерактивные компоненты
- 📊 Таблицы с данными
- 📑 Табы для детального просмотра
- 🎯 Спиннеры для прогресса операций
- 🎨 Emoji для визуального разделения

---

## 🧪 Тестирование

```bash
# Запуск всех тестов
go test ./...

# Тесты с покрытием
go test -cover ./...

# Конкретные тесты
go test ./internal/validator -v

# Интеграционные тесты
go test ./cmd -v
```

---

## 🛠️ Разработка

### Требования

- Go 1.24.3+
- Docker (для тестирования сборки образов)
- Git

### Сборка

```bash
# Установка зависимостей
go mod tidy

# Сборка для текущей платформы
go build -o bin/ai-agents-cli .

# Сборка для разных платформ
make build-all
```

### Добавление новых команд

1. Создайте новый пакет в `cmd/`
2. Реализуйте команду с помощью Cobra
3. Добавьте API методы в `internal/api/`
4. Напишите тесты
5. Обновите документацию

---

## 📚 Документация

- 📖 [Руководство по валидации](TESTING.md)
- 📥 [Установка и настройка](.github/INSTALL.md)
- 💡 [Примеры использования](examples/usage.md)
- 🔌 [API документация](service.swagger.json)
- 🌐 [Cloud.ru AI Agents Docs](https://cloud.ru/docs/ai-agents/ug/index?source-platform=Evolution)

---

## 🤝 Вклад в проект

Мы приветствуем вклад в проект! 

1. Форкните репозиторий
2. Создайте ветку для вашей функции (`git checkout -b feature/amazing-feature`)
3. Внесите изменения и добавьте тесты
4. Закоммитьте изменения (`git commit -m 'Add some amazing feature'`)
5. Запушьте в ветку (`git push origin feature/amazing-feature`)
6. Откройте Pull Request

### Документация для разработчиков

- [Guide to Contributing](CONTRIBUTION_GUIDE.md)
- [Testing Guide](TESTING.md)

---

## 📄 Лицензия

Этот проект лицензирован под MIT License - см. файл [LICENSE](LICENSE) для деталей.

---

## 🆘 Поддержка

- 📧 **Email**: support@cloud.ru
- 📖 **Документация**: [Cloud.ru AI Agents](https://cloud.ru/docs/ai-agents/ug/index?source-platform=Evolution)
- 🐛 **Баги**: [GitHub Issues](https://github.com/cloud-ru/evo-ai-agents-cli/issues)
- 💬 **Обсуждения**: [GitHub Discussions](https://github.com/cloud-ru/evo-ai-agents-cli/discussions)
- 💬 **Вопросы**: [GitHub Q&A](https://github.com/cloud-ru/evo-ai-agents-cli/discussions/categories/q-a)

---

## 🎉 Благодарности

- Всем контрибьюторам проекта
- Команде Cloud.ru за поддержку
- Сообществу пользователей за обратную связь

---

<div align="center">

**[⬆ Вернуться к началу](#-ai-agents-cli)**

Made with ❤️ by [Cloud.ru](https://cloud.ru)

</div>
