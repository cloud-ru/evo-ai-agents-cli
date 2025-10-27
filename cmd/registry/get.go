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
	getIdentifier string
)

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:   "get [registry-name-or-id]",
	Short: "Получить информацию о реестре",
	Long:  "Выводит подробную информацию о реестре по имени или ID",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()

		// Определяем идентификатор реестра
		identifier := getIdentifier
		if len(args) > 0 {
			identifier = args[0]
		}

		if identifier == "" {
			log.Fatal("Registry name or ID is required")
		}

		// Получаем API клиент из DI контейнера
		container := di.GetContainer()
		apiClient, err := container.GetAPI()
		if err != nil {
			log.Fatal("Failed to get API client", "error", err)
		}

		// Получаем информацию о реестре
		registry, err := apiClient.Registries.Get(ctx, identifier)
		if err != nil {
			log.Fatal("Failed to get registry", "error", err)
		}

		// Создаем стили для вывода
		headerStyle := lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("205")).
			Border(lipgloss.RoundedBorder()).
			Padding(0, 1)

		labelStyle := lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("99"))

		valueStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("252"))

		// Выводим информацию
		fmt.Println(headerStyle.Render("📦 Информация о реестре"))
		fmt.Println()
		fmt.Printf("%s: %s\n", labelStyle.Render("ID"), valueStyle.Render(registry.ID))
		fmt.Printf("%s: %s\n", labelStyle.Render("Название"), valueStyle.Render(registry.Name))
		fmt.Printf("%s: %s\n", labelStyle.Render("Тип"), valueStyle.Render(string(registry.RegistryType)))

		// Статус с иконкой
		statusText := string(registry.Status)
		switch registry.Status {
		case api.RegistryStatusCreating:
			statusText = "🟡 " + statusText
		case api.RegistryStatusActive:
			statusText = "🟢 " + statusText
		case api.RegistryStatusError:
			statusText = "🔴 " + statusText
		}
		fmt.Printf("%s: %s\n", labelStyle.Render("Статус"), valueStyle.Render(statusText))

		fmt.Printf("%s: %s\n", labelStyle.Render("Публичный"), valueStyle.Render(func() string {
			if registry.IsPublic {
				return "Да"
			}
			return "Нет"
		}()))

		fmt.Printf("%s: %s\n", labelStyle.Render("Уровень карантина"), valueStyle.Render(string(registry.QuarantineMode)))
		fmt.Printf("%s: %s\n", labelStyle.Render("Создан"), valueStyle.Render(registry.CreatedAt))

		if registry.RetentionPolicyIsEnabled {
			fmt.Printf("%s: %s\n", labelStyle.Render("Политика удаления"), valueStyle.Render("Включена"))
			fmt.Printf("%s: %s\n", labelStyle.Render("Настройки политики"), valueStyle.Render(registry.RetentionPolicy))
		} else {
			fmt.Printf("%s: %s\n", labelStyle.Render("Политика удаления"), valueStyle.Render("Отключена"))
		}
	},
}

func init() {
	getCmd.Flags().StringVarP(&getIdentifier, "identifier", "i", "", "Имя или ID реестра")
}
