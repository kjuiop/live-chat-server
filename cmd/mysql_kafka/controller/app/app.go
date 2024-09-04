package app

import (
	"context"
	"live-chat-server/api/controller"
	"live-chat-server/api/route"
	"live-chat-server/config"
	mysql2 "live-chat-server/internal/database/mysql"
	"live-chat-server/internal/domain/room/repository"
	"live-chat-server/internal/domain/room/usecase"
	server2 "live-chat-server/internal/server"
	"log"
	"sync"
	"time"
)

type App struct {
	cfg *config.EnvConfig
	srv server2.Client
	db  mysql2.Client
}

func NewApplication(ctx context.Context, cfg *config.EnvConfig) *App {

	db, err := mysql2.NewMysqlSingleClient(ctx, cfg.Mysql)
	if err != nil {
		log.Fatalf("fail to connect redis client")
	}

	srv := server2.NewGinServer(cfg)

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

	// repository
	roomRepository := repository.NewRoomMysqlRepository(a.db)

	// use_case
	roomUseCase := usecase.NewRoomUseCase(roomRepository, timeout)

	// controller
	systemController := controller.NewSystemController()
	roomController := controller.NewRoomController(a.cfg.Policy, roomUseCase)

	router := route.RouterConfig{
		Engine:           a.srv.GetEngine(),
		SystemController: systemController,
		RoomController:   roomController,
	}
	router.ApiSetup()
}
