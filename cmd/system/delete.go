package system

import (
	"context"
	"fmt"

	"github.com/charmbracelet/log"
	"github.com/cloud-ru/evo-ai-agents-cli/internal/di"
	"github.com/spf13/cobra"
)

var (
	systemDeleteForce bool
)

// deleteCmd represents the delete command
var deleteCmd = &cobra.Command{
	Use:   "delete <system-id>",
	Short: "Удаление системы агентов",
	Long:  "Удаляет систему агентов. Используйте --force для принудительного удаления",
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

		// Подтверждение удаления
		if !systemDeleteForce {
			fmt.Printf("⚠️  Вы уверены, что хотите удалить систему %s? (y/N): ", systemID)
			var confirmation string
			fmt.Scanln(&confirmation)
			if confirmation != "y" && confirmation != "Y" {
				fmt.Println("❌ Удаление отменено")
				return
			}
		}

		// Удаляем систему

		if err = apiClient.AgentSystems.Delete(ctx, systemID); err != nil {
			log.Fatal("Failed to delete system", "error", err, "system_id", systemID)
		}

		fmt.Printf("✅ Система агентов %s удалена успешно!\n", systemID)
	},
}

func init() {
	RootCMD.AddCommand(deleteCmd)

	deleteCmd.Flags().BoolVarP(&systemDeleteForce, "force", "f", false, "Принудительное удаление без подтверждения")
}
