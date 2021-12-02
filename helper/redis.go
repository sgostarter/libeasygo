package helper

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

func NewRedisClient(dsn string) (redisCli *redis.Client, err error) {
	options, err := redis.ParseURL(dsn)
	if err != nil {
		return
	}

	redisCli = redis.NewClient(options)

	ctx, cf := context.WithTimeout(context.Background(), 3*time.Second)
	defer cf()

	err = redisCli.Ping(ctx).Err()
	if err != nil {
		return
	}

	return
}
