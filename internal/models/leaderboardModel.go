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

type Leaderboard struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Quiz      primitive.ObjectID `json:"quiz" bson:"quiz" validate:"required"`
	Users     []LeaderboardUser  `json:"users" bson:"users"`
	CreatedAt *time.Time         `json:"created_at,omitempty" bson:"created_at,omitempty"`
	UpdatedAt *time.Time         `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
}

type LeaderboardUser struct {
	User  primitive.ObjectID `json:"user"`
	Score int                `json:"score"`
}

func (l *Leaderboard) PreSave() {
	if l.Users == nil {
		l.Users = []LeaderboardUser{}
	}

	now := time.Now()
	if l.CreatedAt == nil {
		l.CreatedAt = &now
	}
	l.UpdatedAt = &now
}

func createQuizIndex(leaderboards *mongo.Collection) {
	indexModel := mongo.IndexModel{
		Keys:    bson.D{{Key: "quiz", Value: 1}},
		Options: options.Index().SetUnique(true),
	}

	_, err := leaderboards.Indexes().CreateOne(context.Background(), indexModel)
	if err != nil {
		log.Fatalf("Failed to create unique index for email: %v", err)
	}
}

func GetLeaderboardCollection() *mongo.Collection {
	collection := services.GetDatabase().Collection("leaderboards")
	createQuizIndex(collection)
	return collection
}
