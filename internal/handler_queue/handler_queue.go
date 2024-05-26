package handler_queue

import (
	queueModel "anonymous_chat/internal/models/queue"
	chatService "anonymous_chat/internal/services/chat"
	"fmt"
)

var Queue queueModel.UserQueue = queueModel.NewUserQueue()

func MustLoad() {
	go start()
}

func start() {
	for {
		if Queue.GetCountUsersIntoQueue() >= 2 {
			bindUser()
		}
	}
}

func bindUser() {
	fmt.Println("There are two users")

	users := Queue.GetRandomUsers(2)

	c := chatService.NewChatService(users)

	go c.Start()
}
