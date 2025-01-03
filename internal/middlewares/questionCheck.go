package middlewares

import (
	"context"

	"github.com/anuj-thakur-513/quizz/internal/models"
	"github.com/anuj-thakur-513/quizz/pkg/core"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func QuestionInQuiz() gin.HandlerFunc {
	return func(c *gin.Context) {
		quizId := c.Param("quizId")
		questionId := c.Param("questionId")
		if quizId == "" || questionId == "" {
			c.JSON(400, core.NewAppError(400, "Invalid Request", "quizId and questionId are required"))
			c.Abort()
			return
		}

		quizzes := models.GetQuizzesCollection()

		var quiz models.Quiz
		id, err := primitive.ObjectIDFromHex(quizId)
		if err != nil {
			c.JSON(400, core.NewAppError(400, "Invalid Request", "quizId is invalid"))
			c.Abort()
			return
		}
		qId, err := primitive.ObjectIDFromHex(questionId)
		if err != nil {
			c.JSON(400, core.NewAppError(400, "Invalid Request", "questionId is invalid"))
			c.Abort()
			return
		}

		if err := quizzes.FindOne(context.Background(), bson.M{"_id": id, "questions": qId}).Decode(&quiz); err != nil {
			c.JSON(400, core.NewAppError(400, "Invalid Request", "wrong question or quiz"))
			c.Abort()
			return
		}
		if quiz.ID == primitive.NilObjectID {
			c.JSON(400, core.NewAppError(400, "Invalid Request", "wrong question or quiz"))
			c.Abort()
			return
		}
		c.Set("isQuizLive", &quiz.IsLive)
		c.Next()
	}
}
