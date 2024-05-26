package chat

import (
	chatModel "anonymous_chat/internal/models/chat"
	messageModel "anonymous_chat/internal/models/message"
	userModel "anonymous_chat/internal/models/user"
	chatRepository "anonymous_chat/internal/repositories/chat"
	hashUtil "anonymous_chat/internal/utils/hash"
	"encoding/json"
	"fmt"
	"time"
)

type ChatService struct {
	Chat           *chatModel.Chat
	chatRepository *chatRepository.ChatRepository
}

func NewChatService(users []*userModel.User) *ChatService {
	chat := newChat(users)

	return &ChatService{
		Chat:           chat,
		chatRepository: chatRepository.NewChatRepository(chat.Hash),
	}
}

func newChat(users []*userModel.User) *chatModel.Chat {
	return &chatModel.Chat{
		Hash:  hashUtil.CreateUniqueModelHash(chatModel.RedisList),
		Users: users,
	}
}

func (c *ChatService) Start() {
	fmt.Println("chat is starting")
	c.registerUsers()
	c.notifyChatStart()
	c.startHandlers()
}

func (c *ChatService) registerUsers() {
	for _, user := range c.Chat.Users {
		c.chatRepository.RegisterUserToStream(user.Hash)
		user.AddChat(c.Chat.Hash)
	}
}

func (c *ChatService) startHandlers() {
	for _, user := range c.Chat.Users {
		go c.listeningUser(user)
		go c.listeningChat(user)
	}
}

func (c *ChatService) listeningUser(user *userModel.User) {
	for {
		fmt.Println("READ MESSAGE FOR ", user.Hash, "chat: ", c.Chat.Hash)

		message, opened := user.GetFromChat(c.Chat.Hash)
		if !opened {
			fmt.Printf("readUserMessages Channel closed for user %s\n", user.Hash)
			return
		}

		message.Payload.ChatHash = c.Chat.Hash
		c.chatRepository.AddMessage(*message)
	}
}

func (c *ChatService) listeningChat(user *userModel.User) {
	for {
		fmt.Println("listeningChat FOR ", user.Hash, "chat: ", c.Chat.Hash)

		chatState, founded := user.GetChatState(c.Chat.Hash)
		if !founded {
			return
		}

		select {
		case <-chatState:
			fmt.Println("Received done signal, exiting goroutine")
			return
		default:
			c.tryLoadMessages(user)
			time.Sleep(10 * time.Millisecond)
		}
	}
}

func (c *ChatService) tryLoadMessages(user *userModel.User) {
	streams := c.chatRepository.GetNewMessages(user.Hash)
	fmt.Println("SEND MESSAGE FOR ", user.Hash, " CHAT ", c.Chat.Hash)

	for _, stream := range streams {
		for _, message := range stream.Messages {

			messageData := message.Values["message"].(string)

			var message messageModel.Message
			json.Unmarshal([]byte(messageData), &message)

			fmt.Println("Read message from redis for: ", user.Hash, "message data : ", messageData)

			if user.Hash == message.Payload.UserHash {
				fmt.Println("CONTINUE for: ", user.Hash)
				continue
			}

			user.PutIntoChannel(&message)

			if message.Category == messageModel.ExitCategory {
				user.CloseChat(message.Payload.ChatHash)
			}
		}
	}

}

func (c *ChatService) notifyChatStart() {
	message := messageModel.NewMessage(
		messageModel.StartCategory,
		c.Chat.Hash,
		"",
		c.Chat.Hash,
	)

	c.chatRepository.AddMessage(*message)
}
