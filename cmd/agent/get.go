package agent

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/charmbracelet/log"
	"github.com/cloudru/ai-agents-cli/internal/di"
	"github.com/cloudru/ai-agents-cli/internal/ui"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

var (
	agentOutputFormat string
)

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:   "get <agent-id>",
	Short: "Получить информацию об агенте",
	Long:  "Показывает подробную информацию о конкретном агенте",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		agentID := args[0]

		// Получаем API клиент из DI контейнера
		container := di.GetContainer()
		apiClient := container.GetAPI()

		// Получаем информацию об агенте
		agent, err := apiClient.Agents.Get(ctx, agentID)
		if err != nil {
			log.Fatal("Failed to get agent", "error", err, "agent_id", agentID)
		}

		if agentOutputFormat == "json" {
			// Выводим в JSON формате
			encoder := json.NewEncoder(os.Stdout)
			encoder.SetIndent("", "  ")
			if err := encoder.Encode(agent); err != nil {
				log.Fatal("Failed to encode JSON", "error", err)
			}
			return
		}

		// Показываем простую статичную версию
		result := ui.RenderAgentDetails(agent, ctx, container)
		fmt.Println(result)
	},
}

// getCreatedByInfo получает информацию о создателе агента
func getCreatedByInfo(ctx context.Context, container *di.Container, userID string) string {
	if userID == "" {
		return "Не указан"
	}

	config := container.GetConfig()
	if config.CustomerID == "" {
		// Если нет customerID, возвращаем ID с пояснением
		return fmt.Sprintf("ID: %s (CUSTOMER_ID не указан)", userID)
	}

	apiClient := container.GetAPI()
	user, err := apiClient.Users.Get(ctx, config.CustomerID, userID)
	if err != nil {
		// При ошибке API тоже показываем ID
		return fmt.Sprintf("ID: %s (ошибка получения данных)", userID)
	}

	return ui.FormatUserName(user.ID, user.FirstName, user.LastName, user.Email)
}

// getUpdatedByInfo получает информацию об изменяющем агента
func getUpdatedByInfo(ctx context.Context, container *di.Container, userID string) string {
	if userID == "" {
		return "Не указан"
	}

	config := container.GetConfig()
	if config.CustomerID == "" {
		// Если нет customerID, возвращаем ID с пояснением
		return fmt.Sprintf("ID: %s (CUSTOMER_ID не указан)", userID)
	}

	apiClient := container.GetAPI()
	user, err := apiClient.Users.Get(ctx, config.CustomerID, userID)
	if err != nil {
		// При ошибке API тоже показываем ID
		return fmt.Sprintf("ID: %s (ошибка получения данных)", userID)
	}

	return ui.FormatUserName(user.ID, user.FirstName, user.LastName, user.Email)
}

// isTerminal проверяет, является ли терминал терминалом
func isTerminal() bool {
	return term.IsTerminal(int(os.Stdout.Fd()))
}

func init() {
	RootCMD.AddCommand(getCmd)

	getCmd.Flags().StringVarP(&agentOutputFormat, "output", "o", "table", "Формат вывода (table, json)")
}
