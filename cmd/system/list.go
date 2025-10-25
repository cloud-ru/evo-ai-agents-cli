package system

import (
	"context"
	"fmt"
	"os"

	"github.com/cloud-ru/evo-ai-agents-cli/internal/di"
	"github.com/cloud-ru/evo-ai-agents-cli/internal/errors"
	"github.com/cloud-ru/evo-ai-agents-cli/internal/ui"
	"github.com/spf13/cobra"
)

var (
	systemOutputFormat string
	systemLimit        int
	systemOffset       int
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Просмотр списка систем агентов",
	Long:  "Показывает список всех агентных систем с возможностью пагинации",
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()

		// Создаем обработчик ошибок
		errorHandler := errors.NewHandler()

		// Получаем API клиент из DI контейнера
		container := di.GetContainer()
		apiClient, err := container.GetAPI()
		if err != nil {
			appErr := errorHandler.WrapAPIError(err, "API_CLIENT_ERROR", "Ошибка получения API клиента")
			appErr = appErr.WithSuggestions(
				"Проверьте переменные окружения: IAM_KEY_ID, IAM_SECRET_KEY, IAM_ENDPOINT",
				"Убедитесь что вы авторизованы: ai-agents-cli auth login",
				"Проверьте доступность API: curl -I $IAM_ENDPOINT",
				"Обратитесь к администратору для получения учетных данных",
				"📚 Подробная документация: https://cloud.ru/docs/ai-agents/ug/index?source-platform=Evolution",
			)
			fmt.Println(errorHandler.HandlePlain(appErr))
			os.Exit(1)
		}

		if systemOutputFormat == "json" {
			// Выводим в JSON формате
			systems, err := apiClient.AgentSystems.List(ctx, systemLimit, systemOffset)
			if err != nil {
				appErr := errorHandler.WrapAPIError(err, "SYSTEMS_LIST_FAILED", "Ошибка получения списка систем")
				appErr = appErr.WithSuggestions(
					"Проверьте переменные окружения: IAM_KEY_ID, IAM_SECRET_KEY, IAM_ENDPOINT",
					"Убедитесь что вы авторизованы: ai-agents-cli auth login",
					"Проверьте доступность API: curl -I $IAM_ENDPOINT",
					"Обратитесь к администратору для получения учетных данных",
					"📚 Подробная документация: https://cloud.ru/docs/ai-agents/ug/index?source-platform=Evolution",
				)
				fmt.Println(errorHandler.HandlePlain(appErr))
				os.Exit(1)
			}

			// Выводим JSON
			fmt.Printf(`{"data":%d,"total":%d}`, len(systems.Data), systems.Total)
			return
		}

		// Показываем интерактивную таблицу
		if err = ui.ShowAgentSystemsListFromAPI(ctx, systemLimit, systemOffset); err != nil {
			appErr := errorHandler.WrapAPIError(err, "SYSTEMS_TABLE_ERROR", "Ошибка отображения таблицы систем")
			appErr = appErr.WithSuggestions(
				"Проверьте переменные окружения: IAM_KEY_ID, IAM_SECRET_KEY, IAM_ENDPOINT",
				"Убедитесь что вы авторизованы: ai-agents-cli auth login",
				"Проверьте доступность API: curl -I $IAM_ENDPOINT",
				"Обратитесь к администратору для получения учетных данных",
				"📚 Подробная документация: https://cloud.ru/docs/ai-agents/ug/index?source-platform=Evolution",
			)
			fmt.Println(errorHandler.HandlePlain(appErr))
			os.Exit(1)
		}
	},
}

func init() {
	RootCMD.AddCommand(listCmd)

	listCmd.Flags().StringVarP(&systemOutputFormat, "output", "o", "table", "Формат вывода (table, json)")
	listCmd.Flags().IntVarP(&systemLimit, "limit", "l", 20, "Количество систем на странице")
	listCmd.Flags().IntVarP(&systemOffset, "offset", "", 0, "Смещение для пагинации")
}
