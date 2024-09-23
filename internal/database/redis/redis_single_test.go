package redis

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestRedisSingleConnect(t *testing.T) {

	testAssert := assert.New(t)

	client := redis.NewClient(&redis.Options{
		Addr:         "127.0.0.1:6379",
		DialTimeout:  time.Second * 3,
		ReadTimeout:  time.Second * 3,
		WriteTimeout: time.Second * 3,
	})

	err := client.Ping(context.Background()).Err()
	testAssert.NoError(err, "expected connection to succeed")
}
