package main

import (
	"fmt"

	"github.com/caarlos0/env/v9"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

type Config struct {
	GithubToken string `env:"BAR_GH_TOKEN"`

	LogLevel string `env:"BAR_LOG_LEVEL" envDefault:"info"`
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

	logrus.SetLevel(logrus.InfoLevel)
	if config.LogLevel == "debug" {
		logrus.SetLevel(logrus.DebugLevel)
	}

	logrus.Info("Loaded config: ", fmt.Sprintf("%+v", config))
}
