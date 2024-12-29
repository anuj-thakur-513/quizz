package config

import "os"

type Keys struct {
	PORT string
}

func GetEnv() *Keys {
	return &Keys{
		PORT: getEnv("PORT", "8000"),
	}
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
