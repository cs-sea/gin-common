package middlewares

import (
	"net/http"

	"github.com/cs-sea/gin-common/contract"
	"github.com/gin-gonic/gin"
)

func NewAuthMiddleware(jwtService contract.JWT) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("token")
		_, err := jwtService.CheckToken(token)

		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		c.Next()
	}
}
