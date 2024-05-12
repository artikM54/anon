package chat

import (
	"anonymous_chat/internal/redis"
	"context"
	"fmt"
)

type ChatRepository struct {
	redisListKeyParticipants string
	redisListKeyMessages     string
	ctx                      context.Context
}

func NewChatRepository(chatHash string) *ChatRepository {
	return &ChatRepository{
		redisListKeyParticipants: "chat_participants:" + chatHash,
		redisListKeyMessages:     "chat_messages:" + chatHash,
		ctx:                      context.Background(),
	}
}

func (c *ChatRepository) AddParticipant(userHash string) {

	added, err := redis.Client.SAdd(c.ctx, c.redisListKeyParticipants, userHash).Result()
	if err != nil {
		panic(err)
	}

	if added == 0 {
		fmt.Println("Error can't add a userHash to chat")
	}
}

func (c *ChatRepository) AddMessage(message string) {
	err := redis.Client.RPush(c.ctx, c.redisListKeyMessages, message).Err()

	if err != nil {
		fmt.Println("Ошибка при добавлении элемента в messages", err)
	}

	fmt.Println("Элемент успешно добавлен в список messages")
}

func (c *ChatRepository) DeleteParticipants() {
	err := redis.Client.Del(c.ctx, c.redisListKeyParticipants).Err()
	if err != nil {
		fmt.Println("Ошибка при удалении списка participants:", err)
		return
	}
	fmt.Println("Список participants успешно удален")
}

func (c *ChatRepository) DeleteMessages() {
	err := redis.Client.Del(c.ctx, c.redisListKeyMessages).Err()
	if err != nil {
		fmt.Println("Ошибка при удалении списка сообщений:", err)
		return
	}
	fmt.Println("Список сообщений успешно удален")
}

func (c *ChatRepository) GetMessages() []string {
	result, err := redis.Client.LRange(c.ctx, c.redisListKeyMessages, 0, -1).Result()
    if err != nil {
        panic(err)
    }

	return result
}