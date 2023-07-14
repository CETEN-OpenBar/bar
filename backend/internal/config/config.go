package config

import (
	"fmt"

	"github.com/caarlos0/env/v9"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

type Config struct {
	MongoConfig struct {
		ConnectionURI string `env:"URI" envDefault:"mongodb://localhost:27017"`
		Database      string `env:"DB" envDefault:"bar"`
		Timeout       int    `env:"TIMEOUT" envDefault:"30"`
	} `envPrefix:"BAR_MONGO_"`

	ApiConfig struct {
		SessionSecret      string `env:"SESSION_SECRET"`
		AdminSessionSecret string `env:"ADMIN_SESSION_SECRET"`
		Port               string `env:"PORT" envDefault:":8080"`
	} `envPrefix:"BAR_API_"`
}

var config Config

func GetConfig() Config {
	return config
}

func init() {
	godotenv.Load()
	if err := env.Parse(&config); err != nil {
		logrus.Fatal(err)
	}

	logrus.Info("Loaded config: ", fmt.Sprintf("%+v", config))
}
