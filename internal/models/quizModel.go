package models

import (
	"context"
	"log"
	"time"

	"github.com/anuj-thakur-513/quizz/internal/services"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Quiz struct {
	ID              primitive.ObjectID   `json:"_id,omitempty" bson:"_id,omitempty"`
	Category        string               `json:"category" validate:"required"`
	QuestionCount   int                  `json:"question_count" bson:"question_count"`
	Questions       []primitive.ObjectID `json:"questions" bson:"questions"`
	IsLive          *bool                `json:"is_live" bson:"is_live"`
	LiveTime        *time.Time           `json:"live_time" bson:"live_time"`
	DurationSeconds int                  `json:"duration_seconds" bson:"duration_seconds"`
	CreatedAt       *time.Time           `json:"created_at,omitempty" bson:"created_at,omitempty"`
	UpdatedAt       *time.Time           `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
}

func (q *Quiz) PreSave() {
	if q.QuestionCount == 0 {
		q.QuestionCount = 10
	}
	if q.IsLive == nil {
		f := false
		q.IsLive = &f
	}
	if q.DurationSeconds == 0 {
		q.DurationSeconds = q.QuestionCount * 30 // 30 seconds for every question
	}
	now := time.Now()
	if q.CreatedAt == nil {
		q.CreatedAt = &now
	}
	q.UpdatedAt = &now
}

func createCategoryIndex(quizzes *mongo.Collection) {
	indexModel := mongo.IndexModel{
		Keys:    bson.D{{Key: "category", Value: 1}},
		Options: options.Index().SetUnique(true),
	}

	_, err := quizzes.Indexes().CreateOne(context.Background(), indexModel)
	if err != nil {
		log.Fatalf("Failed to create unique index for email: %v", err)
	}
}

func GetQuizzesCollection() *mongo.Collection {
	quizzes := services.GetDatabase().Collection("quizzes")
	createCategoryIndex(quizzes)
	return quizzes
}
