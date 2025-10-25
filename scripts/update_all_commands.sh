#!/bin/bash

# Скрипт для обновления всех команд с новой системой обработки ошибок

echo "🔄 Обновление всех команд с новой системой обработки ошибок..."

# Список файлов для обновления
FILES=(
    "cmd/mcp_server/get.go"
    "cmd/mcp_server/create.go"
    "cmd/mcp_server/delete.go"
    "cmd/mcp_server/deploy.go"
    "cmd/mcp_server/history.go"
    "cmd/mcp_server/resume.go"
    "cmd/mcp_server/suspend.go"
    "cmd/mcp_server/update.go"
    "cmd/system/get.go"
    "cmd/system/create.go"
    "cmd/system/delete.go"
    "cmd/system/deploy.go"
    "cmd/system/resume.go"
    "cmd/system/suspend.go"
    "cmd/system/update.go"
    "cmd/agent/get.go"
    "cmd/agent/deploy.go"
    "cmd/agent/marketplace.go"
    "cmd/ci/logs.go"
    "cmd/ci/status.go"
    "cmd/deploy.go"
    "cmd/validate.go"
)

# Функция для обновления импортов
update_imports() {
    local file="$1"
    echo "📝 Обновление импортов в $file"
    
    # Заменяем импорты
    sed -i '' 's|"github.com/charmbracelet/log"|"github.com/cloud-ru/evo-ai-agents-cli/internal/errors"|g' "$file"
    
    # Добавляем os импорт если его нет
    if ! grep -q '"os"' "$file"; then
        sed -i '' '/import (/a\
	"os"
' "$file"
    fi
}

# Функция для обновления обработки ошибок
update_error_handling() {
    local file="$1"
    echo "🔧 Обновление обработки ошибок в $file"
    
    # Заменяем log.Fatal на новую систему
    sed -i '' 's|log\.Fatal(\([^,]*\), "error", err)|errorHandler := errors.NewHandler()\
			appErr := errorHandler.WrapAPIError(err, "API_ERROR", \1)\
			appErr = appErr.WithSuggestions(\
				"Проверьте переменные окружения: IAM_KEY_ID, IAM_SECRET_KEY, IAM_ENDPOINT",\
				"Убедитесь что вы авторизованы: ai-agents-cli auth login или в папке выполнения команды лежит .env файл с перемнными выше",\
				"Проверьте доступность API: curl -I $IAM_ENDPOINT",\
				"Обратитесь к администратору для получения учетных данных",\
				"📚 Подробная документация: https://cloud.ru/docs/ai-agents/ug/index?source-platform=Evolution",\
			)\
			fmt.Println(errorHandler.HandlePlain(appErr))\
			os.Exit(1)|g' "$file"
    
    # Заменяем log.Error на новую систему
    sed -i '' 's|log\.Error(\([^,]*\), "error", err)|errorHandler := errors.NewHandler()\
			appErr := errorHandler.WrapAPIError(err, "API_ERROR", \1)\
			appErr = appErr.WithSuggestions(\
				"Проверьте переменные окружения: IAM_KEY_ID, IAM_SECRET_KEY, IAM_ENDPOINT",\
				"Убедитесь что вы авторизованы: ai-agents-cli auth login или в папке выполнения команды лежит .env файл с перемнными выше",\
				"Проверьте доступность API: curl -I $IAM_ENDPOINT",\
				"Обратитесь к администратору для получения учетных данных",\
				"📚 Подробная документация: https://cloud.ru/docs/ai-agents/ug/index?source-platform=Evolution",\
			)\
			fmt.Println(errorHandler.HandlePlain(appErr))|g' "$file"
}

# Обновляем все файлы
for file in "${FILES[@]}"; do
    if [ -f "$file" ]; then
        echo "🔄 Обработка $file"
        update_imports "$file"
        update_error_handling "$file"
        echo "✅ $file обновлен"
    else
        echo "⚠️  Файл $file не найден"
    fi
done

echo "🎉 Все команды обновлены с новой системой обработки ошибок!"
