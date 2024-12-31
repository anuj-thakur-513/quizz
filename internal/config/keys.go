package config

import "os"

type Keys struct {
	PORT       string
	MONGO_URL  string
	JWT_SECRET string
}

func GetEnv() *Keys {
	return &Keys{
		PORT:       getEnv("PORT", "8000"),
		MONGO_URL:  getEnv("MONGO_URL", "mongodb://localhost:27017"),
		JWT_SECRET: getEnv("JWT_SECRET", "secret"),
	}
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
