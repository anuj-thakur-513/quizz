package router

import (
	"github.com/anuj-thakur-513/quizz/pkg/core"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173", "https://quizz.anuj-thakur.com", "https://www.quizz.anuj-thakur.com"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
		AllowWebSockets:  true,
	}))

	gin.SetMode(gin.ReleaseMode)

	router.GET("/", func(c *gin.Context) {
		c.JSON(200, core.ApiResponse(200, "Quizz Backend Working Fine", nil))
	})

	v1 := router.Group("/api/v1")
	V1Router(v1)

	return router
}
