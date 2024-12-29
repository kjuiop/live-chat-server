package redis

import (
	"bufio"
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"
	"net"
	"strconv"
	"strings"
	"testing"
	"time"
)

// redis-benchmark
// select 0
// redis 는 초당 10만건 정도는 수용할 수 있다.

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

func TestRedisInfo(t *testing.T) {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Redis INFO 명령어 실행
	info, err := testClient.rdb.Info(ctx).Result()
	assert.NoError(t, err)

	// INFO 명령어가 빈 문자열을 반환하지 않는지 확인
	assert.NotEmpty(t, info)

	t.Logf("redis info : %s", info)

	// 특정 섹션에서 "used_memory"가 포함되어 있는지 확인
	assert.True(t, strings.Contains(info, "used_memory"), "INFO 결과에 'used_memory'가 포함되어야 합니다.")
}

func TestRedisStat(t *testing.T) {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 1초 동안 Redis에 몇 개의 명령어를 발행
	for i := 0; i < 100; i++ {
		err := testClient.rdb.Set(ctx, "key:"+strconv.Itoa(i), "value", 0).Err()
		assert.NoError(t, err)
	}

	// 통계 정보를 얻기 위해 잠시 대기
	time.Sleep(1 * time.Second)

	// Redis INFO 명령어 실행
	info, err := testClient.rdb.Info(ctx).Result()
	assert.NoError(t, err)

	// INFO 결과에서 'instantaneous_ops_per_sec' 확인
	lines := strings.Split(info, "\n")
	var instantaneousOpsPerSec string
	for _, line := range lines {
		if strings.HasPrefix(line, "instantaneous_ops_per_sec:") {
			instantaneousOpsPerSec = strings.TrimPrefix(line, "instantaneous_ops_per_sec:")
			break
		}
	}

	assert.NotEmpty(t, instantaneousOpsPerSec, "INFO 에 'instantaneous_ops_per_sec'가 포함되어야 합니다.")

	// 초당 처리된 명령어 수가 0보다 큰지 확인
	opsPerSec, err := strconv.Atoi(strings.TrimSpace(instantaneousOpsPerSec))
	assert.NoError(t, err)

	t.Logf("ops per sec : %d", opsPerSec)
	assert.Greater(t, opsPerSec, 0, "초당 처리된 명령어 수가 0보다 커야 합니다.")

	// 키 삭제
	for i := 0; i < 100; i++ {
		err := testClient.rdb.Del(ctx, "key:"+strconv.Itoa(i)).Err()
		assert.NoError(t, err)
	}
}

func TestRedisMonitor(t *testing.T) {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := net.Dial("tcp", "127.0.0.1:6379")
	assert.NoError(t, err)
	defer conn.Close()

	// MONITOR 명령 전송
	_, err = conn.Write([]byte("MONITOR\r\n"))
	assert.NoError(t, err)

	// 데이터베이스에 키-값 저장
	err = testClient.rdb.Set(ctx, "key:monitor", "value for monitor", 0).Err()
	assert.NoError(t, err)

	// MONITOR 출력 읽기
	scanner := bufio.NewScanner(conn)
	go func() {
		for scanner.Scan() {
			t.Logf("MONITOR: %s", scanner.Text())
		}
	}()

	time.Sleep(2 * time.Second)

	// 클린업: 데이터 삭제
	err = testClient.rdb.Del(ctx, "key:monitor").Err()
	assert.NoError(t, err)

	// 에러가 발생한 경우 오류 출력
	if err := scanner.Err(); err != nil {
		t.Errorf("Error reading from monitor: %v", err)
	}
}
