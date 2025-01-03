package models

import (
	"time"

	"github.com/anuj-thakur-513/quizz/internal/services"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Question struct {
	ID                primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	QuestionText      string             `json:"question_text" bson:"question_text"`
	IsMultipleCorrect bool               `json:"is_multiple_correct" bson:"is_multiple_correct"`
	Options           []Options          `json:"options"`
	SolutionText      string             `json:"solution_text" bson:"solution_text"`
	Difficulty        string             `json:"difficulty"` // easy[1], medium[2], hard[3]
	Category          string             `json:"category"`
	CreatedAt         *time.Time         `json:"created_at,omitempty" bson:"created_at,omitempty"`
	UpdatedAt         *time.Time         `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
}

type Options struct {
	Option    string `json:"option"`
	IsCorrect bool   `json:"is_correct" bson:"is_correct"`
}

func (q *Question) PreSave() {
	now := time.Now()
	if q.CreatedAt == nil {
		q.CreatedAt = &now
	}
	q.UpdatedAt = &now
}

func GetQuestionsCollection() *mongo.Collection {
	return services.GetDatabase().Collection("questions")
}
