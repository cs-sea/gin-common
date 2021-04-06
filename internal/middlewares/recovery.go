package middlewares

import "github.com/gin-gonic/gin"

func NewRecoveryMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// 异常捕获，panic 才会走到这，gin 其实已经实现了，我们加一层拦截 以便告警
				c.Abort()
			}
		}()
		c.Next()
	}
}
