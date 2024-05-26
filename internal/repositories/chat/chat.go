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
	lastMessageId            string
	ctx                      context.Context
}

func NewChatRepository(chatHash string) *ChatRepository {
	return &ChatRepository{
		redisListKeyParticipants: "chat_participants:" + chatHash,
		redisListKeyMessages:     "chat_messages:" + chatHash,
		redisStream:              "chat_stream:" + chatHash,
		lastMessageId:            "0",
		ctx:                      context.Background(),
	}
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

func (c *ChatRepository) GetNewMessages() []redis.XMessage {
	streams, err := redisUtil.Client.XRead(c.ctx, &redis.XReadArgs{
		Streams: []string{c.redisStream, c.lastMessageId},
		Count:   10,
	}).Result()

	if err != nil {
		fmt.Println("Error reading new messages:", err)
	}

	messages := streams[0].Messages

	if len(messages) > 0 {
		i := len(messages) - 1
		c.updateLastMessageId(messages[i])
	}

	return messages
}

func (c *ChatRepository) updateLastMessageId(message redis.XMessage) {
	c.lastMessageId = message.ID
	fmt.Println("Updated lastMessageId value:", c.lastMessageId)

}
