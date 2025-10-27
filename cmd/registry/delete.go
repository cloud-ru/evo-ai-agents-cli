package registry

import (
	"context"
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
	"github.com/cloud-ru/evo-ai-agents-cli/internal/di"
	"github.com/spf13/cobra"
)

var (
	deleteIdentifier string
	confirmDelete    bool
)

// deleteCmd represents the delete command
var deleteCmd = &cobra.Command{
	Use:   "delete [registry-name-or-id]",
	Short: "Удалить реестр",
	Long:  "Удаляет реестр по имени или ID (требуется подтверждение)",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()

		// Определяем идентификатор реестра
		identifier := deleteIdentifier
		if len(args) > 0 {
			identifier = args[0]
		}

		if identifier == "" {
			log.Fatal("Registry name or ID is required")
		}

		// Если флаг подтверждения не установлен, запрашиваем подтверждение
		if !confirmDelete {
			log.Warn("Внимание! Удаление реестра необратимо и удалит все образы.")
			log.Fatal("Используйте флаг --confirm для подтверждения удаления")
		}

		// Получаем API клиент из DI контейнера
		container := di.GetContainer()
		apiClient, err := container.GetAPI()
		if err != nil {
			log.Fatal("Failed to get API client", "error", err)
		}

		// Удаляем реестр (возвращает операцию)
		operation, err := apiClient.Registries.Delete(ctx, identifier)
		if err != nil {
			log.Fatal("Failed to delete registry", "error", err)
		}

		// Выводим сообщение
		successStyle := lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("2")).
			Border(lipgloss.RoundedBorder()).
			Padding(0, 1)

		fmt.Println(successStyle.Render("✅ Операция удаления реестра запущена"))
		fmt.Println()
		fmt.Printf("ID операции: %s\n", operation.ID)
		
		if operation.Description != "" {
			fmt.Printf("Описание: %s\n", operation.Description)
		}
		
		if operation.Done {
			fmt.Println("Статус: Завершено")
		} else {
			fmt.Println("Статус: В процессе")
		}
	},
}

func init() {
	deleteCmd.Flags().StringVarP(&deleteIdentifier, "identifier", "i", "", "Имя или ID реестра")
	deleteCmd.Flags().BoolVarP(&confirmDelete, "confirm", "y", false, "Подтвердить удаление")
}
