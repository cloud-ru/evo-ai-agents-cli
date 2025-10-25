package system

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/charmbracelet/log"
	"github.com/cloud-ru/evo-ai-agents-cli/internal/api"
	"github.com/cloud-ru/evo-ai-agents-cli/internal/di"
	"github.com/spf13/cobra"
)

var (
	systemUpdateName        string
	systemUpdateDescription string
	systemUpdateAgents      []string
	systemUpdateOptions     string
)

// updateCmd represents the update command
var updateCmd = &cobra.Command{
	Use:   "update <system-id>",
	Short: "Обновление существующей системы агентов",
	Long:  "Обновляет параметры существующей системы агентов",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		systemID := args[0]

		// Получаем API клиент из DI контейнера
		container := di.GetContainer()
		apiClient, err := container.GetAPI()
	if err != nil {
		log.Fatal("Failed to get API client", "error", err)
	}

		// Парсим опции из JSON
		var options map[string]interface{}
		if systemUpdateOptions != "" {
			if err := json.Unmarshal([]byte(systemUpdateOptions), &options); err != nil {
				log.Fatal("Failed to parse options JSON", "error", err)
			}
		}

		// Создаем запрос
		req := &api.AgentSystemUpdateRequest{
			Name:        systemUpdateName,
			Description: systemUpdateDescription,
			Agents:      systemUpdateAgents,
			Options:     options,
		}

		// Обновляем систему
		system, err := apiClient.AgentSystems.Update(ctx, systemID, req)
		if err != nil {
			log.Fatal("Failed to update system", "error", err, "system_id", systemID)
		}

		fmt.Printf("✅ Система агентов обновлена успешно!\n")
		fmt.Printf("ID: %s\n", system.ID)
		fmt.Printf("Название: %s\n", system.Name)
		if system.Description != "" {
			fmt.Printf("Описание: %s\n", system.Description)
		}
		fmt.Printf("Статус: %s\n", system.Status)
		fmt.Printf("Агентов: %d\n", len(system.Agents))
	},
}

func init() {
	RootCMD.AddCommand(updateCmd)

	updateCmd.Flags().StringVarP(&systemUpdateName, "name", "n", "", "Новое название системы")
	updateCmd.Flags().StringVarP(&systemUpdateDescription, "description", "d", "", "Новое описание системы")
	updateCmd.Flags().StringSliceVarP(&systemUpdateAgents, "agents", "a", []string{}, "Новый список ID агентов для системы")
	updateCmd.Flags().StringVarP(&systemUpdateOptions, "options", "o", "", "Новые опции системы в формате JSON")
}
