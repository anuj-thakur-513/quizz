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

type TestSubmission struct {
	ID          primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	User        primitive.ObjectID `json:"user" bson:"user"`
	Quiz        primitive.ObjectID `json:"quiz" bson:"quiz"`
	IsSubmitted *bool              `json:"is_submitted" bson:"is_submitted"`
	Score       int                `json:"score" bson:"score"`
	CreatedAt   *time.Time         `json:"created_at,omitempty" bson:"created_at,omitempty"`
	UpdatedAt   *time.Time         `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
}

func (ts *TestSubmission) PreSave() {
	now := time.Now()
	if ts.CreatedAt == nil {
		ts.CreatedAt = &now
	}
	ts.UpdatedAt = &now
	if ts.IsSubmitted == nil {
		t := true
		ts.IsSubmitted = &t
	}
}

func createUserQuizIndex(test_submissions *mongo.Collection) {
	indexModel := mongo.IndexModel{
		Keys:    bson.D{{Key: "user", Value: 1}, {Key: "quiz", Value: 1}},
		Options: options.Index().SetUnique(true),
	}

	_, err := test_submissions.Indexes().CreateOne(context.Background(), indexModel)
	if err != nil {
		log.Fatalf("Failed to create unique index for email: %v", err)
	}
}

func GetTestSubmissionCollection() *mongo.Collection {
	collection := services.GetDatabase().Collection("test_submissions")
	createUserQuizIndex(collection)
	return collection
}
