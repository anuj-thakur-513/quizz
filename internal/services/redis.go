package services

import (
	"crypto/tls"

	"github.com/anuj-thakur-513/quizz/internal/config"
	"github.com/go-redis/redis/v8"
)

func ConnectRedis() *redis.Client {
	r := redis.NewClient(&redis.Options{
		Addr:      config.GetEnv().REDIS.Address,
		Username:  config.GetEnv().REDIS.Username,
		Password:  config.GetEnv().REDIS.Password,
		TLSConfig: &tls.Config{InsecureSkipVerify: true},
	})
	return r
}
