package main

import (
	"fmt"

	"try-on/internal/pkg/config"

	"github.com/go-redis/redis"
)

func main() {
	cfg := config.Redis{Host: "localhost", Port: 6379}
	redisClient := redis.NewClient(&redis.Options{
		Addr: cfg.DSN(),
	})

	res, err := redisClient.SPop("123").Result()
	fmt.Println(res, err)
}
