package app

import (
	"context"
	"live-chat-server/api/controller"
	"live-chat-server/api/route"
	"live-chat-server/config"
	"live-chat-server/internal/database/mysql"
	rr "live-chat-server/internal/domain/room/repository"
	ru "live-chat-server/internal/domain/room/usecase"
	sr "live-chat-server/internal/domain/system/repository"
	su "live-chat-server/internal/domain/system/usecase"
	"live-chat-server/internal/logger"
	"live-chat-server/internal/reporter"
	"live-chat-server/internal/server"
	"log"
	"sync"
	"time"
)

type App struct {
	cfg *config.EnvConfig
	srv server.Client
	db  mysql.Client
}

func NewApplication(ctx context.Context) *App {

	cfg, err := config.LoadEnvConfig()
	if err != nil {
		log.Fatalf("fail to read config err : %v", err)
	}

	reporter.NewSlackReporter(cfg.Slack)

	if err := logger.SlogInit(cfg.Logger); err != nil {
		log.Fatalf("fail to init slog err : %v", err)
	}

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

	if err := app.setupRouter(); err != nil {
		log.Fatalf("failed initialize router, err : %v", err)
	}

	return app
}

func (a *App) Start(wg *sync.WaitGroup) {
	a.srv.Run(wg)
}

func (a *App) Stop(ctx context.Context) {
	a.srv.Shutdown(ctx)
	a.db.Close()
}

func (a *App) setupRouter() error {

	timeout := time.Duration(a.cfg.Policy.ContextTimeout) * time.Second

	// repository
	roomRepository := rr.NewRoomMysqlRepository(a.db)
	systemRepository := sr.NewSystemMySqlRepository(a.db)

	// use_case
	roomUseCase := ru.NewRoomUseCase(roomRepository, timeout)
	systemUseCase := su.NewSystemUseCase(systemRepository)

	// controller
	roomController := controller.NewRoomController(a.cfg.Policy, roomUseCase)
	systemController := controller.NewSystemController(systemUseCase)

	router := route.RouterConfig{
		Engine:           a.srv.GetEngine(),
		SystemController: systemController,
		RoomController:   roomController,
	}
	router.ApiSetup()
	return nil
}
