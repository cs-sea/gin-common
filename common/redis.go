package common

import "github.com/go-redis/redis/v8"

type Redis struct {
	*redis.Client
}

var RedisClient Redis
