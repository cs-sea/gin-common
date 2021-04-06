package contract

import (
	"github.com/gin-gonic/gin"
	"github.com/juju/ratelimit"
	"time"
)

type RateLimit interface {
	GetBucket(ctx *gin.Context, key string) (*ratelimit.Bucket, bool)
	AddBuckets(ctx *gin.Context, rules ...LimiterBucketRule) RateLimit
}


type LimiterBucketRule struct {
	Key          string
	FillInterval time.Duration
	Capacity     int64
	Quantum      int64
}