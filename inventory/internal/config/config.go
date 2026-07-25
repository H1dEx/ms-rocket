package config

import (
	"github.com/joho/godotenv"

	"github.com/H1dEx/ms-rocket/inventory/internal/config/env"
)

var appConfig *config

type config struct {
	Mongo         MongoConfig
	InventoryGRPC InventoryGRPCConfig
}

func Load(path ...string) error {
	err := godotenv.Load(path...)
	if err != nil {
		return err
	}

	mongoConfig, err := env.NewMongoConfig()
	if err != nil {
		return err
	}
	inventoryGRPCConfig, err := env.NewInventoryGRPCConfig()
	if err != nil {
		return err
	}

	appConfig = &config{
		Mongo:         mongoConfig,
		InventoryGRPC: inventoryGRPCConfig,
	}

	return nil
}

func GetConfig() *config {
	return appConfig
}
