package services

import (
	"context"
	"fmt"

	"github.com/anuj-thakur-513/quizz/internal/config"
	"github.com/anuj-thakur-513/quizz/pkg/utils"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var mongoClient *mongo.Client

func ConnectDb() *mongo.Client {
	var connectionString string = config.GetEnv().MONGO_URL
	clientOptions := options.Client().ApplyURI(connectionString)
	dbClient, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		panic(err)
	}
	mongoClient = dbClient
	fmt.Println("MongoDB connection successful")
	return mongoClient
}

func GetDatabase() *mongo.Database {
	return mongoClient.Database(utils.DB_NAME)
}
