package config

type MongoConfig interface {
	URI() string
	DatabaseName() string
}

type InventoryGRPCConfig interface {
	Address() string
}

type LoggerConfig interface {
	Level() string
	AsJSON() bool
}
