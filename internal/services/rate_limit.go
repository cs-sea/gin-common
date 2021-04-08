package services

import (
	"fmt"
	"time"

	"github.com/cs-sea/gin-common/common"
	"github.com/cs-sea/gin-common/contract"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

const rateLimitBucket = "rate_limit_bucket"

type bucket struct {
	fillInterval time.Duration
	// 上一次填充时间
	preTime  time.Time
	quantum  uint64
	capacity uint64
}

type RateLimitService struct {
	redis *common.Redis
}

var _ contract.RateLimit = &RateLimitService{}

func NewRateLimitService(redis *common.Redis) *RateLimitService {
	return &RateLimitService{redis}
}

func (r *RateLimitService) AddBuckets(ctx *gin.Context, rules ...contract.LimiterBucketRule) contract.RateLimit {
	for _, rule := range rules {
		r.addBuckets(ctx, rule)

		go func(rule contract.LimiterBucketRule) {
			r.interval(ctx, rule)
		}(rule)
	}

	return r
}

func (r *RateLimitService) addBuckets(ctx *gin.Context, rule contract.LimiterBucketRule) {
	err := r.redis.WithContext(ctx).ZIncrBy(ctx, rateLimitBucket, 0, rule.Key).Err()
	if err != nil {
		common.Logger(ctx).WithError(err)
	}
}

func (r *RateLimitService) GetToken(ctx *gin.Context, key string, count int64) error {
	return r.getToken(ctx, key, count)
}

func (r *RateLimitService) getToken(ctx *gin.Context, key string, count int64) error {
	tokenCount := r.redis.ZIncrBy(ctx, rateLimitBucket, -float64(count), key).Val()
	if tokenCount < 0 {
		r.redis.ZIncrBy(ctx, rateLimitBucket, float64(count), key)
		return errors.New("not valid token")
	}

	return nil
}

func (r *RateLimitService) addToken(ctx *gin.Context, rule contract.LimiterBucketRule) {
	script := `
local key = KEYS[1]
local score = ARGV[1]
local member = ARGV[2]
local max = ARGV[3]
local currentVal = redis.call("zincrby", key, score, member)

if (tonumber(max) < tonumber(currentVal))
then
	redis.call("zadd", key, max, member)
end
return currentVal;
`
	fmt.Printf("%+v\n", rule)
	err := r.redis.Eval(ctx, script, []string{rateLimitBucket}, rule.Quantum, rule.Key, rule.Capacity).Val()
	fmt.Println(err)
}

func (r *RateLimitService) interval(ctx *gin.Context, rule contract.LimiterBucketRule) {
	newDuration := time.Second * time.Duration(rule.FillInterval)
	var ticker *time.Ticker = time.NewTicker(newDuration)

	for true {
		select {
		case <-ticker.C:
			r.addToken(ctx, rule)
		}
	}

}
