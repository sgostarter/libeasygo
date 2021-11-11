package redlock

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/sgostarter/libeasygo/cuserror"
)

type RedisLock interface {
	Unlock() (err error)
	RedisKey() string
}

type redisLock struct {
	redisCli   *redis.Client
	key        string
	token      string
	timeout    time.Duration
	fnUnlockOb func(key string)
}

var unlockScript = redis.NewScript(`
	if redis.call("get", KEYS[1]) == ARGV[1]
	then
		return redis.call("del", KEYS[1])
	else
		return 0
	end
`)

func (lock *redisLock) tryLock() (bool, error) {
	return lock.redisCli.SetNX(context.TODO(), lock.key, lock.token, lock.timeout).Result()
}

func (lock *redisLock) Unlock() (err error) {
	sha1, err := unlockScript.Load(context.TODO(), lock.redisCli).Result()
	if err != nil {
		return
	}

	f, err := lock.redisCli.EvalSha(lock.redisCli.Context(), sha1, []string{lock.key}, lock.token).Result()

	if err != nil {
		return
	}

	if lock.fnUnlockOb != nil {
		lock.fnUnlockOb(lock.key)
	}

	if f.(int64) != 1 {
		err = cuserror.NewWithErrorMsg(fmt.Sprintf("%d", f))

		return
	}

	return
}

func (lock *redisLock) RedisKey() string {
	return lock.key
}

func TryLock(key string) (redisCli *redis.Client, lock RedisLock, err error) {
	lock, err = TryLockWithTimeout(redisCli, key, defaultTimeout)
	if err != nil {
		return
	}

	return
}

func TryLockWithTimeout(redisCli *redis.Client, key string, timeout time.Duration) (lock RedisLock, err error) {
	return tryLockWithTimeout(redisCli, key, timeout, nil)
}

func tryLockWithTimeout(redisCli *redis.Client, key string, timeout time.Duration, unlockOb func(key string)) (lock RedisLock, err error) {
	lockImpl := &redisLock{
		redisCli:   redisCli,
		key:        key,
		token:      uuid.New().String(),
		timeout:    timeout,
		fnUnlockOb: unlockOb,
	}

	ok, err := lockImpl.tryLock()
	if err == nil && ok {
		lock = lockImpl
	}

	return
}
