package chat

import (
	chatModel "anonymous_chat/internal/models/chat"
	messageModel "anonymous_chat/internal/models/message"
	userModel "anonymous_chat/internal/models/user"
	chatRepository "anonymous_chat/internal/repositories/chat"
	hashUtil "anonymous_chat/internal/utils/hash"
	"encoding/json"
	"fmt"

	"log"
	"time"

	"github.com/gorilla/websocket"
)

type ChatService struct {
	chat           *chatModel.Chat
	chatRepository *chatRepository.ChatRepository
	queue          *[]*userModel.User
}

func NewChatService(users []*userModel.User, queue *[]*userModel.User) *ChatService {
	chat := newChat(users)

	return &ChatService{
		chat:           chat,
		chatRepository: chatRepository.NewChatRepository(chat.Hash),
		queue:          queue,
	}
}

func newChat(users []*userModel.User) *chatModel.Chat {
	return &chatModel.Chat{
		Hash:  hashUtil.CreateUniqueModelHash(chatModel.RedisList),
		Users: users,
	}
}

func (c *ChatService) AddUserToQueue(user *userModel.User) {
	slice := append(*c.queue, user)
	*c.queue = slice
}

func (c *ChatService) Start() {
	fmt.Println("chat is starting")
	defer c.closeChat()

	c.registerUsers()
	c.notifyChatStart()
	c.startHandlers()
}

func (c *ChatService) registerUsers() {
	for _, user := range c.chat.Users {
		c.chatRepository.RegisterUserToStream(user.Hash)
	}
}

func (c *ChatService) startHandlers() {
	for _, user := range c.chat.Users {
		go c.readUserMessages(user)
		go c.sendMessages(user)
	}
}

func (c *ChatService) readUserMessages(user *userModel.User) {
	for {
		fmt.Println("READ MESSAGE FOR ", user.Hash)
		fmt.Println("READ MESSAGE within chat ", c.chat.Hash)

		select {
		case message, ok := <-user.ChannelChat:
			if !ok {
				fmt.Printf("readUserMessages Channel closed for user %s\n", user.Hash)
				return
			}

			message = c.handleMessage(user, message)
			c.chatRepository.AddMessage(*message)
		case <-time.After(5 * time.Second):
			fmt.Printf("Timeout reading message for user %s\n", user.Hash)
		}
	}
}

func (c *ChatService) handleMessage(user *userModel.User, message *messageModel.Message) *messageModel.Message {
	switch message.Category {
	case "CHAT":
		message.Payload.UserHash = user.Hash
		message.Payload.ChatHash = c.chat.Hash
		message.Payload.Timestamp = time.Now().Format("2006-01-02 15:04:05")

	case "FRONT:CHAT_EXIT":
		message.Category = "CHAT_EXIT"
		message.Payload.UserHash = user.Hash
		message.Payload.ChatHash = c.chat.Hash
		message.Payload.Timestamp = time.Now().Format("2006-01-02 15:04:05")
	}

	return message
}

func (c *ChatService) sendMessages(user *userModel.User) {
	for {
		streams := c.chatRepository.GetNewMessages(user.Hash)
		fmt.Println("SEND MESSAGE FOR ", user.Hash)
		fmt.Println("SEND MESSAGE within chat ", c.chat.Hash)

		for _, stream := range streams {
			for _, message := range stream.Messages {
				messageData := message.Values["message"].(string)

				var message messageModel.Message
				json.Unmarshal([]byte(messageData), &message)

				fmt.Println("sendMessages user: ", user.Hash)
				fmt.Println("Received message data:", messageData)

				if user.Hash == message.Payload.UserHash {
					continue
				}

				c.SendMessage(user.Conn, &message)

				if message.Category == "CHAT_EXIT" {
					close(user.ChannelChat)
				}
			}
		}
	}
}

func (c *ChatService) notifyChatStart() {
	message := messageModel.NewMessage(
		"CHAT_START",
		c.chat.Hash,
		time.Now().Format("2006-01-02 15:04:05"),
		"",
		"",
	)

	c.chatRepository.AddMessage(*message)
}

func (c *ChatService) SendMessage(conn *websocket.Conn, message *messageModel.Message) error {
	data, err := json.Marshal(message)
	if err != nil {
		log.Println("Error encoding message to JSON: ", err)
	}

	return conn.WriteMessage(websocket.TextMessage, data)
}

func (c *ChatService) closeChat() {
	fmt.Println("Close chat")
}
