package services

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"time"

	"github.com/anuj-thakur-513/quizz/internal/config"
	"github.com/go-redis/redis/v8"
)

type LeaderboardSetMember struct {
	UserId   string `json:"user_id"`
	Username string `json:"username"`
}

var redisClient *redis.Client

func ConnectRedis() *redis.Client {
	redisClient = redis.NewClient(&redis.Options{
		Addr:      config.GetEnv().REDIS.Address,
		Username:  config.GetEnv().REDIS.Username,
		Password:  config.GetEnv().REDIS.Password,
		TLSConfig: &tls.Config{InsecureSkipVerify: true},
	})
	return redisClient
}

func GetRedisClient() *redis.Client {
	return redisClient
}

func AddToZSet(key string, score float64, userId string, username string) {
	finalKey := "quizz:" + key

	set := redisClient.ZRange(context.Background(), finalKey, 0, -1)
	if len(set.Val()) == 0 {
		member := &LeaderboardSetMember{
			UserId:   userId,
			Username: username,
		}
		json, _ := json.Marshal(member)

		redisClient.ZAdd(context.Background(), finalKey, &redis.Z{
			Score:  score,
			Member: json,
		})
	} else {
		var members []*LeaderboardSetMember
		for _, memberString := range set.Val() {
			var member *LeaderboardSetMember
			if err := json.Unmarshal([]byte(memberString), &member); err != nil {
				continue
			}
			members = append(members, member)
		}
		pipeline := redisClient.Pipeline()
		var userFound bool = false
		for _, member := range members {
			if member.UserId == userId {
				userFound = true
				json, _ := json.Marshal(member)
				res := redisClient.ZScore(context.Background(), finalKey, string(json))
				sc, err := res.Result()
				if err != nil {
					continue
				}
				sc += score
				pipeline.ZAdd(context.Background(), finalKey, &redis.Z{
					Score:  sc,
					Member: json,
				})
			}
		}

		if !userFound {
			member := &LeaderboardSetMember{
				UserId:   userId,
				Username: username,
			}
			json, _ := json.Marshal(member)
			pipeline.ZAdd(context.Background(), finalKey, &redis.Z{
				Score:  score,
				Member: json,
			})
		}

		_, err := pipeline.Exec(context.Background())
		if err != nil {
			fmt.Println(err.Error())
		}
	}
	redisClient.Expire(context.Background(), key, 48*time.Hour)
}

func GetZSet(key string) []string {
	finalKey := "quizz:" + key
	res := redisClient.ZRange(context.Background(), finalKey, 0, -1)
	return res.Val()
}

func GetZScore(key string, member *LeaderboardSetMember) float64 {
	finalKey := "quizz:" + key
	m, err := json.Marshal(member)
	if err != nil {
		return -1
	}
	res := redisClient.ZScore(context.Background(), finalKey, string(m))
	return res.Val()
}

func SetCache(key string, value string) {
	redisClient.Set(context.Background(), key, value, 24*time.Hour)
}

func GetCache(key string) string {
	res := redisClient.Get(context.Background(), key)
	return res.Val()
}

func DeleteCache(key string) {
	redisClient.Del(context.Background(), key)
}
