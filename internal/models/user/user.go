package user

import (
	"github.com/gorilla/websocket"
)

type User struct {
	Conn *websocket.Conn
}

func NewUser(conn *websocket.Conn) *User {
	return &User{Conn: conn}
}
