package redis

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestRedisSentinelConnect(t *testing.T) {

	testAssert := assert.New(t)

	repoConn := redis.NewFailoverClient(&redis.FailoverOptions{
		MasterName:    "primary",
		SentinelAddrs: []string{"127.0.0.1:26379"},
		DialTimeout:   time.Second * 3,
		ReadTimeout:   time.Second * 3,
		WriteTimeout:  time.Second * 3,
	})

	err := repoConn.Ping(context.Background()).Err()
	testAssert.NoError(err, "expected connection to succeed")
}
