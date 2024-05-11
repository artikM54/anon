package user

import (
	userModel "anonymous_chat/internal/models/user"
	hashUtil "anonymous_chat/internal/utils/hash"

	"github.com/gorilla/websocket"
)

func NewUser(conn *websocket.Conn, hash string) *userModel.User {
	if hash == "" {
        hash = hashUtil.CreateUniqueModelHash(userModel.RedisList)
    }
	
	return &userModel.User{
		Hash: hash,
		Conn: conn,
	}
}
