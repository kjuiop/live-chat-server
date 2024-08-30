package app

import (
	"context"
	"live-chat-server/config"
	redis "live-chat-server/internal/redis"
	"live-chat-server/internal/server"
	"log"
	"sync"
)

type App struct {
	cfg config.EnvConfig
	srv server.Client
}

func NewApplication(ctx context.Context, cfg config.EnvConfig) *App {

	redisClient, err := redis.NewRedisSingleClient(ctx, cfg.Redis)
	if err != nil {
		log.Fatalf("fail to connect redis client")
	}

	srv := server.NewGinServer(&cfg, redisClient)

	return &App{
		cfg: cfg,
		srv: srv,
	}
}

func (a *App) Start(wg *sync.WaitGroup) {
	a.srv.Run(wg)
}

func (a *App) Stop(ctx context.Context) {
	a.srv.Shutdown(ctx)
}