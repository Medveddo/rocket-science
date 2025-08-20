package config

type LoggerConfig interface {
	Level() string
	AsJson() bool
}

type HTTPConfig interface {
	Address() string
	Port() string
	Host() string
}

type PaymentGRPCConfig interface {
	Address() string
}
