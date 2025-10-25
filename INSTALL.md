# 🚀 Установка AI Agents CLI

## Быстрая установка

### 1. Автоматическая установка (рекомендуется)
```bash
# Собрать, установить и настроить CLI
make setup
```

### 2. Ручная установка
```bash
# Собрать и установить в пользовательскую директорию
make build-install

# Настроить PATH
./scripts/setup.sh
```

## Способы установки

### 🏠 Пользовательская установка (без sudo)
```bash
# Собрать и установить
make build-install

# Настроить PATH
./scripts/setup.sh

# Перезагрузить shell
source ~/.bashrc  # или ~/.zshrc
```

### 🌐 Системная установка (требует sudo)
```bash
# Собрать и установить в систему
make install

# CLI будет доступен из любой директории
ai-agents-cli --help
```

## Проверка установки

```bash
# Проверить, что CLI работает
ai-agents-cli --version

# Проверить все команды
ai-agents-cli --help

# Проверить команды создания проектов
ai-agents-cli create --help
```

## Использование

### Создание MCP сервера
```bash
# Создать новый MCP сервер
ai-agents-cli create mcp my-mcp-server

# С кастомными параметрами
ai-agents-cli create mcp my-server --author "John Doe" --python-version "3.11"
```

### Создание AI агента
```bash
# Создать нового AI агента
ai-agents-cli create agent my-ai-agent

# С кастомными параметрами
ai-agents-cli create agent my-agent --author "Jane Smith" --python-version "3.10"
```

## Автодополнение

### Bash
```bash
# Генерировать автодополнение для bash
ai-agents-cli completion bash > ~/.local/share/bash-completion/completions/ai-agents-cli

# Добавить в .bashrc
echo 'source ~/.local/share/bash-completion/completions/ai-agents-cli' >> ~/.bashrc
```

### Zsh
```bash
# Создать директорию для автодополнения
mkdir -p ~/.zsh/completions

# Генерировать автодополнение для zsh
ai-agents-cli completion zsh > ~/.zsh/completions/_ai-agents-cli

# Добавить в .zshrc
echo 'fpath=(~/.zsh/completions $fpath)' >> ~/.zshrc
echo 'autoload -U compinit && compinit' >> ~/.zshrc
```

## Удаление

### Пользовательская установка
```bash
# Удалить из пользовательской директории
rm ~/.local/bin/ai-agents-cli

# Удалить из PATH (опционально)
# Отредактируйте ~/.bashrc или ~/.zshrc
```

### Системная установка
```bash
# Удалить из системы
make uninstall
```

## Разработка

### Локальная разработка
```bash
# Собрать для тестирования
make build

# Запустить из локальной директории
./bin/ai-agents-cli --help

# Запустить в режиме разработки
make dev
```

### Тестирование
```bash
# Запустить тесты
make test

# Запустить тесты с покрытием
make test-coverage

# Запустить линтер
make lint
```

## Troubleshooting

### CLI не найден
```bash
# Проверить PATH
echo $PATH

# Проверить, что ~/.local/bin в PATH
ls ~/.local/bin/ai-agents-cli

# Добавить в PATH вручную
export PATH=$HOME/.local/bin:$PATH
```

### Проблемы с автодополнением
```bash
# Перезагрузить shell
source ~/.bashrc  # или ~/.zshrc

# Проверить автодополнение
ai-agents-cli <TAB>
```

### Проблемы с правами
```bash
# Сделать CLI исполняемым
chmod +x ~/.local/bin/ai-agents-cli

# Проверить права
ls -la ~/.local/bin/ai-agents-cli
```

## Обновление

```bash
# Обновить CLI
git pull
make build-install

# Или полная настройка
make setup
```

## Поддержка

- 📖 **Документация**: README.md
- 🐛 **Баг-репорты**: GitHub Issues
- 💬 **Обсуждения**: GitHub Discussions
- 📧 **Контакты**: Cloud.ru Team
