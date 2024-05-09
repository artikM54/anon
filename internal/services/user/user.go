package user

import (
	userModel "anonymous_chat/internal/models/user"
	hashUtil "anonymous_chat/internal/utils/hash"

	"github.com/gorilla/websocket"
)

func NewUser(conn *websocket.Conn) *userModel.User {
	return &userModel.User{
		Hash: hashUtil.CreateUniqueModelHash(userModel.RedisList),
		Conn: conn,
	}
}
