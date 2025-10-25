# Скрипт для установки автодополнения ai-agents-cli для Fish shell
# Запуск: fish scripts/install-fish-completion.fish

# Цвета для вывода
set -g RED '\033[0;31m'
set -g GREEN '\033[0;32m'
set -g YELLOW '\033[1;33m'
set -g BLUE '\033[0;34m'
set -g NC '\033[0m' # No Color

# Функции для вывода сообщений
function log
    echo -e "$BLUE[INFO]$NC $argv"
end

function success
    echo -e "$GREEN[SUCCESS]$NC $argv"
end

function warning
    echo -e "$YELLOW[WARNING]$NC $argv"
end

function error
    echo -e "$RED[ERROR]$NC $argv"
end

# Проверка наличия ai-agents-cli
if not command -v ai-agents-cli > /dev/null
    error "ai-agents-cli не найден в PATH. Установите CLI сначала."
    exit 1
end

log "Установка автодополнения для Fish shell..."

# Создание директории для completion
set completion_dir "$HOME/.config/fish/completions"
set completion_file "$completion_dir/ai-agents-cli.fish"

# Создание директории если не существует
if not test -d "$completion_dir"
    log "Создание директории $completion_dir"
    mkdir -p "$completion_dir"
end

# Генерация completion
log "Генерация автодополнения..."
ai-agents-cli completion fish > "$completion_file"

# Проверка, что файл создался
if test -f "$completion_file"
    success "Автодополнение для fish установлено в $completion_file"
    
    # Проверка размера файла
    set file_size (stat -f%z "$completion_file" 2>/dev/null || stat -c%s "$completion_file" 2>/dev/null || echo "0")
    if test "$file_size" -gt 0
        success "Файл автодополнения создан успешно (размер: $file_size байт)"
    else
        warning "Файл автодополнения пуст, возможно есть проблемы с генерацией"
    end
    
    # Информация о применении
    echo ""
    success "Автодополнение установлено успешно!"
    echo ""
    echo "Для применения изменений:"
    echo "  1. Перезапустите Fish shell, или"
    echo "  2. Выполните: source ~/.config/fish/config.fish"
    echo ""
    echo "Теперь вы можете использовать Tab для автодополнения команд ai-agents-cli"
    echo ""
    echo "Примеры использования:"
    echo "  ai-agents-cli <TAB>          # Показать все команды"
    echo "  ai-agents-cli create <TAB>   # Показать подкоманды create"
    echo "  ai-agents-cli create mcp <TAB> # Показать флаги для mcp"
    
else
    error "Не удалось создать файл автодополнения для fish"
    exit 1
end
