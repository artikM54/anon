package handler_queue

import (
	chatModel "anonymous_chat/internal/models/chat"
	messageModel "anonymous_chat/internal/models/message"
	queueModel "anonymous_chat/internal/models/queue"
	chatService "anonymous_chat/internal/services/chat"
	"fmt"
)

var Queue queueModel.UserQueue = queueModel.NewUserQueue()
var Chats map[string]*chatModel.Chat = make(map[string]*chatModel.Chat)

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
	addChat(c.Chat)

	c.Start()
}

func addChat(chat *chatModel.Chat) {
	Chats[chat.Hash] = chat
}

func ExitChat(hash string) bool {
	_, found := Chats[hash]

	return found
}

func PutIntoChat(message *messageModel.Message) {
	Chats[message.Payload.ChatHash].Channel <- message
}
