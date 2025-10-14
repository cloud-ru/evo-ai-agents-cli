package main

import (
	"context"
	"os"

	"github.com/charmbracelet/log"
	"github.com/cloudru/ai-agents-cli/cmd"
	"github.com/cloudru/ai-agents-cli/internal/config"
	"github.com/cloudru/ai-agents-cli/internal/di"
)

func main() {
	ctx := context.Background()

	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
		return
	}

	shutdown := initTracing(ctx, cfg)
	defer shutdown()

	// Инициализируем DI контейнер
	log.Debug("Инициализация DI контейнера")
	container := di.GetContainer()
	log.Debug("DI контейнер инициализирован успешно")

	defer func() {
		if err := container.Close(); err != nil {
			log.Error("Ошибка закрытия DI контейнера", "error", err)
		}
	}()

	if err := cmd.RootCMD.ExecuteContext(ctx); err != nil {
		os.Exit(1)
	}
}
