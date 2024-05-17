package user

import (
    messageModel "anonymous_chat/internal/models/message"
	"github.com/gorilla/websocket"
	"sync"
)

const (
	RedisList = "unique_users"
)

type User struct {
	Hash   string
	Conn   *websocket.Conn
    ChannelChat chan *messageModel.Message 
	mutex     sync.Mutex
}

func (u *User) ReadMessage() ([]byte, error) {
	u.mutex.Lock()
	defer u.mutex.Unlock()
	_, data, err := u.Conn.ReadMessage()
	return data, err
}
