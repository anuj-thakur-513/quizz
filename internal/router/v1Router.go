package router

import "github.com/gin-gonic/gin"

func V1Router(router *gin.RouterGroup) {
	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"succcess": true,
			"message":  "v1 routes working fine",
		})
	})
}
