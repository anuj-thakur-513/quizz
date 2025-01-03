package router

import (
	"github.com/anuj-thakur-513/quizz/internal/controllers"
	"github.com/anuj-thakur-513/quizz/internal/middlewares"
	"github.com/gin-gonic/gin"
)

func AdminRouter(router *gin.RouterGroup) {
	router.POST("/create-quiz", middlewares.AdminCheck(), controllers.CreateQuiz)
}
