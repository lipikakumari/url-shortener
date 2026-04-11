package main

import (
	"context"
	"fmt"
	"os"

	"github.com/redis/go-redis/v9"
)

func initRedis() *redis.Client {
	addr := os.Getenv("REDIS_URL")

	if addr == "" {
		addr = "localhost:6379"
	}

	client := redis.NewClient(&redis.Options{
		Addr: addr,
	})

	ctx := context.Background()
	_, err := client.Ping(ctx).Result()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to Redis : %v\n", err)
		os.Exit(1)
	}

	fmt.Println("connected to redis")
	return client
}
