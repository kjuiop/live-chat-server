package redis

import (
	"context"
	"github.com/stretchr/testify/assert"
	"live-chat-server/config"
	"testing"
	"time"
)

func TestNewRedisClient(t *testing.T) {

	tests := []struct {
		title    string
		mode     string
		address  string
		master   string
		password string
	}{
		{
			title:   "redis mode",
			mode:    "cluster",
			address: "127.0.0.1:6371",
		},
	}

	for _, tc := range tests {
		t.Run(tc.title, func(t *testing.T) {
			ta := assert.New(t)

			cfg := config.Redis{
				Mode:     tc.mode,
				Addr:     tc.address,
				Masters:  tc.master,
				Password: tc.password,
			}

			ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
			defer cancel()

			_, err := NewUniversalClient(ctx, cfg)
			ta.NoError(err)
		})
	}

}
