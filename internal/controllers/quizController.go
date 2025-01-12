package controllers

import (
	"context"
	"net/http"
	"sync"
	"time"

	"github.com/anuj-thakur-513/quizz/internal/models"
	"github.com/anuj-thakur-513/quizz/internal/repository"
	"github.com/anuj-thakur-513/quizz/internal/services"
	"github.com/anuj-thakur-513/quizz/pkg/core"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetQuizzes(c *gin.Context) {
	ctx := context.Background()
	quizzes := models.GetQuizzesCollection()
	currTime := time.Now()
	cursor, err := quizzes.Find(ctx, bson.M{"live_time": bson.M{"$gt": currTime}}, options.Find().SetProjection(
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
	ctx := context.Background()
	quizId := c.Param("quizId")
	if quizId == "" {
		c.JSON(400, core.NewAppError(400, "Invalid Request", "quizId is required"))
		return
	}

	quizzes := models.GetQuizzesCollection()
	var quiz *models.Quiz
	id, err := primitive.ObjectIDFromHex(quizId)
	if err != nil {
		c.JSON(400, core.NewAppError(400, "Invalid Request", err.Error()))
		return
	}

	if err := quizzes.FindOne(ctx, bson.M{"_id": id}, options.FindOne().SetProjection(
		bson.M{"created_at": 0, "updated_at": 0},
	)).Decode(&quiz); err != nil {
		c.JSON(500, core.NewAppError(500, "Failed to get quiz", err.Error()))
		return
	}

	if quiz.LiveTime.Before(time.Now()) {
		c.JSON(400, core.NewAppError(400, "Invalid Request", "quiz has expired"))
		return
	}

	c.JSON(200, core.ApiResponse(200, "Quiz returned successfully", quiz))
}

func SubmitSolution(c *gin.Context) {
	ctx := context.Background()
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
	if !value.(bool) {
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
	difficulty := question.Difficulty
	score := 0
	for _, option := range options {
		if option.Option == body["solution"] {
			isCorrect = option.IsCorrect
		}
	}
	if isCorrect {
		if difficulty == "Hard" {
			score = 3
		} else if difficulty == "Medium" {
			score = 2
		} else {
			score = 1
		}
	}

	if err != nil {
		c.JSON(400, core.NewAppError(400, "Invalid Request", "quizId is invalid"))
		return
	}

	solution := &models.Solution{}
	solution.User = userId
	solution.Question = qId
	solution.Score = score
	solution.Quiz = quizObjectId
	solution.IsCorrect = isCorrect
	solution.PreSave()

	if _, err := solutions.InsertOne(ctx, solution); err != nil {
		c.JSON(500, core.NewAppError(500, "Failed to submit answer", err.Error()))
		return
	}
	solution.PostSave(user.Name)

	c.JSON(200, core.ApiResponse(200, "Answer submitted successfully", body))
}

func SubmitQuiz(c *gin.Context) {
	ctx := context.Background()
	quizId := c.Param("quizId")
	qId, err := primitive.ObjectIDFromHex(quizId)
	if err != nil {
		c.JSON(400, core.NewAppError(400, "Invalid Request", "quizId is invalid"))
		return
	}

	u, exists := c.Get("user")
	if !exists {
		c.JSON(400, core.NewAppError(400, "Invalid Request", "user not found"))
		return
	}
	user := u.(*models.User)
	userId := user.ID

	TestSubmissions := models.GetTestSubmissionCollection()
	quizzes := models.GetQuizzesCollection()
	solutions := models.GetSolutionsCollection()

	var testSubmission *models.TestSubmission
	if err := TestSubmissions.FindOne(ctx, bson.M{"user": userId, "quiz": qId}).Decode(&testSubmission); err == nil {
		c.JSON(409, core.NewAppError(500, "Failed to submit quiz", "quiz already submitted"))
		return
	}

	var quiz *models.Quiz
	if err := quizzes.FindOne(ctx, bson.M{"_id": qId}).Decode(&quiz); err != nil {
		c.JSON(500, core.NewAppError(500, "Failed to get quiz", err.Error()))
		return
	}

	questions := quiz.Questions
	// solutions -> userId, qId, questionId
	solArr := []primitive.M{}
	cursor, err := solutions.Find(ctx, bson.M{"user": userId, "quiz": qId, "question": bson.M{"$in": questions}})
	if err != nil {
		c.JSON(500, core.NewAppError(500, "Failed to get solutions", err.Error()))
		return
	}
	defer cursor.Close(ctx)

	if err := cursor.All(ctx, &solArr); err != nil {
		c.JSON(500, core.NewAppError(500, "Failed to get solutions", err.Error()))
		return
	}

	finalScore := 0
	for _, solution := range solArr {
		if score32, ok := solution["score"].(int32); ok {
			finalScore += int(score32)
		}
	}

	testSubmission = &models.TestSubmission{}
	testSubmission.User = userId
	testSubmission.Quiz = qId
	testSubmission.Score = finalScore
	testSubmission.PreSave()

	if _, err := TestSubmissions.InsertOne(ctx, testSubmission); err != nil {
		c.JSON(500, core.NewAppError(500, "Failed to submit quiz", err.Error()))
		return
	}
	testSubmission.PostSave()

	c.JSON(201, core.ApiResponse(200, "Quiz submitted successfully", map[string]interface{}{
		"user":         userId,
		"quiz":         qId,
		"is_submitted": true,
		"score":        finalScore,
	}))
}

func StartQuiz(c *gin.Context) {
	wg := &sync.WaitGroup{}
	mu := &sync.RWMutex{}

	quizId := c.Param("quizId")
	qId, err := primitive.ObjectIDFromHex(quizId)
	if err != nil {
		c.JSON(400, core.NewAppError(400, "Invalid Request", "quizId is invalid"))
		return
	}

	data, err := repository.GetQuizWithQuestionDetails(qId)
	if err != nil {
		c.JSON(500, core.NewAppError(500, "Failed to get quiz", err.Error()))
		return
	}

	quiz := (*data)[0]

	quizStartTime := quiz["live_time"].(primitive.DateTime).Time()
	quizDuration := quiz["duration_seconds"].(int32)
	quizEndTime := quizStartTime.Add(time.Duration(quizDuration) * time.Second)
	questionCount := quiz["question_count"].(int32)
	timePerQuestion := quizDuration / questionCount
	questions := quiz["questions"].(primitive.A)

	if time.Now().Truncate(time.Second).Equal(quizStartTime) || time.Now().Truncate(time.Second).After(quizStartTime) {
		c.JSON(400, core.NewAppError(400, "Invalid Request", "quiz has already started"))
		return
	}

	u, exists := c.Get("user")
	if !exists {
		c.JSON(400, core.NewAppError(400, "Invalid Request", "user not found"))
		return
	}
	user := u.(*models.User)
	// upgrade connection to WS
	conn, err := services.UpgradeWsConnection(c.Writer, c.Request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to establish WebSocket connection"})
		return
	}
	defer func() {
		services.RemoveConnection(user.ID.Hex())
		conn.Close()
	}()

	counter := 0
	qIndex := 0

	// Add the connection to the activeConnections map
	services.AddConnection(user.ID.Hex(), conn)
	for {
		currTime := time.Now().Truncate(time.Second)
		// quiz has ended (take a buffer of 5 minutes for this event so that the ended event is sent properly)
		if currTime.Equal(quizEndTime.Add(5*time.Minute)) || currTime.After(quizEndTime.Add(5*time.Minute)) {
			break
		}
		// send questions when quiz has started
		if currTime.Equal(quizStartTime) || currTime.After(quizStartTime) {
			if counter%int(timePerQuestion) == 0 {
				if qIndex < int(questionCount) {
					question := questions[qIndex].(primitive.M)
					wg.Add(1)
					go services.SendQuestion(conn, question, wg, mu)
					if qIndex > 0 {
						wg.Add(1)
						go services.SendLeaderboard(conn, quizId+"_leaderboard", wg, mu)
					}
					wg.Wait()
					qIndex += 1
				} else {
					wg.Add(2)
					go services.SendLeaderboard(conn, quizId+"_leaderboard", wg, mu)
					go services.SendQuizEnded(conn, wg, mu)
					wg.Wait()
					break
				}
			}
			counter += 1
		}
		// quiz has not yet started
		if currTime.Before(quizStartTime.Add(-120 * time.Second)) {
			break
		}
		time.Sleep(1 * time.Second)
	}
}
