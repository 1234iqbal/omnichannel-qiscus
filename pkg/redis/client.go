package redis

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"
)

func NewClient(redisURL string) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     redisURL,
		Password: "",
		DB:       0,
	})

	// Test connection
	ctx := context.Background()
	_, err := client.Ping(ctx).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	fmt.Println("Connected to Redis successfully")
	return client, nil
}
