package main

import (
	"github.com/cs-sea/gin-common/routers"

	"github.com/cs-sea/gin-common/internal/middlewares"

	"github.com/gin-gonic/gin"
)

func main() {
	g := gin.New()

	logMiddleware := middlewares.LogMiddleware()

	routers.SetTestGroup(g, logMiddleware)

	g.Run("localhost:1333")
}
