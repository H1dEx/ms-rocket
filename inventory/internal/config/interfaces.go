package config

type MongoConfig interface {
	URI() string
	DatabaseName() string
}

type InventoryGRPCConfig interface {
	Address() string
}
