package main

import (
	"context"
	"live-chat-server/cmd/mysql_kafka/worker/app"
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

	a := app.NewApplication(ctx)

	wg.Add(1)
	go a.Start(&wg)

	slog.Debug("live chat worker app start", "git_hash", GIT_HASH, "build_time", BUILD_TIME, "app_version", APP_VERSION)

	<-exitSignal()
	a.Stop(ctx)
	cancel()
	wg.Wait()
	slog.Debug("live chat worker app gracefully stopped")
}

func exitSignal() <-chan os.Signal {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	return sig
}
