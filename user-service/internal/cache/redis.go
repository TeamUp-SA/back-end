package cache

import (
	"context"
	"log"
	"os"

	"github.com/redis/go-redis/v9"
)

var (
	Rdb *redis.Client
	Ctx = context.Background()
)

func InitRedis() {
	addr := os.Getenv("REDIS_ADDR")
	if addr == "" {
		addr = "redis:6379"
	}

	Rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: "",
		DB:       0,
	})

	if err := Rdb.Ping(Ctx).Err(); err != nil {
		log.Fatalf("failed to connect redis: %v", err)
	}

	log.Printf("connected to redis at %s", addr)
}
