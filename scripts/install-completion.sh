#!/bin/bash

# Скрипт для установки автодополнения ai-agents-cli
# Поддерживает bash, zsh, fish

set -e

# Цвета для вывода
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Функция для вывода сообщений
log() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Проверка наличия ai-agents-cli
if ! command -v ai-agents-cli &> /dev/null; then
    error "ai-agents-cli не найден в PATH. Установите CLI сначала."
    exit 1
fi

# Определение оболочки
SHELL_NAME=$(basename "$SHELL")

log "Обнаружена оболочка: $SHELL_NAME"

# Функция для установки bash completion
install_bash_completion() {
    local completion_dir="$HOME/.bash_completion.d"
    local completion_file="$completion_dir/ai-agents-cli"
    
    log "Установка автодополнения для bash..."
    
    # Создание директории если не существует
    mkdir -p "$completion_dir"
    
    # Генерация completion
    ai-agents-cli completion bash > "$completion_file"
    
    # Проверка наличия в .bashrc
    if ! grep -q "source.*bash_completion.d" "$HOME/.bashrc" 2>/dev/null; then
        log "Добавление автозагрузки в .bashrc..."
        echo "" >> "$HOME/.bashrc"
        echo "# Автодополнение для ai-agents-cli" >> "$HOME/.bashrc"
        echo "if [ -d ~/.bash_completion.d ]; then" >> "$HOME/.bashrc"
        echo "    for file in ~/.bash_completion.d/*; do" >> "$HOME/.bashrc"
        echo "        [ -f \"\$file\" ] && source \"\$file\"" >> "$HOME/.bashrc"
        echo "    done" >> "$HOME/.bashrc"
        echo "fi" >> "$HOME/.bashrc"
    fi
    
    success "Автодополнение для bash установлено в $completion_file"
}

# Функция для установки zsh completion
install_zsh_completion() {
    local completion_dir="$HOME/.zsh/completions"
    local completion_file="$completion_dir/_ai-agents-cli"
    
    log "Установка автодополнения для zsh..."
    
    # Создание директории если не существует
    mkdir -p "$completion_dir"
    
    # Генерация completion
    ai-agents-cli completion zsh > "$completion_file"
    
    # Проверка наличия в .zshrc
    if ! grep -q "fpath.*completions" "$HOME/.zshrc" 2>/dev/null; then
        log "Добавление автозагрузки в .zshrc..."
        echo "" >> "$HOME/.zshrc"
        echo "# Автодополнение для ai-agents-cli" >> "$HOME/.zshrc"
        echo "fpath=(~/.zsh/completions \$fpath)" >> "$HOME/.zshrc"
        echo "autoload -U compinit && compinit" >> "$HOME/.zshrc"
    fi
    
    success "Автодополнение для zsh установлено в $completion_file"
}

# Функция для установки fish completion
install_fish_completion() {
    local completion_dir="$HOME/.config/fish/completions"
    local completion_file="$completion_dir/ai-agents-cli.fish"
    
    log "Установка автодополнения для fish..."
    
    # Создание директории если не существует
    mkdir -p "$completion_dir"
    
    # Генерация completion
    ai-agents-cli completion fish > "$completion_file"
    
    # Проверка, что файл создался
    if [ -f "$completion_file" ]; then
        success "Автодополнение для fish установлено в $completion_file"
        log "Fish автоматически загрузит completion при следующем запуске"
    else
        error "Не удалось создать файл автодополнения для fish"
        exit 1
    fi
}

# Основная логика
case "$SHELL_NAME" in
    "bash")
        install_bash_completion
        ;;
    "zsh")
        install_zsh_completion
        ;;
    "fish")
        install_fish_completion
        ;;
    *)
        warning "Неподдерживаемая оболочка: $SHELL_NAME"
        echo "Поддерживаемые оболочки: bash, zsh, fish"
        echo ""
        echo "Для ручной установки используйте:"
        echo "  ai-agents-cli completion bash > ~/.bash_completion.d/ai-agents-cli"
        echo "  ai-agents-cli completion zsh > ~/.zsh/completions/_ai-agents-cli"
        echo "  ai-agents-cli completion fish > ~/.config/fish/completions/ai-agents-cli.fish"
        exit 1
        ;;
esac

echo ""
success "Автодополнение установлено успешно!"
echo ""
echo "Для применения изменений:"
echo "  1. Перезапустите терминал, или"
echo "  2. Выполните: source ~/.${SHELL_NAME}rc"
echo ""
echo "Теперь вы можете использовать Tab для автодополнения команд ai-agents-cli"
