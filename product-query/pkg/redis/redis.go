package redis

import (
	"github.com/redis/go-redis/v9"
)

func ConnectRedis(addr string) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: "",
		DB:       0,
		Protocol: 2,
	})
	return rdb
}
