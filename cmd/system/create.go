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
	systemCreateName        string
	systemCreateDescription string
	systemCreateAgents      []string
	systemCreateOptions     string
)

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Создание новой системы агентов",
	Long:  "Создает новую систему агентов с указанными параметрами",
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()

		// Получаем API клиент из DI контейнера
		container := di.GetContainer()
		apiClient := container.GetAPI()

		// Парсим опции из JSON
		var options map[string]interface{}
		if systemCreateOptions != "" {
			if err := json.Unmarshal([]byte(systemCreateOptions), &options); err != nil {
				log.Fatal("Failed to parse options JSON", "error", err)
			}
		}

		// Создаем запрос
		req := &api.AgentSystemCreateRequest{
			Name:        systemCreateName,
			Description: systemCreateDescription,
			Agents:      systemCreateAgents,
			Options:     options,
		}

		// Создаем систему
		system, err := apiClient.AgentSystems.Create(ctx, req)
		if err != nil {
			log.Fatal("Failed to create system", "error", err)
		}

		fmt.Printf("✅ Система агентов создана успешно!\n")
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
	RootCMD.AddCommand(createCmd)

	createCmd.Flags().StringVarP(&systemCreateName, "name", "n", "", "Название системы (обязательно)")
	createCmd.Flags().StringVarP(&systemCreateDescription, "description", "d", "", "Описание системы")
	createCmd.Flags().StringSliceVarP(&systemCreateAgents, "agents", "a", []string{}, "Список ID агентов для добавления в систему")
	createCmd.Flags().StringVarP(&systemCreateOptions, "options", "o", "", "Опции системы в формате JSON")

	createCmd.MarkFlagRequired("name")
}
