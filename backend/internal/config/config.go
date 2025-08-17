package config

import (
	"github.com/caarlos0/env/v9"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

type Config struct {
	MongoConfig struct {
		ConnectionURI string `env:"URI" envDefault:"mongodb://localhost:27017"`
		Database      string `env:"DB" envDefault:"bar"`
		Timeout       uint64 `env:"TIMEOUT" envDefault:"30"`
	} `envPrefix:"BAR_MONGO_"`

	ApiConfig struct {
		SessionSecret      string `env:"SESSION_SECRET"`
		AdminSessionSecret string `env:"ADMIN_SESSION_SECRET"`
		Port               string `env:"PORT" envDefault:":8080"`
		BasePath           string `env:"BASE_PATH" envDefault:"http://localhost:8080"`
		FrontendBasePath   string `env:"FRONTEND_BASE_PATH" envDefault:"http://localhost:5173"`
		LocalToken         string `env:"LOCAL_TOKEN"`
	} `envPrefix:"BAR_API_"`

	OauthConfig struct {
		GoogleClientID     string `env:"GOOGLE_CLIENT_ID"`
		GoogleClientSecret string `env:"GOOGLE_CLIENT_SECRET"`
	} `envPrefix:"BAR_OAUTH_"`

	HelloAssoConfig struct {
		URL          string `env:"URL" envDefault:"https://api.helloasso.com"`
		ClientID     string `env:"CLIENT_ID"`
		ClientSecret string `env:"CLIENT_SECRET"`
		Slug         string `env:"SLUG"`
	} `envPrefix:"HELLOASSO_"`

	StorageConfig struct {
		StoragePath string `env:"STORAGE_PATH" envDefault:"./storage"`
	} `envPrefix:"BAR_STORAGE_"`
	DiscordWebhookURL string `env:"BAR_DISCORD_WEBHOOK_URL" envDefault:""`
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
	switch config.LogLevel {
	case "info":
		logrus.SetLevel(logrus.InfoLevel)
	case "debug":
		logrus.SetLevel(logrus.DebugLevel)
	case "warn":
		logrus.SetLevel(logrus.WarnLevel)
	case "error":
		logrus.SetLevel(logrus.ErrorLevel)
	case "fatal":
		logrus.SetLevel(logrus.FatalLevel)
	case "panic":
		logrus.SetLevel(logrus.PanicLevel)
	case "trace":
		logrus.SetLevel(logrus.TraceLevel)
	}

	// Always add date to logrus
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})

	logrus.Info("Loaded config.")
}
