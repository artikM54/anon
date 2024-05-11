package main

import (
	"anonymous_chat/internal/config"
	"anonymous_chat/internal/server"
	"anonymous_chat/internal/redis"
	"anonymous_chat/internal/handler_queue"
)

// for start "CONFIG_PATH="/home/cat/test_area/go/anonymous_chat" go run main.go"
func main() {
	config.MustLoad()
	redis.MustLoad()
	handler_queue.MustLoad()
	server.MustLoad()
}
