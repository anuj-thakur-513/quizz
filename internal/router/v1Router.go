package router

import (
	"github.com/anuj-thakur-513/quizz/pkg/core"
	"github.com/gin-gonic/gin"
)

func V1Router(router *gin.RouterGroup) {
	router.GET("/", func(c *gin.Context) {
		c.JSON(200, core.ApiResponse(200, "Quizz Backend Working Fine", nil))
	})
}
