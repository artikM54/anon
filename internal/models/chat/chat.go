package chat

import (
	userModel "anonymous_chat/internal/models/user"
	messageModel "anonymous_chat/internal/models/message"
)

const (
	RedisList = "unique_chats"
)

type Chat struct {
	Hash  string
	Users []*userModel.User
	Channel chan *messageModel.Message
}
