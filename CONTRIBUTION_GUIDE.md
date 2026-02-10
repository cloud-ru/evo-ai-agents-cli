# Руководство по внесению вклада в ai-agents-cli

Добро пожаловать в проект ai-agents-cli! Это руководство поможет вам понять архитектуру проекта, паттерны разработки и процессы внесения изменений.

## 📋 Содержание

- [Архитектура проекта](#-архитектура-проекта)
- [Паттерны разработки](#-паттерны-разработки)
- [Система обработки ошибок](#-система-обработки-ошибок)
- [Создание команд](#-создание-команд)
- [TUI компоненты](#-tui-компоненты)
- [Шаблоны проектов](#-шаблоны-проектов)
- [Тестирование](#-тестирование)
- [Git workflow](#-git-workflow)
- [Создание MR](#-создание-mr)
- [Стиль кода](#-стиль-кода)

## 🏗 Архитектура проекта

### Структура директорий

```
ai-agents-cli/
├── cmd/                    # Cobra команды
│   ├── agent/             # Команды для работы с агентами
│   ├── create/            # Команды создания проектов
│   ├── mcp_server/        # Команды для MCP серверов
│   ├── system/            # Команды для систем
│   └── ci/                # CI/CD команды
├── internal/              # Внутренние пакеты
│   ├── api/               # API клиент
│   ├── auth/              # Аутентификация
│   ├── config/            # Конфигурация
│   ├── deployer/          # Развертывание
│   ├── di/                # Dependency Injection
│   ├── errors/            # Система обработки ошибок
│   ├── scaffolder/        # Создание проектов
│   ├── ui/                # TUI компоненты
│   └── validator/         # Валидация
├── templates/             # Шаблоны проектов
├── localizations/        # Локализация
├── schemas/              # JSON схемы
└── scripts/              # Скрипты
```

### Основные компоненты

1. **Cobra CLI** - фреймворк для командной строки
2. **Dependency Injection** - управление зависимостями через `samber/do`
3. **TUI** - терминальный интерфейс с `bubbletea` и `huh`
4. **Error Handling** - структурированная система ошибок
5. **Templates** - система шаблонов с `embed.FS`

## 🎯 Паттерны разработки

### 1. Dependency Injection

Используем `samber/do` для управления зависимостями:

```go
// internal/di/container.go
func (c *Container) GetAPI() (*api.Client, error) {
    client, err := do.Invoke[*api.Client](c.container)
    if err != nil {
        return nil, oops.Wrap(err, "failed to get API client")
    }
    return client, nil
}
```

**Правила:**
- ❌ Никогда не используйте `MustInvoke`
- ✅ Всегда используйте `Invoke` с обработкой ошибок
- ✅ Оборачивайте ошибки через `oops.Wrap`

### 2. Структурированные ошибки

Все ошибки должны быть структурированными:

```go
// Создание ошибки
err := errors.ValidationError("INVALID_EMAIL", "Некорректный формат email").
    WithDetails("Email должен содержать символ @").
    WithContext("input", email)

// Обработка ошибки
errorHandler := errors.NewHandler()
appErr := errorHandler.WrapValidationError(err, "EMAIL_VALIDATION", "Ошибка валидации email")
fmt.Println(errorHandler.HandlePlain(appErr))
```

**Типы ошибок:**
- `ErrorTypeValidation` - ошибки валидации
- `ErrorTypeConfiguration` - ошибки конфигурации
- `ErrorTypeAuthentication` - ошибки аутентификации
- `ErrorTypeAPI` - ошибки API
- `ErrorTypeNetwork` - сетевые ошибки
- `ErrorTypeFileSystem` - ошибки файловой системы
- `ErrorTypeTemplate` - ошибки шаблонов
- `ErrorTypeUser` - пользовательские ошибки
- `ErrorTypeSystem` - системные ошибки

### 3. Логирование

Используем структурированное логирование:

```go
// Получение логгера
logger := errorHandler.GetLogger()

// Логирование с контекстом
logger.WithContext(map[string]interface{}{
    "user_id": "123",
    "operation": "create_project",
}).LogError(err, "Failed to create project")

// Логирование структурированной ошибки
logger.LogAppError(appErr, "Structured error occurred")
```

## 🚨 Система обработки ошибок

### Создание ошибок

```go
// Простая ошибка
err := errors.ValidationError("MISSING_FIELD", "Обязательное поле не заполнено")

// Ошибка с контекстом
err = err.WithContext("field", "project_name").
    WithDetails("Поле project_name обязательно для заполнения")

// Оборачивание существующей ошибки
wrappedErr := errorHandler.WrapFileSystemError(originalErr, "FILE_NOT_FOUND", "Файл не найден")
```

### Обработка ошибок

```go
// Создание обработчика
errorHandler := errors.NewHandler()

// Простое отображение (рекомендуется для CLI)
message := errorHandler.HandlePlain(err)

// Стилизованное отображение
message := errorHandler.HandleSimple(err)

// Полное отображение с рамками
message := errorHandler.Handle(err)
```

### Отображение в UI

```go
// Простой текст
fmt.Println(errors.FormatPlainError(err))

// Стилизованный текст
fmt.Println(errors.FormatSimpleError(err))

// Полное отображение
fmt.Println(errors.FormatError(err))
```

## 🛠 Создание команд

### Структура команды

```go
// cmd/example/root.go
package example

import (
    "fmt"
    
    "github.com/cloud-ru/evo-ai-agents-cli/internal/errors"
    "github.com/spf13/cobra"
)

var exampleCmd = &cobra.Command{
    Use:   "example",
    Short: "Пример команды",
    Long:  "Подробное описание команды",
    Run: func(cmd *cobra.Command, args []string) {
        // Создаем обработчик ошибок
        errorHandler := errors.NewHandler()
        
        // Выполняем операцию
        if err := doSomething(); err != nil {
            appErr := errorHandler.WrapUserError(err, "OPERATION_FAILED", "Ошибка выполнения операции")
            fmt.Println(errorHandler.HandlePlain(appErr))
            return
        }
        
        fmt.Println("Операция выполнена успешно")
    },
}

func init() {
    // Добавляем флаги
    exampleCmd.Flags().StringP("param", "p", "", "Параметр команды")
}

// Регистрируем команду в cmd/root.go
func init() {
    rootCmd.AddCommand(exampleCmd)
}
```

### Правила создания команд

1. **Всегда создавайте обработчик ошибок** в начале функции `Run`
2. **Используйте структурированные ошибки** вместо `log.Fatal`
3. **Обрабатывайте ошибки** через `errorHandler.HandlePlain()`
4. **Добавляйте валидацию** входных параметров
5. **Используйте контекст** для отладки

## 🎨 TUI компоненты

### Создание форм с huh

```go
// internal/ui/example_form.go
package ui

import (
    "github.com/charmbracelet/huh"
    "github.com/cloud-ru/evo-ai-agents-cli/internal/errors"
)

type ExampleFormData struct {
    Name        string
    Description string
    Options     []string
}

func RunExampleForm() (*ExampleFormData, error) {
    var formData ExampleFormData
    
    form := huh.NewForm(
        huh.NewGroup(
            huh.NewInput().
                Title("📝 Название").
                Description("Введите название").
                Value(&formData.Name).
                Validate(func(str string) error {
                    if str == "" {
                        return errors.ValidationError("MISSING_NAME", "Название обязательно")
                    }
                    return nil
                }),
                
            huh.NewText().
                Title("📄 Описание").
                Description("Введите описание").
                Value(&formData.Description),
                
            huh.NewMultiSelect[string]().
                Title("⚙️ Опции").
                Description("Выберите опции").
                Options(huh.NewOptions("option1", "option2", "option3")...).
                Value(&formData.Options),
        ),
    ).WithTheme(huh.ThemeCharm()).
        WithWidth(120).
        WithHeight(40)
    
    if err := form.Run(); err != nil {
        return nil, errors.Wrap(err, errors.ErrorTypeUser, errors.SeverityMedium, "FORM_ERROR", "Ошибка заполнения формы")
    }
    
    return &formData, nil
}
```

### Создание таблиц

```go
// internal/ui/example_table.go
func ShowExampleTable(data []ExampleItem) {
    table := table.New().
        Header("ID", "Name", "Status").
        Rows(func() []table.Row {
            var rows []table.Row
            for _, item := range data {
                rows = append(rows, table.Row{
                    item.ID,
                    item.Name,
                    item.Status,
                })
            }
            return rows
        }()).
        Border(table.RoundedBorder()).
        BorderStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("62")))
    
    fmt.Println(table)
}
```

## 📄 Шаблоны проектов

### Структура шаблонов

```
templates/
├── mcp/                   # Шаблоны для MCP проектов
│   ├── README.md.tmpl
│   ├── Makefile.tmpl
│   ├── Dockerfile.tmpl
│   └── src/
│       └── main.py.tmpl
├── agent/                 # Базовые шаблоны агентов
└── agent-frameworks/      # Шаблоны для фреймворков
    ├── adk/
    ├── langgraph/
    └── crewai/
```

### Создание шаблона

```go
// internal/scaffolder/templates.go
//go:embed templates/*
var TemplatesFS embed.FS

// internal/scaffolder/scaffolder.go
func (s *Scaffolder) processTemplate(templateContent string, data *ProjectData) (string, error) {
    tmpl, err := template.New("template").Parse(templateContent)
    if err != nil {
        return "", errors.Wrap(err, errors.ErrorTypeTemplate, errors.SeverityMedium, "TEMPLATE_PARSE_FAILED", "Ошибка парсинга шаблона")
    }
    
    var buf bytes.Buffer
    if err := tmpl.Execute(&buf, data); err != nil {
        return "", errors.Wrap(err, errors.ErrorTypeTemplate, errors.SeverityMedium, "TEMPLATE_EXECUTE_FAILED", "Ошибка выполнения шаблона")
    }
    
    return buf.String(), nil
}
```

### Переменные шаблонов

```go
type ProjectData struct {
    ProjectName string    // Название проекта
    ProjectType string    // Тип проекта (mcp, agent)
    Framework   string    // Фреймворк (adk, langgraph, crewai)
    Author      string    // Автор
    Year        string    // Год
    CICDType    string    // Тип CI/CD
    Description string    // Описание
}
```

## 🧪 Тестирование

### Unit тесты

```go
// internal/errors/errors_test.go
func TestAppError(t *testing.T) {
    err := ValidationError("TEST_ERROR", "Тестовая ошибка").
        WithDetails("Детали ошибки").
        WithContext("key", "value")
    
    assert.Equal(t, ErrorTypeValidation, err.Type)
    assert.Equal(t, SeverityMedium, err.Severity)
    assert.Equal(t, "TEST_ERROR", err.Code)
    assert.Equal(t, "Тестовая ошибка", err.Message)
    assert.Equal(t, "Детали ошибки", err.Details)
    assert.Equal(t, "value", err.Context["key"])
}
```

### Integration тесты

```go
// cmd/integration_test.go
func TestCreateProject(t *testing.T) {
    // Создаем временную директорию
    tempDir := t.TempDir()
    
    // Тестируем создание проекта
    scaffolder := scaffolder.NewScaffolder()
    err := scaffolder.CreateProject("mcp", "test-project", tempDir, "both")
    
    assert.NoError(t, err)
    
    // Проверяем созданные файлы
    assert.FileExists(t, filepath.Join(tempDir, "README.md"))
    assert.FileExists(t, filepath.Join(tempDir, "Makefile"))
}
```

## 🌿 Git Workflow

### Создание веток

```bash
# Создание feature ветки
git checkout -b feature/error-handling-improvements

# Создание bugfix ветки
git checkout -b bugfix/fix-template-validation

# Создание hotfix ветки
git checkout -b hotfix/critical-security-fix
```

### Именование веток

- `feature/` - новые функции
- `bugfix/` - исправления багов
- `hotfix/` - критические исправления
- `refactor/` - рефакторинг
- `docs/` - документация
- `test/` - тесты

### Коммиты

```bash
# Структура коммита
<type>(<scope>): <description>

# Примеры
feat(create): add agent framework selection
fix(ui): resolve nested error boxes issue
refactor(errors): improve error handling system
docs(readme): update installation instructions
test(scaffolder): add unit tests for template processing
```

**Типы коммитов:**
- `feat` - новая функция
- `fix` - исправление бага
- `refactor` - рефакторинг
- `docs` - документация
- `test` - тесты
- `style` - форматирование
- `perf` - производительность
- `ci` - CI/CD

## 📝 Создание MR

### Подготовка MR

1. **Создайте ветку** от `main`
2. **Внесите изменения** согласно паттернам
3. **Добавьте тесты** для новой функциональности
4. **Обновите документацию** при необходимости
5. **Проверьте линтер** и исправьте ошибки

### Описание MR

```markdown
## 📋 Описание

Краткое описание изменений

## 🎯 Тип изменений

- [ ] Новая функция
- [ ] Исправление бага
- [ ] Рефакторинг
- [ ] Документация
- [ ] Тесты

## 🔧 Изменения

- Добавлена система обработки ошибок
- Улучшено отображение TUI форм
- Добавлены unit тесты

## 🧪 Тестирование

- [ ] Unit тесты пройдены
- [ ] Integration тесты пройдены
- [ ] Ручное тестирование выполнено

## 📸 Скриншоты

(если применимо)

## ✅ Checklist

- [ ] Код соответствует стилю проекта
- [ ] Добавлены тесты для новой функциональности
- [ ] Документация обновлена
- [ ] Линтер не показывает ошибок
```

### Review процесс

1. **Автоматические проверки** (CI/CD)
2. **Code review** от коллег
3. **Тестирование** функциональности
4. **Merge** после одобрения

## 🎨 Стиль кода

### Go

```go
// Именование
const MaxRetries = 3
var errorHandler *Handler
type ProjectData struct {}

// Функции
func (s *Scaffolder) CreateProject(projectType, projectName string) error {
    // Реализация
}

// Обработка ошибок
if err != nil {
    return errors.Wrap(err, errors.ErrorTypeSystem, errors.SeverityHigh, "OPERATION_FAILED", "Операция не выполнена")
}
```

### Документация

```go
// Package scaffolder provides project scaffolding functionality.
package scaffolder

// Scaffolder represents the project scaffolding functionality.
// It handles template processing, directory creation, and file generation.
type Scaffolder struct {
    templates embed.FS
    config    *ScaffolderConfig
}

// CreateProject creates a new project from templates.
// It validates inputs, processes templates, and generates project files.
func (s *Scaffolder) CreateProject(projectType, projectName, targetPath, cicdType string) error {
    // Реализация
}
```

## 🚀 Быстрый старт

### 1. Клонирование и настройка

```bash
git clone <repository-url>
cd ai-agents-cli
go mod download
```

### 2. Создание feature ветки

```bash
git checkout -b feature/my-feature
```

### 3. Разработка

```bash
# Создание команды
mkdir cmd/my-command
# Создание тестов
go test ./...
# Проверка линтера
golangci-lint run
```

### 4. Коммит и push

```bash
git add .
git commit -m "feat(my-command): add new command"
git push origin feature/my-feature
```

### 5. Создание MR

Создайте MR через GitLab интерфейс с подробным описанием.

## 📚 Дополнительные ресурсы

- [📚 Cloud.ru AI Agents Documentation](https://cloud.ru/docs/ai-agents/ug/index?source-platform=Evolution) - Официальная документация по AI Agents
- [Cobra CLI Documentation](https://cobra.dev/)
- [Bubble Tea Documentation](https://github.com/charmbracelet/bubbletea)
- [Huh Documentation](https://github.com/charmbracelet/huh)
- [Go Embed Documentation](https://pkg.go.dev/embed)
- [Samber/Do Documentation](https://github.com/samber/do)

## 🤝 Поддержка

Если у вас есть вопросы:

1. Проверьте документацию
2. Посмотрите примеры в коде
3. Создайте issue в репозитории
4. Обратитесь к команде разработки

---

**Спасибо за вклад в проект! 🎉**
