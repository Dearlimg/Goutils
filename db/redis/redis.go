package redis

import (
	"context"
	"github.com/go-redis/redis/v8"
)

func RedisInit(Addr, Password string, PoolSize, DB int) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     Addr,
		Password: Password,
		DB:       DB,
		PoolSize: PoolSize,
	})
	_, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		panic(err)
	}
	return rdb
}
