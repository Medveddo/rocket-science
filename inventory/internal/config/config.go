package config

import (
	"os"

	"github.com/joho/godotenv"

	"github.com/Medveddo/rocket-science/inventory/internal/config/env"
)

var appConfig *config

type config struct {
	Logger  LoggerConfig
	InventoryGRPC InventoryGRPCConfig
	Mongo   MongoConfig
	HTTP    HTTPConfig
}

func Load(path ...string) error {
	err := godotenv.Load(path...)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	loggerCfg, err := env.NewLoggerConfig()
	if err != nil {
		return err
	}

	inventoryGRPCCfg, err := env.NewInventoryGRPCConfig()
	if err != nil {
		return err
	}

	mongoCfg, err := env.NewMongoConfig()
	if err != nil {
		return err
	}

	httpCfg, err := env.NewHTTPConfig()
	if err != nil {
		return err
	}

	appConfig = &config{
		Logger:  loggerCfg,
		InventoryGRPC: inventoryGRPCCfg,
		Mongo:   mongoCfg,
		HTTP:    httpCfg,
	}

	return nil
}

func AppConfig() *config {
	return appConfig
}
