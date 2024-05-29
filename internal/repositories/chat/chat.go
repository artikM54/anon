package chat

import (
	messageModel "anonymous_chat/internal/models/message"
	redisUtil "anonymous_chat/internal/redis"
	"context"
	"encoding/json"
	"fmt"
)

type ChatRepository struct {
	name_list string
	ctx       context.Context
}

func NewChatRepository(chatHash string) *ChatRepository {
	return &ChatRepository{
		name_list: "list_messages:" + chatHash,
		ctx:       context.Background(),
	}
}

func (c *ChatRepository) AddMessage(message messageModel.Message) {
	data, err := json.Marshal(message)
	if err != nil {
		fmt.Println("AddMessage Redis Error encoding message to JSON: ", err)
	}

	err = redisUtil.Client.RPush(c.ctx, c.name_list, string(data)).Err()
	if err != nil {
		fmt.Println("AddMessage Redis Error: ", err)
	}
}

func (c *ChatRepository) DeleteChat() {
	err := redisUtil.Client.Del(c.ctx, c.name_list).Err()
	if err != nil {
		fmt.Println("Error delete the list:", err)
	} else {
		fmt.Println("Success delete the list")
	}
}
