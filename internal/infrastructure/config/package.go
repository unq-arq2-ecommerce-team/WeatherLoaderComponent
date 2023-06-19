package config

import (
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	loggerPkg "github.com/unq-arq2-ecommerce-team/WeatherLoaderComponent/internal/infrastructure/logger"
	"strings"
	"time"
)

const (
	OtlServiceName = "weather-loader"
	ServiceName    = "WeatherLoaderComponent"

	EnvDockerCompose = "docker-compose"
)

type Config struct {
	Environment    string        `required:"true" default:"development"`
	Port           int           `required:"true" default:"8080"`
	LogLevel       string        `split_words:"true" default:"DEBUG"`
	LokiHost       string        `split_words:"true" required:"true"`
	Otel           OtelConfig    `split_words:"true" required:"true"`
	Mongo          MongoConfig   `split_words:"true" required:"true"`
	TickerLoopTime time.Duration `split_words:"true" default:"60m"`
	Weather        Weather       `required:"true"`
}

// IsIntegrationEnv return true if Environment is equal to EnvDockerCompose (no case sensitive)
func (c Config) IsIntegrationEnv() bool {
	return strings.EqualFold(c.Environment, EnvDockerCompose)
}

type Weather struct {
	ApiUrl       string     `split_words:"true" required:"true"`
	GeocodingUrl string     `split_words:"true" required:"true"`
	ApiKey       string     `split_words:"true" required:"true"`
	Lat          string     `required:"true"`
	Long         string     `required:"true"`
	HttpConfig   HttpConfig `split_words:"true"`
}

type HttpConfig struct {
	OtelEnabled bool          `required:"true" default:"false"`
	Timeout     time.Duration `default:"10s"`
	Retries     int           `default:"0"`
	RetryWait   time.Duration `split_words:"true" default:"15s"`
}

type MongoConfig struct {
	URI      string        `split_words:"true" required:"true"`
	Database string        `split_words:"true" required:"true"`
	Timeout  time.Duration `split_words:"true" required:"true"`
}

type OtelConfig struct {
	URL string `split_words:"true" required:"true"`
}

func LoadConfig() Config {
	defaultLogger := loggerPkg.DefaultLogger(ServiceName, loggerPkg.JsonFormat)

	// Auto load ".env" file
	err := godotenv.Load()
	if err != nil {
		defaultLogger.Error("error loading .env file")
	}
	var config Config
	if err := envconfig.Process("", &config); err != nil {
		panic(err.Error())
	}
	return config
}
