package user

import (
	messageModel "anonymous_chat/internal/models/message"
)

const (
	RedisList = "unique_users"
)

type User struct {
	Hash        string
	channel     chan *messageModel.Message
	stopChannel chan struct{}
}

func NewUser() *User {
	return &User{
		channel:     make(chan *messageModel.Message),
		stopChannel: make(chan struct{}),
	}
}

func (u *User) SetToken(hash string) {
	u.Hash = hash
}

func (u *User) PutIntoChannel(message *messageModel.Message) {
	u.channel <- message
}

func (u *User) GetFromChannel() (*messageModel.Message, bool) {
	message, closed := <-u.channel
	return message, closed
}

func (u *User) CloseChannel() {
	close(u.stopChannel)
	close(u.channel)
}


func (u *User) GetState() chan struct{} {
	return u.stopChannel
}
