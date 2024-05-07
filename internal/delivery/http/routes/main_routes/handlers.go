package main_routes

import (
	userModel "anonymous_chat/internal/models/user"
	chatService "anonymous_chat/internal/services/chat"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
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
	fmt.Println("try to connect")
	conn, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		log.Println("Error upgrading to WebSocket:", err)
		return
	}

	user := userModel.NewUser(conn)
	chatService.TryCreateChat(user)
}
