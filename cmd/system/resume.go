package system

import (
	"context"
	"fmt"

	"github.com/charmbracelet/log"
	"github.com/cloud-ru/evo-ai-agents-cli/internal/di"
	"github.com/spf13/cobra"
)

// resumeCmd represents the resume command
var resumeCmd = &cobra.Command{
	Use:   "resume <system-id>",
	Short: "Возобновление работы системы агентов",
	Long:  "Возобновляет работу приостановленной системы агентов",
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

		// Возобновляем работу системы

		if err := apiClient.AgentSystems.Resume(ctx, systemID); err != nil {
			log.Fatal("Failed to resume system", "error", err, "system_id", systemID)
		}

		fmt.Printf("✅ Система агентов %s возобновлена успешно!\n", systemID)
	},
}

func init() {
	RootCMD.AddCommand(resumeCmd)
}
