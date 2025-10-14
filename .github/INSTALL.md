# Установка AI Agents CLI

## 🚀 Быстрая установка

### Windows (winget)

```bash
# Установка через winget
winget install CloudRu.AIAgentsCLI

# Обновление
winget upgrade CloudRu.AIAgentsCLI

# Удаление
winget uninstall CloudRu.AIAgentsCLI
```

### macOS/Linux (Homebrew)

```bash
# Добавление tap (если нужно)
brew tap cloud-ru/evo-ai-agents-cli

# Установка
brew install ai-agents-cli

# Обновление
brew upgrade ai-agents-cli

# Удаление
brew uninstall ai-agents-cli
```

## 📦 Ручная установка

### 1. Скачивание бинарного файла

Перейдите на [страницу релизов](https://github.com/cloud-ru/evo-ai-agents-cli/releases) и скачайте подходящий файл для вашей платформы:

- **Windows**: `ai-agents-cli-windows-amd64.zip`
- **macOS (Intel)**: `ai-agents-cli-darwin-amd64.tar.gz`
- **macOS (Apple Silicon)**: `ai-agents-cli-darwin-arm64.tar.gz`
- **Linux (Intel)**: `ai-agents-cli-linux-amd64.tar.gz`
- **Linux (ARM)**: `ai-agents-cli-linux-arm64.tar.gz`

### 2. Распаковка и установка

#### Windows

```powershell
# Распакуйте архив
Expand-Archive ai-agents-cli-windows-amd64.zip -DestinationPath C:\ai-agents-cli

# Добавьте в PATH (временно)
$env:PATH += ";C:\ai-agents-cli"

# Или добавьте в PATH постоянно через системные настройки
```

#### macOS/Linux

```bash
# Распакуйте архив
tar -xzf ai-agents-cli-darwin-amd64.tar.gz

# Переместите в системную директорию
sudo mv ai-agents-cli /usr/local/bin/

# Сделайте исполняемым
sudo chmod +x /usr/local/bin/ai-agents-cli
```

### 3. Проверка установки

```bash
# Проверьте версию
ai-agents-cli --version

# Проверьте справку
ai-agents-cli --help
```

## 🔧 Настройка

### 1. Переменные окружения

Создайте файл `.env` или установите переменные окружения:

```bash
# .env файл
export IAM_KEY_ID="your-iam-key-id"
export IAM_SECRET="your-iam-secret"
export PROJECT_ID="your-project-id"
export IAM_ENDPOINT="https://iam.api.cloud.ru"
export PUBLIC_API_ENDPOINT="ai-agents.api.cloud.ru"
```

### 2. Проверка подключения

```bash
# Проверьте статус системы
ai-agents-cli ci status
```

## 🐳 Docker

### Использование Docker образа

```bash
# Запуск через Docker
docker run --rm -it \
  -e API_KEY="your-api-key" \
  -e PROJECT_ID="your-project-id" \
  cloudru/ai-agents-cli:latest --help
```

### Создание собственного образа

```dockerfile
FROM golang:1.21-alpine AS builder

WORKDIR /app
COPY . .
RUN go build -o ai-agents-cli .

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/ai-agents-cli .
ENTRYPOINT ["./ai-agents-cli"]
```

## 🔄 Обновление

### Автоматическое обновление

```bash
# Windows (winget)
winget upgrade CloudRu.AIAgentsCLI

# macOS/Linux (Homebrew)
brew upgrade ai-agents-cli
```

### Ручное обновление

1. Скачайте новую версию с [GitHub Releases](https://github.com/cloudru/ai-agents-cli/releases)
2. Замените старый бинарный файл новым
3. Перезапустите терминал

## 🛠️ Разработка

### Установка из исходного кода

```bash
# Клонируйте репозиторий
git clone https://github.com/cloud-ru/evo-ai-agents-cli.git
cd evo-ai-agents-cli

# Установите зависимости
go mod download

# Соберите проект
make build

# Или используйте go build
go build -o ai-agents-cli .
```

### Установка через go install

```bash
# Установка последней версии
go install github.com/cloud-ru/evo-ai-agents-cli@latest

# Установка конкретной версии
go install github.com/cloud-ru/evo-ai-agents-cli@v1.0.0
```

## 🔍 Устранение неполадок

### Проблема: "command not found"

**Решение:**
- Убедитесь, что бинарный файл находится в PATH
- Перезапустите терминал
- Проверьте права доступа к файлу

### Проблема: "IAM_KEY_ID environment variable is required"

**Решение:**
- Установите переменную окружения IAM_KEY_ID
- Установите переменную окружения IAM_SECRET
- Проверьте правильность IAM ключей
- Убедитесь, что переменные экспортированы

### Проблема: "Permission denied"

**Решение:**
```bash
# Сделайте файл исполняемым
chmod +x ai-agents-cli

# Или запустите с sudo (не рекомендуется)
sudo ./ai-agents-cli
```

## 📚 Дополнительные ресурсы

- [Документация](https://github.com/cloud-ru/evo-ai-agents-cli/blob/main/README.md)
- [Примеры использования](https://github.com/cloud-ru/evo-ai-agents-cli/tree/main/examples)
- [CI/CD интеграция](https://github.com/cloud-ru/evo-ai-agents-cli/tree/main/.github/workflows)
- [Сообщить об ошибке](https://github.com/cloud-ru/evo-ai-agents-cli/issues)

## 🆘 Поддержка

Если у вас возникли проблемы с установкой:

1. Проверьте [Issues](https://github.com/cloud-ru/evo-ai-agents-cli/issues)
2. Создайте новый Issue с подробным описанием проблемы
3. Обратитесь к команде разработки

---

**AI Agents CLI** - ваш надежный помощник в управлении AI агентами! 🚀
