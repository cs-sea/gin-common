package middlewares

import (
	"context"
	"github.com/gin-gonic/gin"
	"time"
)

// 超时控制中间件
func NewTimeoutMiddleware(timeout time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c.Request.Context(), timeout)
		defer cancel()

		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}