package chat

import (
	chatModel "anonymous_chat/internal/models/chat"
	messageModel "anonymous_chat/internal/models/message"
	userModel "anonymous_chat/internal/models/user"
	chatRepository "anonymous_chat/internal/repositories/chat"
	hashUtil "anonymous_chat/internal/utils/hash"
	"encoding/json"
	"fmt"
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
		Hash:    hashUtil.CreateUniqueModelHash(chatModel.RedisList),
		Users:   users,
		Channel: make(chan *messageModel.Message, 20),
	}
}

func (c *ChatService) Start() {
	fmt.Println("chat is starting")

	c.notifyChatStart()
	c.startHandlers()
}

func (c *ChatService) startHandlers() {
	go c.listening()
	go c.sending()
}

func (c *ChatService) listening() {
	for {
		fmt.Println("READ MESSAGE FOR chat: ", c.Chat.Hash)

		message, opened := <-c.Chat.Channel
		if !opened {
			fmt.Printf("readUserMessages Channel closed for user %s\n", c.Chat.Hash)
			return
		}

		message.Payload.ChatHash = c.Chat.Hash
		c.chatRepository.AddMessage(*message)
	}
}

func (c *ChatService) sending() {
	for {
		messages := c.chatRepository.GetNewMessages()
		fmt.Println("SEND MESSAGE FOR CHAT ", c.Chat.Hash)

		for _, message := range messages {

			messageData := message.Values["message"].(string)

			var message messageModel.Message
			json.Unmarshal([]byte(messageData), &message)

			fmt.Println("Read message from redis message data : ", messageData)

			for _, user := range c.Chat.Users {
				if user.Hash == message.Payload.UserHash {
					fmt.Println("CONTINUE for: ", user.Hash)
					continue
				}

				user.PutIntoChannel(&message)
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
