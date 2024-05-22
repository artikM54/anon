package user

import (
	messageModel "anonymous_chat/internal/models/message"
)

const (
	RedisList = "unique_users"
)

type UserChat struct {
	in   chan *messageModel.Message
	stop chan struct{}
}

type User struct {
	Hash       string
	outChannel chan *messageModel.Message
	chats      map[string]*UserChat
}

func NewUser() *User {
	return &User{
		outChannel: make(chan *messageModel.Message),
		chats:      make(map[string]*UserChat),
	}
}

func (u *User) SetToken(hash string) {
	u.Hash = hash
}

func (u *User) PutToOutChannel(message *messageModel.Message) {
	u.outChannel <- message
}

func (u *User) GetFromOutChannel() (*messageModel.Message, bool) {
	message, closed := <-u.outChannel
	return message, closed
}

func (u *User) CloseOutChannel() {
	close(u.outChannel)
}

func (u *User) AddChat(hash string) {
	userChat := UserChat{
		in:   make(chan *messageModel.Message),
		stop: make(chan struct{}),
	}

	u.chats[hash] = &userChat
}

func (u *User) PutToInChat(message *messageModel.Message) {
	u.chats[message.Payload.ChatHash].in <- message
}

func (u *User) GetChatState(hash string) chan struct{} {
	return u.chats[hash].stop
}

func (u *User) GetFromInChat(hash string) (*messageModel.Message, bool) {
	message, closed := <-u.chats[hash].in
	return message, closed
}

func (u *User) CloseChat(hash string) {
	close(u.chats[hash].in)
	close(u.chats[hash].stop)
}

func (u *User) ExitChat(hash string) bool {
	_, found := u.chats[hash]

	return found
}
