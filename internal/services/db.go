package services

import (
	"fmt"

	"github.com/anuj-thakur-513/quizz/internal/config"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func ConnectDb() *mongo.Client {
	var connectionString string = config.GetEnv().MONGO_URL
	const DB_NAME = "quizz"
	clientOptions := options.Client().ApplyURI(connectionString)
	client, err := mongo.Connect(clientOptions)
	if err != nil {
		panic(err)
	}
	fmt.Println("MongoDB connection successful")
	return client
}
