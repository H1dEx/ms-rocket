package config

import (
	"github.com/joho/godotenv"

	"github.com/H1dEx/ms-rocket/order/internal/config/env"
)

var appConfig *config

type config struct {
	Postgres      PostgresConfig
	OrderHTTP     OrderHTTPConfig
	InventoryGRPC InventoryGRPCConfig
	PaymentGRPC   PaymentGRPCConfig
}

func Load(path ...string) error {
	err := godotenv.Load(path...)
	if err != nil {
		return err
	}

	postgresConfig, err := env.NewPostgresConfig()
	if err != nil {
		return err
	}
	orderHTTPConfig, err := env.NewOrderHTTPConfig()
	if err != nil {
		return err
	}
	inventoryGRPCConfig, err := env.NewInventoryGRPCConfig()
	if err != nil {
		return err
	}
	paymentGRPCConfig, err := env.NewPaymentGRPCConfig()
	if err != nil {
		return err
	}

	appConfig = &config{
		Postgres:      postgresConfig,
		OrderHTTP:     orderHTTPConfig,
		InventoryGRPC: inventoryGRPCConfig,
		PaymentGRPC:   paymentGRPCConfig,
	}

	return nil
}

func GetConfig() *config {
	return appConfig
}
