package server

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"live-chat-server/api/controller"
	"live-chat-server/api/route"
	"live-chat-server/config"
	"live-chat-server/internal/domain/system/usecase/mocks"
	"net/http"
	"sync"
	"testing"
	"time"
)

func TestRunAndShutdown(t *testing.T) {

	wg := &sync.WaitGroup{}
	cfg, _ := config.LoadEnvConfig()
	cfg.Server = config.Server{
		Mode:           "test",
		Port:           "8090",
		TrustedProxies: "127.0.0.1/32",
	}

	s := NewGinServer(cfg)
	sm := mocks.NewSystemUseCaseStub()
	router := route.RouterConfig{
		Engine:           s.GetEngine(),
		SystemController: controller.NewSystemController(sm),
	}
	router.SetupSystemRouter(router.Engine.Group("/api"))

	wg.Add(1)
	go s.Run(wg)

	time.Sleep(10 * time.Millisecond)

	resp, err := http.Get(fmt.Sprintf("http://localhost:%s/api/system/health-check", cfg.Server.Port))
	if err != nil {
		t.Fatalf("Failed to make GET request: %v", err)
	}

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Shutdown the server
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	s.Shutdown(ctx)
	wg.Wait()
}
