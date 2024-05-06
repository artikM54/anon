package main

import (
	"anonymous_chat/internal/config"
	"anonymous_chat/internal/server"

	"anonymous_chat/internal/redis_client"

	"context"
)

// for start "CONFIG_PATH="/home/cat/test_area/go/anonymous_chat" go run main.go"
func main() {
	config.MustLoad()
	server.MustLoad()
	redis_client.MustLoad()

	err := redis_client.RedisClient.Set(context.Background(), "key", "value", 0).Err()
	if err != nil {
		panic(err)
	}
}
