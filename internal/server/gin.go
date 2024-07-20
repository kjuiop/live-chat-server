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

func NewGinServer(serverCfg config.Server) Client {

	router := getGinEngine(serverCfg.Mode)

	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	systemController := controller.NewSystemController()
	router.GET("/ping", systemController.GetHealth)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", serverCfg.Port),
		Handler: router,
	}

	return &Gin{
		srv: srv,
		cfg: serverCfg,
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
