package redlock

import (
	"context"
	"testing"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"
)

func testGetRedis() *redis.Client {
	options, err := redis.ParseURL("redis://:redis_default_pass@127.0.0.1:8900/1")
	// options, cuserror := redis.ParseURL("redis://:@127.0.0.1:8901/0")
	if err != nil {
		panic(err)
	}

	return redis.NewClient(options)
}

func TestBase(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	rlm := NewRedisLockManager(ctx, testGetRedis(), nil)
	lock, err := rlm.TryLock("lock-1")
	assert.Nil(t, err)
	assert.NotNil(t, lock)

	err = lock.Unlock()
	assert.Nil(t, err)

	rlm.Terminal()
	rlm.Wait()
}

func TestReentrantBase(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	rlm := NewRedisLockManager(ctx, testGetRedis(), nil)
	lock, err := rlm.TryReentrantLock("lock-1", "gid1")
	assert.Nil(t, err)
	assert.NotNil(t, lock)

	lock2, err := rlm.TryReentrantLock("lock-1", "gid1")
	assert.Nil(t, err)
	assert.NotNil(t, lock2)

	err = lock.Unlock()
	assert.Nil(t, err)

	err = lock2.Unlock()
	assert.Nil(t, err)

	rlm.Terminal()
	rlm.Wait()
}

func TestTTL(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	rlm := NewRedisLockManager(ctx, testGetRedis(), nil)
	lock, err := rlm.TryLock("lock-1")
	assert.Nil(t, err)
	assert.NotNil(t, lock)

	time.Sleep(time.Second * 10)

	err = lock.Unlock()
	assert.Nil(t, err)

	rlm.Terminal()
	rlm.Wait()
}
