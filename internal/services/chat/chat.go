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
}

func NewChatService(users []*userModel.User) *ChatService {
	chat := newChat(users)

	return &ChatService{
		chat:           chat,
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

	for _, user := range c.chat.Users {
		c.NotifyChatStart(user)
	}

	go c.HandleStreamMessages(c.chat.Users[0], c.chat.Users[1])
	go c.HandleStreamMessages(c.chat.Users[1], c.chat.Users[0])
}

func (c *ChatService) HandleStreamMessages(user1 *userModel.User, user2 *userModel.User) {
	for {
		_, textMessage, err := user1.Conn.ReadMessage()
		if err != nil {
			log.Println("Error reading message from client 1:", err)
			c.closeUserConn(user2)
			return
		}

		var message messageModel.Message
		json.Unmarshal(textMessage, &message)

		message.Category = "CHAT"
		message.Payload.Timestamp = time.Now().Format("2000-01-01 00:00:00")
		message.Payload.UserHash = user1.Hash
		message.Payload.ChatHash = c.chat.Hash

		if err := c.SendMessage(user2.Conn, &message); err != nil {
			log.Println("Error sending message to client 2: ", err)
			return
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
