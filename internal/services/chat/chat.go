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
	Chat           *chatModel.Chat
	chatRepository *chatRepository.ChatRepository
	queue          *map[string]*userModel.User
}

func NewChatService(users []*userModel.User, queue *map[string]*userModel.User) *ChatService {
	chat := newChat(users)

	return &ChatService{
		Chat:           chat,
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
	slice := *c.queue
	slice[user.Hash] = user

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
	for _, user := range c.Chat.Users {
		c.chatRepository.RegisterUserToStream(user.Hash)
		user.AddChat(c.Chat.Hash)
	}
}

func (c *ChatService) startHandlers() {
	for _, user := range c.Chat.Users {
		go c.readUserMessages(user)
		go c.sendMessages(user)
	}
}

func (c *ChatService) readUserMessages(user *userModel.User) {
	for {
		fmt.Println("READ MESSAGE FOR ", user.Hash, "chat: ", c.Chat.Hash)

		message, opened := user.GetFromInChat(c.Chat.Hash)
		if !opened {
			fmt.Printf("readUserMessages Channel closed for user %s\n", user.Hash)
			return
		}

		message.Payload.ChatHash = c.Chat.Hash
		c.chatRepository.AddMessage(*message)
	}
}

func (c *ChatService) sendMessages(user *userModel.User) {
	for {
		chanUser := user.GetChatState(c.Chat.Hash)

		select {
		case <-chanUser:
			fmt.Println("Received done signal, exiting goroutine")
			return
		default:
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

					user.PutToOutChannel(&message)

					if message.Category == "CHAT_EXIT" {
						user.CloseChat(message.Payload.ChatHash)
					}
				}
			}
		}
	}
}

func (c *ChatService) notifyChatStart() {
	now := time.Now()
	unixTimestamp := now.Unix()

	message := messageModel.NewMessage(
		"CHAT_START",
		c.Chat.Hash,
		unixTimestamp,
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
