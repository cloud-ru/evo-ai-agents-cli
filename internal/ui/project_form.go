package ui

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/charmbracelet/huh"
)

// ProjectFormData represents the data collected from the form
type ProjectFormData struct {
	ProjectName string
	ProjectPath string
	Author      string
	Framework   string // New field for agent framework selection
	CICDType    string
	GitInit     bool
	CreateEnv   bool
	InstallDeps bool
}

// RunProjectForm runs the project creation form using huh
func RunProjectForm(projectType string) (*ProjectFormData, error) {
	// Get default author from git config
	defaultAuthor := getGitAuthorFromConfig()
	if defaultAuthor == "" {
		defaultAuthor = "Cloud.ru Team"
	}

	// Form data with default values
	formData := ProjectFormData{
		Author:      defaultAuthor,
		CICDType:    "both",
		GitInit:     true,
		CreateEnv:   true,
		InstallDeps: false,
	}

	// Set default framework and project names
	if projectType == "agent" {
		formData.Framework = "adk"
		formData.ProjectName = "my-awesome-agent"
		formData.ProjectPath = "./"
	} else {
		formData.ProjectName = "my-awesome-mcp"
		formData.ProjectPath = "./"
	}

	// Create form fields based on project type
	var form *huh.Form

	if projectType == "agent" {
		// Agent project form with framework selection
		form = huh.NewForm(
			huh.NewGroup(
				// Project name
				huh.NewInput().
					Title("🚀 Название проекта").
					Description("Введите название вашего проекта").
					Placeholder("my-awesome-agent").
					Value(&formData.ProjectName).
					Validate(func(str string) error {
						// Accept non-empty strings or default values
						if str == "" {
							return fmt.Errorf("название проекта обязательно")
						}
						return nil
					}),

				// Project path
				huh.NewInput().
					Title("📁 Путь для создания").
					Description("Директория для создания проекта (оставьте пустым для текущей)").
					Placeholder("./").
					Value(&formData.ProjectPath).
					Validate(func(str string) error {
						// Path is optional, accept empty or default values
						return nil
					}),

				// Author
				huh.NewInput().
					Title("👤 Автор проекта").
					Description("Автор проекта (автоматически определен из git config)").
					Value(&formData.Author).
					Placeholder(defaultAuthor).
					Validate(func(str string) error {
						// Author is optional, accept empty or default values
						return nil
					}),

				// Framework selection
				huh.NewSelect[string]().
					Title("🤖 Фреймворк агента").
					Description("Выберите фреймворк для разработки AI агента").
					Options(
						huh.NewOption("ADK (Agent Development Kit)", "adk"),
						huh.NewOption("LangGraph", "langgraph"),
						huh.NewOption("CrewAI", "crewai"),
					).
					Value(&formData.Framework),

				// CI/CD system
				huh.NewSelect[string]().
					Title("🔧 CI/CD система").
					Description("Выберите систему CI/CD для проекта").
					Options(
						huh.NewOption("GitLab CI", "gitlab"),
						huh.NewOption("GitHub Actions", "github"),
						huh.NewOption("Оба варианта", "both"),
						huh.NewOption("Без CI/CD", "none"),
					).
					Value(&formData.CICDType),

				// Git initialization
				huh.NewConfirm().
					Title("📦 Инициализировать Git репозиторий").
					Description("Создать git репозиторий и сделать первый коммит").
					Value(&formData.GitInit).
					Affirmative("Да").
					Negative("Нет"),

				// Create .env file
				huh.NewConfirm().
					Title("⚙️ Создать .env файл").
					Description("Скопировать .env.example в .env для локальной разработки").
					Value(&formData.CreateEnv).
					Affirmative("Да").
					Negative("Нет"),

				// Install dependencies
				huh.NewConfirm().
					Title("📚 Установить зависимости").
					Description("Запустить uv sync для установки зависимостей").
					Value(&formData.InstallDeps).
					Affirmative("Да").
					Negative("Нет"),
			),
		)
	} else {
		// MCP project form without framework selection
		form = huh.NewForm(
			huh.NewGroup(
				// Project name
				huh.NewInput().
					Title("🚀 Название проекта").
					Description("Введите название вашего проекта").
					Placeholder("my-awesome-mcp").
					Value(&formData.ProjectName).
					Validate(func(str string) error {
						if str == "" {
							return fmt.Errorf("название проекта обязательно")
						}
						return nil
					}),

				// Project path
				huh.NewInput().
					Title("📁 Путь для создания").
					Description("Директория для создания проекта (оставьте пустым для текущей)").
					Placeholder("./").
					Value(&formData.ProjectPath).
					Validate(func(str string) error {
						// Path is optional, accept empty or default values
						return nil
					}),

				// Author
				huh.NewInput().
					Title("👤 Автор проекта").
					Description("Автор проекта (автоматически определен из git config)").
					Value(&formData.Author).
					Placeholder(defaultAuthor).
					Validate(func(str string) error {
						// Author is optional, accept empty or default values
						return nil
					}),

				// CI/CD system
				huh.NewSelect[string]().
					Title("🔧 CI/CD система").
					Description("Выберите систему CI/CD для проекта").
					Options(
						huh.NewOption("GitLab CI", "gitlab"),
						huh.NewOption("GitHub Actions", "github"),
						huh.NewOption("Оба варианта", "both"),
						huh.NewOption("Без CI/CD", "none"),
					).
					Value(&formData.CICDType),

				// Git initialization
				huh.NewConfirm().
					Title("📦 Инициализировать Git репозиторий").
					Description("Создать git репозиторий и сделать первый коммит").
					Value(&formData.GitInit).
					Affirmative("Да").
					Negative("Нет"),

				// Create .env file
				huh.NewConfirm().
					Title("⚙️ Создать .env файл").
					Description("Скопировать .env.example в .env для локальной разработки").
					Value(&formData.CreateEnv).
					Affirmative("Да").
					Negative("Нет"),

				// Install dependencies
				huh.NewConfirm().
					Title("📚 Установить зависимости").
					Description("Запустить uv sync для установки зависимостей").
					Value(&formData.InstallDeps).
					Affirmative("Да").
					Negative("Нет"),
			),
		)
	}

	// Configure form with fullscreen mode
	form = form.
		WithTheme(huh.ThemeCharm()).
		WithAccessible(os.Getenv("ACCESSIBLE") != "").
		WithWidth(120).
		WithHeight(40)

	// Run the form
	if err := form.Run(); err != nil {
		return nil, fmt.Errorf("failed to run form: %w", err)
	}

	return &formData, nil
}

// getGitAuthorFromConfig retrieves author information from git config
func getGitAuthorFromConfig() string {
	// Get git user.name
	nameCmd := exec.Command("git", "config", "user.name")
	nameOutput, err := nameCmd.Output()
	if err != nil {
		return ""
	}

	name := strings.TrimSpace(string(nameOutput))
	if name == "" {
		return ""
	}

	// Get git user.email
	emailCmd := exec.Command("git", "config", "user.email")
	emailOutput, err := emailCmd.Output()
	if err != nil {
		return name
	}

	email := strings.TrimSpace(string(emailOutput))
	if email == "" {
		return name
	}

	return fmt.Sprintf("%s <%s>", name, email)
}
