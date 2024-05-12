package main_routes

import (
	"anonymous_chat/internal/handler_queue"
	messageModel "anonymous_chat/internal/models/message"
	userService "anonymous_chat/internal/services/user"
	"anonymous_chat/internal/utils/sender"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"time"
)

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			// Allow all origins
			return true
		},
	}
)

func homeHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Тест"))
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("try to connect at", time.Now().Format("2006-01-02 15:04:05"))
	conn, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		log.Println("Error upgrading to WebSocket:", err)
		return
	}

	go handleConn(conn)
}

func handleConn(conn *websocket.Conn) {
	message := messageModel.NewMessage(
		"CONNECT",
		"SUCCESS",
		time.Now().Format("2006-01-02 15:04:05"),
		"",
		"",
	)

	if err := sender.SendMessage(conn, message); err != nil {
		log.Println("Error sending message to client 2: ", err)
		return
	}

	for {
		_, data, err := conn.ReadMessage()
		if err != nil {
			log.Println("Error reading:", err)
			return
		}

		var message messageModel.Message

		err = json.Unmarshal(data, &message)
		if err != nil {
			log.Println("Error reading message from client 1:", err)
		}

		fmt.Println(message)

		switch message.Category {
		case "FRONT:GET_TOKEN":
			fmt.Println("FRONT:GET_TOKEN")
			user := userService.NewUser(conn, "")
			sender.NotifyToken(user)
			handler_queue.AddUserToQueue(user)
			return
		case "FRONT:GIVE_TOKEN":
			fmt.Println("FRONT:GIVE_TOKEN")
			user := userService.NewUser(conn, message.Payload.Text)
			handler_queue.AddUserToQueue(user)
			return

		default:
			fmt.Println("Undefined command")
		}

	}
}
