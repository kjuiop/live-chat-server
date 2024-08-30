package server

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"live-chat-server/config"
	"net/http"
	"sync"
	"testing"
	"time"
)

func TestRunAndShutdown(t *testing.T) {

	wg := &sync.WaitGroup{}
	cfg, _ := config.LoadEnvConfig()
	cfg.Server = config.Server{
		Mode: "test",
		Port: "8090",
	}

	s := NewGinServer(cfg)

	wg.Add(1)
	go s.Run(wg)

	time.Sleep(10 * time.Millisecond)

	resp, err := http.Get(fmt.Sprintf("http://localhost:%s/api/health-check", cfg.Server.Port))
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
