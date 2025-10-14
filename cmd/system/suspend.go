package system

import (
	"context"
	"fmt"

	"github.com/charmbracelet/log"
	"github.com/cloud-ru/evo-ai-agents-cli/internal/di"
	"github.com/spf13/cobra"
)

// suspendCmd represents the suspend command
var suspendCmd = &cobra.Command{
	Use:   "suspend <system-id>",
	Short: "Приостановка работы системы агентов",
	Long:  "Приостанавливает работу системы агентов",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		systemID := args[0]

		// Получаем API клиент из DI контейнера
		container := di.GetContainer()
		apiClient := container.GetAPI()

		// Приостанавливаем работу системы
		err := apiClient.AgentSystems.Suspend(ctx, systemID)
		if err != nil {
			log.Fatal("Failed to suspend system", "error", err, "system_id", systemID)
		}

		fmt.Printf("✅ Система агентов %s приостановлена успешно!\n", systemID)
	},
}

func init() {
	RootCMD.AddCommand(suspendCmd)
}
