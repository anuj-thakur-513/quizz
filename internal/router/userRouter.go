package router

import (
	"github.com/anuj-thakur-513/quizz/internal/controllers"
	"github.com/anuj-thakur-513/quizz/internal/middlewares"
	"github.com/gin-gonic/gin"
)

func UserRouter(router *gin.RouterGroup) {
	router.POST("/signup", controllers.Signup)
	router.GET("/login", controllers.Login)
	router.GET("/authCheck", middlewares.AuthCheck(), controllers.AuthCheck)
}
