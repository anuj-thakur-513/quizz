package models

import (
	"time"

	"github.com/anuj-thakur-513/quizz/internal/services"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Solution struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	User      primitive.ObjectID `json:"user"`
	Quiz      primitive.ObjectID `json:"quiz"`
	Question  primitive.ObjectID `json:"question"`
	IsCorrect bool               `json:"is_correct" bson:"is_correct"`
	Score     int                `json:"score"`
	CreatedAt *time.Time         `json:"created_at,omitempty" bson:"created_at,omitempty"`
	UpdatedAt *time.Time         `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
}

func (s *Solution) PreSave() {
	now := time.Now()
	if s.CreatedAt == nil {
		s.CreatedAt = &now
	}
	s.UpdatedAt = &now
}

func (s *Solution) PostSave(username string) {
	key := s.Quiz.Hex() + "_leaderboard"
	services.AddToZSet(key, float64(s.Score), s.User.String(), username)
}

func GetSolutionsCollection() *mongo.Collection {
	return services.GetDatabase().Collection("solutions")
}
