package services

import (
	"github.com/cs-sea/gin-common/contract"
	"github.com/gin-gonic/gin"
	"github.com/juju/ratelimit"
	"github.com/pkg/errors"
)

type Limiter struct {
	limiterBuckets map[string]*ratelimit.Bucket
}

type RateLimitService struct {
	*Limiter
}

var _ contract.RateLimit = &RateLimitService{}

func NewRateLimitService() *RateLimitService {
	return &RateLimitService{&Limiter{
		limiterBuckets: make(map[string]*ratelimit.Bucket),
	}}
}

func (r *RateLimitService) GetBucket(ctx *gin.Context, key string) (*ratelimit.Bucket, bool) {
	bucket, ok := r.limiterBuckets[key]
	return bucket, ok
}

func (r *RateLimitService) AddBuckets(ctx *gin.Context, rules ...contract.LimiterBucketRule) contract.RateLimit {
	for _, rule := range rules {
		if _, ok := r.limiterBuckets[rule.Key]; !ok {
			r.limiterBuckets[rule.Key] = ratelimit.NewBucketWithQuantum(rule.FillInterval, rule.Capacity, rule.Quantum)
		}
	}

	return r
}

func (r *RateLimitService) addBuckets(ctx *gin.Context, rule contract.LimiterBucketRule) {
	r.limiterBuckets[rule.Key] = ratelimit.NewBucketWithQuantum(rule.FillInterval, rule.Capacity, rule.Quantum)
}

func (r *RateLimitService) GetToken(ctx *gin.Context, bucket *ratelimit.Bucket) error {
	tokenCount := bucket.TakeAvailable(1)
	if tokenCount <= 0 {
		return errors.New("not valid token")
	}

	return nil
}
