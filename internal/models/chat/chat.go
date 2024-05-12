package chat

import (
	userModel "anonymous_chat/internal/models/user"
)

const (
	RedisList = "unique_chats"
)

type Chat struct {
	Hash  string
	Users []*userModel.User
}
