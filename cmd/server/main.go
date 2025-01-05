package main

import (
	"context"
	"fmt"
	"log"

	"github.com/anuj-thakur-513/quizz/internal/config"
	"github.com/anuj-thakur-513/quizz/internal/router"
	"github.com/anuj-thakur-513/quizz/internal/services"
)

func main() {
	db := services.ConnectDb()
	redis := services.ConnectRedis()
	fmt.Println("REDIS", redis.Ping(context.TODO()))
	fmt.Println(db.Ping(context.TODO(), nil))
	router := router.SetupRouter()
	port := config.GetEnv().PORT
	err := router.Run(":" + port)
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
