package config

import "time"

const (
	BearerType    string = "Bearer"
	Authorization        = "Authorization"
	AccountUUID          = "accountUUID"
	PrincipalType        = "principalType"
	AuthToken            = "authToken"
	SYSTEM               = "SYSTEM"
)

type Configuration struct {
	Ports                 Ports
	Flags                 Flags
	Auth                  Auth
	DefaultHierarchy      DefaultHierarchy
	MongoConfig           MongoConfig
	MaxQueryExecutionTime time.Duration `split_words:"true" default:"4s"`
	Conferences           DefaultConferences
	SlackConfig           SlackConfig
	SkipSecureVerify      bool `split_words:"true" default:"true"`
}

type SlackConfig struct {
	EnableSlackClient bool   `envconfig:"ENABLE_SLACK_CLIENT" split_words:"true"`
	ConnectionMode    string `envconfig:"CONNECTION_MODE" split_words:"true"`
	SlackAppToken     string `envconfig:"SLACK_APP_TOKEN" split_words:"true"`
	SlackBotToken     string `envconfig:"SLACK_BOT_TOKEN" split_words:"true"`
	IncidentChannelID string `envconfig:"INCIDENT_CHANNEL_ID" split_words:"true"`
}

type DefaultConferences struct {
	ZoomLink string `envconfig:"ZOOM_LINK" split_words:"true"`
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

type Auth struct {
	PasswordHashCost    int    `envconfig:"PASSWORD_HASH_COST" default:"10"`
	DefaultUserID       string `envconfig:"DEFAULT_USER_ID" default:"admin"`
	DefaultUserEmail    string `envconfig:"DEFAULT_USER_EMAIL" default:"admin@respondnow.io"`
	DefaultUserName     string `envconfig:"DEFAULT_USER_NAME" default:"Admin"`
	DefaultUserPassword string `envconfig:"DEFAULT_USER_PASSWORD" default:"respondnow"`
	JWTSecret           string `envconfig:"JWT_SECRET" default:"83820c2001b55fb6c401d08cfd050973ff31084d5b1e66478cf0f951cc4a9e60"`
}

type DefaultHierarchy struct {
	DefaultAccountId   string `envconfig:"DEFAULT_ACCOUNT_ID" default:"default_account_id"`
	DefaultAccountName string `envconfig:"DEFAULT_ACCOUNT_NAME" default:"Default Account"`
	DefaultOrgId       string `envconfig:"DEFAULT_ORG_ID" default:"default_org_id"`
	DefaultOrgName     string `envconfig:"DEFAULT_ORG_NAME" default:"Default Org"`
	DefaultProjectId   string `envconfig:"DEFAULT_PROJECT_ID" default:"default_project_id"`
	DefaultProjectName string `envconfig:"DEFAULT_PROJECT_NAME" default:"Default Project"`
}

var EnvConfig Configuration
