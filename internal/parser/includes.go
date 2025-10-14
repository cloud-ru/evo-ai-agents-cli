package parser

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// IncludeError представляет ошибку при обработке includes
type IncludeError struct {
	File    string
	Message string
	Err     error
}

func (e *IncludeError) Error() string {
	return fmt.Sprintf("include error in %s: %s: %v", e.File, e.Message, e.Err)
}

// IncludeProcessor обрабатывает includes в YAML файлах
type IncludeProcessor struct {
	processedFiles map[string]bool
	baseDir        string
}

// NewIncludeProcessor создает новый процессор includes
func NewIncludeProcessor(baseDir string) *IncludeProcessor {
	return &IncludeProcessor{
		processedFiles: make(map[string]bool),
		baseDir:        baseDir,
	}
}

// ProcessIncludes обрабатывает includes в YAML файле
func (p *IncludeProcessor) ProcessIncludes(filePath string) (map[string]interface{}, error) {
	// Проверяем на циклические зависимости
	if p.processedFiles[filePath] {
		return nil, &IncludeError{
			File:    filePath,
			Message: "circular dependency detected",
			Err:     fmt.Errorf("file %s is already being processed", filePath),
		}
	}

	// Отмечаем файл как обрабатываемый
	p.processedFiles[filePath] = true
	defer func() {
		delete(p.processedFiles, filePath)
	}()

	// Читаем файл
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, &IncludeError{
			File:    filePath,
			Message: "failed to read file",
			Err:     err,
		}
	}

	// Парсим YAML
	var content map[string]interface{}
	if err := yaml.Unmarshal(data, &content); err != nil {
		return nil, &IncludeError{
			File:    filePath,
			Message: "failed to parse YAML",
			Err:     err,
		}
	}

	// Обрабатываем includes
	processedContent, err := p.processNode(content, filePath)
	if err != nil {
		return nil, err
	}

	// Приводим к нужному типу
	if result, ok := processedContent.(map[string]interface{}); ok {
		return result, nil
	}

	return nil, &IncludeError{
		File:    filePath,
		Message: "processed content is not a map",
		Err:     fmt.Errorf("expected map[string]interface{}, got %T", processedContent),
	}
}

// processNode рекурсивно обрабатывает узел YAML на предмет includes
func (p *IncludeProcessor) processNode(node interface{}, currentFile string) (interface{}, error) {
	switch v := node.(type) {
	case map[string]interface{}:
		result := make(map[string]interface{})
		for key, value := range v {
			// Проверяем на include директиву
			if key == "!include" {
				// Обрабатываем include
				includePath, ok := value.(string)
				if !ok {
					return nil, &IncludeError{
						File:    currentFile,
						Message: "include path must be a string",
						Err:     fmt.Errorf("expected string, got %T", value),
					}
				}

				// Разрешаем путь
				resolvedPath, err := p.resolvePath(includePath, currentFile)
				if err != nil {
					return nil, &IncludeError{
						File:    currentFile,
						Message: "failed to resolve include path",
						Err:     err,
					}
				}

				// Обрабатываем включаемый файл
				includedContent, err := p.ProcessIncludes(resolvedPath)
				if err != nil {
					return nil, err
				}

				// Возвращаем содержимое включаемого файла
				return includedContent, nil
			}

			// Обрабатываем значение рекурсивно
			processedValue, err := p.processNode(value, currentFile)
			if err != nil {
				return nil, err
			}
			result[key] = processedValue
		}
		return result, nil

	case []interface{}:
		result := make([]interface{}, len(v))
		for i, item := range v {
			processedItem, err := p.processNode(item, currentFile)
			if err != nil {
				return nil, err
			}
			result[i] = processedItem
		}
		return result, nil

	default:
		// Примитивные типы возвращаем как есть
		return node, nil
	}
}

// resolvePath разрешает путь include относительно текущего файла
func (p *IncludeProcessor) resolvePath(includePath, currentFile string) (string, error) {
	// Если путь абсолютный, используем его как есть
	if filepath.IsAbs(includePath) {
		return includePath, nil
	}

	// Получаем директорию текущего файла
	currentDir := filepath.Dir(currentFile)

	// Если baseDir не установлен, используем директорию текущего файла
	baseDir := p.baseDir
	if baseDir == "" {
		baseDir = currentDir
	}

	// Разрешаем путь относительно baseDir
	resolvedPath := filepath.Join(baseDir, includePath)

	// Проверяем, что файл существует
	if _, err := os.Stat(resolvedPath); os.IsNotExist(err) {
		return "", fmt.Errorf("included file does not exist: %s", resolvedPath)
	}

	return resolvedPath, nil
}

// ProcessYAMLFile обрабатывает YAML файл с includes
func ProcessYAMLFile(filePath string) (map[string]interface{}, error) {
	// Получаем абсолютный путь к файлу
	absPath, err := filepath.Abs(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to get absolute path: %w", err)
	}

	// Получаем базовую директорию
	baseDir := filepath.Dir(absPath)

	// Создаем процессор includes
	processor := NewIncludeProcessor(baseDir)

	// Обрабатываем файл
	return processor.ProcessIncludes(absPath)
}

// ValidateIncludes проверяет, что все includes в файле корректны
func ValidateIncludes(filePath string) error {
	_, err := ProcessYAMLFile(filePath)
	return err
}

// GetIncludeDependencies возвращает список всех файлов, от которых зависит данный файл
func GetIncludeDependencies(filePath string) ([]string, error) {
	absPath, err := filepath.Abs(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to get absolute path: %w", err)
	}

	baseDir := filepath.Dir(absPath)
	processor := NewIncludeProcessor(baseDir)

	var dependencies []string

	// Создаем функцию для сбора зависимостей
	var collectDeps func(node interface{}, currentFile string) error
	collectDeps = func(node interface{}, currentFile string) error {
		switch v := node.(type) {
		case map[string]interface{}:
			for key, value := range v {
				if key == "!include" {
					if includePath, ok := value.(string); ok {
						resolvedPath, err := processor.resolvePath(includePath, currentFile)
						if err != nil {
							return err
						}
						dependencies = append(dependencies, resolvedPath)

						// Рекурсивно собираем зависимости включаемого файла
						data, err := ioutil.ReadFile(resolvedPath)
						if err != nil {
							return err
						}

						var content map[string]interface{}
						if err := yaml.Unmarshal(data, &content); err != nil {
							return err
						}

						return collectDeps(content, resolvedPath)
					}
				} else {
					if err := collectDeps(value, currentFile); err != nil {
						return err
					}
				}
			}
		case []interface{}:
			for _, item := range v {
				if err := collectDeps(item, currentFile); err != nil {
					return err
				}
			}
		}
		return nil
	}

	// Читаем и парсим файл
	data, err := ioutil.ReadFile(absPath)
	if err != nil {
		return nil, err
	}

	var content map[string]interface{}
	if err := yaml.Unmarshal(data, &content); err != nil {
		return nil, err
	}

	// Собираем зависимости
	if err := collectDeps(content, absPath); err != nil {
		return nil, err
	}

	return dependencies, nil
}
