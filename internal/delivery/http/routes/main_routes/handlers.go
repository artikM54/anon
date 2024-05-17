package main_routes

import (
	messageModel "anonymous_chat/internal/models/message"
	userModel "anonymous_chat/internal/models/user"
	userService "anonymous_chat/internal/services/user"
	hashUtil "anonymous_chat/internal/utils/hash"
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
	fmt.Println("ws handler end", time.Now().Format("2006-01-02 15:04:05"))
}

func handleConn(conn *websocket.Conn) {
	fmt.Println("handler conn start", time.Now().Format("2006-01-02 15:04:05"))
	sender.NotifyConnect(conn)

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

		fmt.Println("handler conn message", message)
		switch message.Category {
		case "FRONT:GET_TOKEN":
			token := hashUtil.CreateUniqueModelHash(userModel.RedisList)
			sender.NotifyToken(conn, token)

			user := userService.NewUser(conn, token)
			userService := userService.NewUserService(user)
			go userService.HandleUsersCommand()

			return
		case "FRONT:GIVE_TOKEN":
			token := message.Payload.Text

			user := userService.NewUser(conn, token)
			userService := userService.NewUserService(user)
			go userService.HandleUsersCommand()

			return
		default:
			fmt.Println("Undefined command 1")
		}

	}
}
