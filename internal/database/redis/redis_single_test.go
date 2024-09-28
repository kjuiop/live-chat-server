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

func TestRedisSetGet(t *testing.T) {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := testClient.rdb.Set(ctx, "test_key", "test_value", 0).Err()
	assert.NoError(t, err)

	val, err := testClient.rdb.Get(ctx, "test_key").Result()
	assert.NoError(t, err)
	assert.Equal(t, "test_value", val)

	t.Logf("test_key: %s", val)

	err = testClient.rdb.Del(ctx, "test_key").Err()
	assert.NoError(t, err)
}
