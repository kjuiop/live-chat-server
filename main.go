package main

import (
	"live-chat-server/config"
	"live-chat-server/logger"
	"log"
	"log/slog"
)

var BUILD_TIME = "no flag of BUILD_TIME"
var GIT_HASH = "no flag of GIT_HASH"
var APP_VERSION = "no flag of APP_VERSION"

func main() {

	cfg, err := config.LoadEnvConfig()
	if err != nil {
		log.Fatalf("fail to read config err : %v\n", err)
	}

	if err := cfg.CheckValid(); err != nil {
		log.Fatalf("fail to invalid config, err : %v\n", err)
	}

	if err := logger.SlogInit(cfg.Logger); err != nil {
		log.Fatalf("fail to init slog err : %v\n", err)
	}

	slog.Debug("live chat server app start", "git_hash", GIT_HASH, "build_time", BUILD_TIME, "app_version", APP_VERSION)

}
