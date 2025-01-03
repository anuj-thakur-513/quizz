package router

import (
	"github.com/anuj-thakur-513/quizz/internal/controllers"
	"github.com/anuj-thakur-513/quizz/internal/middlewares"
	"github.com/gin-gonic/gin"
)

func QuizRouter(router *gin.RouterGroup) {
	router.GET("/", middlewares.AuthCheck(), controllers.GetQuizzes)
	router.GET("/:quizId", middlewares.AuthCheck(), controllers.GetQuiz)

	router.POST("/:quizId/:questionId", middlewares.AuthCheck(), middlewares.QuestionInQuiz(), controllers.SubmitSolution)
}
