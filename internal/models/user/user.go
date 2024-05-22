package user

import (
	messageModel "anonymous_chat/internal/models/message"
)

const (
	RedisList = "unique_users"
)

type UserChat struct {
	channel chan *messageModel.Message
	stop    chan struct{}
}

type User struct {
	Hash        string
	channel     chan *messageModel.Message
	stopChannel chan struct{}
	chats       map[string]*UserChat
}

func NewUser() *User {
	return &User{
		channel:     make(chan *messageModel.Message),
		stopChannel: make(chan struct{}),
		chats:       make(map[string]*UserChat),
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

func (u *User) AddChat(hash string) {
	userChat := UserChat{
		channel: make(chan *messageModel.Message),
		stop:    make(chan struct{}),
	}

	u.chats[hash] = &userChat
}

func (u *User) PutIntoChat(message *messageModel.Message) {
	u.chats[message.Payload.ChatHash].channel <- message
}

func (u *User) GetChatState(hash string) (chan struct{}, bool) {
	found := u.ExitChat(hash)

	return u.chats[hash].stop, found
}

func (u *User) GetState() chan struct{} {
	return u.stopChannel
}

func (u *User) GetFromChat(hash string) (*messageModel.Message, bool) {
	message, closed := <-u.chats[hash].channel
	return message, closed
}

func (u *User) CloseChat(hash string) {
	close(u.chats[hash].channel)
	close(u.chats[hash].stop)
}

func (u *User) ExitChat(hash string) bool {
	_, found := u.chats[hash]

	return found
}
