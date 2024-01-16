package db

import (
	"encoding/json"
	"os"

	"github.com/go-redis/redis"
	"github.com/matx64/go_xat/models"
)

func StartRedis() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_HOST") + ":" + os.Getenv("REDIS_PORT"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0,
	})
}

func StoreMessage(msg models.Message, rdb *redis.Client) {
	json, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}

	if err := rdb.RPush("room:"+msg.RoomId+":messages", json).Err(); err != nil {
		panic(err)
	}
}
