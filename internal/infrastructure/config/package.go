package config

import (
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	loggerPkg "github.com/unq-arq2-ecommerce-team/WeatherLoaderComponent/internal/infrastructure/logger"
	"time"
)

const ServiceName = "WeatherLoaderComponent"

type Config struct {
	Environment    string        `required:"true" default:"development"`
	Port           int           `required:"true" default:"8080"`
	LogLevel       string        `split_words:"true" default:"DEBUG"`
	MongoURI       string        `split_words:"true" required:"true"`
	MongoDatabase  string        `split_words:"true" required:"true"`
	MongoTimeout   time.Duration `split_words:"true" required:"true"`
	TickerLoopTime time.Duration `split_words:"true" default:"60m"`
	Weather        Weather       `required:"true"`
}

type Weather struct {
	ApiUrl     string     `split_words:"true" required:"true"`
	ApiKey     string     `split_words:"true" required:"true"`
	Lat        string     `required:"true"`
	Long       string     `required:"true"`
	HttpConfig HttpConfig `split_words:"true"`
}

type HttpConfig struct {
	Timeout   time.Duration `default:"10s"`
	Retries   int           `default:"0"`
	RetryWait time.Duration `split_words:"true" default:"15s"`
}

func LoadConfig() Config {
	primitiveLogger := loggerPkg.New(&loggerPkg.Config{
		ServiceName: ServiceName,
		LogFormat:   loggerPkg.JsonFormat,
	})

	// Auto load ".env" file
	err := godotenv.Load()
	if err != nil {
		primitiveLogger.Error("error loading .env file")
	}
	var config Config
	if err := envconfig.Process("", &config); err != nil {
		panic(err.Error())
	}
	return config
}
