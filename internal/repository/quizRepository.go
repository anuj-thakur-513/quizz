package repository

import (
	"context"

	"github.com/anuj-thakur-513/quizz/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetQuizWithQuestionDetails(quizId primitive.ObjectID) (*[]bson.M, error) {
	ctx := context.Background()
	quizzes := models.GetQuizzesCollection()
	pipeline := []bson.D{
		{
			{Key: "$match", Value: bson.M{"_id": quizId}},
		},
		{
			{Key: "$lookup", Value: bson.M{
				"from":         "questions",
				"localField":   "questions",
				"foreignField": "_id",
				"as":           "questions",
			}},
		},
		{
			{Key: "$project", Value: bson.M{
				"category":                      1,
				"question_count":                1,
				"live_time":                     1,
				"duration_seconds":              1,
				"questions.question_text":       1,
				"questions.is_multiple_correct": 1,
				"questions.options":             1,
				"questions._id":                 1,
			}},
		},
	}

	var data []primitive.M
	cursor, err := quizzes.Aggregate(ctx, mongo.Pipeline(pipeline))
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	if err := cursor.All(ctx, &data); err != nil {
		return nil, err
	}

	return &data, nil
}
