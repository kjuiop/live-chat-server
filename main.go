package main

import (
	"context"
	"live-chat-server/config"
	"live-chat-server/internal/server"
	"live-chat-server/logger"
	"live-chat-server/repository"
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
		log.Fatalf("fail to read config err : %v\n", err)
	}

	if err := cfg.CheckValid(); err != nil {
		log.Fatalf("fail to invalid config, err : %v\n", err)
	}

	if err := logger.SlogInit(cfg.Logger); err != nil {
		log.Fatalf("fail to init slog err : %v\n", err)
	}

	slog.Debug("live chat server app start", "git_hash", GIT_HASH, "build_time", BUILD_TIME, "app_version", APP_VERSION)

	redisClient, err := repository.NewRedisClient(ctx, cfg.Redis)
	if err != nil {
		log.Fatalf("fail to connect redis client")
	}

	srv := server.NewGinServer(cfg, redisClient)

	wg.Add(1)
	go srv.Run(&wg)

	<-exitSignal()
	srv.Shutdown(ctx)
	cancel()
	wg.Wait()
	slog.Debug("live chat server app gracefully stopped")
}

func exitSignal() <-chan os.Signal {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	return sig
}
