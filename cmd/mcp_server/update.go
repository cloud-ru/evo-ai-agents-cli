package mcp_server

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
	"github.com/cloudru/ai-agents-cli/internal/api"
	"github.com/spf13/cobra"
)

var (
	updateName        string
	updateDescription string
	updateConfigFile  string
)

// updateCmd represents the update command
var updateCmd = &cobra.Command{
	Use:   "update <server-id>",
	Short: "Обновить MCP сервер",
	Long:  "Обновляет существующий MCP сервер с новыми параметрами",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		serverID := args[0]

		var req *api.MCPServerUpdateRequest

		if updateConfigFile != "" {
			// Загружаем конфигурацию из файла
			data, err := ioutil.ReadFile(updateConfigFile)
			if err != nil {
				log.Fatal("Failed to read config file", "error", err, "file", updateConfigFile)
			}

			var config struct {
				Name        string                 `json:"name"`
				Description string                 `json:"description"`
				Options     map[string]interface{} `json:"options"`
			}

			if err := json.Unmarshal(data, &config); err != nil {
				log.Fatal("Failed to parse config file", "error", err)
			}

			req = &api.MCPServerUpdateRequest{
				Name:        config.Name,
				Description: config.Description,
				Options:     config.Options,
			}
		} else {
			// Используем параметры командной строки
			req = &api.MCPServerUpdateRequest{}

			if updateName != "" {
				req.Name = updateName
			}
			if updateDescription != "" {
				req.Description = updateDescription
			}
		}

		// Обновляем MCP сервер
		server, err := apiClient.MCPServers.Update(ctx, serverID, req)
		if err != nil {
			log.Fatal("Failed to update MCP server", "error", err, "server_id", serverID)
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
		fmt.Println(successStyle.Render("✅ MCP сервер успешно обновлен"))
		fmt.Println()
		fmt.Printf("%s: %s\n", labelStyle.Render("ID"), valueStyle.Render(server.ID))
		fmt.Printf("%s: %s\n", labelStyle.Render("Название"), valueStyle.Render(server.Name))

		if server.Description != "" {
			fmt.Printf("%s: %s\n", labelStyle.Render("Описание"), valueStyle.Render(server.Description))
		}

		fmt.Printf("%s: %s\n", labelStyle.Render("Статус"), valueStyle.Render(server.Status))
		fmt.Printf("%s: %s\n", labelStyle.Render("Обновлен"), valueStyle.Render(server.UpdatedAt.Format("02.01.2006 15:04:05")))
	},
}

func init() {
	RootCMD.AddCommand(updateCmd)

	updateCmd.Flags().StringVarP(&updateName, "name", "n", "", "Новое название MCP сервера")
	updateCmd.Flags().StringVarP(&updateDescription, "description", "d", "", "Новое описание MCP сервера")
	updateCmd.Flags().StringVarP(&updateConfigFile, "config", "c", "", "Путь к файлу конфигурации (JSON)")
}
