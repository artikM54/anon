package sender

import (
	messageModel "anonymous_chat/internal/models/message"
	userModel "anonymous_chat/internal/models/user"
	"encoding/json"
	"github.com/gorilla/websocket"
	"log"
	"time"
)

func SendMessage(conn *websocket.Conn, message *messageModel.Message) error {
	data, err := json.Marshal(message)
	if err != nil {
		log.Println("Error encoding message to JSON: ", err)
	}

	return conn.WriteMessage(websocket.TextMessage, data)
}

func NotifyToken(user *userModel.User) {
	message := messageModel.NewMessage(
		"TOKEN",
		user.Hash,
		time.Now().Format("2006-01-02 15:04:05"),
		"",
		"",
	)

	if err := SendMessage(user.Conn, message); err != nil {
		log.Println("Error sending message to client 2: ", err)
		return
	}
}
