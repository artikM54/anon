package user

import (
	messageModel "anonymous_chat/internal/models/message"
	hashUtil "anonymous_chat/internal/utils/hash"
	"github.com/gorilla/websocket"
	"sync"
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
	Conn       *websocket.Conn
	outChannel chan *messageModel.Message
	chats      map[string]*UserChat
	mutex      sync.Mutex
}

func NewUser(conn *websocket.Conn, hash string) *User {
	if hash == "" {
		hash = hashUtil.CreateUniqueModelHash(RedisList)
	}

	return &User{
		Hash:       hash,
		Conn:       conn,
		outChannel: make(chan *messageModel.Message),
		chats:      make(map[string]*UserChat),
	}
}

func (u *User) ReadMessage() ([]byte, error) {
	u.mutex.Lock()
	defer u.mutex.Unlock()
	_, data, err := u.Conn.ReadMessage()
	return data, err
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
