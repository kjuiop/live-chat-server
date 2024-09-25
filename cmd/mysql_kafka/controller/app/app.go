package app

import (
	"context"
	"live-chat-server/api/controller"
	"live-chat-server/api/route"
	"live-chat-server/config"
	"live-chat-server/internal/database/mysql"
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
	"log/slog"
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

	mq, err := kafka.NewKafkaConsumerClient(cfg.Kafka)
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

	if err := app.setupRouter(ctx); err != nil {
		log.Fatalf("failed initialize router, err : %v", err)
	}

	if err := app.initProcess(); err != nil {
		log.Fatalf("failed initialized process, err : %v", err)
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

func (a *App) setupRouter(ctx context.Context) error {

	timeout := time.Duration(a.cfg.Policy.ContextTimeout) * time.Second

	// mq
	systemPubSub := spq.NewSystemPubSub(a.cfg.Kafka, a.mq)

	// repository
	roomRepository := rr.NewRoomMysqlRepository(a.db)
	systemRepository := sr.NewSystemMySqlRepository(a.db)

	// use_case
	roomUseCase := ru.NewRoomUseCase(roomRepository, timeout)
	systemUseCase := su.NewSystemUseCase(ctx, systemRepository, systemPubSub)

	// controller
	roomController := controller.NewRoomController(a.cfg.Policy, roomUseCase)
	systemController := controller.NewSystemController(systemUseCase)

	router := route.RouterConfig{
		Engine:           a.srv.GetEngine(),
		SystemController: systemController,
		RoomController:   roomController,
	}
	router.ApiSetup()
	a.systemUseCase = systemUseCase
	return nil
}

func (a *App) initProcess() error {

	if err := a.systemUseCase.RegisterSubTopic("chat"); err != nil {
		return err
	}

	return nil
}

func (a *App) LoopServerInfo(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		select {
		case <-ctx.Done():
			slog.Debug("close Loop Sub Kafka goroutine")
			return
		default:
			event, err := a.systemUseCase.LoopSubKafka(a.cfg.Kafka.ConsumerTimeout)
			if err != nil {
				slog.Error("received event error", "error", err)
				continue
			}

			if event == nil {
				continue
			}

			slog.Debug("received event", "event", string(event.Value))
		}
	}

}
