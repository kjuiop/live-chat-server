package redis

import (
	"github.com/go-redis/redis/v8"
	"live-chat-server/config"
	"live-chat-server/internal/utils"
	"log"
	"os"
	"testing"
	"time"
)

var testClient *TestClient

type TestClient struct {
	cfg *config.EnvConfig
	rdb *redis.Client
}

func TestMain(m *testing.M) {

	if err := utils.LoadTestEnv(); err != nil {
		log.Fatalf("failed load test env, err : %v", err)
	}

	cfg, err := config.LoadEnvConfig()
	if err != nil {
		log.Fatalf("fail to read config err : %v", err)
	}

	client := redis.NewClient(&redis.Options{
		Addr:         "127.0.0.1:6379",
		DialTimeout:  time.Second * 3,
		ReadTimeout:  time.Second * 3,
		WriteTimeout: time.Second * 3,
	})

	testClient = &TestClient{
		cfg: cfg,
		rdb: client,
	}

	exitCode := m.Run()

	if err := client.Close(); err != nil {
		log.Printf("failed close redis err : %v\n", err)
	}

	os.Exit(exitCode)
}
