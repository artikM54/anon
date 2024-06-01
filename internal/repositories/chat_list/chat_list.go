package chat_list

import (
	"anonymous_chat/internal/memory_storages/chat_list"
	chatModel "anonymous_chat/internal/models/chat"
	messageModel "anonymous_chat/internal/models/message"
	"fmt"
)

type ChatListRepository struct {
	store *chat_list.InMemoryChatListStore
}

func NewChatListRepository() *ChatListRepository {
	return &ChatListRepository{
		store: chat_list.GetInstance(),
	}
}

func (r *ChatListRepository) Add(chat *chatModel.Chat) {
	r.store.Add(chat)
}

func (r *ChatListRepository) Get(hash string) *chatModel.Chat {
	chat, _ := r.store.Get(hash)
	return chat
}

func (r *ChatListRepository) Delete(hash string) {
	r.store.Delete(hash)
}

func (r *ChatListRepository) ExistChat(hash string) bool {
	_, found := r.store.Get(hash)

	return found
}

func (r *ChatListRepository) PutMessageIntoChat(message *messageModel.Message) {
	chat, found := r.store.Get(message.Payload.ChatHash)

	if !found {
		fmt.Println("Error repository Put message not found")
	}

	chat.Channel <- message
}

func (r *ChatListRepository) ExitUserFromChat(chatHash string, userHash string) {
	chat, found := r.store.Get(chatHash)

	if !found {
		fmt.Println("Error repository exit user not found")
	}

	if chat.ExistUser(userHash) {
		delete(chat.Users, userHash)
	}
}
