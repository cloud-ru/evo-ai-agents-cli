package mcp_server

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
	"github.com/cloud-ru/evo-ai-agents-cli/internal/api"
	"github.com/cloud-ru/evo-ai-agents-cli/internal/di"
	"github.com/spf13/cobra"
)

var (
	name        string
	description string
	configFile  string
)

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Создать новый MCP сервер",
	Long:  "Создает новый MCP сервер с указанными параметрами",
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()

		var req *api.MCPServerCreateRequest

		if configFile != "" {
			// Загружаем конфигурацию из файла
			data, err := ioutil.ReadFile(configFile)
			if err != nil {
				log.Fatal("Failed to read config file", "error", err, "file", configFile)
			}

			var config struct {
				Name        string                 `json:"name"`
				Description string                 `json:"description"`
				Options     map[string]interface{} `json:"options"`
			}

			if err := json.Unmarshal(data, &config); err != nil {
				log.Fatal("Failed to parse config file", "error", err)
			}

			req = &api.MCPServerCreateRequest{
				Name:        config.Name,
				Description: config.Description,
				Options:     config.Options,
			}
		} else {
			// Используем параметры командной строки
			if name == "" {
				log.Fatal("Name is required. Use --name flag or --config file")
			}

			req = &api.MCPServerCreateRequest{
				Name:        name,
				Description: description,
				Options:     make(map[string]interface{}),
			}
		}

		// Получаем API клиент из DI контейнера
		container := di.GetContainer()
		apiClient, err := container.GetAPI()
	if err != nil {
		log.Fatal("Failed to get API client", "error", err)
	}

		// Создаем MCP сервер
		server, err := apiClient.MCPServers.Create(ctx, req)
		if err != nil {
			log.Fatal("Failed to create MCP server", "error", err)
		}

		// Создаем стили для вывода
		successStyle := lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("2")).
			Border(lipgloss.RoundedBorder()).
			Padding(0, 1)

		labelStyle := lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("99"))

		valueStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("252"))

		// Выводим результат
		fmt.Println(successStyle.Render("✅ MCP сервер успешно создан"))
		fmt.Println()
		fmt.Printf("%s: %s\n", labelStyle.Render("ID"), valueStyle.Render(server.ID))
		fmt.Printf("%s: %s\n", labelStyle.Render("Название"), valueStyle.Render(server.Name))

		if server.Description != "" {
			fmt.Printf("%s: %s\n", labelStyle.Render("Описание"), valueStyle.Render(server.Description))
		}

		fmt.Printf("%s: %s\n", labelStyle.Render("Статус"), valueStyle.Render(server.Status))
		fmt.Printf("%s: %s\n", labelStyle.Render("Создан"), valueStyle.Render(server.CreatedAt.Time.Format("02.01.2006 15:04:05")))
	},
}

func init() {
	RootCMD.AddCommand(createCmd)

	createCmd.Flags().StringVarP(&name, "name", "n", "", "Название MCP сервера")
	createCmd.Flags().StringVarP(&description, "description", "d", "", "Описание MCP сервера")
	createCmd.Flags().StringVarP(&configFile, "config", "c", "", "Путь к файлу конфигурации (JSON)")
}
