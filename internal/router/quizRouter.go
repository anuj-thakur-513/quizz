package router

import (
	"github.com/anuj-thakur-513/quizz/internal/controllers"
	"github.com/anuj-thakur-513/quizz/internal/middlewares"
	"github.com/gin-gonic/gin"
)

func QuizRouter(router *gin.RouterGroup) {
	router.POST("/", middlewares.AdminCheck(), controllers.CreateQuiz)
	router.GET("/", middlewares.AdminCheck(), controllers.GetQuizzes)
	router.GET("/:quizId", controllers.GetQuiz) // add middleware for auth check
}
