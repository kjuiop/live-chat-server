package app

import (
	"context"
	"fmt"
	"live-chat-server/api/controller"
	"live-chat-server/api/route"
	"live-chat-server/config"
	"live-chat-server/internal/database/mysql"
	cu "live-chat-server/internal/domain/chat/usecase"
	"live-chat-server/internal/domain/room/repository"
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
	cfg           *config.EnvConfig
	srv           server.Client
	db            mysql.Client
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

	db, err := mysql.NewMysqlSingleClient(ctx, cfg.Mysql)
	if err != nil {
		log.Fatalf("fail to connect redis client")
	}

	mq, err := kafka.NewKafkaProducerClient(cfg.Kafka)
	if err != nil {
		log.Fatalf("fail to connect kafka client")
	}

	srv := server.NewGinServer(cfg)

	app := &App{
		cfg: cfg,
		srv: srv,
		db:  db,
		mq:  mq,
	}

	if err := app.setupRouter(); err != nil {
		log.Fatalf("failed initialize router, err : %v", err)
	}

	app.registerServer()

	return app
}

func (a *App) Start(wg *sync.WaitGroup) {
	a.srv.Run(wg)
}

func (a *App) Stop(ctx context.Context) {
	a.srv.Shutdown(ctx)
	a.db.Close()
	a.mq.Close("producer")
}

func (a *App) setupRouter() error {

	timeout := time.Duration(a.cfg.Policy.ContextTimeout) * time.Second

	// mq
	systemPubSub := spq.NewSystemPubSub(a.cfg.Kafka, a.mq)

	// repository
	roomRepository := repository.NewRoomMysqlRepository(a.db)
	systemRepository := sr.NewSystemMySqlRepository(a.db)

	// use_case
	roomUseCase := ru.NewRoomUseCase(roomRepository, timeout)
	chatUseCase := cu.NewChatUseCase(roomUseCase, timeout)
	systemUseCase := su.NewSystemUseCase(context.TODO(), systemRepository, systemPubSub)

	chatController := controller.NewChatController(chatUseCase)

	router := route.RouterConfig{
		Engine:         a.srv.GetEngine(),
		ChatController: chatController,
	}
	router.WsSetup()
	a.systemUseCase = systemUseCase
	return nil
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
	if err := a.systemUseCase.SetChatServerInfo(addr, true); err != nil {
		log.Fatalf("failed register server info, address : %s, err : %v", addr, err)
	}

	a.systemUseCase.PublishServerStatusEvent(addr, true)
}
