package main

import (
	"log"

	"github.com/anuj-thakur-513/quizz/internal/config"
	"github.com/anuj-thakur-513/quizz/internal/router"
)

func main() {
	router := router.SetupRouter()
	port := config.GetEnv().PORT
	err := router.Run(":" + port)
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
