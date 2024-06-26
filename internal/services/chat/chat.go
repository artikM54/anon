package chat

import (
	chatModel "anonymous_chat/internal/models/chat"
	messageModel "anonymous_chat/internal/models/message"
	userModel "anonymous_chat/internal/models/user"
	chatRepository "anonymous_chat/internal/repositories/chat"
	"anonymous_chat/internal/repositories/chat_list"
	hashUtil "anonymous_chat/internal/utils/hash"
	"fmt"
	"sync"
)

type ChatService struct {
	Chat               *chatModel.Chat
	chatRepository     *chatRepository.ChatRepository
	chatListRepository *chat_list.ChatListRepository
	wg                 sync.WaitGroup
}

func NewChatService(users map[string]*userModel.User) *ChatService {
	chat := newChat(users)

	return &ChatService{
		Chat:               chat,
		chatRepository:     chatRepository.NewChatRepository(chat.Hash),
		chatListRepository: chat_list.NewChatListRepository(),
	}
}

func newChat(users map[string]*userModel.User) *chatModel.Chat {
	return &chatModel.Chat{
		Hash:    hashUtil.CreateUniqueModelHash(chatModel.RedisList),
		Users:   users,
		Channel: make(chan *messageModel.Message, 20),
	}
}

func (c *ChatService) Start() {
	fmt.Println("chat is starting")
	c.chatListRepository.Add(c.Chat)

	c.notifyChatStart()

	c.wg.Add(1)
	go c.handler()
	c.wg.Wait()

	c.chatRepository.DeleteChat()
	c.chatListRepository.Delete(c.Chat.Hash)
}

func (c *ChatService) handler() {
	defer c.wg.Done()

	for message := range c.Chat.Channel {
		fmt.Println("READ MESSAGE FOR chat: ", c.Chat.Hash)
		if c.Chat.IsEmpty() {
			fmt.Println("close chat channel")

			close(c.Chat.Channel)
		}

		message.Payload.ChatHash = c.Chat.Hash
		c.chatRepository.AddMessage(*message)

		for _, user := range c.Chat.Users {
			if user.Hash == message.Payload.UserHash {
				fmt.Println("CONTINUE for: ", user.Hash)
				continue
			}

			user.PutIntoChannel(message)
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

	c.Chat.Channel <- message
}
