package routers

import (
	"fmt"

	"github.com/cs-sea/gin-common/contract"

	"github.com/cs-sea/gin-common/internal/services"

	"github.com/cs-sea/gin-common/common"
	"github.com/cs-sea/gin-common/models"

	"github.com/gin-gonic/gin"
)

func SetTestGroup(r gin.IRouter, middlewares ...gin.HandlerFunc) {
	test := r.Group("test").Use(middlewares...)

	test.GET("v1", func(c *gin.Context) {
		db := common.NewDB()

		s := &models.ApiUser{}
		db.WithContext(c).Debug().Where("sdf = ?", "sd").Find(s)

		redis := common.NewRedis()

		rs := services.NewRateLimitService(redis)
		rs.AddBuckets(c, contract.LimiterBucketRule{
			Key:          "test",
			FillInterval: 1,
			Capacity:     100,
			Quantum:      10,
		})

		fmt.Println("222")
	})
}
