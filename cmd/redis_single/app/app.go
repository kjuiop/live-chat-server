package app

import (
	"context"
	"live-chat-server/api/controller"
	"live-chat-server/api/route"
	"live-chat-server/config"
	"live-chat-server/internal/database/redis"
	cu "live-chat-server/internal/domain/chat/usecase"
	"live-chat-server/internal/domain/room/repository"
	ru "live-chat-server/internal/domain/room/usecase"
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
	db  redis.Client
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

	db, err := redis.NewRedisSingleClient(ctx, cfg.Redis)
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

	// repository
	roomRepository := repository.NewRoomRedisRepository(a.db)

	// use_case
	roomUseCase := ru.NewRoomUseCase(roomRepository, timeout)
	chatUseCase := cu.NewChatUseCase(roomUseCase, timeout)

	// controller
	systemController := controller.NewSystemController()
	roomController := controller.NewRoomController(a.cfg.Policy, roomUseCase)
	chatController := controller.NewChatController(chatUseCase)

	router := route.RouterConfig{
		Engine:           a.srv.GetEngine(),
		SystemController: systemController,
		RoomController:   roomController,
		ChatController:   chatController,
	}
	router.Setup()
}
