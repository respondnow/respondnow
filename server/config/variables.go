package config

import "time"

type Configuration struct {
	Ports                 Ports
	Flags                 Flags
	MongoConfig           MongoConfig
	MaxQueryExecutionTime time.Duration `split_words:"true" default:"4s"`
	SkipSecureVerify      bool          `split_words:"true" default:"true"`
}

type Ports struct {
	HttpPort   string `envconfig:"HTTP_PORT" default:"8080"`
	MetricPort string `envconfig:"METRIC_PORT" default:"8890"`
}

type Flags struct {
	EnableSLIMetrics bool `envconfig:"ENABLE_SLI_METRICS" default:"true"`
}

type MongoConfig struct {
	MongoURL        string `envconfig:"MONGO_URL" default:"mongodb://mongo1:27018,mongo2:27019,mongo3:27020"`
	MongoDbName     string `envconfig:"MONGO_DB_NAME" default:"respondnow"`
	MongoDbUsername string `envconfig:"MONGO_DB_USERNAME" required:"true" split_words:"true"`
	MongoDbPassword string `envconfig:"MONGO_DB_PASSWORD" required:"true" split_words:"true"`
}

var EnvConfig Configuration
