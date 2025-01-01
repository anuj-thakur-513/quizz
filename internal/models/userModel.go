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
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name      string             `json:"name" bson:"name" validate:"required,min=2,max=100"`
	Email     string             `json:"email" bson:"email" validate:"email,required"`
	Role      Role               `json:"role" bson:"role"`
	Password  string             `json:"password" bson:"password" validate:"required,min=6,max=100"`
	CreatedAt *time.Time         `json:"created_at" bson:"created_at"`
	UpdatedAt *time.Time         `json:"updated_at" bson:"updated_at"`
}

type Role string

const (
	Admin   Role = "admin"
	Default Role = "default"
)

func (u *User) PreSave() {
	if u.Password != "" {
		hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), 10)
		if err != nil {
			panic(err)
		}
		u.Password = string(hash)
	}

	now := time.Now()
	if u.CreatedAt == nil {
		u.CreatedAt = &now
	}
	u.UpdatedAt = &now

	if u.Role == "" {
		u.Role = Default
	}
}

func (u *User) ComparePassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}

func createUniqueEmailIndex(users *mongo.Collection) {
	indexModel := mongo.IndexModel{
		Keys:    bson.D{{Key: "email", Value: 1}},
		Options: options.Index().SetUnique(true),
	}

	_, err := users.Indexes().CreateOne(context.Background(), indexModel)
	if err != nil {
		log.Fatalf("Failed to create unique index for email: %v", err)
	}
}

func GetUsersCollection() *mongo.Collection {
	users := services.GetDatabase().Collection("users")
	createUniqueEmailIndex(users)
	return users
}
