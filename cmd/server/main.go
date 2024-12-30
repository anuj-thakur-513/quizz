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
	client := services.ConnectDb()
	fmt.Println(client.Ping(context.TODO(), nil))
	router := router.SetupRouter()
	port := config.GetEnv().PORT
	err := router.Run(":" + port)
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
