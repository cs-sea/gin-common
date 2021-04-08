package common

import (
	"github.com/go-redis/redis/v8"
)

type Redis struct {
	*redis.Client
}

var RedisClient Redis

func NewRedis() *Redis {

	redis.SetLogger(new(Writer))
	redisClient := redis.NewClient(&redis.Options{})
	return &Redis{redisClient}
}
