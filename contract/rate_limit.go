package contract

import (
	"time"

	"github.com/gin-gonic/gin"
)

type RateLimit interface {
	AddBuckets(ctx *gin.Context, rules ...LimiterBucketRule) RateLimit
	GetToken(ctx *gin.Context, key string, count int64) error
}

type LimiterBucketRule struct {
	Key          string
	FillInterval time.Duration
	Capacity     int64
	Quantum      int64
}
