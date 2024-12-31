package router

import (
	"github.com/anuj-thakur-513/quizz/internal/controllers"
	"github.com/gin-gonic/gin"
)

func UserRouter(router *gin.RouterGroup) {
	router.POST("/signup", controllers.CreateUser)
}
