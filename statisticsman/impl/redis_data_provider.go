package impl

import (
	"context"
	"strconv"

	"github.com/go-redis/redis/v8"
	"github.com/sgostarter/i/commerr"
	"github.com/sgostarter/libeasygo/statisticsman/inters"
)

func NewRedisDataProvider(redisCli *redis.Client) inters.DataProvider {
	return &redisDataProviderImpl{
		redisCli: redisCli,
	}
}

type redisDataProviderImpl struct {
	redisCli *redis.Client
}

func (impl *redisDataProviderImpl) Exists(k string) (exists bool, err error) {
	n, err := impl.redisCli.Exists(context.Background(), k).Result()

	if err != nil {
		return
	}

	exists = n > 0

	return
}

func (impl *redisDataProviderImpl) ScanEx(k string, cb inters.DataScannerCB, reset bool) error {
	return impl.scanImpl(k, cb, reset)
}

func (impl *redisDataProviderImpl) Scan(k string, cb inters.DataScannerCB) error {
	return impl.scanImpl(k, cb, false)
}

func (impl *redisDataProviderImpl) scanImpl(k string, cb inters.DataScannerCB, reset bool) error {
	if cb == nil {
		return commerr.ErrInvalidArgument
	}

	var cursor uint64

	var keys []string

	var err error

	for {
		keys, cursor, err = impl.redisCli.HScan(context.Background(), k, cursor, "", 100).Result()
		// nolint: nestif
		if err != nil {
			err = cb(k, "", 0, err)
		} else {
			for idx := 0; idx < len(keys); idx += 2 {
				var v int64
				v, err = strconv.ParseInt(keys[idx+1], 0, 64)
				err = cb(k, keys[idx], v, err)
				if err != nil {
					break
				}
				if reset {
					err = impl.redisCli.HIncrBy(context.Background(), k, keys[idx], -v).Err()
					if err != nil {
						break
					}
				}
			}
		}

		if err != nil {
			break
		}

		if cursor <= 0 {
			break
		}
	}

	return err
}

func (impl *redisDataProviderImpl) Delete(k string) error {
	return impl.redisCli.Del(context.Background(), k).Err()
}
