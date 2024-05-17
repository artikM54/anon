package sender

import (
	messageModel "anonymous_chat/internal/models/message"
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

func NotifyConnect(conn *websocket.Conn) {
	message := messageModel.NewMessage(
		"CONNECT",
		"SUCCESS",
		time.Now().Format("2006-01-02 15:04:05"),
		"",
		"",
	)

	if err := SendMessage(conn, message); err != nil {
		log.Println("Error sending message to client 2: ", err)
		return
	}
}

func NotifyToken(conn *websocket.Conn, token string) {
	message := messageModel.NewMessage(
		"TOKEN",
		token,
		time.Now().Format("2006-01-02 15:04:05"),
		"",
		"",
	)

	if err := SendMessage(conn, message); err != nil {
		log.Println("Error sending message to client 2: ", err)
		return
	}
}
