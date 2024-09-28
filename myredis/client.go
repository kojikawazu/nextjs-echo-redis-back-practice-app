package myredis

import (
	"context"
	"log"
	"os"

	"github.com/redis/go-redis/v9"
)

var rdb *redis.Client

func InitRedis() error {
	log.Println("Initializing Redis client...")

	rdb = redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_HOST") + ":" + os.Getenv("REDIS_PORT"),
		Password: "", // Redisにパスワードが設定されている場合はここで指定
		DB:       0,  // defaultのDBを使用
	})

	// 接続確認
	log.Println("Pinging Redis server...")
	_, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		log.Printf("Failed to connect to Redis: %v", err)
		return err
	}

	log.Println("Connected to Redis successfully")
	return nil
}

// Redisクライアントを取得する関数
func GetRedisClient() *redis.Client {
	return rdb
}
