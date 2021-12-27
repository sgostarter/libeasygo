package impl

import (
	"context"

	"github.com/go-redis/redis/v8"
	"github.com/sgostarter/i/l"
	"github.com/sgostarter/libeasygo/statisticsman/inters"
)

func NewRedisCounterStorage(redisCli *redis.Client, logger l.Wrapper) inters.Storage {
	if logger == nil {
		logger = l.NewNopLoggerWrapper()
	}

	return &redisCounterStorageImpl{
		redisCli: redisCli,
		logger:   logger,
	}
}

type redisCounterStorageImpl struct {
	redisCli *redis.Client
	logger   l.Wrapper
}

func (cs *redisCounterStorageImpl) Inc(key, field string, incV int64) {
	err := cs.redisCli.HIncrBy(context.Background(), key, field, incV).Err()
	if err != nil {
		cs.logger.WithFields(l.StringField("key", key), l.StringField("field", field),
			l.Int64Field("incV", incV)).Error("data_counter_failed")
	}
}
