package server

import (
	"context"
	"errors"
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"live-chat-server/api/controller"
	"live-chat-server/api/middleware"
	"live-chat-server/config"
	redis "live-chat-server/internal/redis"
	"live-chat-server/repository"
	"live-chat-server/usecase"
	"log/slog"
	"net/http"
	"sync"
)

type Gin struct {
	srv *http.Server
	cfg config.Server
}

func NewGinServer(cfg *config.EnvConfig, redis redis.Client) Client {

	serverCfg := cfg.Server
	router := getGinEngine(serverCfg.Mode)

	router.Use(middleware.LoggingMiddleware)
	router.Use(middleware.RecoveryErrorReport())
	router.Use(cors.New(cors.Config{
		AllowWebSockets:  true,
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"*"},
		AllowCredentials: true,
	}))

	api := router.Group("/api")

	setupSystemGroup(api)
	setupRoomGroup(api, cfg.RoomPolicy, redis)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", serverCfg.Port),
		Handler: router,
	}

	return &Gin{
		srv: srv,
		cfg: cfg.Server,
	}
}

func (g *Gin) Run(wg *sync.WaitGroup) {
	defer wg.Done()

	err := g.srv.ListenAndServe()
	if errors.Is(err, http.ErrServerClosed) {
		slog.Debug("server close")
	} else {
		slog.Error("run server error", "error", err)
	}
}

func (g *Gin) Shutdown(ctx context.Context) {
	if err := g.srv.Shutdown(ctx); err != nil {
		slog.Error("error during server shutdown", "error", err)
	}
}

func getGinEngine(mode string) *gin.Engine {
	switch mode {
	case "prod":
		return gin.New()
	case "test":
		gin.SetMode(gin.TestMode)
		return gin.Default()
	default:
		return gin.Default()
	}
}

func setupSystemGroup(router *gin.RouterGroup) {
	systemController := controller.NewSystemController()
	router.GET("/health-check", systemController.GetHealth)
	router.GET("/panic-test", systemController.OccurPanic)
}

func setupRoomGroup(router *gin.RouterGroup, cfg config.RoomPolicy, redis redis.Client) {
	rr := repository.NewRoomRepository(redis)
	ur := usecase.NewRoomUseCase(rr)
	roomController := controller.NewRoomController(cfg, ur)

	router.POST("/room", roomController.CreateRoom)
}
