package config

import (
	"github.com/kelseyhightower/envconfig"
)

type EnvConfig struct {
	Server Server
	Logger Logger
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

func LoadEnvConfig() (*EnvConfig, error) {

	var config EnvConfig
	if err := envconfig.Process("fua", &config); err != nil {
		return nil, err
	}

	return &config, nil
}

func (c *EnvConfig) CheckValid() error {
	return nil
}
