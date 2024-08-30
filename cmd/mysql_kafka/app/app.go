package app

import (
	"context"
	"live-chat-server/config"
	"live-chat-server/database"
	"live-chat-server/server"
	"log"
	"sync"
)

type App struct {
	cfg *config.EnvConfig
	srv server.Client
	db  database.Client
}

func NewApplication(ctx context.Context, cfg *config.EnvConfig) *App {

	db, err := database.NewMysqlSingleClient(ctx, cfg.Mysql)
	if err != nil {
		log.Fatalf("fail to connect redis client")
	}

	srv := server.NewGinServer(cfg, db)

	return &App{
		cfg: cfg,
		srv: srv,
		db:  db,
	}
}

func (a *App) Start(wg *sync.WaitGroup) {
	a.srv.Run(wg)
}

func (a *App) Stop(ctx context.Context) {
	a.srv.Shutdown(ctx)
}
