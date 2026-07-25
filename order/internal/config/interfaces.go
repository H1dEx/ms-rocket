package config

type PostgresConfig interface {
	URI() string
	MigrationDir() string
}

type OrderHTTPConfig interface {
	Address() string
}

type InventoryGRPCConfig interface {
	Address() string
}

type PaymentGRPCConfig interface {
	Address() string
}
