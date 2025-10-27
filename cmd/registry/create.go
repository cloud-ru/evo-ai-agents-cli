package registry

import (
	"context"
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
	"github.com/cloud-ru/evo-ai-agents-cli/internal/api"
	"github.com/cloud-ru/evo-ai-agents-cli/internal/di"
	"github.com/spf13/cobra"
)

var (
	createName     string
	createType     string
	createIsPublic bool
)

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Создать новый реестр образов контейнеров",
	Long:  "Создает новый реестр в Artifact Registry для хранения образов контейнеров",
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()

		// Проверяем обязательные параметры
		if createName == "" {
			log.Fatal("Name is required. Use --name flag")
		}

		// Получаем API клиент из DI контейнера
		container := di.GetContainer()
		apiClient, err := container.GetAPI()
		if err != nil {
			log.Fatal("Failed to get API client", "error", err)
		}

		// Определяем тип реестра
		var registryType api.RegistryType
		switch createType {
		case "debian":
			registryType = api.RegistryTypeDebian
		case "rpm":
			registryType = api.RegistryTypeRPM
		default:
			registryType = api.RegistryTypeDocker
		}

		// Создаем запрос
		req := &api.RegistryCreateRequest{
			Name:         createName,
			RegistryType: registryType,
			IsPublic:     createIsPublic,
		}

		// Создаем реестр (возвращает операцию)
		operation, err := apiClient.Registries.Create(ctx, req)
		if err != nil {
			log.Fatal("Failed to create registry", "error", err)
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
		fmt.Println(successStyle.Render("✅ Операция создания реестра запущена"))
		fmt.Println()
		fmt.Printf("%s: %s\n", labelStyle.Render("ID операции"), valueStyle.Render(operation.ID))

		if operation.ResourceID != "" {
			fmt.Printf("%s: %s\n", labelStyle.Render("ID реестра"), valueStyle.Render(operation.ResourceID))
		}

		if operation.ResourceName != "" {
			fmt.Printf("%s: %s\n", labelStyle.Render("Название"), valueStyle.Render(operation.ResourceName))
		}

		if operation.Description != "" {
			fmt.Printf("%s: %s\n", labelStyle.Render("Описание"), valueStyle.Render(operation.Description))
		}

		fmt.Printf("%s: %s\n", labelStyle.Render("Статус"), valueStyle.Render(func() string {
			if operation.Done {
				return "Завершено"
			}
			return "В процессе"
		}()))

		// Выводим инструкции по использованию
		fmt.Println()
		fmt.Println("💡 После завершения операции используйте:")
		fmt.Println("  ai-agents-cli registry list")
		fmt.Println("  ai-agents-cli registry get <registry-id>")
	},
}

func init() {
	createCmd.Flags().StringVarP(&createName, "name", "n", "", "Название реестра (обязательно)")
	createCmd.Flags().StringVarP(&createType, "type", "t", "docker", "Тип реестра (docker, debian, rpm)")
	createCmd.Flags().BoolVar(&createIsPublic, "public", false, "Сделать реестр публичным")
}
