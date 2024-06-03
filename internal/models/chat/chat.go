package chat

import (
	messageModel "anonymous_chat/internal/models/message"
	userModel "anonymous_chat/internal/models/user"
	"sync"
)

const (
	RedisList = "unique_chats"
)

type Chat struct {
	Hash    string
	Users   map[string]*userModel.User
	Channel chan *messageModel.Message
	Mu      sync.Mutex
}

func (c *Chat) IsEmpty() bool {
	var result bool

	if len(c.Users) < 2 {
		result = true
	} else {
		result = false
	}

	return result
}

func (c *Chat) ExistUser(userHash string) bool {
	_, found := c.Users[userHash]

	return found
}
