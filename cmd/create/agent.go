package create

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/cloud-ru/evo-ai-agents-cli/internal/errors"
	"github.com/cloud-ru/evo-ai-agents-cli/internal/scaffolder"
	"github.com/cloud-ru/evo-ai-agents-cli/internal/ui"
	"github.com/spf13/cobra"
)

var (
	agentProjectPath string
	agentAuthor      string
)

// createAgentCmd represents the agent create command
var createAgentCmd = &cobra.Command{
	Use:   "agent [project-name]",
	Short: "Создать новый проект AI агента из шаблона",
	Long: `Создает новый проект AI агента из готового шаблона.

AI агент предназначен для автоматизации задач и интеллектуальной обработки данных.
Шаблон включает полную структуру Python проекта с FastAPI, Docker конфигурацией,
CI/CD пайплайнами и документацией.

Команда запускает интерактивную форму для настройки всех параметров проекта.
`,
	Args: cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// Создаем обработчик ошибок
		errorHandler := errors.NewHandler()

		// Always use full-screen TUI form for better UX
		formData, err := ui.RunProjectForm("agent")
		if err != nil {
			appErr := errorHandler.WrapUserError(err, "FORM_ERROR", "Ошибка при заполнении формы")
			fmt.Println(errorHandler.Handle(appErr))
			os.Exit(1)
		}

		projectName := formData.ProjectName
		if projectName == "" {
			appErr := errors.ValidationError("MISSING_PROJECT_NAME", "Название проекта обязательно")
			fmt.Println(errorHandler.Handle(appErr))
			os.Exit(1)
		}

		targetPath := formData.ProjectPath
		if targetPath == "" {
			targetPath = projectName
		} else {
			targetPath = filepath.Join(targetPath, projectName)
		}

		author := formData.Author
		framework := formData.Framework
		cicdTypeStr := formData.CICDType

		// Collect options
		var options []string
		if formData.GitInit {
			options = append(options, "git_init")
		}
		if formData.CreateEnv {
			options = append(options, "create_env")
		}
		if formData.InstallDeps {
			options = append(options, "install_deps")
		}

		// Show project info
		fmt.Println(lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("205")).
			Render("🤖 Создание проекта AI агента"))

		fmt.Printf("Название проекта: %s\n", projectName)
		fmt.Printf("Путь: %s\n", targetPath)
		fmt.Printf("Автор: %s\n", author)
		fmt.Printf("Фреймворк: %s\n", framework)
		fmt.Printf("CI/CD система: %s\n", cicdTypeStr)
		if len(options) > 0 {
			fmt.Printf("Дополнительные опции: %s\n", strings.Join(options, ", "))
		}
		fmt.Println()

		// Create scaffolder with custom config
		config := &scaffolder.ScaffolderConfig{
			Author:      author,
			DefaultCICD: cicdTypeStr,
		}
		scaffolderInstance := scaffolder.NewScaffolderWithConfig(config)

		// Validate template
		if err := scaffolderInstance.ValidateTemplate("agent"); err != nil {
			appErr := errorHandler.WrapTemplateError(err, "TEMPLATE_VALIDATION_FAILED", "Ошибка валидации шаблона агента")
			fmt.Println(errorHandler.Handle(appErr))
			os.Exit(1)
		}

		// Create project
		fmt.Println(lipgloss.NewStyle().
			Foreground(lipgloss.Color("240")).
			Render("Создание проекта..."))

		if err := scaffolderInstance.CreateProjectWithOptions("agent", projectName, targetPath, cicdTypeStr, framework, options); err != nil {
			appErr := errorHandler.WrapFileSystemError(err, "PROJECT_CREATION_FAILED", "Ошибка создания проекта агента")
			appErr = appErr.WithSuggestions(
				"Проверьте права доступа к директории: ls -la "+targetPath,
				"Убедитесь что директория существует: mkdir -p "+targetPath,
				"Проверьте свободное место на диске: df -h",
				"Попробуйте создать проект в другой директории",
				"📚 Документация: https://cloud.ru/docs/ai-agents/ug/index?source-platform=Evolution",
			)
			fmt.Println(errorHandler.Handle(appErr))
			os.Exit(1)
		}

		// Show success message
		successStyle := lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("2")).
			Border(lipgloss.RoundedBorder()).
			Padding(1, 2)

		fmt.Println(successStyle.Render("✅ Проект AI агента создан успешно!"))
		fmt.Println()

		// Show next steps
		nextSteps := lipgloss.NewStyle().
			Foreground(lipgloss.Color("252")).
			Render(fmt.Sprintf(`
Следующие шаги:

1. Перейдите в директорию проекта:
   cd %s

2. Установите зависимости:
   make install
   # или
   uv sync
   # или (fallback)
   pip install -r requirements.txt

3. Настройте переменные окружения:
   cp env.example .env
   # Отредактируйте .env файл
   # Добавьте API ключи для AI сервисов

4. Запустите агента:
   make run
   # или
   python src/agent.py

5. Для Docker:
   make docker-build
   make docker-run

📖 Документация: README.md
🔧 Команды разработки: make help
🤖 Тестирование агента: POST /agent/process
`, targetPath))

		fmt.Println(nextSteps)
	},
}

func init() {
	RootCMD.AddCommand(createAgentCmd)

	createAgentCmd.Flags().StringVarP(&agentProjectPath, "path", "p", "", "Путь для создания проекта (по умолчанию: текущая директория)")
	createAgentCmd.Flags().StringVarP(&agentAuthor, "author", "a", "", "Автор проекта (по умолчанию: из git config или 'Cloud.ru Team')")
}
