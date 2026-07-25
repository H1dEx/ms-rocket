package config

import (
	"github.com/joho/godotenv"

	"github.com/H1dEx/ms-rocket/payment/internal/config/env"
)

var appConfig *config

type config struct {
	PaymentGRPC PaymentGRPCConfig
}

func Load(path ...string) error {
	err := godotenv.Load(path...)
	if err != nil {
		return err
	}
	paymentGRPCConfig, err := env.NewPaymentGRPCConfig()
	if err != nil {
		return err
	}
	appConfig = &config{
		PaymentGRPC: paymentGRPCConfig,
	}
	return nil
}

func GetConfig() *config {
	return appConfig
}
