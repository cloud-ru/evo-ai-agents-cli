package mcp_server

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
	"github.com/cloudru/ai-agents-cli/internal/api"
	"github.com/cloudru/ai-agents-cli/localizations"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var (
	deployFile string
	dryRun     bool
)

// deployCmd represents the deploy command
var deployCmd = &cobra.Command{
	Use:   "deploy [config-file]",
	Short: localizations.Localization.Get("deploy_short"),
	Long:  localizations.Localization.Get("deploy_long"),
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()

		// Определяем файл конфигурации
		configFile := deployFile
		if len(args) > 0 {
			configFile = args[0]
		}

		if configFile == "" {
			// Ищем файл конфигурации по умолчанию
			defaultFiles := []string{
				"mcp-servers.yaml",
				"mcp-servers.yml",
				"mcp-servers.json",
			}

			for _, file := range defaultFiles {
				if _, err := os.Stat(file); err == nil {
					configFile = file
					break
				}
			}

			if configFile == "" {
				log.Fatal("No configuration file specified and no default file found")
			}
		}

		// Читаем файл конфигурации
		data, err := ioutil.ReadFile(configFile)
		if err != nil {
			log.Fatal("Failed to read config file", "error", err, "file", configFile)
		}

		// Парсим конфигурацию
		var config struct {
			MCPServers []struct {
				Name        string                 `json:"name" yaml:"name"`
				Description string                 `json:"description" yaml:"description"`
				Options     map[string]interface{} `json:"options" yaml:"options"`
			} `json:"mcp-servers" yaml:"mcp-servers"`
		}

		ext := filepath.Ext(configFile)
		switch ext {
		case ".json":
			err = json.Unmarshal(data, &config)
		case ".yaml", ".yml":
			err = yaml.Unmarshal(data, &config)
		default:
			log.Fatal("Unsupported file format", "extension", ext)
		}

		if err != nil {
			log.Fatal("Failed to parse config file", "error", err)
		}

		if len(config.MCPServers) == 0 {
			log.Fatal("No MCP servers found in configuration")
		}

		// Создаем стили для вывода
		headerStyle := lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("205")).
			Border(lipgloss.RoundedBorder()).
			Padding(0, 1)

		successStyle := lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("2"))

		errorStyle := lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("1"))

		// Выводим заголовок
		fmt.Println(headerStyle.Render("🚀 Развертывание MCP серверов"))
		fmt.Printf("Файл конфигурации: %s\n", configFile)
		fmt.Printf("Количество серверов: %d\n", len(config.MCPServers))
		if dryRun {
			fmt.Println("🔍 Режим предварительного просмотра (dry-run)")
		}
		fmt.Println()

		// Развертываем каждый сервер
		successCount := 0
		errorCount := 0

		for i, serverConfig := range config.MCPServers {
			fmt.Printf("[%d/%d] ", i+1, len(config.MCPServers))

			if dryRun {
				fmt.Printf("🔍 %s (dry-run)\n", serverConfig.Name)
				continue
			}

			// Создаем MCP сервер
			req := &api.MCPServerCreateRequest{
				Name:        serverConfig.Name,
				Description: serverConfig.Description,
				Options:     serverConfig.Options,
			}

			server, err := apiClient.MCPServers.Create(ctx, req)
			if err != nil {
				fmt.Printf("%s %s: %v\n", errorStyle.Render("❌"), serverConfig.Name, err)
				errorCount++
				continue
			}

			fmt.Printf("%s %s (ID: %s)\n", successStyle.Render("✅"), serverConfig.Name, server.ID[:8]+"...")
			successCount++
		}

		// Выводим итоги
		fmt.Println()
		if dryRun {
			fmt.Println(headerStyle.Render("🔍 Предварительный просмотр завершен"))
		} else {
			fmt.Println(headerStyle.Render("🎉 Развертывание завершено"))
			fmt.Printf("Успешно: %d\n", successCount)
			if errorCount > 0 {
				fmt.Printf("Ошибок: %d\n", errorCount)
			}
		}
	},
}

func init() {
	RootCMD.AddCommand(deployCmd)

	deployCmd.Flags().StringVarP(&deployFile, "file", "f", "", "Путь к файлу конфигурации")
	deployCmd.Flags().BoolVarP(&dryRun, "dry-run", "d", false, "Режим предварительного просмотра без создания ресурсов")
}
