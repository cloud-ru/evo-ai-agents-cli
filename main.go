package main

import (
	"context"
	"github.com/charmbracelet/log"
	"github.com/cloudru/ai-agents-cli/cmd"
	"github.com/cloudru/ai-agents-cli/internal/config"
	"os"
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

	if err := cmd.RootCMD.ExecuteContext(ctx); err != nil {
		os.Exit(1)
	}
}
