package config

import (
	"github.com/caarlos0/env/v11"
	"github.com/charmbracelet/log"
	_ "github.com/joho/godotenv/autoload"
)

type ServiceName string

const (
	SenderApp    ServiceName = "ai-assistant-billing-sender"
	HandlerApp   ServiceName = "ai-assistant-billing"
	SchedulerApp ServiceName = "ai-assistant-billing-scheduler"
	ConsumerApp  ServiceName = "ai-assistant-billing-consumer"
)

// Config
// .go:generate go run github.com/g4s8/envdoc@latest -types='*' -output .env.example -format dotenv
type Config struct {
	ServiceConfig                   ServiceConfig `envPrefix:"SERVICE_"`
	ServiceAccountConfig            IAMConfig     `envPrefix:"IAM_"`
	BulkOperationsConcurrencyFactor int           `env:"BULK_OPERATIONS_CONCURRENCY" envDefault:"20"`

	IntegrationApiGrpcAddr string `env:"PUBLIC_API_ENDPOINT"          envDefault:"ai-agents.api.cloud.ru"`
}

func Load() (*Config, error) {
	cfg := &Config{}
	err := env.Parse(cfg)
	if err != nil {
		log.Fatalf("%+v", err)
		return nil, err
	}
	return cfg, err
}
