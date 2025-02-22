package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file", err)
	}
}

type Keys struct {
	PORT       string
	MONGO_URL  string
	JWT_SECRET string
	QUIZ_API   QuizAPIKeys
	REDIS      RedisKeys
}
type QuizAPIKeys struct {
	URL   string
	TOKEN string
}

type RedisKeys struct {
	Address  string
	Username string
	Password string
}

func GetEnv() *Keys {
	return &Keys{
		PORT:       getEnv("PORT", "8000"),
		MONGO_URL:  getEnv("MONGO_URL", "mongodb://localhost:27017"),
		JWT_SECRET: getEnv("JWT_SECRET", "secret"),
		QUIZ_API: QuizAPIKeys{
			URL:   getEnv("QUIZ_API_URL", "https://quizapi.io/api/v1/questions"),
			TOKEN: getEnv("QUIZ_API_TOKEN", "secret"),
		},
		REDIS: RedisKeys{
			Address:  getEnv("REDIS_ADDRESS", "localhost:6379"),
			Username: getEnv("REDIS_USERNAME", "default"),
			Password: getEnv("REDIS_PASSWORD", ""),
		},
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
