package routers

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func SetTestGroup(r gin.IRouter, middlewares ...gin.HandlerFunc) {
	test := r.Group("test").Use(middlewares...)

	test.GET("v1", func(c *gin.Context) {
		fmt.Println("222")
	})
}
