package delayqueue

import (
	"context"
	"errors"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/sgostarter/libeasygo/helper"
)

func NewRedisReadyQueue(ctx context.Context, redisCli *redis.Client) ReadyPool {
	return &redisReadyQueueImpl{
		ctx:      ctx,
		redisCli: redisCli,
	}
}

type redisReadyQueueImpl struct {
	ctx      context.Context
	redisCli *redis.Client
}

func (impl *redisReadyQueueImpl) NewReadyJob(topic, jobID string) (err error) {
	return helper.RunWithTimeout4Redis(impl.ctx, func(ctx context.Context) error {
		return impl.redisCli.RPush(ctx, topic, jobID).Err()
	})
}

func (impl *redisReadyQueueImpl) GetReadyJob(timeout time.Duration, topics ...string) (jid *JobIdentify, err error) {
	vs, err := impl.redisCli.BLPop(impl.ctx, timeout, topics...).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			err = nil
		}

		return
	}

	if len(vs) == 0 {
		return
	}

	jid = &JobIdentify{
		Topic: vs[0],
		ID:    vs[1],
	}

	return
}
