package controllers

import (
	"github.com/anuj-thakur-513/quizz/internal/models"
	"github.com/anuj-thakur-513/quizz/pkg/core"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func CreateQuiz(c *gin.Context) {
	var quizData *models.Quiz
	if err := c.BindJSON(&quizData); err != nil {
		c.JSON(400, core.NewAppError(400, "Invalid JSON body", err.Error()))
		return
	}
	quizData.PreSave()
	if err := validate.Struct(quizData); err != nil {
		c.JSON(400, core.NewAppError(400, "Invalid JSON body", err.Error()))
		return
	}

	quizzes := models.GetQuizzesCollection()

	if _, err := quizzes.InsertOne(ctx, quizData); err != nil {
		c.JSON(500, core.NewAppError(500, "Failed to create quiz", err.Error()))
		return
	}

	c.JSON(201, core.ApiResponse(200, "Create Quiz", map[string]interface{}{
		"category":       quizData.Category,
		"question_count": quizData.QuestionCount,
	}))
}

func GetQuizzes(c *gin.Context) {
	quizzes := models.GetQuizzesCollection()
	cursor, err := quizzes.Find(ctx, bson.M{}, options.Find().SetProjection(
		bson.M{
			"created_at": 0,
			"updated_at": 0,
		},
	))
	if err != nil {
		c.JSON(500, core.NewAppError(500, "Failed to get quizzes", err.Error()))
		return
	}
	defer cursor.Close(ctx)

	var quizzesData []models.Quiz
	if err := cursor.All(ctx, &quizzesData); err != nil {
		c.JSON(500, core.NewAppError(500, "Failed to get quizzes", err.Error()))
		return
	}
	c.JSON(200, core.ApiResponse(200, "All Quizzes returned successfully", quizzesData))
}

func GetQuiz(c *gin.Context) {
	quizId := c.Param("quizId")
	if quizId == "" {
		c.JSON(400, core.NewAppError(400, "Invalid Request", "quizId is required"))
		return
	}

	quizzes := models.GetQuizzesCollection()
	var quiz *models.Quiz
	id, err := primitive.ObjectIDFromHex(quizId)
	if err != nil {
		c.JSON(400, core.NewAppError(400, "Invalid Request", "quizId is invalid"))
	}

	if err := quizzes.FindOne(ctx, bson.M{"_id": id}, options.FindOne().SetProjection(bson.M{"created_at": 0, "updated_at": 0})).Decode(&quiz); err != nil {
		c.JSON(500, core.NewAppError(500, "Failed to get quiz", err.Error()))
		return
	}

	c.JSON(200, core.ApiResponse(200, "Quiz returned successfully", quiz))
}
