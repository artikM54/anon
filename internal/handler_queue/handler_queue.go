package handler_queue

import (
	"anonymous_chat/internal/services/user_queue"
)

func MustLoad() {
	userQueueService := user_queue.NewUserQueueService()
	userQueueService.Start()
}
