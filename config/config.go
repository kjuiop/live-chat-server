package config

import (
	"fmt"
	"github.com/kelseyhightower/envconfig"
)

type EnvConfig struct {
	Server Server
	Logger Logger
	Slack  Slack
	Redis  Redis
	Policy Policy
	Mysql  Mysql
	Kafka  Kafka
}

type Server struct {
	Mode           string `envconfig:"LCS_ENV" default:"dev"`
	Port           string `envconfig:"LCS_SERVER_PORT" default:"8090"`
	TrustedProxies string `envconfig:"LCS_TRUSTED_PROXIES" default:"127.0.0.1/32"`
}

type Logger struct {
	Level       string `envconfig:"LCS_LOG_LEVEL" default:"debug"`
	Path        string `envconfig:"LCS_LOG_PATH" default:"./logs/access.log"`
	PrintStdOut bool   `envconfig:"LOG_STDOUT" default:"true"`
}

type Slack struct {
	WebhookReportUrl string `envconfig:"LCS_SLACK_WEBHOOK_REPORT_URL"`
}

type Redis struct {
	Mode     string `envconfig:"LCS_REDIS_MODE" default:"single"`
	Addr     string `envconfig:"LCS_REDIS_ADDR" default:":6379"`
	Password string `envconfig:"KUS_REDIS_PASSWORD"`
	Masters  string `envconfig:"KUS_REDIS_MASTERS"`
	PoolSize int    `envconfig:"KUS_REDIS_POOL_SIZE" default:"100"`
}

type Mysql struct {
	Host     string `envconfig:"LCS_MYSQL_HOST" default:"localhost:3306"`
	Driver   string `envconfig:"LCS_MYSQL_DATABASE" default:"mysql"`
	User     string `envconfig:"LCS_MYSQL_USER" default:"root"`
	Password string `envconfig:"LCS_MYSQL_PASSWORD" default:"1234"`
	Database string `envconfig:"LCS_MYSQL_DATABASE" default:"chatting"`
}

type Policy struct {
	Prefix         string `envconfig:"LCS_ROOM_PREFIX" default:"N1,N2"`
	ContextTimeout int    `envconfig:"LCS_CONTEXT_TIMEOUT" default:"60"`
}

type Kafka struct {
	URL             string `envconfig:"LCS_KAFKA_URL" default:"localhost:9292"`
	GroupID         string `envconfig:"LCS_KAFKA_GROUP_ID" default:"chat-consumer-1"`
	ClientID        string `envconfig:"LCS_KAFKA_CLIENT_ID" default:"chat-producer-1"`
	ConsumerTimeout int    `envconfig:"LCS_KAFKA_CONSUMER_TIMEOUT" default:"1000"`
}

func LoadEnvConfig() (*EnvConfig, error) {
	var config EnvConfig
	if err := envconfig.Process("lcs", &config); err != nil {
		return nil, err
	}

	if err := config.CheckValid(); err != nil {
		return nil, err
	}

	return &config, nil
}

func (c *EnvConfig) CheckValid() error {

	if c.Redis.Addr == "" {
		return fmt.Errorf("check, redis addr is empty")
	}

	return nil
}
