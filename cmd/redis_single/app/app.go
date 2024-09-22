package app

import (
	"context"
	"fmt"
	"live-chat-server/api/controller"
	"live-chat-server/api/route"
	"live-chat-server/config"
	"live-chat-server/internal/database/redis"
	cu "live-chat-server/internal/domain/chat/usecase"
	rr "live-chat-server/internal/domain/room/repository"
	ru "live-chat-server/internal/domain/room/usecase"
	"live-chat-server/internal/domain/system"
	spq "live-chat-server/internal/domain/system/pubsub"
	sr "live-chat-server/internal/domain/system/repository"
	su "live-chat-server/internal/domain/system/usecase"
	"live-chat-server/internal/logger"
	"live-chat-server/internal/mq/kafka"
	"live-chat-server/internal/reporter"
	"live-chat-server/internal/server"
	"log"
	"net"
	"sync"
	"time"
)

type App struct {
	cfg           config.EnvConfig
	srv           server.Client
	db            redis.Client
	mq            kafka.Client
	systemUseCase system.UseCase
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

	mq := kafka.NewMemoryClient()

	srv := server.NewGinServer(cfg)

	app := &App{
		cfg: *cfg,
		srv: srv,
		db:  db,
		mq:  mq,
	}

	app.setupRouter()
	app.registerServer()

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

	// mq
	systemPubSub := spq.NewMemorySystemPubSub(a.cfg.Kafka, nil)

	// repository
	roomRepository := rr.NewRoomRedisRepository(a.db)
	systemRepository := sr.NewSystemRedisRepository(a.cfg.Redis)
	// use_case
	roomUseCase := ru.NewRoomUseCase(roomRepository, timeout)
	chatUseCase := cu.NewChatUseCase(roomUseCase, timeout)
	systemUseCase := su.NewSystemUseCase(context.TODO(), systemRepository, systemPubSub)

	// controller
	systemController := controller.NewSystemController(systemUseCase)
	roomController := controller.NewRoomController(a.cfg.Policy, roomUseCase)
	chatController := controller.NewChatController(chatUseCase)

	router := route.RouterConfig{
		Engine:           a.srv.GetEngine(),
		SystemController: systemController,
		RoomController:   roomController,
		ChatController:   chatController,
	}
	router.Setup()
	a.systemUseCase = systemUseCase
}

func (a *App) registerServer() {

	addrs, err := net.InterfaceAddrs()
	if err != nil {
		log.Fatalf("failed parsing ip address, err : %v", err)
	}

	var ip net.IP
	for _, addr := range addrs {
		if ipNet, ok := addr.(*net.IPNet); ok {
			if !ipNet.IP.IsLoopback() && ipNet.IP.To4() != nil {
				ip = ipNet.IP
				break
			}
		}
	}

	if ip == nil {
		log.Fatalln("no ip address found")
	}

	addr := fmt.Sprintf("%s:%s", ip.String(), a.cfg.Server.Port)
	if err := a.systemUseCase.ChatServerSet(addr, true); err != nil {
		log.Fatalf("failed register server info, address : %s, err : %v", addr, err)
	}

}
