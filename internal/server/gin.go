package server

import (
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"live-chat-server/api/controller"
	"live-chat-server/config"
	"log/slog"
	"net/http"
	"sync"
)

type Gin struct {
	srv *http.Server
	cfg config.Server
}

func NewGinServer(cfg *config.EnvConfig) Client {

	serverCfg := cfg.Server
	router := getGinEngine(serverCfg.Mode)

	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	api := router.Group("/api")

	systemController := controller.NewSystemController()
	setupSystemGroup(api, systemController)

	roomController := controller.NewRoomController(cfg.RoomPolicy)
	setupRoomGroup(api, roomController)

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

func setupSystemGroup(router *gin.RouterGroup, systemController *controller.SystemController) {
	router.GET("/health-check", systemController.GetHealth)
}

func setupRoomGroup(router *gin.RouterGroup, roomController *controller.RoomController) {
	router.POST("/room", roomController.CreateRoom)
}
