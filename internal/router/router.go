package router

import "github.com/gin-gonic/gin"

func SetupRouter() *gin.Engine {
	router := gin.Default()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	gin.SetMode(gin.ReleaseMode)

	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"succcess": true,
			"message":  "Server working fine",
		})
	})

	v1 := router.Group("/api/v1")
	V1Router(v1)

	return router
}
