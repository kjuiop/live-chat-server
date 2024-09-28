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

func TestRedisSlowLog(t *testing.T) {

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	originalValue, err := testClient.rdb.ConfigGet(ctx, "slowlog-log-slower-than").Result()
	assert.NoError(t, err)

	// 모든 명령어를 slow log 기록
	err = testClient.rdb.ConfigSet(ctx, "slowlog-log-slower-than", "0").Err()
	assert.NoError(t, err)

	// 1초 sleep 명령어 실행
	_, err = testClient.rdb.Do(ctx, "DEBUG", "SLEEP", 1).Result() // 1초 Sleep
	assert.NoError(t, err)

	// 모든 slow log 를 조회
	slowLogs, err := testClient.rdb.SlowLogGet(ctx, -1).Result()
	assert.NoError(t, err)

	assert.True(t, len(slowLogs) > 0, "Slowlog 기록이 없습니다.")

	for _, log := range slowLogs {
		t.Logf("Slowlog ID: %d, Command: %v, Execution Time: %d µs\n",
			log.ID, log.Args, log.Duration)
	}

	err = testClient.rdb.ConfigSet(ctx, "slowlog-log-slower-than", originalValue[1].(string)).Err()
	assert.NoError(t, err)
}
