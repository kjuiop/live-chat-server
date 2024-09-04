package main

import (
	"context"
	"live-chat-server/cmd/mysql_kafka/worker/app"
	"live-chat-server/config"
	"live-chat-server/internal/logger"
	"live-chat-server/internal/reporter"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

var BUILD_TIME = "no flag of BUILD_TIME"
var GIT_HASH = "no flag of GIT_HASH"
var APP_VERSION = "no flag of APP_VERSION"

func main() {

	wg := sync.WaitGroup{}
	ctx, cancel := context.WithCancel(context.Background())

	cfg, err := config.LoadEnvConfig()
	if err != nil {
		log.Fatalf("fail to read config err : %v", err)
	}

	if err := cfg.CheckValid(); err != nil {
		log.Fatalf("fail to invalid config, err : %v", err)
	}

	reporter.NewSlackReporter(cfg.Slack)

	if err := logger.SlogInit(cfg.Logger); err != nil {
		log.Fatalf("fail to init slog err : %v", err)
	}

	slog.Debug("live chat server app start", "git_hash", GIT_HASH, "build_time", BUILD_TIME, "app_version", APP_VERSION)

	a := app.NewApplication(ctx, cfg)

	wg.Add(1)
	go a.Start(&wg)

	<-exitSignal()
	a.Stop(ctx)
	cancel()
	wg.Wait()
	slog.Debug("live chat server app gracefully stopped")
}

func exitSignal() <-chan os.Signal {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	return sig
}
