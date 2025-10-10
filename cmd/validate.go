package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
	"github.com/cloudru/ai-agents-cli/internal/validator"
	"github.com/spf13/cobra"
)

var (
	validateFile string
	validateDir  string
)

// validateCmd represents the validate command
var validateCmd = &cobra.Command{
	Use:   "validate [config-file]",
	Short: "Валидация конфигурационных файлов",
	Long:  "Проверяет корректность конфигурационных файлов перед развертыванием",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// Создаем валидатор
		configValidator := validator.NewConfigValidator()

		// Загружаем схемы
		if err := loadSchemas(configValidator); err != nil {
			log.Warn("Failed to load schemas", "error", err)
		}

		// Определяем файлы для валидации
		var files []string

		if len(args) > 0 {
			files = []string{args[0]}
		} else if validateFile != "" {
			files = []string{validateFile}
		} else if validateDir != "" {
			// Находим все конфигурационные файлы в директории
			var err error
			files, err = findConfigFiles(validateDir)
			if err != nil {
				log.Fatal("Failed to find config files", "error", err, "dir", validateDir)
			}
		} else {
			// Ищем файлы в текущей директории
			var err error
			files, err = findConfigFiles(".")
			if err != nil {
				log.Fatal("Failed to find config files", "error", err)
			}
		}

		if len(files) == 0 {
			log.Fatal("No configuration files found")
		}

		// Валидируем каждый файл
		allValid := true
		headerStyle := lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("205")).
			Border(lipgloss.RoundedBorder()).
			Padding(0, 1)

		fmt.Println(headerStyle.Render("🔍 Валидация конфигурационных файлов"))
		fmt.Println()

		for _, file := range files {
			fmt.Printf("Проверка %s...\n", file)

			result, err := configValidator.ValidateFile(file)
			if err != nil {
				log.Errorf("Ошибка при валидации %s: %v", file, err)
				allValid = false
				continue
			}

			if !result.Valid {
				allValid = false
				configValidator.PrintErrors(result)
			} else {
				log.Infof("✅ %s валиден", file)
			}
			fmt.Println()
		}

		// Выводим итоговый результат
		if allValid {
			fmt.Println(headerStyle.Copy().Foreground(lipgloss.Color("2")).Render("🎉 Все файлы валидны!"))
			os.Exit(0)
		} else {
			fmt.Println(headerStyle.Copy().Foreground(lipgloss.Color("1")).Render("❌ Обнаружены ошибки валидации"))
			os.Exit(1)
		}
	},
}

func loadSchemas(validator *validator.ConfigValidator) error {
	schemas := map[string]string{
		"mcp-servers":   "schemas/mcp.schema.json",
		"agents":        "schemas/agent.schema.json",
		"agent-systems": "schemas/systems.schema.json",
	}

	for name, path := range schemas {
		if err := validator.LoadSchema(name, path); err != nil {
			log.Warn("Failed to load schema", "schema", name, "path", path, "error", err)
		}
	}

	return nil
}

func findConfigFiles(dir string) ([]string, error) {
	var files []string

	// Ищем файлы с расширениями .yaml, .yml, .json
	extensions := []string{".yaml", ".yml", ".json"}

	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		name := entry.Name()
		for _, ext := range extensions {
			if len(name) > len(ext) && name[len(name)-len(ext):] == ext {
				// Добавляем полный путь к файлу
				files = append(files, filepath.Join(dir, name))
				break
			}
		}
	}

	return files, nil
}

func init() {
	RootCMD.AddCommand(validateCmd)

	validateCmd.Flags().StringVarP(&validateFile, "file", "f", "", "Файл для валидации")
	validateCmd.Flags().StringVarP(&validateDir, "dir", "d", "", "Директория с файлами для валидации")
}
