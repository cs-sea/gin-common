package services

import (
	"fmt"
	"time"

	"github.com/cs-sea/gin-common/common"
	"github.com/cs-sea/gin-common/contract"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

const rateLimitPrefix = "rate_limit_prefix:"

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
		rule := rule
		go func() {
			r.interval(ctx, rule.Key)
		}()
	}

	return r
}

func (r *RateLimitService) addBuckets(ctx *gin.Context, rule contract.LimiterBucketRule) {
	r.redis.SetNX(ctx, r.buildKey(rule.Key), 0, 0)
}

func (r *RateLimitService) GetToken(ctx *gin.Context, key string, count int64) error {
	tokenCount := r.redis.DecrBy(ctx, r.buildKey(key), count).Val()
	if tokenCount < 0 {
		r.redis.IncrBy(ctx, r.buildKey(key), count)
		return errors.New("not valid token")
	}

	return nil
}

func (r *RateLimitService) addToken(ctx *gin.Context, key string, count int64) {
	r.redis.IncrBy(ctx, r.buildKey(key), count)
}

func (r *RateLimitService) buildKey(key string) string {
	return rateLimitPrefix + key
}

func (r *RateLimitService) interval(ctx *gin.Context, key string) {
	var ticker *time.Ticker = time.NewTicker(time.Second)

	for true {
		select {
		case <-ticker.C:
			fmt.Println(11)
			r.addToken(ctx, key, 10)
		}
	}

}
