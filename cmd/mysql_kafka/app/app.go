package app

import (
	"context"
	"live-chat-server/api/route"
	"live-chat-server/config"
	"live-chat-server/database/mysql"
	"live-chat-server/repository/room"
	"live-chat-server/server"
	"live-chat-server/usecase"
	"log"
	"sync"
	"time"
)

type App struct {
	cfg *config.EnvConfig
	srv server.Client
	db  mysql.Client
}

func NewApplication(ctx context.Context, cfg *config.EnvConfig) *App {

	db, err := mysql.NewMysqlSingleClient(ctx, cfg.Mysql)
	if err != nil {
		log.Fatalf("fail to connect redis client")
	}

	srv := server.NewGinServer(cfg)

	app := &App{
		cfg: cfg,
		srv: srv,
		db:  db,
	}

	app.setupRouter()

	return app
}

func (a *App) Start(wg *sync.WaitGroup) {
	a.srv.Run(wg)
}

func (a *App) Stop(ctx context.Context) {
	a.srv.Shutdown(ctx)
}

func (a *App) setupRouter() {

	timeout := time.Duration(a.cfg.Policy.ContextTimeout) * time.Second

	api := a.srv.GetRouterGroup("/api")
	ws := a.srv.GetRouterGroup("/ws")

	route.SetupSystemGroup(api)

	rr := room.NewRoomMysqlRepository(a.db)
	ur := usecase.NewRoomUseCase(rr, timeout)
	route.SetupRoomGroup(api, a.cfg.Policy, ur)

	cu := usecase.NewChatUseCase(ur, timeout)
	route.SetupChatGroup(ws, cu)
}
