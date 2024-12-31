package services

import (
	"fmt"

	"github.com/anuj-thakur-513/quizz/internal/config"
	"github.com/anuj-thakur-513/quizz/pkg/utils"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

var client *mongo.Client

func ConnectDb() *mongo.Client {
	var connectionString string = config.GetEnv().MONGO_URL
	clientOptions := options.Client().ApplyURI(connectionString)
	dbClient, err := mongo.Connect(clientOptions)
	if err != nil {
		panic(err)
	}
	client = dbClient
	fmt.Println("MongoDB connection successful")
	return client
}

func GetDatabase() *mongo.Database {
	return client.Database(utils.DB_NAME)
}
