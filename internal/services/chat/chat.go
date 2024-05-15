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

	for _, user := range c.chat.Users {
		c.NotifyChatStart(user)
		c.chatRepository.RegisterUserToStream(user.Hash)
	}

	for _, user := range c.chat.Users {
		go c.write(user)
		go c.read(user)
	}
}

func (c *ChatService) write(user *userModel.User) {
	for {
		_, textMessage, err := user.Conn.ReadMessage()
		if err != nil {
			log.Println("Error reading message from user connection 1:", err)

			// TODO notify other users
			// c.closeUserConn(user2)
			return
		}

		var message messageModel.Message
		json.Unmarshal(textMessage, &message)

		switch message.Category {
		case "CHAT":
			message.Category = "CHAT"
			message.Payload.Timestamp = time.Now().Format("2006-01-02 15:04:05")
			message.Payload.UserHash = user.Hash
			message.Payload.ChatHash = c.chat.Hash

			c.chatRepository.AddMessage(message)
		case "FRONT:CHAT_EXIT":
			fmt.Println("FRONT:CHAT_EXIT HandleStreamMessages ", time.Now().Format("2006-01-02 15:04:05"))
			message.Category = "CHAT_EXIT"
			message.Payload.Timestamp = time.Now().Format("2006-01-02 15:04:05")
			message.Payload.UserHash = ""
			message.Payload.ChatHash = ""

			//TODO handler this case within the redis stream 
		case "FRONT:START_QUEUE":
			fmt.Println("FRONT:START_QUEUE HandleStreamMessages ", time.Now().Format("2006-01-02 15:04:05"))

			c.AddUserToQueue(user)

			return
		}
	}
}

func (c *ChatService) read(user *userModel.User) {
	for {
		streams := c.chatRepository.GetNewMessages(user.Hash)

		for _, stream := range streams {
			for _, message := range stream.Messages {
				messageData := message.Values["message"].(string)

				var message messageModel.Message
				json.Unmarshal([]byte(messageData), &message)

				fmt.Println("Received message data:", messageData)

				if user.Hash == message.Payload.UserHash {
					continue
				}

				c.SendMessage(user.Conn, &message)
			}
		}
	}
}

func (c *ChatService) NotifyChatStart(user *userModel.User) {
	message := messageModel.NewMessage(
		"CHAT_START",
		c.chat.Hash,
		time.Now().Format("2006-01-02 15:04:05"),
		"",
		"",
	)

	if err := c.SendMessage(user.Conn, message); err != nil {
		log.Println("Error sending message to client 2: ", err)
		return
	}
}

func (c *ChatService) SendMessage(conn *websocket.Conn, message *messageModel.Message) error {
	data, err := json.Marshal(message)
	if err != nil {
		log.Println("Error encoding message to JSON: ", err)
	}

	return conn.WriteMessage(websocket.TextMessage, data)
}

func (c *ChatService) closeUserConn(user *userModel.User) {
	message := messageModel.NewMessage(
		"CLOSE",
		string("Собеседнкик покинул чат"),
		time.Now().Format("2000-01-01 00:00:00"),
		"",
		"",
	)

	if err := c.SendMessage(user.Conn, message); err != nil {
		log.Println("Error sending message to client 2: ", err)
	}

	user.Conn.Close()
}
