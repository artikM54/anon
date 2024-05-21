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
	now := time.Now()
	unixTimestamp := now.Unix()

	message := messageModel.NewMessage(
		messageModel.ConnectCategory,
		"SUCCESS",
		unixTimestamp,
		"",
		"",
	)

	if err := SendMessage(conn, message); err != nil {
		log.Println("Error sending message to client 2: ", err)
		return
	}
}

func NotifyToken(conn *websocket.Conn, token string) {
	now := time.Now()
	unixTimestamp := now.Unix()

	message := messageModel.NewMessage(
		messageModel.TokenCategory,
		token,
		unixTimestamp,
		"",
		"",
	)

	if err := SendMessage(conn, message); err != nil {
		log.Println("Error sending message to client 2: ", err)
		return
	}
}
