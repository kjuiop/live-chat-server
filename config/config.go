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
}

type Server struct {
	Mode string `envconfig:"LCS_ENV" default:"dev"`
	Port string `envconfig:"LCS_SERVER_PORT" default:"8090"`
}

type Logger struct {
	Level       string `envconfig:"LCS_LOG_LEVEL" default:"debug"`
	Path        string `envconfig:"LCS_LOG_PATH" default:"./logs/access.log"`
	PrintStdOut bool   `envconfig:"LOG_STDOUT" default:"false"`
}

type Slack struct {
	WebhookReportUrl string `envconfig:"LCS_SLACK_WEBHOOK_REPORT_URL" default:""`
}

type Redis struct {
	Addr string `envconfig:"LCS_REDIS_ADDR" default:":6379"`
}

type Policy struct {
	Prefix         string `envconfig:"LCS_ROOM_PREFIX" default:"N1,N2"`
	ContextTimeout int    `envconfig:"LCS_CONTEXT_TIMEOUT" default:"3"`
}

func LoadEnvConfig() (*EnvConfig, error) {
	var config EnvConfig
	if err := envconfig.Process("fua", &config); err != nil {
		return nil, err
	}
	return &config, nil
}

func (c *EnvConfig) CheckValid() error {

	if c.Redis.Addr == "" {
		return fmt.Errorf("check, redis addr is empty")
	}

	if c.Slack.WebhookReportUrl == "" {
		return fmt.Errorf("check, slack webhook report url is empty")
	}

	return nil
}
