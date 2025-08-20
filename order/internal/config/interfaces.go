package config

import "time"

type LoggerConfig interface {
	Level() string
	AsJson() bool
}

type HTTPConfig interface {
	Address() string
	Port() string
	Host() string
	ReadHeaderTimeout() time.Duration
	ShutdownTimeout() time.Duration
}

type PostgreSQLConfig interface {
	URI() string
	DatabaseName() string
}

type InventoryGRPCConfig interface {
	Address() string
}

type PaymentGRPCConfig interface {
	Address() string
}
