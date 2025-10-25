package create

import (
	"github.com/spf13/cobra"
)

// RootCMD represents the create command
var RootCMD = &cobra.Command{
	Use:   "create",
	Short: "Создание проектов из шаблонов",
	Long: `Создание новых проектов MCP серверов и AI агентов из готовых шаблонов.

Команды создания проектов:
• mcp [project-name] - Создать проект MCP сервера
• agent [project-name] - Создать проект AI агента

Каждый шаблон включает:
• Полную структуру Python проекта
• Dockerfile и docker-compose.yml
• CI/CD конфигурации (GitLab CI и GitHub Actions)
• Makefile с командами разработки
• README.md с документацией
• .gitignore, .editorconfig, LICENSE
• Базовые зависимости и настройки`,
}
