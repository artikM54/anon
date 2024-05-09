package chat

import (
	chatModel "anonymous_chat/internal/models/chat"
	messageModel "anonymous_chat/internal/models/message"
	userModel "anonymous_chat/internal/models/user"
	chatRepository "anonymous_chat/internal/repositories/chat"
	hashUtil "anonymous_chat/internal/utils/hash"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"math/rand"
)

var userQueue []*userModel.User

type ChatService struct {
	chat           *chatModel.Chat
	chatRepository *chatRepository.ChatRepository
}

func NewChatService() *ChatService {
	return &ChatService{}
}

func (c *ChatService) addUserToQueue(user *userModel.User) {
	userQueue = append(userQueue, user)
}

func (c *ChatService) AddUserAndTryUpChat(user *userModel.User) {
	c.addUserToQueue(user)
	c.notifyConnect(user)
	c.chat = c.NewChat()

	c.chatRepository = chatRepository.NewChatRepository(c.chat.Hash)

	if c.getCountUsersIntoQueue() >= 2 {
		c.bindUser()
	} else {
		c.notifyWait(user)
	}
}

func (c *ChatService) NewChat() *chatModel.Chat {
	return &chatModel.Chat{
		Hash: hashUtil.CreateUniqueModelHash(chatModel.RedisList),
	}
}

func (c *ChatService) getCountUsersIntoQueue() int {
	return len(userQueue)
}

func (c *ChatService) notifyWait(user *userModel.User) {
	message := messageModel.NewMessage(
		"WAIT",
		string("Нет свободных участников, пожалуйста, дождитесь свободного участника"),
	)

	if err := c.sendMessage(user, message); err != nil {
		log.Println("Error sending message to client 2: ", err)
		return
	}
}

func (c *ChatService) notifyConnect(user *userModel.User) {
	message := messageModel.NewMessage(
		"TOKEN",
		user.Hash,
	)

	if err := c.sendMessage(user, message); err != nil {
		log.Println("Error sending message to client 2: ", err)
		return
	}
}

func (c *ChatService) notifyChatStart(user *userModel.User) {
	message := messageModel.NewMessage(
		"CHAT_START",
		"Собеседник найден",
	)

	if err := c.sendMessage(user, message); err != nil {
		log.Println("Error sending message to client 2: ", err)
		return
	}
}

func (c *ChatService) bindUser() {
	fmt.Println("There are two users")

	user1, user2, err := c.chooseRandomPair()
	if err != nil {
		log.Println("Error choosing random pair:", err)
		return
	}

	userQueue = c.removeClientFromSlice(userQueue, user1)
	userQueue = c.removeClientFromSlice(userQueue, user2)

	c.notifyChatStart(user1)
	c.notifyChatStart(user2)

	go c.HandleStreamMessages(user1, user2)
	go c.HandleStreamMessages(user2, user1)

	c.chatRepository.AddParticipant(user1.Hash)
	c.chatRepository.AddParticipant(user2.Hash)
}

func (c *ChatService) chooseRandomPair() (*userModel.User, *userModel.User, error) {
	idx1 := rand.Intn(len(userQueue))
	idx2 := rand.Intn(len(userQueue))

	for idx2 == idx1 {
		idx2 = rand.Intn(len(userQueue))
	}

	return userQueue[idx1], userQueue[idx2], nil
}

func (c *ChatService) removeClientFromSlice(slice []*userModel.User, user *userModel.User) []*userModel.User {
	for i, c := range slice {
		if c == user {
			return append(slice[:i], slice[i+1:]...)
		}
	}

	return slice
}

func (c *ChatService) HandleStreamMessages(user1 *userModel.User, user2 *userModel.User) {
	for {
		_, textMessage, err := user1.Conn.ReadMessage()
		if err != nil {
			log.Println("Error reading message from client 1:", err)
			c.closeUserConn(user2)
			return
		}

		message := messageModel.NewMessage(
			"CHAT",
			string(textMessage),
		)

		if err := c.sendMessage(user2, message); err != nil {
			log.Println("Error sending message to client 2: ", err)
			return
		}
	}
}

func (c *ChatService) sendMessage(user *userModel.User, message *messageModel.Message) error {
	data, err := json.Marshal(message)
	if err != nil {
		log.Println("Error encoding message to JSON: ", err)
	}

	return user.Conn.WriteMessage(websocket.TextMessage, data)
}

func (c *ChatService) closeUserConn(user *userModel.User) {
	message := messageModel.NewMessage(
		"CLOSE",
		string("Собеседнкик покинул чат"),
	)

	if err := c.sendMessage(user, message); err != nil {
		log.Println("Error sending message to client 2: ", err)
	}

	user.Conn.Close()
}
