package router

import (
	"github.com/gin-gonic/gin"
)

func V1Router(router *gin.RouterGroup) {
	user := router.Group("/user")
	UserRouter(user)

	quiz := router.Group("/quiz")
	QuizRouter(quiz)
}
