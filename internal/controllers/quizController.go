package controllers

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"github.com/anuj-thakur-513/quizz/internal/config"
	"github.com/anuj-thakur-513/quizz/internal/models"
	"github.com/anuj-thakur-513/quizz/pkg/core"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func extractOptions(answers map[string]interface{}, correctAnswers map[string]interface{}) []models.Options {
	var options []models.Options
	for key, value := range answers {
		if value != nil {
			isCorrect := false
			if correctAnswers[key+"_correct"] != nil && correctAnswers[key+"_correct"].(string) == "true" {
				isCorrect = true
			}
			options = append(options, models.Options{
				Option:    value.(string),
				IsCorrect: isCorrect,
			})
		}
	}
	return options
}

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

	// API call to get all the questions and then add them to questions collection and quiz as reference
	apiUrl := config.GetEnv().QUIZ_API.URL + "?apiKey=" + config.GetEnv().QUIZ_API.TOKEN +
		"&limit=" + strconv.Itoa(quizData.QuestionCount) + "&category=" + quizData.Category

	res, err := http.Get(apiUrl)
	if err != nil {
		c.JSON(500, core.NewAppError(500, "Failed to get questions", err.Error()))
		return
	}

	if res.StatusCode != 200 {
		c.JSON(res.StatusCode, core.NewAppError(res.StatusCode, "Failed to get questions", "Failed to get questions"))
		return
	}
	defer res.Body.Close()

	responseBytes, err := io.ReadAll(res.Body)
	if err != nil {
		c.JSON(500, core.NewAppError(500, "Failed to get questions", err.Error()))
	}

	var jsonResponse []map[string]interface{}
	if err := json.Unmarshal(responseBytes, &jsonResponse); err != nil {
		c.JSON(500, core.NewAppError(500, "Failed to get questions", err.Error()))
		return
	}

	var questionData []interface{}
	for _, question := range jsonResponse {
		q := models.Question{
			QuestionText:      question["question"].(string),
			IsMultipleCorrect: question["multiple_correct_answers"].(string) == "true",
			Options: extractOptions(question["answers"].(map[string]interface{}),
				question["correct_answers"].(map[string]interface{})),
			SolutionText: question["explanation"].(string),
			Difficulty:   question["difficulty"].(string),
			Category:     question["category"].(string),
		}
		q.PreSave()
		questionData = append(questionData, q)
	}
	questions := models.GetQuestionsCollection()
	quizzes := models.GetQuizzesCollection()

	r, dbErr := questions.InsertMany(ctx, questionData)
	if dbErr != nil {
		c.JSON(500, core.NewAppError(500, "Failed to create questions", dbErr.Error()))
		return
	}

	for _, id := range r.InsertedIDs {
		quizData.Questions = append(quizData.Questions, id.(primitive.ObjectID))
	}

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

	if err := quizzes.FindOne(ctx, bson.M{"_id": id}, options.FindOne().SetProjection(
		bson.M{"created_at": 0, "updated_at": 0},
	)).Decode(&quiz); err != nil {
		c.JSON(500, core.NewAppError(500, "Failed to get quiz", err.Error()))
		return
	}

	c.JSON(200, core.ApiResponse(200, "Quiz returned successfully", quiz))
}

func SubmitSolution(c *gin.Context) {
	quizId := c.Param("quizId")
	questionId := c.Param("questionId")
	if quizId == "" || questionId == "" {
		c.JSON(400, core.NewAppError(400, "Invalid Request", "quizId and questionId are required"))
		return
	}

	var body map[string]string
	if err := c.BindJSON(&body); err != nil {
		c.JSON(400, core.NewAppError(400, "Invalid JSON body", err.Error()))
		return
	}

	// var quiz *models.Quiz
	value, _ := c.Get("isQuizLive")
	if value == nil {
		c.JSON(400, core.NewAppError(400, "Invalid Request", "quiz live status is not set"))
		return
	}
	// Dereference **bool to get the actual bool value
	if isLive, ok := value.(**bool); ok && *isLive != nil && !**isLive {
		c.JSON(400, core.NewAppError(400, "Invalid Request", "quiz is not live"))
		return
	}

	// find question
	var question models.Question
	qId, err := primitive.ObjectIDFromHex(questionId)
	if err != nil {
		c.JSON(400, core.NewAppError(400, "Invalid Request", "questionId is invalid"))
		return
	}
	questions := models.GetQuestionsCollection()
	if err := questions.FindOne(ctx, bson.M{"_id": qId}).Decode(&question); err != nil {
		c.JSON(500, core.NewAppError(500, "Failed to get question", err.Error()))
		return
	}

	u, exists := c.Get("user")
	if !exists {
		c.JSON(400, core.NewAppError(400, "Invalid Request", "user not found"))
		return
	}
	user := u.(*models.User)
	userId := user.ID
	quizObjectId, err := primitive.ObjectIDFromHex(quizId)

	// checks if the user has already submitted the answer for given question
	var prevSolution *models.Solution
	solutions := models.GetSolutionsCollection()
	if err := solutions.FindOne(ctx, bson.M{"quiz": quizObjectId, "question": qId, "user": userId}).Decode(&prevSolution); err == nil {
		c.JSON(409, core.NewAppError(500, "Failed to submit answer", "answer already submitted"))
		return
	}

	var isCorrect bool
	options := question.Options
	for _, option := range options {
		if option.Option == body["solution"] {
			isCorrect = option.IsCorrect
		}
	}

	if err != nil {
		c.JSON(400, core.NewAppError(400, "Invalid Request", "quizId is invalid"))
		return
	}

	solution := &models.Solution{}
	solution.User = userId
	solution.Question = qId
	solution.Quiz = quizObjectId
	solution.IsCorrect = isCorrect
	solution.PreSave()

	if _, err := solutions.InsertOne(ctx, solution); err != nil {
		c.JSON(500, core.NewAppError(500, "Failed to submit answer", err.Error()))
		return
	}

	c.JSON(200, core.ApiResponse(200, "Answer submitted successfully", body))
}
