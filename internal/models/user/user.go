package user

import (
	"github.com/gorilla/websocket"
)

const(
    RedisList = "unique_users"
)

type User struct {
	Hash string
	Conn *websocket.Conn
}
