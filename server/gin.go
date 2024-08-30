package server

import (
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"live-chat-server/api/middleware"
	"live-chat-server/api/route"
	"live-chat-server/config"
	"live-chat-server/database"
	"log/slog"
	"net/http"
	"sync"
	"time"
)

type Gin struct {
	srv *http.Server
	cfg config.Server
}

func NewGinServer(cfg *config.EnvConfig, db database.Client) Client {

	serverCfg := cfg.Server
	router := getGinEngine(serverCfg.Mode)

	router.Use(middleware.LoggingMiddleware)
	router.Use(middleware.RecoveryErrorReport())
	router.Use(middleware.SetCorsPolicy())

	timeout := time.Duration(cfg.Policy.ContextTimeout) * time.Second

	api := router.Group("/api")
	ws := router.Group("/ws")
	route.Setup(api, ws, cfg.Policy, timeout, db)

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
