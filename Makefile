# AI Agents CLI Makefile

# Переменные
BINARY_NAME=ai-agents-cli
BUILD_DIR=bin
VERSION?=1.0.0
BUILD_TIME=$(shell date +%Y-%m-%d_%H:%M:%S)
GIT_COMMIT=$(shell git rev-parse --short HEAD)
LDFLAGS=-ldflags "-X main.Version=$(VERSION) -X main.BuildTime=$(BUILD_TIME) -X main.GitCommit=$(GIT_COMMIT)"

# Цвета для вывода
GREEN=\033[0;32m
YELLOW=\033[1;33m
RED=\033[0;31m
NC=\033[0m # No Color

.PHONY: help build clean test lint run install deps validate examples

# Помощь
help: ## Показать справку
	@echo "$(GREEN)AI Agents CLI - Makefile$(NC)"
	@echo ""
	@echo "$(YELLOW)Доступные команды:$(NC)"
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  $(GREEN)%-15s$(NC) %s\n", $$1, $$2}' $(MAKEFILE_LIST)

# Сборка
build: ## Собрать бинарный файл
	@echo "$(YELLOW)Сборка $(BINARY_NAME)...$(NC)"
	@mkdir -p $(BUILD_DIR)
	@go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) .
	@echo "$(GREEN)✅ Сборка завершена: $(BUILD_DIR)/$(BINARY_NAME)$(NC)"

# Сборка для разных платформ
build-all: ## Собрать для всех платформ
	@echo "$(YELLOW)Сборка для всех платформ...$(NC)"
	@mkdir -p $(BUILD_DIR)
	@GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 .
	@GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64 .
	@GOOS=darwin GOARCH=arm64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-arm64 .
	@GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe .
	@echo "$(GREEN)✅ Сборка для всех платформ завершена$(NC)"

# Очистка
clean: ## Очистить собранные файлы
	@echo "$(YELLOW)Очистка...$(NC)"
	@rm -rf $(BUILD_DIR)
	@go clean
	@echo "$(GREEN)✅ Очистка завершена$(NC)"

# Тестирование
test: ## Запустить тесты
	@echo "$(YELLOW)Запуск тестов...$(NC)"
	@go test -v ./...
	@echo "$(GREEN)✅ Тесты завершены$(NC)"

# Тестирование с покрытием
test-coverage: ## Запустить тесты с покрытием
	@echo "$(YELLOW)Запуск тестов с покрытием...$(NC)"
	@go test -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "$(GREEN)✅ Отчет о покрытии создан: coverage.html$(NC)"

# Линтинг
lint: ## Запустить линтер
	@echo "$(YELLOW)Запуск линтера...$(NC)"
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run; \
	else \
		echo "$(RED)❌ golangci-lint не установлен. Установите: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest$(NC)"; \
		exit 1; \
	fi
	@echo "$(GREEN)✅ Линтинг завершен$(NC)"

# Форматирование кода
fmt: ## Форматировать код
	@echo "$(YELLOW)Форматирование кода...$(NC)"
	@go fmt ./...
	@echo "$(GREEN)✅ Форматирование завершено$(NC)"

# Валидация
validate: ## Валидировать конфигурационные файлы
	@echo "$(YELLOW)Валидация конфигурационных файлов...$(NC)"
	@go run . validate
	@echo "$(GREEN)✅ Валидация завершена$(NC)"

# Запуск
run: build ## Собрать и запустить CLI
	@echo "$(YELLOW)Запуск CLI...$(NC)"
	@./$(BUILD_DIR)/$(BINARY_NAME) --help

# Установка зависимостей
deps: ## Установить зависимости
	@echo "$(YELLOW)Установка зависимостей...$(NC)"
	@go mod download
	@go mod tidy
	@echo "$(GREEN)✅ Зависимости установлены$(NC)"

# Установка CLI
install: build ## Установить CLI в систему
	@echo "$(YELLOW)Установка CLI...$(NC)"
	@sudo cp $(BUILD_DIR)/$(BINARY_NAME) /usr/local/bin/
	@echo "$(GREEN)✅ CLI установлен в /usr/local/bin/$(NC)"

# Удаление из системы
uninstall: ## Удалить CLI из системы
	@echo "$(YELLOW)Удаление CLI...$(NC)"
	@sudo rm -f /usr/local/bin/$(BINARY_NAME)
	@echo "$(GREEN)✅ CLI удален$(NC)"

# Примеры
examples: ## Запустить примеры
	@echo "$(YELLOW)Запуск примеров...$(NC)"
	@echo "$(GREEN)Примеры использования CLI:$(NC)"
	@echo ""
	@echo "$(YELLOW)1. Проверка статуса:$(NC)"
	@echo "   ./$(BUILD_DIR)/$(BINARY_NAME) ci status"
	@echo ""
	@echo "$(YELLOW)2. Список MCP серверов:$(NC)"
	@echo "   ./$(BUILD_DIR)/$(BINARY_NAME) mcp-servers list"
	@echo ""
	@echo "$(YELLOW)3. Валидация конфигурации:$(NC)"
	@echo "   ./$(BUILD_DIR)/$(BINARY_NAME) validate"
	@echo ""
	@echo "$(YELLOW)4. Развертывание:$(NC)"
	@echo "   ./$(BUILD_DIR)/$(BINARY_NAME) mcp-servers deploy --dry-run"

# Разработка
dev: ## Запустить в режиме разработки
	@echo "$(YELLOW)Запуск в режиме разработки...$(NC)"
	@go run . --help

# Проверка перед коммитом
pre-commit: fmt lint test ## Выполнить все проверки перед коммитом
	@echo "$(GREEN)✅ Все проверки пройдены$(NC)"

# Создание релиза
release: clean build-all ## Создать релиз
	@echo "$(YELLOW)Создание релиза...$(NC)"
	@mkdir -p release
	@cp $(BUILD_DIR)/* release/
	@echo "$(GREEN)✅ Релиз создан в директории release/$(NC)"

# Показать версию
version: ## Показать версию
	@echo "$(GREEN)Версия: $(VERSION)$(NC)"
	@echo "$(GREEN)Git commit: $(GIT_COMMIT)$(NC)"
	@echo "$(GREEN)Build time: $(BUILD_TIME)$(NC)"

# По умолчанию
.DEFAULT_GOAL := help
