package router

import (
	"time"

	"github.com/anuj-thakur-513/quizz/internal/controllers"
	"github.com/anuj-thakur-513/quizz/internal/middlewares"
	"github.com/gin-gonic/gin"
)

func UserRouter(router *gin.RouterGroup) {
	rl := middlewares.NewRateLimiter(5, 10*time.Minute)

	router.POST("/signup", rl.Limit(), controllers.Signup)
	router.POST("/login", rl.Limit(), controllers.Login)
	router.GET("/authCheck", middlewares.AuthCheck(), controllers.AuthCheck)
}
