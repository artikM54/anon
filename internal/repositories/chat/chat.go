package chat

import (
	messageModel "anonymous_chat/internal/models/message"
	redisUtil "anonymous_chat/internal/redis"
	"context"
	"encoding/json"
	"fmt"
	"github.com/redis/go-redis/v9"
)

type ChatRepository struct {
	redisListKeyParticipants string
	redisListKeyMessages     string
	redisStream              string
	ctx                      context.Context
}

func NewChatRepository(chatHash string) *ChatRepository {
	return &ChatRepository{
		redisListKeyParticipants: "chat_participants:" + chatHash,
		redisListKeyMessages:     "chat_messages:" + chatHash,
		redisStream:              "chat_stream:" + chatHash,
		ctx:                      context.Background(),
	}
}

func (c *ChatRepository) AddParticipant(userHash string) {

	added, err := redisUtil.Client.SAdd(c.ctx, c.redisListKeyParticipants, userHash).Result()
	if err != nil {
		panic(err)
	}

	if added == 0 {
		fmt.Println("Error can't add a userHash to chat")
	}
}

func (c *ChatRepository) DeleteParticipants() {
	err := redisUtil.Client.Del(c.ctx, c.redisListKeyParticipants).Err()
	if err != nil {
		fmt.Println("Ошибка при удалении списка participants:", err)
		return
	}
	fmt.Println("Список participants успешно удален")
}

func (c *ChatRepository) RegisterUserToStream(userHash string) {
	redisUtil.Client.XGroupCreateMkStream(c.ctx, c.redisStream, userHash, "0")
}

func (c *ChatRepository) AddMessage(message messageModel.Message) {
	data, err := json.Marshal(message)
	if err != nil {
		fmt.Println("AddMessage Redis Error encoding message to JSON: ", err)
	}

	redisUtil.Client.XAdd(c.ctx, &redis.XAddArgs{
		Stream: c.redisStream,
		Values: map[string]interface{}{
			"message": data,
		},
	})
}

func (c *ChatRepository) GetNewMessages(userHash string) []redis.XStream {
	messages, err := redisUtil.Client.XReadGroup(c.ctx, &redis.XReadGroupArgs{
		Group:    userHash,
		Consumer: "consumer1",
		Streams:  []string{c.redisStream, ">"},
		Count:    10,
	}).Result()

	if err != nil {
		fmt.Println("Error GetNewMessages:")
	}

	return messages
}

func (c *ChatRepository) PointMessageAsRead(userHash string, messageId string) {
	redisUtil.Client.XAck(c.ctx, c.redisStream, userHash, messageId)
}
